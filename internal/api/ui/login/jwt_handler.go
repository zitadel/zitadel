package login

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
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

func (l *Login) handleJWTExtraction(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *query.IDPTemplate) {
	var token string
	var err error
	//token, err := getToken(r, idpConfig.JWTHeaderName) // TODO: template
	if err != nil {
		emtpyTokens := &oidc.Tokens{Token: &oauth2.Token{}}
		if _, actionErr := l.runPostExternalAuthenticationActions(&domain.ExternalUser{}, emtpyTokens, authReq, r, err); actionErr != nil {
			logging.WithError(err).Error("both external user authentication and action post authentication failed")
		}

		l.renderError(w, r, authReq, err)
		return
	}
	tokenClaims, err := validateToken(r.Context(), token, idpConfig)
	tokens := &oidc.Tokens{IDToken: token, IDTokenClaims: tokenClaims, Token: &oauth2.Token{}}
	if err != nil {
		if _, actionErr := l.runPostExternalAuthenticationActions(&domain.ExternalUser{}, tokens, authReq, r, err); actionErr != nil {
			logging.WithError(err).Error("both external user authentication and action post authentication failed")
		}
		l.renderError(w, r, authReq, err)
		return
	}
	var externalUser *domain.ExternalUser
	//externalUser := l.mapTokenToLoginUser(tokens, idpConfig)
	externalUser, err = l.runPostExternalAuthenticationActions(externalUser, tokens, authReq, r, nil)
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
	redirect, err := l.redirectToJWTCallback(r.Context(), authReq)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (l *Login) jwtExtractionUserNotFound(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *query.IDPTemplate, tokens *oidc.Tokens, err error) {
	if errors.IsNotFound(err) {
		err = nil
	}
	if !idpConfig.IsAutoCreation {
		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, err)
		return
	}
	authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	resourceOwner := l.getOrgID(r, authReq)
	orgIamPolicy, err := l.getOrgDomainPolicy(r, resourceOwner)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	user, externalIDP, metadata := mapExternalUserToLoginUser(authReq.LinkingUsers[len(authReq.LinkingUsers)-1], orgIamPolicy.UserLoginMustBeDomain)
	user, metadata, err = l.runPreCreationActions(authReq, r, user, metadata, resourceOwner, domain.FlowTypeExternalAuthentication)
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
	userGrants, err := l.runPostCreationActions(authReq.UserID, authReq, r, resourceOwner, domain.FlowTypeExternalAuthentication)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	err = l.appendUserGrants(r.Context(), userGrants, resourceOwner)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	redirect, err := l.redirectToJWTCallback(r.Context(), authReq)
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
	if idpConfig.Type != domain.IDPTypeJWT { // TODO: ?
		l.renderLogin(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}

func validateToken(ctx context.Context, token string, config *query.IDPTemplate) (oidc.IDTokenClaims, error) {
	logging.Log("LOGIN-ADf42").Debug("begin token validation")
	offset := 3 * time.Second
	maxAge := time.Hour
	claims := oidc.EmptyIDTokenClaims()
	payload, err := oidc.ParseToken(token, claims)
	if err != nil {
		return nil, err
	}

	//if err = oidc.CheckIssuer(claims, config.JWTIssuer); err != nil { //TODO: check
	//	return nil, err
	//}

	logging.Debug("begin signature validation")
	_ = payload
	//keySet := rp.NewRemoteKeySet(http.DefaultClient, config.JWTKeysEndpoint)
	//if err = oidc.CheckSignature(ctx, token, payload, claims, nil, keySet); err != nil {
	//	return nil, err
	//}

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
