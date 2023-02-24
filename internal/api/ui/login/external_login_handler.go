package login

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
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

//
//func (l *Login) handleExternalLoginStep(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, selectedIDPConfigID string) {
//	for _, idp := range authReq.AllowedExternalIDPs {
//		if idp.IDPConfigID == selectedIDPConfigID {
//			l.handleIDP(w, r, authReq, selectedIDPConfigID)
//			return
//		}
//	}
//	l.renderLogin(w, r, authReq, errors.ThrowInvalidArgument(nil, "VIEW-Fsj7f", "Errors.User.ExternalIDP.NotAllowed"))
//}
//
//func (l *Login) handleExternalLogin(w http.ResponseWriter, r *http.Request) {
//	data := new(externalIDPData)
//	authReq, err := l.getAuthRequestAndParseData(r, data)
//	if err != nil {
//		l.renderError(w, r, authReq, err)
//		return
//	}
//	if authReq == nil {
//		l.defaultRedirect(w, r)
//		return
//	}
//	l.handleIDP(w, r, authReq, data.IDPConfigID)
//}
//
//func (l *Login) handleIDP(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, selectedIDPConfigID string) {
//	idpConfig, err := l.getIDPConfigByID(r, selectedIDPConfigID)
//	if err != nil {
//		l.renderError(w, r, authReq, err)
//		return
//	}
//	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
//	err = l.authRepo.SelectExternalIDP(r.Context(), authReq.ID, idpConfig.IDPConfigID, userAgentID)
//	if err != nil {
//		l.renderLogin(w, r, authReq, err)
//		return
//	}
//	if !idpConfig.IsOIDC {
//		l.handleJWTAuthorize(w, r, authReq, idpConfig)
//		return
//	}
//	l.handleOIDCAuthorize(w, r, authReq, idpConfig, EndpointExternalLoginCallback)
//}
//
//func (l *Login) handleOIDCAuthorize(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *iam_model.IDPConfigView, callbackEndpoint string) {
//	provider, err := l.getRPConfig(r.Context(), idpConfig, callbackEndpoint)
//	if err != nil {
//		l.renderLogin(w, r, authReq, err)
//		return
//	}
//	http.Redirect(w, r, rp.AuthURL(authReq.ID, provider, rp.WithPrompt(oidc.PromptSelectAccount)), http.StatusFound)
//}

func (l *Login) handleJWTAuthorize(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *iam_model.IDPConfigView) {
	redirect, err := url.Parse(idpConfig.JWTEndpoint)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	q := redirect.Query()
	q.Set(QueryAuthRequestID, authReq.ID)
	userAgentID, ok := http_mw.UserAgentIDFromCtx(r.Context())
	if !ok {
		l.renderLogin(w, r, authReq, errors.ThrowPreconditionFailed(nil, "LOGIN-dsgg3", "Errors.AuthRequest.UserAgentNotFound"))
		return
	}
	nonce, err := l.idpConfigAlg.Encrypt([]byte(userAgentID))
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	q.Set(queryUserAgentID, base64.RawURLEncoding.EncodeToString(nonce))
	redirect.RawQuery = q.Encode()
	http.Redirect(w, r, redirect.String(), http.StatusFound)
}

