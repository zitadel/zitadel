package login

import (
	"context"
	"net/http"
	"strings"

	"github.com/crewjam/saml/samlsp"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
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
	"github.com/zitadel/zitadel/internal/idp/providers/saml"
	"github.com/zitadel/zitadel/internal/idp/providers/saml/requesttracker"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	queryIDPConfigID           = "idpConfigID"
	tmplExternalNotFoundOption = "externalnotfoundoption"
)

type externalIDPData struct {
	IDPConfigID string `schema:"idpConfigID"`
}

type externalIDPCallbackData struct {
	State string `schema:"state"`
	Code  string `schema:"code"`

	RelayState string `schema:"RelayState"`
	Method     string `schema:"Method"`

	// Apple returns a user on first registration
	User string `schema:"user"`
}

type externalNotFoundOptionFormData struct {
	externalRegisterFormData
	Link         bool `schema:"linkbutton"`
	AutoRegister bool `schema:"autoregisterbutton"`
	ResetLinking bool `schema:"resetlinking"`
	TermsConfirm bool `schema:"terms-confirm"`
}

type externalNotFoundOptionData struct {
	baseData
	externalNotFoundOptionFormData
	IsLinkingAllowed           bool
	IsCreationAllowed          bool
	ExternalIDPID              string
	ExternalIDPUserID          string
	ExternalIDPUserDisplayName string
	ShowUsername               bool
	ShowUsernameSuffix         bool
	OrgRegister                bool
	ExternalEmail              domain.EmailAddress
	ExternalEmailVerified      bool
	ExternalPhone              domain.PhoneNumber
	ExternalPhoneVerified      bool
	ProviderName               string
}

type externalRegisterFormData struct {
	ExternalIDPConfigID    string              `schema:"external-idp-config-id"`
	ExternalIDPExtUserID   string              `schema:"external-idp-ext-user-id"`
	ExternalIDPDisplayName string              `schema:"external-idp-display-name"`
	ExternalEmail          domain.EmailAddress `schema:"external-email"`
	ExternalEmailVerified  bool                `schema:"external-email-verified"`
	Email                  domain.EmailAddress `schema:"email"`
	Username               string              `schema:"username"`
	Firstname              string              `schema:"firstname"`
	Lastname               string              `schema:"lastname"`
	Nickname               string              `schema:"nickname"`
	ExternalPhone          domain.PhoneNumber  `schema:"external-phone"`
	ExternalPhoneVerified  bool                `schema:"external-phone-verified"`
	Phone                  domain.PhoneNumber  `schema:"phone"`
	Language               string              `schema:"language"`
	TermsConfirm           bool                `schema:"terms-confirm"`
}

// handleExternalLoginStep is called as nextStep
func (l *Login) handleExternalLoginStep(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, selectedIDPID string) {
	for _, idp := range authReq.AllowedExternalIDPs {
		if idp.IDPConfigID == selectedIDPID {
			l.handleIDP(w, r, authReq, selectedIDPID)
			return
		}
	}
	l.renderLogin(w, r, authReq, zerrors.ThrowInvalidArgument(nil, "VIEW-Fsj7f", "Errors.User.ExternalIDP.NotAllowed"))
}

// handleExternalLogin is called when a user selects the idp on the login page
func (l *Login) handleExternalLogin(w http.ResponseWriter, r *http.Request) {
	data := new(externalIDPData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if authReq == nil {
		l.defaultRedirect(w, r)
		return
	}
	l.handleIDP(w, r, authReq, data.IDPConfigID)
}

// handleExternalRegister is called when a user selects the idp on the register options page
func (l *Login) handleExternalRegister(w http.ResponseWriter, r *http.Request) {
	data := new(externalIDPData)
	authReq, err := l.ensureAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.handleIDP(w, r, authReq, data.IDPConfigID)
}

// handleIDP start the authentication of the selected IDP
// it will redirect to the IDPs auth page
func (l *Login) handleIDP(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, id string) {
	identityProvider, err := l.getIDPByID(r, id)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	err = l.authRepo.SelectExternalIDP(r.Context(), authReq.ID, identityProvider.ID, userAgentID)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	var provider idp.Provider

	switch identityProvider.Type {
	case domain.IDPTypeOAuth:
		provider, err = l.oauthProvider(r.Context(), identityProvider)
	case domain.IDPTypeOIDC:
		provider, err = l.oidcProvider(r.Context(), identityProvider)
	case domain.IDPTypeJWT:
		provider, err = l.jwtProvider(identityProvider)
	case domain.IDPTypeAzureAD:
		provider, err = l.azureProvider(r.Context(), identityProvider)
	case domain.IDPTypeGitHub:
		provider, err = l.githubProvider(r.Context(), identityProvider)
	case domain.IDPTypeGitHubEnterprise:
		provider, err = l.githubEnterpriseProvider(r.Context(), identityProvider)
	case domain.IDPTypeGitLab:
		provider, err = l.gitlabProvider(r.Context(), identityProvider)
	case domain.IDPTypeGitLabSelfHosted:
		provider, err = l.gitlabSelfHostedProvider(r.Context(), identityProvider)
	case domain.IDPTypeGoogle:
		provider, err = l.googleProvider(r.Context(), identityProvider)
	case domain.IDPTypeApple:
		provider, err = l.appleProvider(r.Context(), identityProvider)
	case domain.IDPTypeLDAP:
		provider, err = l.ldapProvider(r.Context(), identityProvider)
	case domain.IDPTypeSAML:
		provider, err = l.samlProvider(r.Context(), identityProvider)
	case domain.IDPTypeUnspecified:
		fallthrough
	default:
		l.renderLogin(w, r, authReq, zerrors.ThrowInvalidArgument(nil, "LOGIN-AShek", "Errors.ExternalIDP.IDPTypeNotImplemented"))
		return
	}
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	params := l.sessionParamsFromAuthRequest(r.Context(), authReq, identityProvider.ID)
	session, err := provider.BeginAuth(r.Context(), authReq.ID, params...)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}

	content, redirect := session.GetAuth(r.Context())
	if redirect {
		http.Redirect(w, r, content, http.StatusFound)
		return
	}
	_, err = w.Write([]byte(content))
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
}

