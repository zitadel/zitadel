package idp

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	z_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/form"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/azuread"
	"github.com/zitadel/zitadel/internal/idp/providers/github"
	"github.com/zitadel/zitadel/internal/idp/providers/gitlab"
	"github.com/zitadel/zitadel/internal/idp/providers/google"
	"github.com/zitadel/zitadel/internal/idp/providers/jwt"
	"github.com/zitadel/zitadel/internal/idp/providers/ldap"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	"github.com/zitadel/zitadel/internal/query"
)

const (
	HandlerPrefix    = "/idps"
	callbackPath     = "/callback"
	ldapCallbackPath = callbackPath + "/ldap"

	paramIntentID         = "id"
	paramToken            = "token"
	paramUserID           = "user"
	paramError            = "error"
	paramErrorDescription = "error_description"
)

type Handler struct {
	commands            *command.Commands
	queries             *query.Queries
	parser              *form.Parser
	encryptionAlgorithm crypto.EncryptionAlgorithm
	callbackURL         func(ctx context.Context) string
}

type externalIDPCallbackData struct {
	State            string `schema:"state"`
	Code             string `schema:"code"`
	Error            string `schema:"error"`
	ErrorDescription string `schema:"error_description"`
}

// CallbackURL generates the instance specific URL to the IDP callback handler
func CallbackURL(externalSecure bool) func(ctx context.Context) string {
	return func(ctx context.Context) string {
		return http_utils.BuildOrigin(authz.GetInstance(ctx).RequestedHost(), externalSecure) + HandlerPrefix + callbackPath
	}
}

func NewHandler(
	commands *command.Commands,
	queries *query.Queries,
	encryptionAlgorithm crypto.EncryptionAlgorithm,
	externalSecure bool,
	instanceInterceptor func(next http.Handler) http.Handler,
) http.Handler {
	h := &Handler{
		commands:            commands,
		queries:             queries,
		parser:              form.NewParser(),
		encryptionAlgorithm: encryptionAlgorithm,
		callbackURL:         CallbackURL(externalSecure),
	}

	router := mux.NewRouter()
	router.Use(instanceInterceptor)
	router.HandleFunc(callbackPath, h.handleCallback)
	return router
}

func (h *Handler) handleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data, err := h.parseCallbackRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	intent, err := h.commands.GetActiveIntent(ctx, data.State)
	if err != nil {
		if z_errs.IsNotFound(err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		redirectToFailureURLErr(w, r, intent, err)
		return
	}

	// the provider might have returned an error
	if data.Error != "" {
		cmdErr := h.commands.FailIDPIntent(ctx, intent, reason(data.Error, data.ErrorDescription))
		logging.WithFields("intent", intent.AggregateID).OnError(cmdErr).Error("failed to push failed event on idp intent")
		redirectToFailureURL(w, r, intent, data.Error, data.ErrorDescription)
		return
	}

	provider, err := h.commands.GetProvider(ctx, intent.IDPID, h.callbackURL(ctx))
	if err != nil {
		cmdErr := h.commands.FailIDPIntent(ctx, intent, err.Error())
		logging.WithFields("intent", intent.AggregateID).OnError(cmdErr).Error("failed to push failed event on idp intent")
		redirectToFailureURLErr(w, r, intent, err)
		return
	}

	idpUser, idpSession, err := h.fetchIDPUser(ctx, provider, data.Code)
	if err != nil {
		cmdErr := h.commands.FailIDPIntent(ctx, intent, err.Error())
		logging.WithFields("intent", intent.AggregateID).OnError(cmdErr).Error("failed to push failed event on idp intent")
		redirectToFailureURLErr(w, r, intent, err)
		return
	}
	userID, err := h.checkExternalUser(ctx, intent.IDPID, idpUser.GetID())
	logging.WithFields("intent", intent.AggregateID).OnError(err).Error("could not check if idp user already exists")

	if userID == "" {
		userID, err = h.tryMigrateExternalUser(ctx, intent.IDPID, idpUser, idpSession)
		logging.WithFields("intent", intent.AggregateID).OnError(err).Error("migration check failed")
	}

	token, err := h.commands.SucceedIDPIntent(ctx, intent, idpUser, idpSession, userID)
	if err != nil {
		redirectToFailureURLErr(w, r, intent, z_errs.ThrowInternal(err, "IDP-JdD3g", "Errors.Intent.TokenCreationFailed"))
		return
	}
	redirectToSuccessURL(w, r, intent, token, userID)
}

func (h *Handler) tryMigrateExternalUser(ctx context.Context, idpID string, idpUser idp.User, idpSession idp.Session) (userID string, err error) {
	migration, ok := idpSession.(idp.SessionSupportsMigration)
	if !ok {
		return "", nil
	}
	previousID, err := migration.RetrievePreviousID()
	if err != nil || previousID == "" {
		return "", err
	}
	userID, err = h.checkExternalUser(ctx, idpID, previousID)
	if err != nil {
		return "", err
	}
	return userID, h.commands.MigrateUserIDP(ctx, userID, "", idpID, previousID, idpUser.GetID())
}

