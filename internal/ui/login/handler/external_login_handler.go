package handler

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/oidc"
	"golang.org/x/oauth2"

	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
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
	OrgRegister                bool
	ExternalEmail              string
	ExternalEmailVerified      bool
	ExternalPhone              string
	ExternalPhoneVerified      bool
}

func (l *Login) handleExternalLoginStep(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, selectedIDPConfigID string) {
	for _, idp := range authReq.AllowedExternalIDPs {
		if idp.IDPConfigID == selectedIDPConfigID {
			l.handleIDP(w, r, authReq, selectedIDPConfigID)
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

func (l *Login) handleIDP(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, selectedIDPConfigID string) {
	idpConfig, err := l.getIDPConfigByID(r, selectedIDPConfigID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	err = l.authRepo.SelectExternalIDP(r.Context(), authReq.ID, idpConfig.IDPConfigID, userAgentID)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	if !idpConfig.IsOIDC {
		l.handleJWTAuthorize(w, r, authReq, idpConfig)
		return
	}
	l.handleOIDCAuthorize(w, r, authReq, idpConfig, EndpointExternalLoginCallback)
}

func (l *Login) handleOIDCAuthorize(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *iam_model.IDPConfigView, callbackEndpoint string) {
	provider, err := l.getRPConfig(idpConfig, callbackEndpoint)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	http.Redirect(w, r, rp.AuthURL(authReq.ID, provider, rp.WithPrompt(oidc.PromptSelectAccount)), http.StatusFound)
}

func (l *Login) handleJWTAuthorize(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *iam_model.IDPConfigView) {
	redirect, err := url.Parse(idpConfig.JWTEndpoint)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	q := redirect.Query()
	q.Set(queryAuthRequestID, authReq.ID)
	userAgentID, ok := http_mw.UserAgentIDFromCtx(r.Context())
	if !ok {
		l.renderLogin(w, r, authReq, caos_errors.ThrowPreconditionFailed(nil, "LOGIN-dsgg3", "Errors.AuthRequest.UserAgentNotFound"))
		return
	}
	nonce, err := l.IDPConfigAesCrypto.Encrypt([]byte(userAgentID))
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	q.Set(queryUserAgentID, base64.RawURLEncoding.EncodeToString(nonce))
	redirect.RawQuery = q.Encode()
	http.Redirect(w, r, redirect.String(), http.StatusFound)
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
	idpConfig, err := l.authRepo.GetIDPConfigByID(r.Context(), authReq.SelectedIDPConfigID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if idpConfig.IsOIDC {
		provider, err := l.getRPConfig(idpConfig, EndpointExternalLoginCallback)
		if err != nil {
			l.renderLogin(w, r, authReq, err)
			return
		}
		tokens, err := rp.CodeExchange(r.Context(), data.Code, provider)
		if err != nil {
			l.renderLogin(w, r, authReq, err)
			return
		}
		l.handleExternalUserAuthenticated(w, r, authReq, idpConfig, userAgentID, tokens)
		return
	}
	l.renderError(w, r, authReq, caos_errors.ThrowPreconditionFailed(nil, "RP-asff2", "Errors.ExternalIDP.IDPTypeNotImplemented"))
}

func (l *Login) getRPConfig(idpConfig *iam_model.IDPConfigView, callbackEndpoint string) (rp.RelyingParty, error) {
	oidcClientSecret, err := crypto.DecryptString(idpConfig.OIDCClientSecret, l.IDPConfigAesCrypto)
	if err != nil {
		return nil, err
	}
	if idpConfig.OIDCIssuer != "" {
		return rp.NewRelyingPartyOIDC(idpConfig.OIDCIssuer, idpConfig.OIDCClientID, oidcClientSecret, l.baseURL+callbackEndpoint, idpConfig.OIDCScopes, rp.WithVerifierOpts(rp.WithIssuedAtOffset(3*time.Second)))
	}
	if idpConfig.OAuthAuthorizationEndpoint == "" || idpConfig.OAuthTokenEndpoint == "" {
		return nil, caos_errors.ThrowPreconditionFailed(nil, "RP-4n0fs", "Errors.IdentityProvider.InvalidConfig")
	}
	oauth2Config := &oauth2.Config{
		ClientID:     idpConfig.OIDCClientID,
		ClientSecret: oidcClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  idpConfig.OAuthAuthorizationEndpoint,
			TokenURL: idpConfig.OAuthTokenEndpoint,
		},
		RedirectURL: l.baseURL + callbackEndpoint,
		Scopes:      idpConfig.OIDCScopes,
	}
	return rp.NewRelyingPartyOAuth(oauth2Config, rp.WithVerifierOpts(rp.WithIssuedAtOffset(3*time.Second)))
}

func (l *Login) handleExternalUserAuthenticated(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *iam_model.IDPConfigView, userAgentID string, tokens *oidc.Tokens) {
	externalUser := l.mapTokenToLoginUser(tokens, idpConfig)
	externalUser, err := l.customExternalUserMapping(r.Context(), externalUser, tokens, authReq, idpConfig)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	err = l.authRepo.CheckExternalUserLogin(setContext(r.Context(), ""), authReq.ID, userAgentID, externalUser, domain.BrowserInfoFromRequest(r))
	if err != nil {
		if errors.IsNotFound(err) {
			err = nil
		}
		iam, err := l.query.IAMByID(r.Context(), domain.IAMID)
		if err != nil {
			l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, nil, err)
			return
		}

		resourceOwner := iam.GlobalOrgID

		if authReq.RequestedOrgID != "" && authReq.RequestedOrgID != iam.GlobalOrgID {
			resourceOwner = authReq.RequestedOrgID
		}

		orgIAMPolicy, err := l.getOrgIamPolicy(r, resourceOwner)
		if err != nil {
			l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, nil, err)
			return
		}

		human, idpLinking, _ := l.mapExternalUserToLoginUser(orgIAMPolicy, externalUser, idpConfig)
		if !idpConfig.AutoRegister {
			l.renderExternalNotFoundOption(w, r, authReq, iam, orgIAMPolicy, human, idpLinking, err)
			return
		}
		authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, userAgentID)
		if err != nil {
			l.renderExternalNotFoundOption(w, r, authReq, iam, orgIAMPolicy, human, idpLinking, err)
			return
		}
		l.handleAutoRegister(w, r, authReq)
		return
	}
	if len(externalUser.Metadatas) > 0 {
		authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, userAgentID)
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

func (l *Login) renderExternalNotFoundOption(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, iam *query.IAM, orgIAMPolicy *query.OrgIAMPolicy, human *domain.Human, externalIDP *domain.UserIDPLink, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	if orgIAMPolicy == nil {
		iam, err = l.query.IAMByID(r.Context(), domain.IAMID)
		if err != nil {
			l.renderError(w, r, authReq, err)
			return
		}
		resourceOwner := iam.GlobalOrgID

		if authReq.RequestedOrgID != "" && authReq.RequestedOrgID != iam.GlobalOrgID {
			resourceOwner = authReq.RequestedOrgID
		}

		orgIAMPolicy, err = l.getOrgIamPolicy(r, resourceOwner)
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

	data := externalNotFoundOptionData{
		baseData: l.getBaseData(r, authReq, "ExternalNotFoundOption", errID, errMessage),
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
		OrgRegister:                orgIAMPolicy.UserLoginMustBeDomain,
	}
	if human.Phone != nil {
		data.Phone = human.PhoneNumber
		data.ExternalPhone = human.PhoneNumber
		data.ExternalPhoneVerified = human.IsPhoneVerified
	}
	translator := l.getTranslator(authReq)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplExternalNotFoundOption], data, nil)
}

func (l *Login) handleExternalNotFoundOptionCheck(w http.ResponseWriter, r *http.Request) {
	data := new(externalNotFoundOptionFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, nil, err)
		return
	}
	if data.Link {
		l.renderLogin(w, r, authReq, nil)
		return
	} else if data.ResetLinking {
		userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
		err = l.authRepo.ResetLinkingUsers(r.Context(), authReq.ID, userAgentID)
		if err != nil {
			l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, nil, err)
		}
		l.handleLogin(w, r)
		return
	}
	l.handleAutoRegister(w, r, authReq)
}