// handleExternalLoginCallbackForm handles the callback from a IDP with form_post.
// It will redirect to the "normal" callback endpoint with the form data as query parameter.
// This way cookies will be handled correctly (same site = lax).
func (l *Login) handleExternalLoginCallbackForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		l.renderLogin(w, r, nil, err)
		return
	}
	r.Form.Add("Method", http.MethodPost)
	http.Redirect(w, r, HandlerPrefix+EndpointExternalLoginCallback+"?"+r.Form.Encode(), 302)
}

// handleExternalLoginCallback handles the callback from a IDP
// and tries to extract the user with the provided data
func (l *Login) handleExternalLoginCallback(w http.ResponseWriter, r *http.Request) {
	data := new(externalIDPCallbackData)
	err := l.getParseData(r, data)
	if err != nil {
		l.renderLogin(w, r, nil, err)
		return
	}
	if data.State == "" {
		data.State = data.RelayState
	}
	// workaround because of CSRF on external identity provider flows
	if data.Method == http.MethodPost {
		r.Method = http.MethodPost
		r.PostForm = r.Form
	}

	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	authReq, err := l.authRepo.AuthRequestByID(r.Context(), data.State, userAgentID)
	if err != nil {
		l.externalAuthFailed(w, r, authReq, nil, nil, err)
		return
	}
	identityProvider, err := l.getIDPByID(r, authReq.SelectedIDPConfigID)
	if err != nil {
		l.externalAuthFailed(w, r, authReq, nil, nil, err)
		return
	}
	var provider idp.Provider
	var session idp.Session
	switch identityProvider.Type {
	case domain.IDPTypeOAuth:
		provider, err = l.oauthProvider(r.Context(), identityProvider)
		if err != nil {
			l.externalAuthFailed(w, r, authReq, nil, nil, err)
			return
		}
		session = &oauth.Session{Provider: provider.(*oauth.Provider), Code: data.Code}
	case domain.IDPTypeOIDC:
		provider, err = l.oidcProvider(r.Context(), identityProvider)
		if err != nil {
			l.externalAuthFailed(w, r, authReq, nil, nil, err)
			return
		}
		session = &openid.Session{Provider: provider.(*openid.Provider), Code: data.Code}
	case domain.IDPTypeAzureAD:
		provider, err = l.azureProvider(r.Context(), identityProvider)
		if err != nil {
			l.externalAuthFailed(w, r, authReq, nil, nil, err)
			return
		}
		session = &azuread.Session{Provider: provider.(*azuread.Provider), Code: data.Code}
	case domain.IDPTypeGitHub:
		provider, err = l.githubProvider(r.Context(), identityProvider)
		if err != nil {
			l.externalAuthFailed(w, r, authReq, nil, nil, err)
			return
		}
		session = &oauth.Session{Provider: provider.(*github.Provider).Provider, Code: data.Code}
	case domain.IDPTypeGitHubEnterprise:
		provider, err = l.githubEnterpriseProvider(r.Context(), identityProvider)
		if err != nil {
			l.externalAuthFailed(w, r, authReq, nil, nil, err)
			return
		}
		session = &oauth.Session{Provider: provider.(*github.Provider).Provider, Code: data.Code}
	case domain.IDPTypeGitLab:
		provider, err = l.gitlabProvider(r.Context(), identityProvider)
		if err != nil {
			l.externalAuthFailed(w, r, authReq, nil, nil, err)
			return
		}
		session = &openid.Session{Provider: provider.(*gitlab.Provider).Provider, Code: data.Code}
	case domain.IDPTypeGitLabSelfHosted:
		provider, err = l.gitlabSelfHostedProvider(r.Context(), identityProvider)
		if err != nil {
			l.externalAuthFailed(w, r, authReq, nil, nil, err)
			return
		}
		session = &openid.Session{Provider: provider.(*gitlab.Provider).Provider, Code: data.Code}
	case domain.IDPTypeGoogle:
		provider, err = l.googleProvider(r.Context(), identityProvider)
		if err != nil {
			l.externalAuthFailed(w, r, authReq, nil, nil, err)
			return
		}
		session = &openid.Session{Provider: provider.(*google.Provider).Provider, Code: data.Code}
	case domain.IDPTypeApple:
		provider, err = l.appleProvider(r.Context(), identityProvider)
		if err != nil {
			l.externalAuthFailed(w, r, authReq, nil, nil, err)
			return
		}
		session = &apple.Session{Session: &openid.Session{Provider: provider.(*apple.Provider).Provider, Code: data.Code}, UserFormValue: data.User}
	case domain.IDPTypeSAML:
		provider, err = l.samlProvider(r.Context(), identityProvider)
		if err != nil {
			l.externalAuthFailed(w, r, authReq, nil, nil, err)
			return
		}
		session, err = saml.NewSession(provider.(*saml.Provider), authReq.SAMLRequestID, r)
		if err != nil {
			l.externalAuthFailed(w, r, authReq, nil, nil, err)
			return
		}
	case domain.IDPTypeJWT,
		domain.IDPTypeLDAP,
		domain.IDPTypeUnspecified:
		fallthrough
	default:
		l.renderLogin(w, r, authReq, zerrors.ThrowInvalidArgument(nil, "LOGIN-SFefg", "Errors.ExternalIDP.IDPTypeNotImplemented"))
		return
	}

	user, err := session.FetchUser(r.Context())
	if err != nil {
		logging.WithFields(
			"instance", authz.GetInstance(r.Context()).InstanceID(),
			"providerID", identityProvider.ID,
		).WithError(err).Info("external authentication failed")
		l.externalAuthFailed(w, r, authReq, tokens(session), user, err)
		return
	}
	l.handleExternalUserAuthenticated(w, r, authReq, identityProvider, session, user, l.renderNextStep)
}

func (l *Login) tryMigrateExternalUserID(r *http.Request, session idp.Session, authReq *domain.AuthRequest, externalUser *domain.ExternalUser) (previousIDMatched bool, err error) {
	migration, ok := session.(idp.SessionSupportsMigration)
	if !ok {
		return false, nil
	}
	previousID, err := migration.RetrievePreviousID()
	if err != nil {
		return false, err
	}
	return l.migrateExternalUserID(r, authReq, externalUser, previousID)
}

