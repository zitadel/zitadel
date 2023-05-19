package idp

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/form"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/azuread"
	"github.com/zitadel/zitadel/internal/idp/providers/github"
	"github.com/zitadel/zitadel/internal/idp/providers/gitlab"
	"github.com/zitadel/zitadel/internal/idp/providers/google"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	"github.com/zitadel/zitadel/internal/query"
)

const (
	HandlerPrefix = "/idps"
	callbackPath  = "/callback"

	paramIntentID = "id"
	paramToken    = "token"
	paramUserID   = "user"
	paramError    = "error"
)

type Handler struct {
	commands            *command.Commands
	queries             *query.Queries
	parser              *form.Parser
	encryptionAlgorithm crypto.EncryptionAlgorithm
	callbackURL         func(ctx context.Context) string
}

type externalIDPCallbackData struct {
	State string `schema:"state"`
	Code  string `schema:"code"`
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
	data := new(externalIDPCallbackData)
	err := h.parser.Parse(r, data)
	if err != nil {
		// TODO: ?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	// get intent
	intent, err := h.commands.GetIntentWriteModel(ctx, data.State, "")
	if err != nil {
		// TODO: ?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if intent.State != domain.IDPIntentStateStarted {
		redirectToFailureURL(w, r, intent, errors.ThrowPreconditionFailed(nil, "IDP-Sfrgs", "Errors.Intent.NotStarted"))
		return
	}
	// get idp
	idpTemplate, err := h.queries.IDPTemplateByID(ctx, false, intent.IDPID, false)
	if err != nil {
		// TODO: set failed?
		redirectToFailureURL(w, r, intent, err)
		return
	}

	// fetch user
	idpUser, err := h.fetchIDPUser(w, r, idpTemplate, data.Code)
	if err != nil {
		// on error -> redirect to failed url
		// TODO: set failed?
		redirectToFailureURL(w, r, intent, err)
		return
	}
	// on success
	userID, err := h.checkExternalUser(ctx, idpTemplate.ID, idpUser.GetID())
	if err != nil {
		// TODO: ?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// create intent token
	token, err := h.commands.SucceedIDPIntent(ctx, intent, idpUser, userID)
	if err != nil {
		// TODO: ?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// redirect to success url with internal user id, intent id, intent token
	redirectToSuccessURL(w, r, intent, token, userID)
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

func redirectToFailureURL(w http.ResponseWriter, r *http.Request, i *command.IDPIntentWriteModel, err error) {
	queries := i.FailureURL.Query()
	queries.Set(paramIntentID, i.AggregateID)
	queries.Set(paramError, err.Error())
	i.FailureURL.RawQuery = queries.Encode()
	http.Redirect(w, r, i.FailureURL.String(), http.StatusFound)
}

func (h *Handler) fetchIDPUser(w http.ResponseWriter, r *http.Request, identityProvider *query.IDPTemplate, code string) (user idp.User, err error) {
	var provider idp.Provider
	var session idp.Session
	callback := h.callbackURL(r.Context())
	switch identityProvider.Type {
	case domain.IDPTypeOAuth:
		provider, err = oauth.NewFromQueryTemplate(identityProvider, callback, h.encryptionAlgorithm)
		if err != nil {
			return nil, err
		}
		session = &oauth.Session{Provider: provider.(*oauth.Provider), Code: code}
	case domain.IDPTypeOIDC:
		provider, err = openid.NewFromQueryTemplate(identityProvider, callback, h.encryptionAlgorithm)
		if err != nil {
			return nil, err
		}
		session = &openid.Session{Provider: provider.(*openid.Provider), Code: code}
	case domain.IDPTypeAzureAD:
		provider, err = azuread.NewFromQueryTemplate(identityProvider, callback, h.encryptionAlgorithm)
		if err != nil {
			return nil, err
		}
		session = &oauth.Session{Provider: provider.(*azuread.Provider).Provider, Code: code}
	case domain.IDPTypeGitHub:
		provider, err = github.NewFromQueryTemplate(identityProvider, callback, h.encryptionAlgorithm)
		if err != nil {
			return nil, err
		}
		session = &oauth.Session{Provider: provider.(*github.Provider).Provider, Code: code}
	case domain.IDPTypeGitHubEnterprise:
		provider, err = github.NewCustomFromQueryTemplate(identityProvider, callback, h.encryptionAlgorithm)
		if err != nil {
			return nil, err
		}
		session = &oauth.Session{Provider: provider.(*github.Provider).Provider, Code: code}
	case domain.IDPTypeGitLab:
		provider, err = gitlab.NewFromQueryTemplate(identityProvider, callback, h.encryptionAlgorithm)
		if err != nil {
			return nil, err
		}
		session = &openid.Session{Provider: provider.(*gitlab.Provider).Provider, Code: code}
	case domain.IDPTypeGitLabSelfHosted:
		provider, err = gitlab.NewCustomFromQueryTemplate(identityProvider, callback, h.encryptionAlgorithm)
		if err != nil {
			return nil, err
		}
		session = &openid.Session{Provider: provider.(*gitlab.Provider).Provider, Code: code}
	case domain.IDPTypeGoogle:
		provider, err = google.NewFromQueryTemplate(identityProvider, callback, h.encryptionAlgorithm)
		if err != nil {
			return nil, err
		}
		session = &openid.Session{Provider: provider.(*google.Provider).Provider, Code: code}
	case domain.IDPTypeJWT,
		domain.IDPTypeLDAP,
		domain.IDPTypeUnspecified:
		fallthrough
	default:
		return nil, errors.ThrowInvalidArgument(nil, "IDP-SSDg", "Errors.ExternalIDP.IDPTypeNotImplemented")
	}

	return session.FetchUser(r.Context())
}

func CallbackURL(externalSecure bool) func(ctx context.Context) string {
	return func(ctx context.Context) string {
		return http_utils.BuildOrigin(authz.GetInstance(ctx).RequestedHost(), externalSecure) + HandlerPrefix + callbackPath
	}
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
