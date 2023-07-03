package oidc

import (
	"time"

	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/command"
)

const IDPrefix = "V2_"

type AuthRequestV2 struct {
	*command.AuthRequest
	UserID string
}

func (a *AuthRequestV2) GetID() string {
	return IDPrefix + a.ID
}

func (a *AuthRequestV2) GetACR() string {
	return "" //PLANNED: impl
}

func (a *AuthRequestV2) GetAMR() []string {
	//TODO: get from linked session?
	return nil
}

func (a *AuthRequestV2) GetAudience() []string {
	return a.Audience
}

func (a *AuthRequestV2) GetAuthTime() time.Time {
	return time.Time{} //TODO: get from linked session?
}

func (a *AuthRequestV2) GetClientID() string {
	return a.ClientID
}

func (a *AuthRequestV2) GetCodeChallenge() *oidc.CodeChallenge {
	return CodeChallengeToOIDC(a.CodeChallenge)
}

func (a *AuthRequestV2) GetNonce() string {
	return a.Nonce
}

func (a *AuthRequestV2) GetRedirectURI() string {
	return a.RedirectURI
}

func (a *AuthRequestV2) GetResponseType() oidc.ResponseType {
	return ResponseTypeToOIDC(a.ResponseType)
}

func (a *AuthRequestV2) GetResponseMode() oidc.ResponseMode {
	return ""
}

func (a *AuthRequestV2) GetScopes() []string {
	return a.Scope
}

func (a *AuthRequestV2) GetState() string {
	return a.State
}

func (a *AuthRequestV2) GetSubject() string {
	return a.UserID
}

// Done implements the [op.AuthRequest] interface and will be used to determine if the user has authenticated
// and the auth request can be redirected back to the client.
// Since AuthRequestV2 handles the redirect directly in the gRPC OIDC Service, it will always return false.
func (a *AuthRequestV2) Done() bool {
	return false
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