func (l *Login) migrateExternalUserID(r *http.Request, authReq *domain.AuthRequest, externalUser *domain.ExternalUser, previousID string) (previousIDMatched bool, err error) {
	if previousID == "" {
		return false, nil
	}
	// save the currentID, so we're able to reset to it later on if the user is not found with the old ID as well
	externalUserID := externalUser.ExternalUserID
	externalUser.ExternalUserID = previousID
	if err = l.authRepo.CheckExternalUserLogin(setContext(r.Context(), ""), authReq.ID, authReq.AgentID, externalUser, domain.BrowserInfoFromRequest(r), true); err != nil {
		// always reset to the mapped ID
		externalUser.ExternalUserID = externalUserID
		// but ignore the error if the user was just not found with the previousID
		if zerrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	previousIDMatched = true
	if err = l.authRepo.ResetLinkingUsers(r.Context(), authReq.ID, authReq.AgentID); err != nil {
		return previousIDMatched, err
	}
	// read current auth request state (incl. authorized user)
	authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
	if err != nil {
		return previousIDMatched, err
	}
	return previousIDMatched, l.command.MigrateUserIDP(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, externalUser.IDPConfigID, previousID, externalUserID)
}

// handleExternalUserAuthenticated maps the IDP user, checks for a corresponding externalID and that the IDP is allowed
func (l *Login) handleExternalUserAuthenticated(
	w http.ResponseWriter,
	r *http.Request,
	authReq *domain.AuthRequest,
	provider *query.IDPTemplate,
	session idp.Session,
	user idp.User,
	callback func(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest),
) {
	externalUser := mapIDPUserToExternalUser(user, provider.ID)
	// ensure the linked IDP is added to the login policy
	if err := l.authRepo.SelectExternalIDP(r.Context(), authReq.ID, provider.ID, authReq.AgentID); err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	// check and fill in local linked user
	externalErr := l.authRepo.CheckExternalUserLogin(setContext(r.Context(), ""), authReq.ID, authReq.AgentID, externalUser, domain.BrowserInfoFromRequest(r), false)
	if externalErr != nil && !zerrors.IsNotFound(externalErr) {
		l.renderError(w, r, authReq, externalErr)
		return
	}
	if externalErr != nil && zerrors.IsNotFound(externalErr) {
		previousIDMatched, err := l.tryMigrateExternalUserID(r, session, authReq, externalUser)
		if err != nil {
			l.renderError(w, r, authReq, err)
			return
		}
		// if the old ID matched, ignore the not found error from the current ID
		if previousIDMatched {
			externalErr = nil
		}
	}
	var err error
	// read current auth request state (incl. authorized user)
	authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	externalUser, externalUserChange, err := l.runPostExternalAuthenticationActions(externalUser, tokens(session), authReq, r, user, nil)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	// if action is done and no user linked then link or register
	if zerrors.IsNotFound(externalErr) {
		l.externalUserNotExisting(w, r, authReq, provider, externalUser, externalUserChange)
		return
	}
	if provider.IsAutoUpdate || externalUserChange {
		err = l.updateExternalUser(r.Context(), authReq, externalUser)
		if err != nil {
			l.renderError(w, r, authReq, err)
			return
		}
	}
	if len(externalUser.Metadatas) > 0 {
		_, err = l.command.BulkSetUserMetadata(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, externalUser.Metadatas...)
		if err != nil {
			l.renderError(w, r, authReq, err)
			return
		}
	}
	callback(w, r, authReq)
}

// checkAutoLinking checks if a user with the provided information (username or email) already exists within ZITADEL.
// The decision, which information will be checked is based on the IdP template option.
// The function returns a boolean whether a user was found or not.
// If single a user was found, it will be automatically linked.
func (l *Login) checkAutoLinking(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, provider *query.IDPTemplate, externalUser *domain.ExternalUser) bool {
	queries := make([]query.SearchQuery, 0, 2)
	switch provider.AutoLinking {
	case domain.AutoLinkingOptionUnspecified:
		// is auto linking is disable, we shouldn't even get here, but in case we do we can directly return
		return false
	case domain.AutoLinkingOptionUsername:
		// if we're checking for usernames there are to options:
		//
		// If no specific org has been requested (by id or domain scope), we'll check the provided username against
		// all existing loginnames and directly use that result to either prompt or continue with other idp options.
		if authReq.RequestedOrgID == "" {
			user, err := l.query.GetNotifyUserByLoginName(r.Context(), false, externalUser.PreferredUsername)
			if err != nil {
				return false
			}
			l.autoLinkUser(w, r, authReq, user)
			return true
		}
		// If a specific org has been requested, we'll check the provided username against usernames (of that org).
		usernameQuery, err := query.NewUserUsernameSearchQuery(externalUser.PreferredUsername, query.TextEqualsIgnoreCase)
		if err != nil {
			return false
		}
		queries = append(queries, usernameQuery)
	case domain.AutoLinkingOptionEmail:
		// Email will always be checked against verified email addresses.
		emailQuery, err := query.NewUserVerifiedEmailSearchQuery(string(externalUser.Email))
		if err != nil {
			return false
		}
		queries = append(queries, emailQuery)
	}
	// restrict the possible organization if needed (for email and usernames)
	if authReq.RequestedOrgID != "" {
		resourceOwnerQuery, err := query.NewUserResourceOwnerSearchQuery(authReq.RequestedOrgID, query.TextEquals)
		if err != nil {
			return false
		}
		queries = append(queries, resourceOwnerQuery)
	}
	user, err := l.query.GetNotifyUser(r.Context(), false, queries...)
	if err != nil {
		return false
	}
	l.autoLinkUser(w, r, authReq, user)
	return true
}

func (l *Login) autoLinkUser(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, user *query.NotifyUser) {
	if err := l.authRepo.SelectUser(r.Context(), authReq.ID, user.ID, authReq.AgentID); err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if err := l.authRepo.LinkExternalUsers(r.Context(), authReq.ID, authReq.AgentID, domain.BrowserInfoFromRequest(r)); err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}

// externalUserNotExisting is called if an externalAuthentication couldn't find a corresponding externalID
// possible solutions are:
//
// * auto creation
// * external not found overview:
//   - creation by user
//   - linking to existing user
func (l *Login) externalUserNotExisting(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, provider *query.IDPTemplate, externalUser *domain.ExternalUser, changed bool) {
	resourceOwner := determineResourceOwner(r.Context(), authReq)
	orgIAMPolicy, err := l.getOrgDomainPolicy(r, resourceOwner)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, err)
		return
	}

	human, idpLink, _ := mapExternalUserToLoginUser(externalUser, orgIAMPolicy.UserLoginMustBeDomain)
	// let's check if auto-linking is enabled and if the user would be found by the corresponding option
	if provider.AutoLinking != domain.AutoLinkingOptionUnspecified {
		if l.checkAutoLinking(w, r, authReq, provider, externalUser) {
			return
		}
	}

	// if auto creation is disabled, send the user to the notFoundOption
	// where they can either link or create an account (based on the available options)
	if !provider.IsAutoCreation {
		l.renderExternalNotFoundOption(w, r, authReq, orgIAMPolicy, human, idpLink, nil)
		return
	}

	// reload auth request, to ensure current state (checked external login)
	authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, orgIAMPolicy, human, idpLink, err)
		return
	}
	if changed || len(externalUser.Metadatas) > 0 {
		if err := l.authRepo.SetLinkingUser(r.Context(), authReq, externalUser); err != nil {
			l.renderError(w, r, authReq, err)
			return
		}
	}
	l.autoCreateExternalUser(w, r, authReq)
}

