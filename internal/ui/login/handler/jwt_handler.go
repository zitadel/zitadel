package handler

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/oidc"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
)

type jwtRequest struct {
	AuthRequestID string `schema:"authRequestID"`
	UserAgentID   string `schema:"userAgentID"`
}

func (l *Login) handleJWTRequest(w http.ResponseWriter, r *http.Request) {
	data := new(jwtRequest)
	err := l.getParseData(r, data)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	if data.AuthRequestID == "" || data.UserAgentID == "" {
		l.renderError(w, r, nil, errors.ThrowInvalidArgument(nil, "LOGIN-adfzz", "Errors.AuthRequest.MissingParameters"))
		return
	}
	id, err := base64.RawURLEncoding.DecodeString(data.UserAgentID)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	userAgentID, err := l.IDPConfigAesCrypto.DecryptString(id, l.IDPConfigAesCrypto.EncryptionKeyID())
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
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
		if err != nil {
			l.renderError(w, r, nil, err)
			return
		}
	}
	l.handleJWTExtraction(w, r, authReq, idpConfig)
}

func (l *Login) handleJWTExtraction(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *iam_model.IDPConfigView) {
	token, err := getToken(r, idpConfig.JWTHeaderName)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	tokenClaims, err := validateToken(r.Context(), token, idpConfig)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	tokens := &oidc.Tokens{IDToken: token, IDTokenClaims: tokenClaims}
	externalUser := l.mapTokenToLoginUser(tokens, idpConfig)
	externalUser, err = l.customExternalUserMapping(r.Context(), externalUser, tokens, authReq, idpConfig)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	metadata := externalUser.Metadatas
	err = l.authRepo.CheckExternalUserLogin(setContext(r.Context(), ""), authReq.ID, authReq.AgentID, externalUser, domain.BrowserInfoFromRequest(r))
	if err != nil {
		l.jwtExtractionUserNotFound(w, r, authReq, idpConfig, tokens, err)
		return
	}
	if len(metadata) > 0 {
		authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
		if err != nil {
			l.renderError(w, r, authReq, err)
			return
		}
		_, err = l.command.BulkSetUserMetadata(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, metadata...)
		if err != nil {
			l.renderError(w, r, authReq, err)
			return
		}
	}
	redirect, err := l.redirectToJWTCallback(authReq)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (l *Login) jwtExtractionUserNotFound(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *iam_model.IDPConfigView, tokens *oidc.Tokens, err error) {
	if errors.IsNotFound(err) {
		err = nil
	}
	if !idpConfig.AutoRegister {
		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, nil, err)
		return
	}
	authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	resourceOwner := l.getOrgID(r, authReq)
	orgIamPolicy, err := l.getOrgIamPolicy(r, resourceOwner)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	user, externalIDP, metadata := l.mapExternalUserToLoginUser(orgIamPolicy, authReq.LinkingUsers[len(authReq.LinkingUsers)-1], idpConfig)
	user, metadata, err = l.customExternalUserToLoginUserMapping(user, tokens, authReq, idpConfig, metadata, resourceOwner)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	err = l.authRepo.AutoRegisterExternalUser(setContext(r.Context(), resourceOwner), user, externalIDP, nil, authReq.ID, authReq.AgentID, resourceOwner, metadata, domain.BrowserInfoFromRequest(r))
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	userGrants, err := l.customGrants(authReq.UserID, tokens, authReq, idpConfig, resourceOwner)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	err = l.appendUserGrants(r.Context(), userGrants, resourceOwner)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	redirect, err := l.redirectToJWTCallback(authReq)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	http.Redirect(w, r, redirect, http.StatusFound)
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

func (l *Login) redirectToJWTCallback(authReq *domain.AuthRequest) (string, error) {
	redirect, err := url.Parse(l.baseURL + EndpointJWTCallback)
	if err != nil {
		return "", err
	}
	q := redirect.Query()
	q.Set(queryAuthRequestID, authReq.ID)
	nonce, err := l.IDPConfigAesCrypto.Encrypt([]byte(authReq.AgentID))
	if err != nil {
		return "", err
	}
	q.Set(queryUserAgentID, base64.RawURLEncoding.EncodeToString(nonce))
	redirect.RawQuery = q.Encode()
	return redirect.String(), nil
}

func (l *Login) handleJWTCallback(w http.ResponseWriter, r *http.Request) {
	data := new(jwtRequest)
	err := l.getParseData(r, data)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	id, err := base64.RawURLEncoding.DecodeString(data.UserAgentID)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	userAgentID, err := l.IDPConfigAesCrypto.DecryptString(id, l.IDPConfigAesCrypto.EncryptionKeyID())
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
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
	l.renderNextStep(w, r, authReq)
}

func validateToken(ctx context.Context, token string, config *iam_model.IDPConfigView) (oidc.IDTokenClaims, error) {
	logging.Log("LOGIN-ADf42").Debug("begin token validation")
	offset := 3 * time.Second
	maxAge := time.Hour
	claims := oidc.EmptyIDTokenClaims()
	payload, err := oidc.ParseToken(token, claims)
	if err != nil {
		return nil, err
	}

	if err = oidc.CheckIssuer(claims, config.JWTIssuer); err != nil {
		return nil, err
	}

	logging.Log("LOGIN-Dfg22").Debug("begin signature validation")
	keySet := rp.NewRemoteKeySet(http.DefaultClient, config.JWTKeysEndpoint)
	if err = oidc.CheckSignature(ctx, token, payload, claims, nil, keySet); err != nil {
		return nil, err
	}

	if !claims.GetExpiration().IsZero() {
		if err = oidc.CheckExpiration(claims, offset); err != nil {
			return nil, err
		}
	}

	if !claims.GetIssuedAt().IsZero() {
		if err = oidc.CheckIssuedAt(claims, maxAge, offset); err != nil {
			return nil, err
		}
	}
	return claims, nil
}

func getToken(r *http.Request, headerName string) (string, error) {
	if headerName == "" {
		headerName = http_util.Authorization
	}
	auth := r.Header.Get(headerName)
	if auth == "" {
		return "", errors.ThrowInvalidArgument(nil, "LOGIN-adh42", "Errors.AuthRequest.TokenNotFound")
	}
	return strings.TrimPrefix(auth, oidc.PrefixBearer), nil
}
