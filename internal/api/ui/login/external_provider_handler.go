package login

import (
	"context"
	"net/http"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/google"
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
	ExternalIDPID              string
	ExternalIDPUserID          string
	ExternalIDPUserDisplayName string
	ShowUsername               bool
	ShowUsernameSuffix         bool
	OrgRegister                bool
	ExternalEmail              string
	ExternalEmailVerified      bool
	ExternalPhone              string
	ExternalPhoneVerified      bool
}

type externalRegisterFormData struct {
	ExternalIDPConfigID    string `schema:"external-idp-config-id"`
	ExternalIDPExtUserID   string `schema:"external-idp-ext-user-id"`
	ExternalIDPDisplayName string `schema:"external-idp-display-name"`
	ExternalEmail          string `schema:"external-email"`
	ExternalEmailVerified  bool   `schema:"external-email-verified"`
	Email                  string `schema:"email"`
	Username               string `schema:"username"`
	Firstname              string `schema:"firstname"`
	Lastname               string `schema:"lastname"`
	Nickname               string `schema:"nickname"`
	ExternalPhone          string `schema:"external-phone"`
	ExternalPhoneVerified  bool   `schema:"external-phone-verified"`
	Phone                  string `schema:"phone"`
	Language               string `schema:"language"`
	TermsConfirm           bool   `schema:"terms-confirm"`
}

func (l *Login) handleExternalLoginStep(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, selectedIDPID string) {
	for _, idp := range authReq.AllowedExternalIDPs {
		if idp.IDPConfigID == selectedIDPID {
			l.handleIDP(w, r, authReq, selectedIDPID)
			return
		}
	}
	l.renderLogin(w, r, authReq, errors.ThrowInvalidArgument(nil, "VIEW-Fsj7f", "Errors.User.ExternalIDP.NotAllowed"))
}

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

func (l *Login) handleExternalRegister(w http.ResponseWriter, r *http.Request) {
	data := new(externalIDPData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.handleIDP(w, r, authReq, data.IDPConfigID)
	//l.handleExternalRegisterByConfigID(w, r, authReq, data.IDPConfigID)
}
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
	case domain.IDPTypeOIDC:
		//oidc.New(provider.Name, provider.)
	case domain.IDPTypeJWT:
	case domain.IDPTypeOAuth:
	case domain.IDPTypeLDAP:
	case domain.IDPTypeAzureAD:
	case domain.IDPTypeGitHub:
	case domain.IDPTypeGitHubEE:
	case domain.IDPTypeGitLab:
	case domain.IDPTypeGitLabSelfHosted:
	case domain.IDPTypeGoogle:
		provider, err = l.googleProvider(r.Context(), identityProvider)
	}
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	session, err := provider.BeginAuth(r.Context(), authReq.ID)
	if err != nil {

	}
	http.Redirect(w, r, session.GetAuthURL(), http.StatusFound)
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
	return google.New(identityProvider.GoogleIDPTemplate.ClientID, secret, l.baseURL(ctx)+EndpointExternalLoginCallback, identityProvider.GoogleIDPTemplate.Scopes)
}

func providerOptions(identityProvider *query.IDPTemplate) {
	if identityProvider.IsCreationAllowed {
		openid.WithCreationAllowed()
	}
}