// autoCreateExternalUser takes the externalUser and creates it automatically (without user interaction)
func (l *Login) autoCreateExternalUser(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	if len(authReq.LinkingUsers) == 0 {
		l.renderError(w, r, authReq, zerrors.ThrowPreconditionFailed(nil, "LOGIN-asfg3", "Errors.ExternalIDP.NoExternalUserData"))
		return
	}

	// TODO (LS): how do we get multiple and why do we use the last of them (taken as is)?
	linkingUser := authReq.LinkingUsers[len(authReq.LinkingUsers)-1]

	l.registerExternalUser(w, r, authReq, linkingUser)
}

// renderExternalNotFoundOption renders a page, where the user is able to edit the IDP data,
// create a new externalUser of link to existing on (based on the IDP template)
func (l *Login) renderExternalNotFoundOption(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, orgIAMPolicy *query.DomainPolicy, human *domain.Human, idpLink *domain.UserIDPLink, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	resourceOwner := determineResourceOwner(r.Context(), authReq)
	if orgIAMPolicy == nil {
		orgIAMPolicy, err = l.getOrgDomainPolicy(r, resourceOwner)
		if err != nil {
			l.renderError(w, r, authReq, err)
			return
		}
	}

	if human == nil || idpLink == nil {
		// TODO (LS): how do we get multiple and why do we use the last of them (taken as is)?
		linkingUser := authReq.LinkingUsers[len(authReq.LinkingUsers)-1]
		human, idpLink, _ = mapExternalUserToLoginUser(linkingUser, orgIAMPolicy.UserLoginMustBeDomain)
	}

	labelPolicy, err := l.getLabelPolicy(r, resourceOwner)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	idpTemplate, err := l.getIDPByID(r, idpLink.IDPConfigID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if !idpTemplate.IsCreationAllowed && !idpTemplate.IsLinkingAllowed {
		l.renderError(w, r, authReq, zerrors.ThrowPreconditionFailed(nil, "LOGIN-3kl44", "Errors.User.ExternalIDP.NoOptionAllowed"))
		return
	}

	translator := l.getTranslator(r.Context(), authReq)
	data := externalNotFoundOptionData{
		baseData: l.getBaseData(r, authReq, translator, "ExternalNotFound.Title", "ExternalNotFound.Description", errID, errMessage),
		externalNotFoundOptionFormData: externalNotFoundOptionFormData{
			externalRegisterFormData: externalRegisterFormData{
				Email:     human.EmailAddress,
				Username:  human.Username,
				Firstname: human.FirstName,
				Lastname:  human.LastName,
				Nickname:  human.NickName,
				Language:  human.PreferredLanguage.String(),
			},
		},
		IsLinkingAllowed:           idpTemplate.IsLinkingAllowed,
		IsCreationAllowed:          idpTemplate.IsCreationAllowed,
		ExternalIDPID:              idpLink.IDPConfigID,
		ExternalIDPUserID:          idpLink.ExternalUserID,
		ExternalIDPUserDisplayName: idpLink.DisplayName,
		ExternalEmail:              human.EmailAddress,
		ExternalEmailVerified:      human.IsEmailVerified,
		ShowUsername:               orgIAMPolicy.UserLoginMustBeDomain,
		ShowUsernameSuffix:         !labelPolicy.HideLoginNameSuffix,
		OrgRegister:                orgIAMPolicy.UserLoginMustBeDomain,
		ProviderName:               domain.IDPName(idpTemplate.Name, idpTemplate.Type),
	}
	if human.Phone != nil {
		data.Phone = human.PhoneNumber
		data.ExternalPhone = human.PhoneNumber
		data.ExternalPhoneVerified = human.IsPhoneVerified
	}
	funcs := map[string]interface{}{
		"selectedLanguage": func(l string) bool {
			return data.Language == l
		},
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplExternalNotFoundOption], data, funcs)
}

// handleExternalNotFoundOptionCheck takes the data from the submitted externalNotFound page
// and either links or creates an externalUser
func (l *Login) handleExternalNotFoundOptionCheck(w http.ResponseWriter, r *http.Request) {
	data := new(externalNotFoundOptionFormData)
	authReq, err := l.ensureAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, err)
		return
	}

	idpTemplate, err := l.getIDPByID(r, authReq.SelectedIDPConfigID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	// if the user click on the cancel button / back icon
	if data.ResetLinking {
		userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
		err = l.authRepo.ResetLinkingUsers(r.Context(), authReq.ID, userAgentID)
		if err != nil {
			l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, err)
		}
		l.handleLogin(w, r)
		return
	}
	// if the user selects the linking button
	if data.Link {
		if !idpTemplate.IsLinkingAllowed {
			l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, zerrors.ThrowPreconditionFailed(nil, "LOGIN-AS3ff", "Errors.ExternalIDP.LinkingNotAllowed"))
			return
		}
		l.renderLogin(w, r, authReq, nil)
		return
	}
	// if the user selects the creation button
	if !idpTemplate.IsCreationAllowed {
		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, zerrors.ThrowPreconditionFailed(nil, "LOGIN-dsfd3", "Errors.ExternalIDP.CreationNotAllowed"))
		return
	}
	linkingUser := mapExternalNotFoundOptionFormDataToLoginUser(data)
	l.registerExternalUser(w, r, authReq, linkingUser)
}

