package login

import (
	"context"
	"net/http"
	"strings"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/azuread"
	"github.com/zitadel/zitadel/internal/idp/providers/github"
	"github.com/zitadel/zitadel/internal/idp/providers/gitlab"
	"github.com/zitadel/zitadel/internal/idp/providers/google"
	"github.com/zitadel/zitadel/internal/idp/providers/jwt"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	"github.com/zitadel/zitadel/internal/query"
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
	l.renderLogin(w, r, authReq, errors.ThrowInvalidArgument(nil, "VIEW-Fsj7f", "Errors.User.ExternalIDP.NotAllowed"))
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
	authReq, err := l.getAuthRequestAndParseData(r, data)
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
	case domain.IDPTypeLDAP,
		domain.IDPTypeUnspecified:
		fallthrough
	default:
		l.renderLogin(w, r, authReq, errors.ThrowInvalidArgument(nil, "LOGIN-AShek", "Errors.ExternalIDP.IDPTypeNotImplemented"))
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
	http.Redirect(w, r, session.GetAuthURL(), http.StatusFound)
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
		session = &oauth.Session{Provider: provider.(*azuread.Provider).Provider, Code: data.Code}
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
	case domain.IDPTypeJWT,
		domain.IDPTypeLDAP,
		domain.IDPTypeUnspecified:
		fallthrough
	default:
		l.renderLogin(w, r, authReq, errors.ThrowInvalidArgument(nil, "LOGIN-SFefg", "Errors.ExternalIDP.IDPTypeNotImplemented"))
		return
	}

	user, err := session.FetchUser(r.Context())
	if err != nil {
		l.externalAuthFailed(w, r, authReq, tokens(session), user, err)
		return
	}
	l.handleExternalUserAuthenticated(w, r, authReq, identityProvider, session, user, l.renderNextStep)
}

// handleExternalUserAuthenticated maps the IDP user, checks for a corresponding externalID
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
	externalUser, err := l.runPostExternalAuthenticationActions(externalUser, tokens(session), authReq, r, user, nil)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	err = l.authRepo.CheckExternalUserLogin(setContext(r.Context(), ""), authReq.ID, authReq.AgentID, externalUser, domain.BrowserInfoFromRequest(r))
	if err != nil {
		if !errors.IsNotFound(err) {
			l.renderError(w, r, authReq, err)
			return
		}
		l.externalUserNotExisting(w, r, authReq, provider, externalUser)
		return
	}
	if provider.IsAutoUpdate || len(externalUser.Metadatas) > 0 {
		// read current auth request state (incl. authorized user)
		authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
		if err != nil {
			l.renderError(w, r, authReq, err)
			return
		}
	}
	if provider.IsAutoUpdate {
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

// externalUserNotExisting is called if an externalAuthentication couldn't find a corresponding externalID
// possible solutions are:
//
// * auto creation
// * external not found overview:
//   - creation by user
//   - linking to existing user
func (l *Login) externalUserNotExisting(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, provider *query.IDPTemplate, externalUser *domain.ExternalUser) {
	resourceOwner := authz.GetInstance(r.Context()).DefaultOrganisationID()

	if authReq.RequestedOrgID != "" && authReq.RequestedOrgID != resourceOwner {
		resourceOwner = authReq.RequestedOrgID
	}

	orgIAMPolicy, err := l.getOrgDomainPolicy(r, resourceOwner)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, err)
		return
	}

	human, idpLink, _ := mapExternalUserToLoginUser(externalUser, orgIAMPolicy.UserLoginMustBeDomain)
	// if auto creation or creation itself is disabled, send the user to the notFoundOption
	if !provider.IsCreationAllowed || !provider.IsAutoCreation {
		l.renderExternalNotFoundOption(w, r, authReq, orgIAMPolicy, human, idpLink, err)
		return
	}

	// reload auth request, to ensure current state (checked external login)
	authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, orgIAMPolicy, human, idpLink, err)
		return
	}
	l.autoCreateExternalUser(w, r, authReq)
}