func (l *Login) handleExternalLoginCallback(w http.ResponseWriter, r *http.Request) {
	data := new(externalIDPCallbackData)
	err := l.getParseData(r, data)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	authReq, err := l.authRepo.AuthRequestByID(r.Context(), data.State, userAgentID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	identityProvider, err := l.getIDPByID(r, authReq.SelectedIDPConfigID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	var provider idp.Provider
	var session idp.Session
	switch identityProvider.Type {
	case domain.IDPTypeOIDC:
		//oidc.New(provider.Name, provider.)
	case domain.IDPTypeJWT:
	case domain.IDPTypeOAuth:
	case domain.IDPTypeLDAP:
	case domain.IDPTypeAzureAD:
	case domain.IDPTypeGitHub:
	case domain.IDPTypeGitHubEE:
	case domain.IDPTypeGitLab:
	case domain.IDPTypeGitLabSelfHosted:
	case domain.IDPTypeGoogle:
		provider, err = l.googleProvider(r.Context(), identityProvider)
		if err != nil {

		}
		session = &openid.Session{Provider: provider.(*google.Provider).Provider, AuthURL: "", Code: data.Code}
	}

	user, err := session.FetchUser(r.Context())
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	l.handleExternalUserAuthenticated(w, r, authReq, identityProvider, session, user)
	return
}

// tokens extracts the oidc.Tokens for backwards compatibility of PostExternalAuthenticationActions
func tokens(session idp.Session) *oidc.Tokens {
	if s, ok := session.(*openid.Session); ok {
		return s.Tokens
	}
	return nil
}

func (l *Login) handleExternalUserAuthenticated(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, provider *query.IDPTemplate, session idp.Session, user idp.User) {
	externalUser := mapIDPUserToExternalUser(user, provider.ID)
	externalUser, err := l.runPostExternalAuthenticationActions(externalUser, tokens(session), authReq, r, nil)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	err = l.authRepo.CheckExternalUserLogin(setContext(r.Context(), ""), authReq.ID, authReq.AgentID, externalUser, domain.BrowserInfoFromRequest(r))
	if err != nil {
		if errors.IsNotFound(err) { // TODO: handle error
			err = nil
		}
		l.externalUserNotExisting(w, r, authReq, provider, externalUser)
		return
	}
	if len(externalUser.Metadatas) > 0 {
		authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.ID)
		if err != nil {
			return
		}
		_, err = l.command.BulkSetUserMetadata(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, externalUser.Metadatas...)
		if err != nil {
			l.renderError(w, r, authReq, err)
			return
		}
	}
	l.renderNextStep(w, r, authReq)
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
	human := &domain.Human{
		Username: externalUser.PreferredUsername,
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
		DisplayName:    externalUser.DisplayName,
		NewVersion:     true,
	}
	return human, externalIDP, externalUser.Metadatas
}

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

	human, idpLinking, _ := mapExternalUserToLoginUser(externalUser, orgIAMPolicy.UserLoginMustBeDomain)
	if !provider.IsAutoCreation {
		l.renderExternalNotFoundOption(w, r, authReq, orgIAMPolicy, human, idpLinking, err)
		return
	}

	// reload auth request, to ensure current state (checked external login)
	authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, orgIAMPolicy, human, idpLinking, err)
		return
	}
	l.autoRegister(w, r, authReq)
	return
}