// registerExternalUser creates an externalUser with the provided data
// incl. execution of pre and post creation actions
//
// it is called from either the [autoCreateExternalUser] or [handleExternalNotFoundOptionCheck]
func (l *Login) registerExternalUser(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, externalUser *domain.ExternalUser) {
	resourceOwner := determineResourceOwner(r.Context(), authReq)

	orgIamPolicy, err := l.getOrgDomainPolicy(r, resourceOwner)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, err)
		return
	}
	user, externalIDP, metadata := mapExternalUserToLoginUser(externalUser, orgIamPolicy.UserLoginMustBeDomain)

	user, metadata, err = l.runPreCreationActions(authReq, r, user, metadata, resourceOwner, domain.FlowTypeExternalAuthentication)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, orgIamPolicy, nil, nil, err)
		return
	}
	err = l.authRepo.AutoRegisterExternalUser(setContext(r.Context(), resourceOwner), user, externalIDP, nil, authReq.ID, authReq.AgentID, resourceOwner, metadata, domain.BrowserInfoFromRequest(r))
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, orgIamPolicy, user, externalIDP, err)
		return
	}
	// read auth request again to get current state including userID
	authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	userGrants, err := l.runPostCreationActions(authReq.UserID, authReq, r, resourceOwner, domain.FlowTypeExternalAuthentication)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	err = l.appendUserGrants(r.Context(), userGrants, resourceOwner)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}

// updateExternalUser will update the existing user (email, phone, profile) with data provided by the IDP
func (l *Login) updateExternalUser(ctx context.Context, authReq *domain.AuthRequest, externalUser *domain.ExternalUser) error {
	user, err := l.query.GetUserByID(ctx, true, authReq.UserID)
	if err != nil {
		return err
	}
	if user.Human == nil {
		return zerrors.ThrowPreconditionFailed(nil, "LOGIN-WLTce", "Errors.User.NotHuman")
	}
	err = l.updateExternalUserEmail(ctx, user, externalUser)
	logging.WithFields("authReq", authReq.ID, "user", authReq.UserID).OnError(err).Error("unable to update email")

	err = l.updateExternalUserPhone(ctx, user, externalUser)
	logging.WithFields("authReq", authReq.ID, "user", authReq.UserID).OnError(err).Error("unable to update phone")

	err = l.updateExternalUserProfile(ctx, user, externalUser)
	logging.WithFields("authReq", authReq.ID, "user", authReq.UserID).OnError(err).Error("unable to update profile")

	err = l.updateExternalUsername(ctx, user, externalUser)
	logging.WithFields("authReq", authReq.ID, "user", authReq.UserID).OnError(err).Error("unable to update external username")

	return nil
}

func (l *Login) updateExternalUserEmail(ctx context.Context, user *query.User, externalUser *domain.ExternalUser) error {
	changed := hasEmailChanged(user, externalUser)
	if !changed {
		return nil
	}
	// if the email has changed and / or was not verified, we change it
	emailCodeGenerator, err := l.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyEmailCode, l.userCodeAlg)
	if err != nil {
		return err
	}
	_, err = l.command.ChangeHumanEmail(setContext(ctx, user.ResourceOwner),
		&domain.Email{
			ObjectRoot:      models.ObjectRoot{AggregateID: user.ID},
			EmailAddress:    externalUser.Email,
			IsEmailVerified: externalUser.IsEmailVerified,
		},
		emailCodeGenerator)
	return err
}

func (l *Login) updateExternalUserPhone(ctx context.Context, user *query.User, externalUser *domain.ExternalUser) error {
	changed, err := hasPhoneChanged(user, externalUser)
	if !changed || err != nil {
		return err
	}
	// if the phone has changed and / or was not verified, we change it
	phoneCodeGenerator, err := l.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyPhoneCode, l.userCodeAlg)
	if err != nil {
		return err
	}
	_, err = l.command.ChangeHumanPhone(setContext(ctx, user.ResourceOwner),
		&domain.Phone{
			ObjectRoot:      models.ObjectRoot{AggregateID: user.ID},
			PhoneNumber:     externalUser.Phone,
			IsPhoneVerified: externalUser.IsPhoneVerified,
		},
		user.ResourceOwner,
		phoneCodeGenerator)
	return err
}

func (l *Login) updateExternalUserProfile(ctx context.Context, user *query.User, externalUser *domain.ExternalUser) error {
	if externalUser.FirstName == user.Human.FirstName &&
		externalUser.LastName == user.Human.LastName &&
		externalUser.NickName == user.Human.NickName &&
		externalUser.DisplayName == user.Human.DisplayName &&
		externalUser.PreferredLanguage == user.Human.PreferredLanguage {
		return nil
	}
	_, err := l.command.ChangeHumanProfile(setContext(ctx, user.ResourceOwner), &domain.Profile{
		ObjectRoot:        models.ObjectRoot{AggregateID: user.ID},
		FirstName:         externalUser.FirstName,
		LastName:          externalUser.LastName,
		NickName:          externalUser.NickName,
		DisplayName:       externalUser.DisplayName,
		PreferredLanguage: externalUser.PreferredLanguage,
		Gender:            user.Human.Gender,
	})
	return err
}

func (l *Login) updateExternalUsername(ctx context.Context, user *query.User, externalUser *domain.ExternalUser) error {
	externalIDQuery, err := query.NewIDPUserLinksExternalIDSearchQuery(externalUser.ExternalUserID)
	if err != nil {
		return err
	}
	idpIDQuery, err := query.NewIDPUserLinkIDPIDSearchQuery(externalUser.IDPConfigID)
	if err != nil {
		return err
	}
	userIDQuery, err := query.NewIDPUserLinksUserIDSearchQuery(user.ID)
	if err != nil {
		return err
	}
	links, err := l.query.IDPUserLinks(ctx, &query.IDPUserLinksSearchQuery{Queries: []query.SearchQuery{externalIDQuery, idpIDQuery, userIDQuery}}, nil)
	if err != nil || len(links.Links) == 0 {
		return err
	}
	if links.Links[0].ProvidedUsername == externalUser.PreferredUsername {
		return nil
	}
	return l.command.UpdateUserIDPLinkUsername(
		setContext(ctx, user.ResourceOwner),
		user.ID,
		user.ResourceOwner,
		externalUser.IDPConfigID,
		externalUser.ExternalUserID,
		externalUser.PreferredUsername,
	)
}

func hasEmailChanged(user *query.User, externalUser *domain.ExternalUser) bool {
	externalUser.Email = externalUser.Email.Normalize()
	if externalUser.Email == "" {
		return false
	}
	// ignore if the same email is not set to verified anymore
	if externalUser.Email == user.Human.Email && user.Human.IsEmailVerified {
		return false
	}
	return externalUser.Email != user.Human.Email || externalUser.IsEmailVerified != user.Human.IsEmailVerified
}