//
//func (l *Login) handleExternalLoginCallback(w http.ResponseWriter, r *http.Request) {
//	data := new(externalIDPCallbackData)
//	err := l.getParseData(r, data)
//	if err != nil {
//		l.renderError(w, r, nil, err)
//		return
//	}
//	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
//	authReq, err := l.authRepo.AuthRequestByID(r.Context(), data.State, userAgentID)
//	if err != nil {
//		l.renderError(w, r, authReq, err)
//		return
//	}
//	idpConfig, err := l.authRepo.GetIDPConfigByID(r.Context(), authReq.SelectedIDPConfigID)
//	if err != nil {
//		l.renderError(w, r, authReq, err)
//		return
//	}
//	if idpConfig.IsOIDC {
//		provider, err := l.getRPConfig(r.Context(), idpConfig, EndpointExternalLoginCallback)
//		if err != nil {
//			emtpyTokens := &oidc.Tokens{Token: &oauth2.Token{}}
//			if _, actionErr := l.runPostExternalAuthenticationActions(&domain.ExternalUser{}, emtpyTokens, authReq, r, idpConfig, err); actionErr != nil {
//				logging.WithError(err).Error("both external user authentication and action post authentication failed")
//			}
//
//			l.renderLogin(w, r, authReq, err)
//			return
//		}
//		tokens, err := rp.CodeExchange(r.Context(), data.Code, provider)
//		if err != nil {
//			emtpyTokens := &oidc.Tokens{Token: &oauth2.Token{}}
//			if _, actionErr := l.runPostExternalAuthenticationActions(&domain.ExternalUser{}, emtpyTokens, authReq, r, idpConfig, err); actionErr != nil {
//				logging.WithError(err).Error("both external user authentication and action post authentication failed")
//			}
//
//			l.renderLogin(w, r, authReq, err)
//			return
//		}
//		l.handleExternalUserAuthenticated(w, r, authReq, idpConfig, userAgentID, tokens)
//		return
//	}
//
//	err = errors.ThrowPreconditionFailed(nil, "RP-asff2", "Errors.ExternalIDP.IDPTypeNotImplemented")
//	emtpyTokens := &oidc.Tokens{Token: &oauth2.Token{}}
//	if _, actionErr := l.runPostExternalAuthenticationActions(&domain.ExternalUser{}, emtpyTokens, authReq, r, idpConfig, err); actionErr != nil {
//		logging.WithError(err).Error("both external user authentication and action post authentication failed")
//	}
//
//	l.renderError(w, r, authReq, err)
//}
//
//func (l *Login) getRPConfig(ctx context.Context, idpConfig *iam_model.IDPConfigView, callbackEndpoint string) (rp.RelyingParty, error) {
//	oidcClientSecret, err := crypto.DecryptString(idpConfig.OIDCClientSecret, l.idpConfigAlg)
//	if err != nil {
//		return nil, err
//	}
//	if idpConfig.OIDCIssuer != "" {
//		return rp.NewRelyingPartyOIDC(idpConfig.OIDCIssuer, idpConfig.OIDCClientID, oidcClientSecret, l.baseURL(ctx)+callbackEndpoint, idpConfig.OIDCScopes, rp.WithVerifierOpts(rp.WithIssuedAtOffset(3*time.Second)))
//	}
//	if idpConfig.OAuthAuthorizationEndpoint == "" || idpConfig.OAuthTokenEndpoint == "" {
//		return nil, errors.ThrowPreconditionFailed(nil, "RP-4n0fs", "Errors.IdentityProvider.InvalidConfig")
//	}
//	oauth2Config := &oauth2.Config{
//		ClientID:     idpConfig.OIDCClientID,
//		ClientSecret: oidcClientSecret,
//		Endpoint: oauth2.Endpoint{
//			AuthURL:  idpConfig.OAuthAuthorizationEndpoint,
//			TokenURL: idpConfig.OAuthTokenEndpoint,
//		},
//		RedirectURL: l.baseURL(ctx) + callbackEndpoint,
//		Scopes:      idpConfig.OIDCScopes,
//	}
//	return rp.NewRelyingPartyOAuth(oauth2Config, rp.WithVerifierOpts(rp.WithIssuedAtOffset(3*time.Second)))
//}

