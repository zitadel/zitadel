package login

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/idp/providers/jwt"
	"github.com/zitadel/zitadel/internal/query"
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
	userAgentID, err := l.idpConfigAlg.DecryptString(id, l.idpConfigAlg.EncryptionKeyID())
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	authReq, err := l.authRepo.AuthRequestByID(r.Context(), data.AuthRequestID, userAgentID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	idpConfig, err := l.getIDPByID(r, authReq.SelectedIDPConfigID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if idpConfig.Type != domain.IDPTypeJWT {
		if err != nil {
			l.renderError(w, r, nil, err)
			return
		}
	}
	l.handleJWTExtraction(w, r, authReq, idpConfig)
}

func (l *Login) handleJWTExtraction(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, identityProvider *query.IDPTemplate) {
	token, err := getToken(r, identityProvider.JWTIDPTemplate.HeaderName)
	if err != nil {
		if _, actionErr := l.runPostExternalAuthenticationActions(new(domain.ExternalUser), nil, authReq, r, nil, err); actionErr != nil {
			logging.WithError(err).Error("both external user authentication and action post authentication failed")
		}

		l.renderError(w, r, authReq, err)
		return
	}
	provider, err := l.jwtProvider(identityProvider)
	if err != nil {
		if _, actionErr := l.runPostExternalAuthenticationActions(new(domain.ExternalUser), nil, authReq, r, nil, err); actionErr != nil {
			logging.WithError(err).Error("both external user authentication and action post authentication failed")
		}
		l.renderError(w, r, authReq, err)
		return
	}
	session := &jwt.Session{Provider: provider, Tokens: &oidc.Tokens{IDToken: token, Token: &oauth2.Token{}}}
	user, err := session.FetchUser(r.Context())
	if err != nil {
		if _, actionErr := l.runPostExternalAuthenticationActions(new(domain.ExternalUser), tokens(session), authReq, r, user, err); actionErr != nil {
			logging.WithError(err).Error("both external user authentication and action post authentication failed")
		}
		l.renderError(w, r, authReq, err)
		return
	}
	l.handleExternalUserAuthenticated(w, r, authReq, identityProvider, session, user, l.jwtCallback)
}

func (l *Login) jwtCallback(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	redirect, err := l.redirectToJWTCallback(r.Context(), authReq)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (l *Login) redirectToJWTCallback(ctx context.Context, authReq *domain.AuthRequest) (string, error) {
	redirect, err := url.Parse(l.baseURL(ctx) + EndpointJWTCallback)
	if err != nil {
		return "", err
	}
	q := redirect.Query()
	q.Set(QueryAuthRequestID, authReq.ID)
	nonce, err := l.idpConfigAlg.Encrypt([]byte(authReq.AgentID))
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
	userAgentID, err := l.idpConfigAlg.DecryptString(id, l.idpConfigAlg.EncryptionKeyID())
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	authReq, err := l.authRepo.AuthRequestByID(r.Context(), data.AuthRequestID, userAgentID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	idpConfig, err := l.getIDPByID(r, authReq.SelectedIDPConfigID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if idpConfig.Type != domain.IDPTypeJWT {
		l.renderLogin(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
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
