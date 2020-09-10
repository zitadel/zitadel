package handler

import (
	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/oidc/pkg/rp"
	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/crypto"
	caos_errors "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"net/http"
	"time"
)

const (
	queryIDPConfigID = "idpConfigID"
)

type externalIDPData struct {
	IDPConfigID string `schema:"idpConfigID"`
}

type externalIDPCallbackData struct {
	State string `schema:"state"`
	Code  string `schema:"code"`
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
		l.renderLogin(w, r, authReq, err)
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

	return &model.ExternalUser{
		IDPConfigID:    idpConfig.IDPConfigID,
		ExternalUserID: tokens.IDTokenClaims.Subject,
		DisplayName:    displayName,
	}
}
