package handler

import (
	"github.com/caos/zitadel/internal/v2/domain"
	"net/http"
	"strings"

	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/oidc/pkg/rp"
	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	caos_errors "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

func (l *Login) handleExternalRegister(w http.ResponseWriter, r *http.Request) {
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
	l.handleOIDCAuthorize(w, r, authReq, idpConfig, EndpointExternalRegisterCallback)
}

func (l *Login) handleExternalRegisterCallback(w http.ResponseWriter, r *http.Request) {
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
	provider := l.getRPConfig(w, r, authReq, idpConfig, EndpointExternalRegisterCallback)
	tokens, err := rp.CodeExchange(r.Context(), data.Code, provider)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	l.handleExternalUserRegister(w, r, authReq, idpConfig, userAgentID, tokens)
}

func (l *Login) handleExternalUserRegister(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *iam_model.IDPConfigView, userAgentID string, tokens *oidc.Tokens) {
	iam, err := l.authRepo.GetIAM(r.Context())
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	resourceOwner := iam.GlobalOrgID
	memberRoles := []string{orgProjectCreatorRole}

	if authReq.RequestedOrgID != "" && authReq.RequestedOrgID != iam.GlobalOrgID {
		memberRoles = nil
		resourceOwner = authReq.RequestedOrgID
	}
	orgIamPolicy, err := l.getOrgIamPolicy(r, resourceOwner)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	user, externalIDP := l.mapTokenToLoginHumanAndExternalIDP(orgIamPolicy, tokens, idpConfig)
	_, err = l.command.RegisterHuman(setContext(r.Context(), resourceOwner), resourceOwner, user, externalIDP, memberRoles)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}

func (l *Login) mapTokenToLoginHumanAndExternalIDP(orgIamPolicy *iam_model.OrgIAMPolicyView, tokens *oidc.Tokens, idpConfig *iam_model.IDPConfigView) (*domain.Human, *domain.ExternalIDP) {
	username := tokens.IDTokenClaims.GetPreferredUsername()
	switch idpConfig.OIDCUsernameMapping {
	case iam_model.OIDCMappingFieldEmail:
		if tokens.IDTokenClaims.IsEmailVerified() && tokens.IDTokenClaims.GetEmail() != "" {
			username = tokens.IDTokenClaims.GetEmail()
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
			FirstName:         tokens.IDTokenClaims.GetGivenName(),
			LastName:          tokens.IDTokenClaims.GetFamilyName(),
			PreferredLanguage: tokens.IDTokenClaims.GetLocale(),
			NickName:          tokens.IDTokenClaims.GetNickname(),
		},
		Email: &domain.Email{
			EmailAddress:    tokens.IDTokenClaims.GetEmail(),
			IsEmailVerified: tokens.IDTokenClaims.IsEmailVerified(),
		},
	}

	if tokens.IDTokenClaims.GetPhoneNumber() != "" {
		human.Phone = &domain.Phone{
			PhoneNumber:     tokens.IDTokenClaims.GetPhoneNumber(),
			IsPhoneVerified: tokens.IDTokenClaims.IsPhoneNumberVerified(),
		}
	}

	displayName := tokens.IDTokenClaims.GetPreferredUsername()
	switch idpConfig.OIDCIDPDisplayNameMapping {
	case iam_model.OIDCMappingFieldEmail:
		if tokens.IDTokenClaims.IsEmailVerified() && tokens.IDTokenClaims.GetEmail() != "" {
			displayName = tokens.IDTokenClaims.GetEmail()
		}
	}

	externalIDP := &domain.ExternalIDP{
		IDPConfigID:    idpConfig.IDPConfigID,
		ExternalUserID: tokens.IDTokenClaims.GetSubject(),
		DisplayName:    displayName,
	}
	return human, externalIDP
}