//
//func (l *Login) handleExternalUserAuthenticated(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *iam_model.IDPConfigView, userAgentID string, tokens *oidc.Tokens) {
//	externalUser := l.mapTokenToLoginUser(tokens, idpConfig)
//	externalUser, err := l.runPostExternalAuthenticationActions(externalUser, tokens, authReq, r, idpConfig, nil)
//	if err != nil {
//		l.renderError(w, r, authReq, err)
//		return
//	}
//
//	err = l.authRepo.CheckExternalUserLogin(setContext(r.Context(), ""), authReq.ID, userAgentID, externalUser, domain.BrowserInfoFromRequest(r))
//	if err != nil {
//		if errors.IsNotFound(err) {
//			err = nil
//		}
//		resourceOwner := authz.GetInstance(r.Context()).DefaultOrganisationID()
//
//		if authReq.RequestedOrgID != "" && authReq.RequestedOrgID != resourceOwner {
//			resourceOwner = authReq.RequestedOrgID
//		}
//
//		orgIAMPolicy, err := l.getOrgDomainPolicy(r, resourceOwner)
//		if err != nil {
//			l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, err)
//			return
//		}
//
//		human, idpLinking, _ := l.mapExternalUserToLoginUser(orgIAMPolicy, externalUser, idpConfig)
//		if !idpConfig.AutoRegister {
//			l.renderExternalNotFoundOption(w, r, authReq, orgIAMPolicy, human, idpLinking, err)
//			return
//		}
//		authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, userAgentID)
//		if err != nil {
//			l.renderExternalNotFoundOption(w, r, authReq, orgIAMPolicy, human, idpLinking, err)
//			return
//		}
//		l.handleAutoRegister(w, r, authReq, false)
//		return
//	}
//	if len(externalUser.Metadatas) > 0 {
//		authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, userAgentID)
//		if err != nil {
//			return
//		}
//		_, err = l.command.BulkSetUserMetadata(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, externalUser.Metadatas...)
//		if err != nil {
//			l.renderError(w, r, authReq, err)
//			return
//		}
//	}
//	l.renderNextStep(w, r, authReq)
//}

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
		idpConfig, err := l.authRepo.GetIDPConfigByID(r.Context(), authReq.SelectedIDPConfigID)
		if err != nil {
			l.renderError(w, r, authReq, err)
			return
		}
		linkingUser := authReq.LinkingUsers[len(authReq.LinkingUsers)-1]
		human, externalIDP, _ = l.mapExternalUserToLoginUser(orgIAMPolicy, linkingUser, idpConfig)
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
	linkingUser := l.mapExternalNotFoundOptionFormDataToLoginUser(data)
	l.registerExternalUser(w, r, authReq, linkingUser)
}

//
//func (l *Login) handleAutoRegister(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userNotFound bool) {
//	resourceOwner := authz.GetInstance(r.Context()).DefaultOrganisationID()
//
//	if authReq.RequestedOrgID != "" && authReq.RequestedOrgID != resourceOwner {
//		resourceOwner = authReq.RequestedOrgID
//	}
//
//	orgIamPolicy, err := l.getOrgDomainPolicy(r, resourceOwner)
//	if err != nil {
//		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, err)
//		return
//	}
//
//	idpConfig, err := l.authRepo.GetIDPConfigByID(r.Context(), authReq.SelectedIDPConfigID)
//	if err != nil {
//		l.renderExternalNotFoundOption(w, r, authReq, orgIamPolicy, nil, nil, err)
//		return
//	}
//
//	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
//	if len(authReq.LinkingUsers) == 0 {
//		l.renderError(w, r, authReq, errors.ThrowPreconditionFailed(nil, "LOGIN-asfg3", "Errors.ExternalIDP.NoExternalUserData"))
//		return
//	}
//
//	linkingUser := authReq.LinkingUsers[len(authReq.LinkingUsers)-1]
//	if userNotFound {
//		data := new(externalNotFoundOptionFormData)
//		err := l.getParseData(r, data)
//		if err != nil {
//			l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, err)
//			return
//		}
//		linkingUser = l.mapExternalNotFoundOptionFormDataToLoginUser(data)
//	}
//
//	user, externalIDP, metadata := l.mapExternalUserToLoginUser(orgIamPolicy, linkingUser, idpConfig)
//
//	user, metadata, err = l.runPreCreationActions(authReq, r, user, metadata, resourceOwner, domain.FlowTypeExternalAuthentication)
//	if err != nil {
//		l.renderExternalNotFoundOption(w, r, authReq, orgIamPolicy, nil, nil, err)
//		return
//	}
//	err = l.authRepo.AutoRegisterExternalUser(setContext(r.Context(), resourceOwner), user, externalIDP, nil, authReq.ID, userAgentID, resourceOwner, metadata, domain.BrowserInfoFromRequest(r))
//	if err != nil {
//		l.renderExternalNotFoundOption(w, r, authReq, orgIamPolicy, user, externalIDP, err)
//		return
//	}
//	authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
//	if err != nil {
//		l.renderError(w, r, authReq, err)
//		return
//	}
//	userGrants, err := l.runPostCreationActions(authReq.UserID, authReq, r, resourceOwner, domain.FlowTypeExternalAuthentication)
//	if err != nil {
//		l.renderError(w, r, authReq, err)
//		return
//	}
//	err = l.appendUserGrants(r.Context(), userGrants, resourceOwner)
//	if err != nil {
//		l.renderError(w, r, authReq, err)
//		return
//	}
//	l.renderNextStep(w, r, authReq)
//}

