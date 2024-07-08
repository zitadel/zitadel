package idp

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/form"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/apple"
	"github.com/zitadel/zitadel/internal/idp/providers/azuread"
	"github.com/zitadel/zitadel/internal/idp/providers/github"
	"github.com/zitadel/zitadel/internal/idp/providers/gitlab"
	"github.com/zitadel/zitadel/internal/idp/providers/google"
	"github.com/zitadel/zitadel/internal/idp/providers/jwt"
	"github.com/zitadel/zitadel/internal/idp/providers/ldap"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	saml2 "github.com/zitadel/zitadel/internal/idp/providers/saml"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	HandlerPrefix = "/idps"

	idpPrefix = "/{" + varIDPID + ":[0-9]+|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}"

	callbackPath    = "/callback"
	metadataPath    = idpPrefix + "/saml/metadata"
	acsPath         = idpPrefix + "/saml/acs"
	certificatePath = idpPrefix + "/saml/certificate"

	paramIntentID         = "id"
	paramToken            = "token"
	paramUserID           = "user"
	paramError            = "error"
	paramErrorDescription = "error_description"
	varIDPID              = "idpid"
)

type Handler struct {
	commands                *command.Commands
	queries                 *query.Queries
	parser                  *form.Parser
	encryptionAlgorithm     crypto.EncryptionAlgorithm
	callbackURL             func(ctx context.Context) string
	samlRootURL             func(ctx context.Context, idpID string) string
	loginUICallbackRedirect func(w http.ResponseWriter, r *http.Request, state string) bool
}

type externalIDPCallbackData struct {
	State            string `schema:"state"`
	Code             string `schema:"code"`
	Error            string `schema:"error"`
	ErrorDescription string `schema:"error_description"`

	// Apple returns a user on first registration
	User string `schema:"user"`
}

type externalSAMLIDPCallbackData struct {
	IDPID      string
	Response   string
	RelayState string
}

// CallbackURL generates the instance specific URL to the IDP callback handler
func CallbackURL(externalSecure bool) func(ctx context.Context) string {
	return func(ctx context.Context) string {
		return http_utils.BuildOrigin(authz.GetInstance(ctx).RequestedHost(), externalSecure) + HandlerPrefix + callbackPath
	}
}

func SAMLRootURL(externalSecure bool) func(ctx context.Context, idpID string) string {
	return func(ctx context.Context, idpID string) string {
		return http_utils.BuildOrigin(authz.GetInstance(ctx).RequestedHost(), externalSecure) + HandlerPrefix + "/" + idpID + "/"
	}
}

func NewHandler(
	commands *command.Commands,
	queries *query.Queries,
	encryptionAlgorithm crypto.EncryptionAlgorithm,
	externalSecure bool,
	instanceInterceptor func(next http.Handler) http.Handler,
	loginUICallbackRedirect func(w http.ResponseWriter, r *http.Request, state string) bool,
) http.Handler {
	h := &Handler{
		commands:                commands,
		queries:                 queries,
		parser:                  form.NewParser(),
		encryptionAlgorithm:     encryptionAlgorithm,
		callbackURL:             CallbackURL(externalSecure),
		loginUICallbackRedirect: loginUICallbackRedirect,
		samlRootURL:             SAMLRootURL(externalSecure),
	}

	router := mux.NewRouter()
	router.Use(instanceInterceptor)
	router.HandleFunc(callbackPath, h.handleCallback)
	router.HandleFunc(metadataPath, h.handleMetadata)
	router.HandleFunc(certificatePath, h.handleCertificate)
	router.HandleFunc(acsPath, h.handleACS)
	return router
}

func parseSAMLRequest(r *http.Request) *externalSAMLIDPCallbackData {
	vars := mux.Vars(r)
	return &externalSAMLIDPCallbackData{
		IDPID:      vars[varIDPID],
		Response:   r.FormValue("SAMLResponse"),
		RelayState: r.FormValue("RelayState"),
	}
}

func (h *Handler) getProvider(ctx context.Context, idpID string) (idp.Provider, error) {
	return h.commands.GetProvider(ctx, idpID, h.callbackURL(ctx), h.samlRootURL(ctx, idpID))
}