func (l *Login) handleAutoRegister(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	iam, err := l.query.IAMByID(r.Context(), domain.IAMID)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, nil, err)
		return
	}

	resourceOwner := iam.GlobalOrgID
	memberRoles := []string{domain.RoleSelfManagementGlobal}

	if authReq.RequestedOrgID != "" && authReq.RequestedOrgID != iam.GlobalOrgID {
		memberRoles = nil
		resourceOwner = authReq.RequestedOrgID
	}

	orgIamPolicy, err := l.getOrgIamPolicy(r, resourceOwner)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, nil, err)
		return
	}

	idpConfig, err := l.authRepo.GetIDPConfigByID(r.Context(), authReq.SelectedIDPConfigID)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, iam, orgIamPolicy, nil, nil, err)
		return
	}

	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	if len(authReq.LinkingUsers) == 0 {
		l.renderError(w, r, authReq, caos_errors.ThrowPreconditionFailed(nil, "LOGIN-asfg3", "Errors.ExternalIDP.NoExternalUserData"))
		return
	}
	linkingUser := authReq.LinkingUsers[len(authReq.LinkingUsers)-1]
	user, externalIDP, metadata := l.mapExternalUserToLoginUser(orgIamPolicy, linkingUser, idpConfig)
	user, metadata, err = l.customExternalUserToLoginUserMapping(user, nil, authReq, idpConfig, metadata, resourceOwner)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, iam, orgIamPolicy, nil, nil, err)
		return
	}
	err = l.authRepo.AutoRegisterExternalUser(setContext(r.Context(), resourceOwner), user, externalIDP, memberRoles, authReq.ID, userAgentID, resourceOwner, metadata, domain.BrowserInfoFromRequest(r))
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, iam, orgIamPolicy, user, externalIDP, err)
		return
	}
	authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	userGrants, err := l.customGrants(authReq.UserID, nil, authReq, idpConfig, resourceOwner)
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
	}

	if tokens.IDTokenClaims.GetPhoneNumber() != "" {
		externalUser.Phone = tokens.IDTokenClaims.GetPhoneNumber()
		externalUser.IsPhoneVerified = tokens.IDTokenClaims.IsPhoneNumberVerified()
	}
	return externalUser
}
func (l *Login) mapExternalUserToLoginUser(orgIamPolicy *query.OrgIAMPolicy, linkingUser *domain.ExternalUser, idpConfig *iam_model.IDPConfigView) (*domain.Human, *domain.UserIDPLink, []*domain.Metadata) {
	username := linkingUser.PreferredUsername
	switch idpConfig.OIDCUsernameMapping {
	case iam_model.OIDCMappingFieldEmail:
		if linkingUser.IsEmailVerified && linkingUser.Email != "" {
			username = linkingUser.Email
		}
	}
	if username == "" {
		username = linkingUser.Email
	}

	if orgIamPolicy.UserLoginMustBeDomain {
		splittedUsername := strings.Split(username, "@")
		if len(splittedUsername) > 1 {
			username = splittedUsername[0]
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
		if linkingUser.IsEmailVerified && linkingUser.Email != "" {
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
