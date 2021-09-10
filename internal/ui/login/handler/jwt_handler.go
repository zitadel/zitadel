package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/caos/oidc/pkg/client/rp"
	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type jwtRequest struct {
	AuthRequestID string `schema:"authRequestID"`
	UserAgentID   string `schema:"userAgentID"`
}

func (l *Login) handleJWTRequest(w http.ResponseWriter, r *http.Request) {
	//TODO: error handling
	data := new(jwtRequest)
	err := l.getParseData(r, data)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	userAgentID, err := l.IDPConfigAesCrypto.DecryptString([]byte(data.UserAgentID), l.IDPConfigAesCrypto.EncryptionKeyID())
	_ = err
	authReq, err := l.authRepo.AuthRequestByID(r.Context(), data.AuthRequestID, userAgentID)
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
		l.renderLogin(w, r, authReq, err)
		return
	}
	l.handleJWTExtraction(w, r, authReq, idpConfig)
}

func (l *Login) handleJWTExtraction(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *iam_model.IDPConfigView) {
	token, err := getToken(r)
	tokenClaims, err := validateToken(r.Context(), token, idpConfig)

	externalUser := l.mapTokenToLoginUser(&oidc.Tokens{IDTokenClaims: tokenClaims}, idpConfig)
	err = l.authRepo.CheckExternalUserLogin(r.Context(), authReq.ID, authReq.AgentID, externalUser, domain.BrowserInfoFromRequest(r))
	if err != nil {
		if errors.IsNotFound(err) {
			err = nil
		}
		l.renderExternalNotFoundOption(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}

func validateToken(ctx context.Context, token string, config *iam_model.IDPConfigView) (oidc.IDTokenClaims, error) {
	offset := 3 * time.Second
	maxAge := time.Hour
	claims := oidc.EmptyIDTokenClaims()
	payload, err := oidc.ParseToken(token, claims)
	if err != nil {
		return nil, err
	}

	if err := oidc.CheckSubject(claims); err != nil {
		return nil, err
	}

	if err = oidc.CheckIssuer(claims, config.JWTIssuer); err != nil {
		return nil, err
	}

	keySet := rp.NewRemoteKeySet(http.DefaultClient, config.JWTKeysEndpoint)
	if err = oidc.CheckSignature(ctx, token, payload, claims, nil, keySet); err != nil {
		return nil, err
	}

	if !claims.GetExpiration().IsZero() {
		if err = oidc.CheckExpiration(claims, offset); err != nil {
			return nil, err
		}
	}

	if err = oidc.CheckIssuedAt(claims, maxAge, offset); err != nil {
		return nil, err
	}
	return claims, nil
}

func getToken(r *http.Request) (string, error) {
	return "", nil
}
