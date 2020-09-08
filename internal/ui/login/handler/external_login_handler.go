package handler

import (
	"fmt"
	"github.com/caos/oidc/pkg/rp"
	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/crypto"
	caos_errors "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"net/http"
	"path"
)

const (
	queryIDPConfigID = "idpConfigID"
	queryAggregateID = "aggregateID"
	queryState       = "state"
)

var (
	scopes = []string{"openid", "profile", "email"}
)

type externalIDPData struct {
	IDPConfigID string `schema:"idpConfigID"`
	AggregateID string `schema:"aggregateID"`
}

type externalIDPCallbackData struct {
	State string `schema:"state"`
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
	l.handleOIDCAuthorize(w, r, data, authReq, idpConfig)
}

func (l *Login) handleOIDCAuthorize(w http.ResponseWriter, r *http.Request, data *externalIDPData, authReq *model.AuthRequest, idpConfig *iam_model.IDPConfigView) {
	oidcClientSecret, err := crypto.DecryptString(idpConfig.OIDCClientSecret, l.IDPConfigAesCrypto)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	rpConfig := &rp.Config{
		ClientID:     idpConfig.OIDCClientID,
		ClientSecret: oidcClientSecret,
		Issuer:       idpConfig.OIDCIssuer,
		CallbackURL:  path.Join(l.renderer.pathPrefix, fmt.Sprintf("%s?%s=%s", EndpointExternalLoginCallback)),
		Scopes:       scopes,
	}

	provider, err := rp.NewDefaultRP(rpConfig)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	http.Redirect(w, r, provider.AuthURL(authReq.ID), http.StatusFound)
}

func (l *Login) handleExternalLoginCallback(w http.ResponseWriter, r *http.Request) {
	data := new(externalIDPCallbackData)
	err := l.getParseData(r, data)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	fmt.Println("Callback DATA: %v", data)
}