func (h *Handler) handleCertificate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := parseSAMLRequest(r)

	provider, err := h.getProvider(ctx, data.IDPID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	samlProvider, ok := provider.(*saml2.Provider)
	if !ok {
		http.Error(w, zerrors.ThrowInvalidArgument(nil, "SAML-lrud8s9coi", "Errors.Intent.IDPInvalid").Error(), http.StatusBadRequest)
		return
	}

	certPem := new(bytes.Buffer)
	if _, err := certPem.Write(samlProvider.Certificate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=idp.crt")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	_, err = io.Copy(w, certPem)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to response with certificate: %w", err).Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) handleMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := parseSAMLRequest(r)

	provider, err := h.getProvider(ctx, data.IDPID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	samlProvider, ok := provider.(*saml2.Provider)
	if !ok {
		http.Error(w, zerrors.ThrowInvalidArgument(nil, "SAML-lrud8s9coi", "Errors.Intent.IDPInvalid").Error(), http.StatusBadRequest)
		return
	}

	sp, err := samlProvider.GetSP()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metadata := sp.ServiceProvider.Metadata()

	buf, _ := xml.MarshalIndent(metadata, "", "  ")
	w.Header().Set("Content-Type", "application/samlmetadata+xml")
	_, err = w.Write(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handler) handleACS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := parseSAMLRequest(r)

	if h.loginUICallbackRedirect(w, r, data.RelayState) {
		return
	}
	provider, err := h.getProvider(ctx, data.IDPID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	samlProvider, ok := provider.(*saml2.Provider)
	if !ok {
		err := zerrors.ThrowInvalidArgument(nil, "SAML-ui9wyux0hp", "Errors.Intent.IDPInvalid")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	intent, err := h.commands.GetActiveIntent(ctx, data.RelayState)
	if err != nil {
		if zerrors.IsNotFound(err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		redirectToFailureURLErr(w, r, intent, err)
		return
	}

	session, err := saml2.NewSession(samlProvider, intent.RequestID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idpUser, err := session.FetchUser(r.Context())
	if err != nil {
		cmdErr := h.commands.FailIDPIntent(ctx, intent, err.Error())
		logging.WithFields("intent", intent.AggregateID).OnError(cmdErr).Error("failed to push failed event on idp intent")
		redirectToFailureURLErr(w, r, intent, err)
		return
	}

	userID, err := h.checkExternalUser(ctx, intent.IDPID, idpUser.GetID())
	logging.WithFields("intent", intent.AggregateID).OnError(err).Error("could not check if idp user already exists")

	token, err := h.commands.SucceedSAMLIDPIntent(ctx, intent, idpUser, userID, session.Assertion)
	if err != nil {
		redirectToFailureURLErr(w, r, intent, zerrors.ThrowInternal(err, "IDP-JdD3g", "Errors.Intent.TokenCreationFailed"))
		return
	}
	redirectToSuccessURL(w, r, intent, token, userID)
}

func (h *Handler) handleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data, err := h.parseCallbackRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if h.loginUICallbackRedirect(w, r, data.State) {
		return
	}
	intent, err := h.commands.GetActiveIntent(ctx, data.State)
	if err != nil {
		if zerrors.IsNotFound(err) {
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

	provider, err := h.getProvider(ctx, intent.IDPID)
	if err != nil {
		cmdErr := h.commands.FailIDPIntent(ctx, intent, err.Error())
		logging.WithFields("intent", intent.AggregateID).OnError(cmdErr).Error("failed to push failed event on idp intent")
		redirectToFailureURLErr(w, r, intent, err)
		return
	}

	idpUser, idpSession, err := h.fetchIDPUserFromCode(ctx, provider, data.Code, data.User)
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
		redirectToFailureURLErr(w, r, intent, zerrors.ThrowInternal(err, "IDP-JdD3g", "Errors.Intent.TokenCreationFailed"))
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
		return nil, zerrors.ThrowInvalidArgument(nil, "IDP-Hk38e", "Errors.Intent.StateMissing")
	}
	return data, nil
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
	zErr := new(zerrors.ZitadelError)
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

func (h *Handler) fetchIDPUserFromCode(ctx context.Context, identityProvider idp.Provider, code string, appleUser string) (user idp.User, idpTokens idp.Session, err error) {
	var session idp.Session
	switch provider := identityProvider.(type) {
	case *oauth.Provider:
		session = &oauth.Session{Provider: provider, Code: code}
	case *openid.Provider:
		session = &openid.Session{Provider: provider, Code: code}
	case *azuread.Provider:
		session = &azuread.Session{Provider: provider, Code: code}
	case *github.Provider:
		session = &oauth.Session{Provider: provider.Provider, Code: code}
	case *gitlab.Provider:
		session = &openid.Session{Provider: provider.Provider, Code: code}
	case *google.Provider:
		session = &openid.Session{Provider: provider.Provider, Code: code}
	case *apple.Provider:
		session = &apple.Session{Session: &openid.Session{Provider: provider.Provider, Code: code}, UserFormValue: appleUser}
	case *jwt.Provider, *ldap.Provider, *saml2.Provider:
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "IDP-52jmn", "Errors.ExternalIDP.IDPTypeNotImplemented")
	default:
		return nil, nil, zerrors.ThrowUnimplemented(nil, "IDP-SSDg", "Errors.ExternalIDP.IDPTypeNotImplemented")
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