// autoCreateExternalUser takes the externalUser and creates it automatically (without user interaction)
func (l *Login) autoCreateExternalUser(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	if len(authReq.LinkingUsers) == 0 {
		l.renderError(w, r, authReq, errors.ThrowPreconditionFailed(nil, "LOGIN-asfg3", "Errors.ExternalIDP.NoExternalUserData"))
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
	if orgIAMPolicy == nil {
		resourceOwner := authz.GetInstance(r.Context()).DefaultOrganisationID()

		if authReq.RequestedOrgID != "" && authReq.RequestedOrgID != resourceOwner {
			resourceOwner = authReq.RequestedOrgID
		}

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

	var resourceOwner string
	if authReq != nil {
		resourceOwner = authReq.RequestedOrgID
	}
	if resourceOwner == "" {
		resourceOwner = authz.GetInstance(r.Context()).DefaultOrganisationID()
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

	translator := l.getTranslator(r.Context(), authReq)
	data := externalNotFoundOptionData{
		baseData: l.getBaseData(r, authReq, "ExternalNotFound.Title", "ExternalNotFound.Description", errID, errMessage),
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
	authReq, err := l.getAuthRequestAndParseData(r, data)
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
			l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, errors.ThrowPreconditionFailed(nil, "LOGIN-AS3ff", "Errors.ExternalIDP.LinkingNotAllowed"))
			return
		}
		l.renderLogin(w, r, authReq, nil)
		return
	}
	// if the user selects the creation button
	if !idpTemplate.IsCreationAllowed {
		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, errors.ThrowPreconditionFailed(nil, "LOGIN-dsfd3", "Errors.ExternalIDP.CreationNotAllowed"))
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
	resourceOwner := authz.GetInstance(r.Context()).DefaultOrganisationID()

	if authReq.RequestedOrgID != "" && authReq.RequestedOrgID != resourceOwner {
		resourceOwner = authReq.RequestedOrgID
	}

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
	user, err := l.query.GetUserByID(ctx, true, authReq.UserID, false)
	if err != nil {
		return err
	}
	if user.Human == nil {
		return errors.ThrowPreconditionFailed(nil, "LOGIN-WLTce", "Errors.User.NotHuman")
	}
	if externalUser.Email != "" && externalUser.Email != user.Human.Email && externalUser.IsEmailVerified != user.Human.IsEmailVerified {
		emailCodeGenerator, err := l.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyEmailCode, l.userCodeAlg)
		logging.WithFields("authReq", authReq.ID, "user", authReq.UserID).OnError(err).Error("unable to update email")
		if err == nil {
			_, err = l.command.ChangeHumanEmail(setContext(ctx, authReq.UserOrgID),
				&domain.Email{
					ObjectRoot:      models.ObjectRoot{AggregateID: authReq.UserID},
					EmailAddress:    externalUser.Email,
					IsEmailVerified: externalUser.IsEmailVerified,
				},
				emailCodeGenerator)
			logging.WithFields("authReq", authReq.ID, "user", authReq.UserID).OnError(err).Error("unable to update email")
		}
	}
	if externalUser.Phone != "" && externalUser.Phone != user.Human.Phone && externalUser.IsPhoneVerified != user.Human.IsPhoneVerified {
		phoneCodeGenerator, err := l.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyPhoneCode, l.userCodeAlg)
		logging.WithFields("authReq", authReq.ID, "user", authReq.UserID).OnError(err).Error("unable to update phone")
		if err == nil {
			_, err = l.command.ChangeHumanPhone(setContext(ctx, authReq.UserOrgID),
				&domain.Phone{
					ObjectRoot:      models.ObjectRoot{AggregateID: authReq.UserID},
					PhoneNumber:     externalUser.Phone,
					IsPhoneVerified: externalUser.IsPhoneVerified,
				},
				authReq.UserOrgID,
				phoneCodeGenerator)
			logging.WithFields("authReq", authReq.ID, "user", authReq.UserID).OnError(err).Error("unable to update phone")
		}
	}
	if externalUser.FirstName != user.Human.FirstName ||
		externalUser.LastName != user.Human.LastName ||
		externalUser.NickName != user.Human.NickName ||
		externalUser.DisplayName != user.Human.DisplayName ||
		externalUser.PreferredLanguage != user.Human.PreferredLanguage {
		_, err = l.command.ChangeHumanProfile(setContext(ctx, authReq.UserOrgID), &domain.Profile{
			ObjectRoot:        models.ObjectRoot{AggregateID: authReq.UserID},
			FirstName:         externalUser.FirstName,
			LastName:          externalUser.LastName,
			NickName:          externalUser.NickName,
			DisplayName:       externalUser.DisplayName,
			PreferredLanguage: externalUser.PreferredLanguage,
			Gender:            user.Human.Gender,
		})
		logging.WithFields("authReq", authReq.ID, "user", authReq.UserID).OnError(err).Error("unable to update profile")
	}
	return nil
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

func (l *Login) externalAuthFailed(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, tokens *oidc.Tokens, user idp.User, err error) {
	if _, actionErr := l.runPostExternalAuthenticationActions(&domain.ExternalUser{}, tokens, authReq, r, user, err); actionErr != nil {
		logging.WithError(err).Error("both external user authentication and action post authentication failed")
	}
	l.renderLogin(w, r, authReq, err)
}

// tokens extracts the oidc.Tokens for backwards compatibility of PostExternalAuthenticationActions
func tokens(session idp.Session) *oidc.Tokens {
	switch s := session.(type) {
	case *openid.Session:
		return s.Tokens
	case *jwt.Session:
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

func (l *Login) sessionParamsFromAuthRequest(ctx context.Context, authReq *domain.AuthRequest, identityProviderID string) []any {
	params := []any{authReq.AgentID}

	if authReq.UserID != "" && identityProviderID != "" {
		links, err := l.getUserLinks(ctx, authReq.UserID, identityProviderID)
		if err != nil {
			logging.WithFields("authReqID", authReq.ID, "userID", authReq.UserID, "providerID", identityProviderID).WithError(err).Warn("failed to get user links for")
			return params
		}
		if len(links.Links) == 1 {
			return append(params, keyAndValueToAuthURLOpt("login_hint", links.Links[0].ProvidedUsername))
		}
	}
	if authReq.UserName != "" {
		return append(params, keyAndValueToAuthURLOpt("login_hint", authReq.UserName))
	}
	if authReq.LoginName != "" {
		return append(params, keyAndValueToAuthURLOpt("login_hint", authReq.LoginName))
	}
	if authReq.LoginHint != "" {
		return append(params, keyAndValueToAuthURLOpt("login_hint", authReq.LoginHint))
	}
	return params
}

func keyAndValueToAuthURLOpt(key, value string) rp.AuthURLOpt {
	return func() []oauth2.AuthCodeOption {
		return []oauth2.AuthCodeOption{oauth2.SetAuthURLParam(key, value)}
	}
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
		}, false,
	)
}