func (h *Handler) parseCallbackRequest(r *http.Request) (*externalIDPCallbackData, error) {
	data := new(externalIDPCallbackData)
	err := h.parser.Parse(r, data)
	if err != nil {
		return nil, err
	}
	if data.State == "" {
		return nil, z_errs.ThrowInvalidArgument(nil, "IDP-Hk38e", "Errors.Intent.StateMissing")
	}
	return data, nil
}

func (h *Handler) getActiveIntent(w http.ResponseWriter, r *http.Request, state string) *command.IDPIntentWriteModel {
	intent, err := h.commands.GetIntentWriteModel(r.Context(), state, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	if intent.State == domain.IDPIntentStateUnspecified {
		http.Error(w, reason("IDP-Hk38e", "Errors.Intent.NotStarted"), http.StatusBadRequest)
		return nil
	}
	if intent.State != domain.IDPIntentStateStarted {
		redirectToFailureURL(w, r, intent, "IDP-Sfrgs", "Errors.Intent.NotStarted")
		return nil
	}
	return intent
}

func redirectToSuccessURL(w http.ResponseWriter, r *http.Request, intent *command.IDPIntentWriteModel, token, userID string) {
	queries := intent.SuccessURL.Query()
	queries.Set(paramIntentID, intent.AggregateID)
	queries.Set(paramToken, token)
	if userID != "" {
		queries.Set(paramUserID, userID)
	}
	intent.SuccessURL.RawQuery = queries.Encode()
	http.Redirect(w, r, intent.SuccessURL.String(), http.StatusFound)
}

func redirectToFailureURLErr(w http.ResponseWriter, r *http.Request, i *command.IDPIntentWriteModel, err error) {
	msg := err.Error()
	var description string
	zErr := new(z_errs.CaosError)
	if errors.As(err, &zErr) {
		msg = zErr.GetID()
		description = zErr.GetMessage() // TODO: i18n?
	}
	redirectToFailureURL(w, r, i, msg, description)
}

func redirectToFailureURL(w http.ResponseWriter, r *http.Request, i *command.IDPIntentWriteModel, err, description string) {
	queries := i.FailureURL.Query()
	queries.Set(paramIntentID, i.AggregateID)
	queries.Set(paramError, err)
	queries.Set(paramErrorDescription, description)
	i.FailureURL.RawQuery = queries.Encode()
	http.Redirect(w, r, i.FailureURL.String(), http.StatusFound)
}

func (h *Handler) fetchIDPUser(ctx context.Context, identityProvider idp.Provider, code string) (user idp.User, idpTokens idp.Session, err error) {
	var session idp.Session
	switch provider := identityProvider.(type) {
	case *oauth.Provider:
		session = &oauth.Session{Provider: provider, Code: code}
	case *openid.Provider:
		session = &openid.Session{Provider: provider, Code: code}
	case *azuread.Provider:
		session = &azuread.Session{Session: &oauth.Session{Provider: provider.Provider, Code: code}}
	case *github.Provider:
		session = &oauth.Session{Provider: provider.Provider, Code: code}
	case *gitlab.Provider:
		session = &openid.Session{Provider: provider.Provider, Code: code}
	case *google.Provider:
		session = &openid.Session{Provider: provider.Provider, Code: code}
	case *jwt.Provider, *ldap.Provider:
		return nil, nil, z_errs.ThrowInvalidArgument(nil, "IDP-52jmn", "Errors.ExternalIDP.IDPTypeNotImplemented")
	default:
		return nil, nil, z_errs.ThrowUnimplemented(nil, "IDP-SSDg", "Errors.ExternalIDP.IDPTypeNotImplemented")
	}

	user, err = session.FetchUser(ctx)
	if err != nil {
		return nil, nil, err
	}
	return user, session, nil
}

func (h *Handler) checkExternalUser(ctx context.Context, idpID, externalUserID string) (userID string, err error) {
	idQuery, err := query.NewIDPUserLinkIDPIDSearchQuery(idpID)
	if err != nil {
		return "", err
	}
	externalIDQuery, err := query.NewIDPUserLinksExternalIDSearchQuery(externalUserID)
	if err != nil {
		return "", err
	}
	queries := []query.SearchQuery{
		idQuery, externalIDQuery,
	}
	links, err := h.queries.IDPUserLinks(ctx, &query.IDPUserLinksSearchQuery{Queries: queries}, false)
	if err != nil {
		return "", err
	}
	if len(links.Links) != 1 {
		return "", nil
	}
	return links.Links[0].UserID, nil
}

func reason(err, description string) string {
	if description == "" {
		return err
	}
	return err + ": " + description
}