func (l *Login) mapExternalNotFoundOptionFormDataToLoginUser(formData *externalNotFoundOptionFormData) *domain.ExternalUser {
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

func (l *Login) mapTokenToLoginUser(tokens *oidc.Tokens, idpConfig *iam_model.IDPConfigView) *domain.ExternalUser {
	displayName := tokens.IDTokenClaims.GetPreferredUsername()
	if displayName == "" && tokens.IDTokenClaims.GetEmail() != "" {
		displayName = tokens.IDTokenClaims.GetEmail()
	}
	switch idpConfig.OIDCIDPDisplayNameMapping {
	case iam_model.OIDCMappingFieldEmail:
		if tokens.IDTokenClaims.IsEmailVerified() && tokens.IDTokenClaims.GetEmail() != "" {
			displayName = tokens.IDTokenClaims.GetEmail()
		}
	}

	externalUser := &domain.ExternalUser{
		IDPConfigID:       idpConfig.IDPConfigID,
		ExternalUserID:    tokens.IDTokenClaims.GetSubject(),
		PreferredUsername: tokens.IDTokenClaims.GetPreferredUsername(),
		DisplayName:       displayName,
		FirstName:         tokens.IDTokenClaims.GetGivenName(),
		LastName:          tokens.IDTokenClaims.GetFamilyName(),
		NickName:          tokens.IDTokenClaims.GetNickname(),
		Email:             tokens.IDTokenClaims.GetEmail(),
		IsEmailVerified:   tokens.IDTokenClaims.IsEmailVerified(),
		PreferredLanguage: tokens.IDTokenClaims.GetLocale(),
	}

	if tokens.IDTokenClaims.GetPhoneNumber() != "" {
		externalUser.Phone = tokens.IDTokenClaims.GetPhoneNumber()
		externalUser.IsPhoneVerified = tokens.IDTokenClaims.IsPhoneNumberVerified()
	}
	return externalUser
}
func (l *Login) mapExternalUserToLoginUser(orgIamPolicy *query.DomainPolicy, linkingUser *domain.ExternalUser, idpConfig *iam_model.IDPConfigView) (*domain.Human, *domain.UserIDPLink, []*domain.Metadata) {
	username := linkingUser.PreferredUsername
	switch idpConfig.OIDCUsernameMapping {
	case iam_model.OIDCMappingFieldEmail:
		if linkingUser.IsEmailVerified && linkingUser.Email != "" && username == "" {
			username = linkingUser.Email
		}
	}
	if username == "" {
		username = linkingUser.Email
	}

	if orgIamPolicy.UserLoginMustBeDomain {
		index := strings.LastIndex(username, "@")
		if index > 1 {
			username = username[:index]
		}
	}

	human := &domain.Human{
		Username: username,
		Profile: &domain.Profile{
			FirstName:         linkingUser.FirstName,
			LastName:          linkingUser.LastName,
			PreferredLanguage: linkingUser.PreferredLanguage,
			NickName:          linkingUser.NickName,
		},
		Email: &domain.Email{
			EmailAddress:    linkingUser.Email,
			IsEmailVerified: linkingUser.IsEmailVerified,
		},
	}
	if linkingUser.Phone != "" {
		human.Phone = &domain.Phone{
			PhoneNumber:     linkingUser.Phone,
			IsPhoneVerified: linkingUser.IsPhoneVerified,
		}
	}

	displayName := linkingUser.PreferredUsername
	switch idpConfig.OIDCIDPDisplayNameMapping {
	case iam_model.OIDCMappingFieldEmail:
		if linkingUser.IsEmailVerified && linkingUser.Email != "" && displayName == "" {
			displayName = linkingUser.Email
		}
	}
	if displayName == "" {
		displayName = linkingUser.Email
	}

	externalIDP := &domain.UserIDPLink{
		IDPConfigID:    idpConfig.IDPConfigID,
		ExternalUserID: linkingUser.ExternalUserID,
		DisplayName:    displayName,
	}
	return human, externalIDP, linkingUser.Metadatas
}
