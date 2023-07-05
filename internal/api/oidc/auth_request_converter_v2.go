package oidc

import (
	"time"

	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
)

type AuthRequestV2 struct {
	id            string
	amr           []string
	audience      []string
	authTime      time.Time
	clientID      string
	codeChallenge *domain.OIDCCodeChallenge
	nonce         string
	redirectURI   string
	responseType  domain.OIDCResponseType
	scope         []string
	state         string
	sessionID     string
	userID        string
}

func (a *AuthRequestV2) GetID() string {
	return a.id
}

func (a *AuthRequestV2) GetACR() string {
	return "" //PLANNED: impl
}

func (a *AuthRequestV2) GetAMR() []string {
	return a.amr
}

func (a *AuthRequestV2) GetAudience() []string {
	return a.audience
}

func (a *AuthRequestV2) GetAuthTime() time.Time {
	return a.authTime
}

func (a *AuthRequestV2) GetClientID() string {
	return a.clientID
}

func (a *AuthRequestV2) GetCodeChallenge() *oidc.CodeChallenge {
	return CodeChallengeToOIDC(a.codeChallenge)
}

func (a *AuthRequestV2) GetNonce() string {
	return a.nonce
}

func (a *AuthRequestV2) GetRedirectURI() string {
	return a.redirectURI
}

func (a *AuthRequestV2) GetResponseType() oidc.ResponseType {
	return ResponseTypeToOIDC(a.responseType)
}

func (a *AuthRequestV2) GetResponseMode() oidc.ResponseMode {
	return ""
}

func (a *AuthRequestV2) GetScopes() []string {
	return a.scope
}

func (a *AuthRequestV2) GetState() string {
	return a.state
}

func (a *AuthRequestV2) GetSubject() string {
	return a.userID
}

func (a *AuthRequestV2) Done() bool {
	return a.userID != "" && a.sessionID != ""
}

type RefreshTokenRequestV2 struct {
	*command.OIDCSessionWriteModel
	RequestedScopes []string
}

func (r *RefreshTokenRequestV2) GetAMR() []string {
	return r.AuthMethodsReferences
}

func (r *RefreshTokenRequestV2) GetAudience() []string {
	return r.Audience
}

func (r *RefreshTokenRequestV2) GetAuthTime() time.Time {
	return r.AuthTime
}

func (r *RefreshTokenRequestV2) GetClientID() string {
	return r.ClientID
}

func (r *RefreshTokenRequestV2) GetScopes() []string {
	return r.Scope
}

func (r *RefreshTokenRequestV2) GetSubject() string {
	return r.UserID
}

func (r *RefreshTokenRequestV2) SetCurrentScopes(scopes []string) {
	r.RequestedScopes = scopes
}