func (l *Login) autoRegister(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	if len(authReq.LinkingUsers) == 0 {
		l.renderError(w, r, authReq, errors.ThrowPreconditionFailed(nil, "LOGIN-asfg3", "Errors.ExternalIDP.NoExternalUserData"))
		return
	}

	// how do we get multiple and why do we use the last of them?
	linkingUser := authReq.LinkingUsers[len(authReq.LinkingUsers)-1]

	l.registerExternalUser(w, r, authReq, linkingUser)

	//resourceOwner := authz.GetInstance(r.Context()).DefaultOrganisationID()
	//
	//if authReq.RequestedOrgID != "" && authReq.RequestedOrgID != resourceOwner {
	//	resourceOwner = authReq.RequestedOrgID
	//}
	//
	//orgIamPolicy, err := l.getOrgDomainPolicy(r, resourceOwner)
	//if err != nil {
	//	l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, err)
	//	return
	//}
	//
	////idpConfig, err := l.authRepo.GetIDPConfigByID(r.Context(), authReq.SelectedIDPConfigID)
	////if err != nil {
	////	l.renderExternalNotFoundOption(w, r, authReq, orgIamPolicy, nil, nil, err)
	////	return
	////}
	//
	//userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	//if len(authReq.LinkingUsers) == 0 {
	//	l.renderError(w, r, authReq, errors.ThrowPreconditionFailed(nil, "LOGIN-asfg3", "Errors.ExternalIDP.NoExternalUserData"))
	//	return
	//}
	//
	//linkingUser := authReq.LinkingUsers[len(authReq.LinkingUsers)-1]
	////if userNotFound {
	////	data := new(externalNotFoundOptionFormData)
	////	err := l.getParseData(r, data)
	////	if err != nil {
	////		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, err)
	////		return
	////	}
	////	linkingUser = l.mapExternalNotFoundOptionFormDataToLoginUser(data)
	////}
	//
	//user, externalIDP, metadata := mapExternalUserToLoginUser(linkingUser, orgIamPolicy.UserLoginMustBeDomain)
	//
	//user, metadata, err = l.runPreCreationActions(authReq, r, user, metadata, resourceOwner, domain.FlowTypeExternalAuthentication)
	//if err != nil {
	//	l.renderExternalNotFoundOption(w, r, authReq, orgIamPolicy, nil, nil, err)
	//	return
	//}
	//err = l.authRepo.AutoRegisterExternalUser(setContext(r.Context(), resourceOwner), user, externalIDP, nil, authReq.ID, userAgentID, resourceOwner, metadata, domain.BrowserInfoFromRequest(r))
	//if err != nil {
	//	l.renderExternalNotFoundOption(w, r, authReq, orgIamPolicy, user, externalIDP, err)
	//	return
	//}
	//authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
	//if err != nil {
	//	l.renderError(w, r, authReq, err)
	//	return
	//}
	//userGrants, err := l.runPostCreationActions(authReq.UserID, authReq, r, resourceOwner, domain.FlowTypeExternalAuthentication)
	//if err != nil {
	//	l.renderError(w, r, authReq, err)
	//	return
	//}
	//err = l.appendUserGrants(r.Context(), userGrants, resourceOwner)
	//if err != nil {
	//	l.renderError(w, r, authReq, err)
	//	return
	//}
	//l.renderNextStep(w, r, authReq)
}

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
	//
	////idpConfig, err := l.authRepo.GetIDPConfigByID(r.Context(), authReq.SelectedIDPConfigID)
	////if err != nil {
	////	l.renderExternalNotFoundOption(w, r, authReq, orgIamPolicy, nil, nil, err)
	////	return
	////}
	//
	//userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	//if len(authReq.LinkingUsers) == 0 {
	//	l.renderError(w, r, authReq, errors.ThrowPreconditionFailed(nil, "LOGIN-asfg3", "Errors.ExternalIDP.NoExternalUserData"))
	//	return
	//}
	//
	//linkingUser := authReq.LinkingUsers[len(authReq.LinkingUsers)-1]
	////if userNotFound {
	////	data := new(externalNotFoundOptionFormData)
	////	err := l.getParseData(r, data)
	////	if err != nil {
	////		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, err)
	////		return
	////	}
	////	linkingUser = l.mapExternalNotFoundOptionFormDataToLoginUser(data)
	////}

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

func (l *Login) renderExternalNotFoundOption(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, orgIAMPolicy *query.DomainPolicy, human *domain.Human, externalIDP *domain.UserIDPLink, err error) {
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

	if human == nil || externalIDP == nil {
		//idpConfig, err := l.getIDPByID(r, authReq.SelectedIDPConfigID)
		//if err != nil {
		//	l.renderError(w, r, authReq, err)
		//	return
		//}
		linkingUser := authReq.LinkingUsers[len(authReq.LinkingUsers)-1]
		human, externalIDP, _ = mapExternalUserToLoginUser(linkingUser, orgIAMPolicy.UserLoginMustBeDomain)
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
		ExternalIDPID:              externalIDP.IDPConfigID,
		ExternalIDPUserID:          externalIDP.ExternalUserID,
		ExternalIDPUserDisplayName: externalIDP.DisplayName,
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

func (l *Login) handleExternalNotFoundOptionCheck(w http.ResponseWriter, r *http.Request) {
	data := new(externalNotFoundOptionFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, err)
		return
	}
	if data.Link {
		l.renderLogin(w, r, authReq, nil)
		return
	} else if data.ResetLinking {
		userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
		err = l.authRepo.ResetLinkingUsers(r.Context(), authReq.ID, userAgentID)
		if err != nil {
			l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, err)
		}
		l.handleLogin(w, r)
		return
	}
	linkingUser := mapExternalNotFoundOptionFormDataToLoginUser(data)
	l.registerExternalUser(w, r, authReq, linkingUser)
}

func mapExternalNotFoundOptionFormDataToLoginUser(formData *externalNotFoundOptionFormData) *domain.ExternalUser {
	isEmailVerified := formData.ExternalEmailVerified && formData.Email == formData.ExternalEmail
	isPhoneVerified := formData.ExternalPhoneVerified && formData.Phone == formData.ExternalPhone
	return &domain.ExternalUser{
		IDPConfigID:       formData.ExternalIDPConfigID,
		ExternalUserID:    formData.ExternalIDPExtUserID,
		PreferredUsername: formData.Username,
		DisplayName:       formData.Email,
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
