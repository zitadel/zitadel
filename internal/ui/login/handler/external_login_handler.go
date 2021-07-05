package handler

import (
	"github.com/caos/oidc/pkg/client/rp"
	"github.com/caos/oidc/pkg/oidc"
	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	caos_errors "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"net/http"
	"strings"
	"time"
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
	Link         bool `schema:"link"`
	AutoRegister bool `schema:"autoregister"`
	ResetLinking bool `schema:"resetlinking"`
}

type externalNotFoundOptionData struct {
	baseData
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
		http.Redirect(w, r, l.zitadelURL, http.StatusFound)
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
		l.renderError(w, r, authReq, caos_errors.ThrowInternal(nil, "LOGIN-Rio9s", "Errors.User.ExternalIDP.IDPTypeNotImplemented"))
		return
	}
	l.handleOIDCAuthorize(w, r, authReq, idpConfig, EndpointExternalLoginCallback)
}

func (l *Login) handleOIDCAuthorize(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *iam_model.IDPConfigView, callbackEndpoint string) {
	provider := l.getRPConfig(w, r, authReq, idpConfig, callbackEndpoint)
	http.Redirect(w, r, rp.AuthURL(authReq.ID, provider, rp.WithPrompt(oidc.PromptSelectAccount)), http.StatusFound)
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
	provider := l.getRPConfig(w, r, authReq, idpConfig, EndpointExternalLoginCallback)
	tokens, err := rp.CodeExchange(r.Context(), data.Code, provider)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	l.handleExternalUserAuthenticated(w, r, authReq, idpConfig, userAgentID, tokens)
}

func (l *Login) getRPConfig(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *iam_model.IDPConfigView, callbackEndpoint string) rp.RelyingParty {
	oidcClientSecret, err := crypto.DecryptString(idpConfig.OIDCClientSecret, l.IDPConfigAesCrypto)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return nil
	}
	provider, err := rp.NewRelyingPartyOIDC(idpConfig.OIDCIssuer, idpConfig.OIDCClientID, oidcClientSecret, l.baseURL+callbackEndpoint, idpConfig.OIDCScopes, rp.WithVerifierOpts(rp.WithIssuedAtOffset(3*time.Second)))
	if err != nil {
		l.renderError(w, r, authReq, err)
		return nil
	}
	return provider
}

func (l *Login) handleExternalUserAuthenticated(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *iam_model.IDPConfigView, userAgentID string, tokens *oidc.Tokens) {
	externalUser := l.mapTokenToLoginUser(tokens, idpConfig)
	err := l.authRepo.CheckExternalUserLogin(r.Context(), authReq.ID, userAgentID, externalUser, domain.BrowserInfoFromRequest(r))
	if err != nil {
		if errors.IsNotFound(err) {
			err = nil
		}
		l.renderExternalNotFoundOption(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}

func (l *Login) renderExternalNotFoundOption(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	data := externalNotFoundOptionData{
		baseData: l.getBaseData(r, authReq, "ExternalNotFoundOption", errID, errMessage),
	}
	translator := l.getTranslator(authReq)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplExternalNotFoundOption], data, nil)
}

func (l *Login) handleExternalNotFoundOptionCheck(w http.ResponseWriter, r *http.Request) {
	data := new(externalNotFoundOptionFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, err)
		return
	}
	if data.Link {
		l.renderLogin(w, r, authReq, nil)
		return
	} else if data.ResetLinking {
		userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
		err = l.authRepo.ResetLinkingUsers(r.Context(), authReq.ID, userAgentID)
		if err != nil {
			l.renderExternalNotFoundOption(w, r, authReq, err)
		}
		l.handleLogin(w, r)
		return
	}
	l.handleAutoRegister(w, r, authReq)
}

func (l *Login) handleAutoRegister(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	iam, err := l.authRepo.GetIAM(r.Context())
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, err)
		return
	}

	resourceOwner := iam.GlobalOrgID
	memberRoles := []string{domain.RoleOrgProjectCreator}

	if authReq.RequestedOrgID != "" && authReq.RequestedOrgID != iam.GlobalOrgID {
		memberRoles = nil
		resourceOwner = authReq.RequestedOrgID
	}

	orgIamPolicy, err := l.getOrgIamPolicy(r, resourceOwner)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, err)
		return
	}

	idpConfig, err := l.authRepo.GetIDPConfigByID(r.Context(), authReq.SelectedIDPConfigID)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, err)
		return
	}

	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	user, externalIDP := l.mapExternalUserToLoginUser(orgIamPolicy, authReq.LinkingUsers[len(authReq.LinkingUsers)-1], idpConfig)
	err = l.authRepo.AutoRegisterExternalUser(setContext(r.Context(), resourceOwner), user, externalIDP, memberRoles, authReq.ID, userAgentID, resourceOwner, domain.BrowserInfoFromRequest(r))
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, err)
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
func (l *Login) mapExternalUserToLoginUser(orgIamPolicy *iam_model.OrgIAMPolicyView, linkingUser *domain.ExternalUser, idpConfig *iam_model.IDPConfigView) (*domain.Human, *domain.ExternalIDP) {
	username := linkingUser.PreferredUsername
	switch idpConfig.OIDCUsernameMapping {
	case iam_model.OIDCMappingFieldEmail:
		if linkingUser.IsEmailVerified && linkingUser.Email != "" {
			username = linkingUser.Email
		}
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

	externalIDP := &domain.ExternalIDP{
		IDPConfigID:    idpConfig.IDPConfigID,
		ExternalUserID: linkingUser.ExternalUserID,
		DisplayName:    displayName,
	}
	return human, externalIDP
}