func hasPhoneChanged(user *query.User, externalUser *domain.ExternalUser) (_ bool, err error) {
	if externalUser.Phone == "" {
		return false, nil
	}
	externalUser.Phone, err = externalUser.Phone.Normalize()
	if err != nil {
		return false, err
	}
	// ignore if the same phone is not set to verified anymore
	if externalUser.Phone == user.Human.Phone && user.Human.IsPhoneVerified {
		return false, nil
	}
	return externalUser.Phone != user.Human.Phone || externalUser.IsPhoneVerified != user.Human.IsPhoneVerified, nil
}

func (l *Login) ldapProvider(ctx context.Context, identityProvider *query.IDPTemplate) (*ldap.Provider, error) {
	password, err := crypto.DecryptString(identityProvider.LDAPIDPTemplate.BindPassword, l.idpConfigAlg)
	if err != nil {
		return nil, err
	}
	var opts []ldap.ProviderOpts
	if !identityProvider.LDAPIDPTemplate.StartTLS {
		opts = append(opts, ldap.WithoutStartTLS())
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.IDAttribute != "" {
		opts = append(opts, ldap.WithCustomIDAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.IDAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.FirstNameAttribute != "" {
		opts = append(opts, ldap.WithFirstNameAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.FirstNameAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.LastNameAttribute != "" {
		opts = append(opts, ldap.WithLastNameAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.LastNameAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.DisplayNameAttribute != "" {
		opts = append(opts, ldap.WithDisplayNameAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.DisplayNameAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.NickNameAttribute != "" {
		opts = append(opts, ldap.WithNickNameAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.NickNameAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.PreferredUsernameAttribute != "" {
		opts = append(opts, ldap.WithPreferredUsernameAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.PreferredUsernameAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.EmailAttribute != "" {
		opts = append(opts, ldap.WithEmailAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.EmailAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.EmailVerifiedAttribute != "" {
		opts = append(opts, ldap.WithEmailVerifiedAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.EmailVerifiedAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.PhoneAttribute != "" {
		opts = append(opts, ldap.WithPhoneAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.PhoneAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.PhoneVerifiedAttribute != "" {
		opts = append(opts, ldap.WithPhoneVerifiedAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.PhoneVerifiedAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.PreferredLanguageAttribute != "" {
		opts = append(opts, ldap.WithPreferredLanguageAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.PreferredLanguageAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.AvatarURLAttribute != "" {
		opts = append(opts, ldap.WithAvatarURLAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.AvatarURLAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.ProfileAttribute != "" {
		opts = append(opts, ldap.WithProfileAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.ProfileAttribute))
	}
	return ldap.New(
		identityProvider.Name,
		identityProvider.Servers,
		identityProvider.BaseDN,
		identityProvider.BindDN,
		password,
		identityProvider.UserBase,
		identityProvider.UserObjectClasses,
		identityProvider.UserFilters,
		identityProvider.Timeout,
		l.baseURL(ctx)+EndpointLDAPLogin+"?"+QueryAuthRequestID+"=",
		opts...,
	), nil
}

func (l *Login) googleProvider(ctx context.Context, identityProvider *query.IDPTemplate) (*google.Provider, error) {
	errorHandler := func(w http.ResponseWriter, r *http.Request, errorType string, errorDesc string, state string) {
		logging.Errorf("token exchanged failed: %s - %s (state: %s)", errorType, errorType, state)
		rp.DefaultErrorHandler(w, r, errorType, errorDesc, state)
	}
	openid.WithRelyingPartyOption(rp.WithErrorHandler(errorHandler))
	secret, err := crypto.DecryptString(identityProvider.GoogleIDPTemplate.ClientSecret, l.idpConfigAlg)
	if err != nil {
		return nil, err
	}
	return google.New(
		identityProvider.GoogleIDPTemplate.ClientID,
		secret,
		l.baseURL(ctx)+EndpointExternalLoginCallback,
		identityProvider.GoogleIDPTemplate.Scopes,
	)
}

func (l *Login) oidcProvider(ctx context.Context, identityProvider *query.IDPTemplate) (*openid.Provider, error) {
	secret, err := crypto.DecryptString(identityProvider.OIDCIDPTemplate.ClientSecret, l.idpConfigAlg)
	if err != nil {
		return nil, err
	}
	opts := make([]openid.ProviderOpts, 1, 2)
	opts[0] = openid.WithSelectAccount()
	if identityProvider.OIDCIDPTemplate.IsIDTokenMapping {
		opts = append(opts, openid.WithIDTokenMapping())
	}
	return openid.New(identityProvider.Name,
		identityProvider.OIDCIDPTemplate.Issuer,
		identityProvider.OIDCIDPTemplate.ClientID,
		secret,
		l.baseURL(ctx)+EndpointExternalLoginCallback,
		identityProvider.OIDCIDPTemplate.Scopes,
		openid.DefaultMapper,
		opts...,
	)
}

func (l *Login) jwtProvider(identityProvider *query.IDPTemplate) (*jwt.Provider, error) {
	return jwt.New(
		identityProvider.Name,
		identityProvider.JWTIDPTemplate.Issuer,
		identityProvider.JWTIDPTemplate.Endpoint,
		identityProvider.JWTIDPTemplate.KeysEndpoint,
		identityProvider.JWTIDPTemplate.HeaderName,
		l.idpConfigAlg,
	)
}

func (l *Login) oauthProvider(ctx context.Context, identityProvider *query.IDPTemplate) (*oauth.Provider, error) {
	secret, err := crypto.DecryptString(identityProvider.OAuthIDPTemplate.ClientSecret, l.idpConfigAlg)
	if err != nil {
		return nil, err
	}
	config := &oauth2.Config{
		ClientID:     identityProvider.OAuthIDPTemplate.ClientID,
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  identityProvider.OAuthIDPTemplate.AuthorizationEndpoint,
			TokenURL: identityProvider.OAuthIDPTemplate.TokenEndpoint,
		},
		RedirectURL: l.baseURL(ctx) + EndpointExternalLoginCallback,
		Scopes:      identityProvider.OAuthIDPTemplate.Scopes,
	}
	return oauth.New(
		config,
		identityProvider.Name,
		identityProvider.OAuthIDPTemplate.UserEndpoint,
		func() idp.User {
			return oauth.NewUserMapper(identityProvider.OAuthIDPTemplate.IDAttribute)
		},
	)
}

func (l *Login) samlProvider(ctx context.Context, identityProvider *query.IDPTemplate) (*saml.Provider, error) {
	key, err := crypto.Decrypt(identityProvider.SAMLIDPTemplate.Key, l.idpConfigAlg)
	if err != nil {
		return nil, err
	}
	opts := make([]saml.ProviderOpts, 0, 6)
	if identityProvider.SAMLIDPTemplate.WithSignedRequest {
		opts = append(opts, saml.WithSignedRequest())
	}
	if identityProvider.SAMLIDPTemplate.Binding != "" {
		opts = append(opts, saml.WithBinding(identityProvider.SAMLIDPTemplate.Binding))
	}
	if identityProvider.SAMLIDPTemplate.NameIDFormat.Valid {
		opts = append(opts, saml.WithNameIDFormat(identityProvider.SAMLIDPTemplate.NameIDFormat.V))
	}
	if identityProvider.SAMLIDPTemplate.TransientMappingAttributeName != "" {
		opts = append(opts, saml.WithTransientMappingAttributeName(identityProvider.SAMLIDPTemplate.TransientMappingAttributeName))
	}
	opts = append(opts,
		saml.WithEntityID(http_utils.DomainContext(ctx).Origin()+"/idps/"+identityProvider.ID+"/saml/metadata"),
		saml.WithCustomRequestTracker(
			requesttracker.New(
				func(ctx context.Context, authRequestID, samlRequestID string) error {
					useragent, _ := http_mw.UserAgentIDFromCtx(ctx)
					return l.authRepo.SaveSAMLRequestID(ctx, authRequestID, samlRequestID, useragent)
				},
				func(ctx context.Context, authRequestID string) (*samlsp.TrackedRequest, error) {
					useragent, _ := http_mw.UserAgentIDFromCtx(ctx)
					auhRequest, err := l.authRepo.AuthRequestByID(ctx, authRequestID, useragent)
					if err != nil {
						return nil, err
					}
					return &samlsp.TrackedRequest{
						SAMLRequestID: auhRequest.SAMLRequestID,
						Index:         authRequestID,
					}, nil
				},
			),
		))
	return saml.New(
		identityProvider.Name,
		l.baseURL(ctx)+EndpointExternalLogin+"/",
		identityProvider.SAMLIDPTemplate.Metadata,
		identityProvider.SAMLIDPTemplate.Certificate,
		key,
		opts...,
	)
}

func (l *Login) azureProvider(ctx context.Context, identityProvider *query.IDPTemplate) (*azuread.Provider, error) {
	secret, err := crypto.DecryptString(identityProvider.AzureADIDPTemplate.ClientSecret, l.idpConfigAlg)
	if err != nil {
		return nil, err
	}
	opts := make([]azuread.ProviderOptions, 0, 2)
	if identityProvider.AzureADIDPTemplate.IsEmailVerified {
		opts = append(opts, azuread.WithEmailVerified())
	}
	if identityProvider.AzureADIDPTemplate.Tenant != "" {
		opts = append(opts, azuread.WithTenant(azuread.TenantType(identityProvider.AzureADIDPTemplate.Tenant)))
	}
	return azuread.New(
		identityProvider.Name,
		identityProvider.AzureADIDPTemplate.ClientID,
		secret,
		l.baseURL(ctx)+EndpointExternalLoginCallback,
		identityProvider.AzureADIDPTemplate.Scopes,
		opts...,
	)
}

func (l *Login) githubProvider(ctx context.Context, identityProvider *query.IDPTemplate) (*github.Provider, error) {
	secret, err := crypto.DecryptString(identityProvider.GitHubIDPTemplate.ClientSecret, l.idpConfigAlg)
	if err != nil {
		return nil, err
	}
	return github.New(
		identityProvider.GitHubIDPTemplate.ClientID,
		secret,
		l.baseURL(ctx)+EndpointExternalLoginCallback,
		identityProvider.GitHubIDPTemplate.Scopes,
	)
}

func (l *Login) githubEnterpriseProvider(ctx context.Context, identityProvider *query.IDPTemplate) (*github.Provider, error) {
	secret, err := crypto.DecryptString(identityProvider.GitHubIDPTemplate.ClientSecret, l.idpConfigAlg)
	if err != nil {
		return nil, err
	}
	return github.NewCustomURL(
		identityProvider.Name,
		identityProvider.GitHubIDPTemplate.ClientID,
		secret,
		l.baseURL(ctx)+EndpointExternalLoginCallback,
		identityProvider.GitHubEnterpriseIDPTemplate.AuthorizationEndpoint,
		identityProvider.GitHubEnterpriseIDPTemplate.TokenEndpoint,
		identityProvider.GitHubEnterpriseIDPTemplate.UserEndpoint,
		identityProvider.GitHubIDPTemplate.Scopes,
	)
}

func (l *Login) gitlabProvider(ctx context.Context, identityProvider *query.IDPTemplate) (*gitlab.Provider, error) {
	secret, err := crypto.DecryptString(identityProvider.GitLabIDPTemplate.ClientSecret, l.idpConfigAlg)
	if err != nil {
		return nil, err
	}
	return gitlab.New(
		identityProvider.GitLabIDPTemplate.ClientID,
		secret,
		l.baseURL(ctx)+EndpointExternalLoginCallback,
		identityProvider.GitLabIDPTemplate.Scopes,
	)
}

func (l *Login) gitlabSelfHostedProvider(ctx context.Context, identityProvider *query.IDPTemplate) (*gitlab.Provider, error) {
	secret, err := crypto.DecryptString(identityProvider.GitLabSelfHostedIDPTemplate.ClientSecret, l.idpConfigAlg)
	if err != nil {
		return nil, err
	}
	return gitlab.NewCustomIssuer(
		identityProvider.Name,
		identityProvider.GitLabSelfHostedIDPTemplate.Issuer,
		identityProvider.GitLabSelfHostedIDPTemplate.ClientID,
		secret,
		l.baseURL(ctx)+EndpointExternalLoginCallback,
		identityProvider.GitLabSelfHostedIDPTemplate.Scopes,
	)
}

func (l *Login) appleProvider(ctx context.Context, identityProvider *query.IDPTemplate) (*apple.Provider, error) {
	privateKey, err := crypto.Decrypt(identityProvider.AppleIDPTemplate.PrivateKey, l.idpConfigAlg)
	if err != nil {
		return nil, err
	}
	return apple.New(
		identityProvider.AppleIDPTemplate.ClientID,
		identityProvider.AppleIDPTemplate.TeamID,
		identityProvider.AppleIDPTemplate.KeyID,
		l.baseURL(ctx)+EndpointExternalLoginCallbackFormPost,
		privateKey,
		identityProvider.AppleIDPTemplate.Scopes,
	)
}

func (l *Login) appendUserGrants(ctx context.Context, userGrants []*domain.UserGrant, resourceOwner string) error {
	if len(userGrants) == 0 {
		return nil
	}
	for _, grant := range userGrants {
		_, err := l.command.AddUserGrant(setContext(ctx, resourceOwner), grant, resourceOwner)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Login) externalAuthFailed(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, tokens *oidc.Tokens[*oidc.IDTokenClaims], user idp.User, err error) {
	if authReq == nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	if _, _, actionErr := l.runPostExternalAuthenticationActions(&domain.ExternalUser{}, tokens, authReq, r, user, err); actionErr != nil {
		logging.WithError(err).Error("both external user authentication and action post authentication failed")
	}
	l.renderLogin(w, r, authReq, err)
}

// tokens extracts the oidc.Tokens for backwards compatibility of PostExternalAuthenticationActions
func tokens(session idp.Session) *oidc.Tokens[*oidc.IDTokenClaims] {
	switch s := session.(type) {
	case *openid.Session:
		return s.Tokens
	case *jwt.Session:
		return s.Tokens
	case *oauth.Session:
		return s.Tokens
	case *azuread.Session:
		return s.Tokens()
	case *apple.Session:
		return s.Tokens
	}
	return nil
}

func mapIDPUserToExternalUser(user idp.User, id string) *domain.ExternalUser {
	return &domain.ExternalUser{
		IDPConfigID:       id,
		ExternalUserID:    user.GetID(),
		PreferredUsername: user.GetPreferredUsername(),
		DisplayName:       user.GetDisplayName(),
		FirstName:         user.GetFirstName(),
		LastName:          user.GetLastName(),
		NickName:          user.GetNickname(),
		Email:             user.GetEmail(),
		IsEmailVerified:   user.IsEmailVerified(),
		PreferredLanguage: user.GetPreferredLanguage(),
		Phone:             user.GetPhone(),
		IsPhoneVerified:   user.IsPhoneVerified(),
	}
}

func mapExternalUserToLoginUser(externalUser *domain.ExternalUser, mustBeDomain bool) (*domain.Human, *domain.UserIDPLink, []*domain.Metadata) {
	username := externalUser.PreferredUsername
	if mustBeDomain {
		index := strings.LastIndex(username, "@")
		if index > 1 {
			username = username[:index]
		}
	}
	human := &domain.Human{
		Username: username,
		Profile: &domain.Profile{
			FirstName:         externalUser.FirstName,
			LastName:          externalUser.LastName,
			PreferredLanguage: externalUser.PreferredLanguage,
			NickName:          externalUser.NickName,
			DisplayName:       externalUser.DisplayName,
		},
		Email: &domain.Email{
			EmailAddress:    externalUser.Email,
			IsEmailVerified: externalUser.IsEmailVerified,
		},
	}
	if externalUser.Phone != "" {
		human.Phone = &domain.Phone{
			PhoneNumber:     externalUser.Phone,
			IsPhoneVerified: externalUser.IsPhoneVerified,
		}
	}
	externalIDP := &domain.UserIDPLink{
		IDPConfigID:    externalUser.IDPConfigID,
		ExternalUserID: externalUser.ExternalUserID,
		DisplayName:    externalUser.PreferredUsername,
	}
	return human, externalIDP, externalUser.Metadatas
}

func mapExternalNotFoundOptionFormDataToLoginUser(formData *externalNotFoundOptionFormData) *domain.ExternalUser {
	isEmailVerified := formData.ExternalEmailVerified && formData.Email == formData.ExternalEmail
	isPhoneVerified := formData.ExternalPhoneVerified && formData.Phone == formData.ExternalPhone
	return &domain.ExternalUser{
		IDPConfigID:       formData.ExternalIDPConfigID,
		ExternalUserID:    formData.ExternalIDPExtUserID,
		PreferredUsername: formData.Username,
		DisplayName:       string(formData.Email),
		FirstName:         formData.Firstname,
		LastName:          formData.Lastname,
		NickName:          formData.Nickname,
		Email:             formData.Email,
		IsEmailVerified:   isEmailVerified,
		Phone:             formData.Phone,
		IsPhoneVerified:   isPhoneVerified,
		PreferredLanguage: language.Make(formData.Language),
	}
}

func (l *Login) sessionParamsFromAuthRequest(ctx context.Context, authReq *domain.AuthRequest, identityProviderID string) []idp.Parameter {
	params := make([]idp.Parameter, 1, 2)
	params[0] = idp.UserAgentID(authReq.AgentID)

	if authReq.UserID != "" && identityProviderID != "" {
		links, err := l.getUserLinks(ctx, authReq.UserID, identityProviderID)
		if err != nil {
			logging.WithFields("authReqID", authReq.ID, "userID", authReq.UserID, "providerID", identityProviderID).WithError(err).Warn("failed to get user links for")
			return params
		}
		if len(links.Links) == 1 {
			return append(params, idp.LoginHintParam(links.Links[0].ProvidedUsername))
		}
	}
	if authReq.UserName != "" {
		return append(params, idp.LoginHintParam(authReq.UserName))
	}
	if authReq.LoginName != "" {
		return append(params, idp.LoginHintParam(authReq.LoginName))
	}
	if authReq.LoginHint != "" {
		return append(params, idp.LoginHintParam(authReq.LoginHint))
	}
	return params
}

func (l *Login) getUserLinks(ctx context.Context, userID, idpID string) (*query.IDPUserLinks, error) {
	userIDQuery, err := query.NewIDPUserLinksUserIDSearchQuery(userID)
	if err != nil {
		return nil, err
	}
	idpIDQuery, err := query.NewIDPUserLinkIDPIDSearchQuery(idpID)
	if err != nil {
		return nil, err
	}
	return l.query.IDPUserLinks(ctx,
		&query.IDPUserLinksSearchQuery{
			Queries: []query.SearchQuery{
				userIDQuery,
				idpIDQuery,
			},
		}, nil,
	)
}
