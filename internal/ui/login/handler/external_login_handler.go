package handler

import (
	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/oidc/pkg/rp"
	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/crypto"
	caos_errors "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
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
}

type externalNotFoundOptionData struct {
	baseData
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
	idpConfig, err := l.getIDPConfigByID(r, data.IDPConfigID)
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

func (l *Login) handleOIDCAuthorize(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, idpConfig *iam_model.IDPConfigView, callbackEndpoint string) {
	provider := l.getRPConfig(w, r, authReq, idpConfig, callbackEndpoint)
	http.Redirect(w, r, provider.AuthURL(authReq.ID), http.StatusFound)
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
	tokens, err := provider.CodeExchange(r.Context(), data.Code)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	l.handleExternalUserAuthenticated(w, r, authReq, idpConfig, userAgentID, tokens)
}

func (l *Login) getRPConfig(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, idpConfig *iam_model.IDPConfigView, callbackEndpoint string) rp.DelegationTokenExchangeRP {
	oidcClientSecret, err := crypto.DecryptString(idpConfig.OIDCClientSecret, l.IDPConfigAesCrypto)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return nil
	}
	rpConfig := &rp.Config{
		ClientID:     idpConfig.OIDCClientID,
		ClientSecret: oidcClientSecret,
		Issuer:       idpConfig.OIDCIssuer,
		CallbackURL:  l.baseURL + callbackEndpoint,
		Scopes:       idpConfig.OIDCScopes,
	}

	provider, err := rp.NewDefaultRP(rpConfig, rp.WithVerifierOpts(rp.WithIssuedAtOffset(3*time.Second)))
	if err != nil {
		l.renderError(w, r, authReq, err)
		return nil
	}
	return provider
}

func (l *Login) handleExternalUserAuthenticated(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, idpConfig *iam_model.IDPConfigView, userAgentID string, tokens *oidc.Tokens) {
	externalUser := l.mapTokenToLoginUser(tokens, idpConfig)
	err := l.authRepo.CheckExternalUserLogin(r.Context(), authReq.ID, userAgentID, externalUser)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, nil)
		return
	}
	l.renderNextStep(w, r, authReq)
}

func (l *Login) renderExternalNotFoundOption(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	var errType, errMessage string
	if err != nil {
		errMessage = l.getErrorMessage(r, err)
	}
	data := externalNotFoundOptionData{
		baseData: l.getBaseData(r, authReq, "ExternalNotFoundOption", errType, errMessage),
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplExternalNotFoundOption], data, nil)
}

func (l *Login) handleExternalNotFoundOptionCheck(w http.ResponseWriter, r *http.Request) {
	data := new(externalNotFoundOptionFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if data.Link {
		l.renderLogin(w, r, authReq, nil)
		return
	}
	l.handleAutoRegister(w, r, authReq)
}

func (l *Login) handleAutoRegister(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest) {
	orgIamPolicy, err := l.getOrgIamPolicy(r, authReq.GetScopeOrgID())
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, err)
		return
	}
	iam, err := l.authRepo.GetIAM(r.Context())
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, err)
		return
	}
	resourceOwner := iam.GlobalOrgID
	member := &org_model.OrgMember{
		ObjectRoot: models.ObjectRoot{AggregateID: iam.GlobalOrgID},
		Roles:      []string{orgProjectCreatorRole},
	}
	if authReq.GetScopeOrgID() != iam.GlobalOrgID && authReq.GetScopeOrgID() != "" {
		member = nil
		resourceOwner = authReq.GetScopeOrgID()
	}

	idpConfig, err := l.authRepo.GetIDPConfigByID(r.Context(), authReq.SelectedIDPConfigID)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, err)
		return
	}

	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	user, externalIDP := l.mapExternalUserToLoginUser(orgIamPolicy, authReq.LinkingUsers[len(authReq.LinkingUsers)-1], idpConfig)
	err = l.authRepo.AutoRegisterExternalUser(setContext(r.Context(), resourceOwner), user, externalIDP, member, authReq.ID, userAgentID, resourceOwner)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}

func (l *Login) mapTokenToLoginUser(tokens *oidc.Tokens, idpConfig *iam_model.IDPConfigView) *model.ExternalUser {
	displayName := tokens.IDTokenClaims.PreferredUsername
	switch idpConfig.OIDCIDPDisplayNameMapping {
	case iam_model.OIDCMappingFieldEmail:
		if tokens.IDTokenClaims.EmailVerified && tokens.IDTokenClaims.Email != "" {
			displayName = tokens.IDTokenClaims.Email
		}
	}

	externalUser := &model.ExternalUser{
		IDPConfigID:       idpConfig.IDPConfigID,
		ExternalUserID:    tokens.IDTokenClaims.Subject,
		PreferredUsername: tokens.IDTokenClaims.PreferredUsername,
		DisplayName:       displayName,
		FirstName:         tokens.IDTokenClaims.GivenName,
		LastName:          tokens.IDTokenClaims.FamilyName,
		NickName:          tokens.IDTokenClaims.Nickname,
		Email:             tokens.IDTokenClaims.Email,
		IsEmailVerified:   tokens.IDTokenClaims.EmailVerified,
	}

	if tokens.IDTokenClaims.PhoneNumber != "" {
		externalUser.Phone = tokens.IDTokenClaims.PhoneNumber
		externalUser.IsPhoneVerified = tokens.IDTokenClaims.PhoneNumberVerified
	}
	return externalUser
}

func (l *Login) mapExternalUserToLoginUser(orgIamPolicy *org_model.OrgIAMPolicy, linkingUser *model.ExternalUser, idpConfig *iam_model.IDPConfigView) (*usr_model.User, *usr_model.ExternalIDP) {
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

	user := &usr_model.User{
		UserName: username,
		Human: &usr_model.Human{
			Profile: &usr_model.Profile{
				FirstName:         linkingUser.FirstName,
				LastName:          linkingUser.LastName,
				PreferredLanguage: linkingUser.PreferredLanguage,
				NickName:          linkingUser.NickName,
			},
			Email: &usr_model.Email{
				EmailAddress:    linkingUser.Email,
				IsEmailVerified: linkingUser.IsEmailVerified,
			},
		},
	}
	if linkingUser.Phone != "" {
		user.Phone = &usr_model.Phone{
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

	externalIDP := &usr_model.ExternalIDP{
		IDPConfigID: idpConfig.IDPConfigID,
		UserID:      linkingUser.ExternalUserID,
		DisplayName: displayName,
	}
	return user, externalIDP
}
