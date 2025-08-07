package azuread

import (
	"context"
	"net/http"
	"time"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
)

var _ idp.Session = (*Session)(nil)

// Session extends the [oauth.Session] to be able to handle the id_token and to implement the [idp.SessionSupportsMigration] functionality
type Session struct {
	*Provider
	Code string

	OAuthSession *oauth.Session
}

func NewSession(provider *Provider, code string) *Session {
	return &Session{Provider: provider, Code: code}
}

// GetAuth implements the [idp.Provider] interface by calling the wrapped [oauth.Session].
func (s *Session) GetAuth(ctx context.Context) (idp.Auth, error) {
	return s.oauth().GetAuth(ctx)
}

// RetrievePreviousID implements the [idp.SessionSupportsMigration] interface by returning the `sub` from the userinfo endpoint
func (s *Session) RetrievePreviousID() (string, error) {
	req, err := http.NewRequest("GET", userinfoEndpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("authorization", s.oauth().Tokens.TokenType+" "+s.oauth().Tokens.AccessToken)
	userinfo := new(oidc.UserInfo)
	if err := httphelper.HttpRequest(s.Provider.HttpClient(), req, &userinfo); err != nil {
		return "", err
	}
	return userinfo.Subject, nil
}

// PersistentParameters implements the [idp.Session] interface.
func (s *Session) PersistentParameters() map[string]any {
	return nil
}

// FetchUser implements the [idp.Session] interface.
// It will execute an OAuth 2.0 code exchange if needed to retrieve the access token,
// call the specified userEndpoint and map the received information into an [idp.User].
// In case of a specific TenantID as [TenantType] it will additionally extract the id_token and validate it.
func (s *Session) FetchUser(ctx context.Context) (user idp.User, err error) {
	user, err = s.oauth().FetchUser(ctx)
	if err != nil {
		return nil, err
	}
	// since azure will sign the id_token always with the issuer of the application it might differ from
	// the issuer the auth and token were based on, e.g. when allowing all account types to login,
	// then the auth endpoint must be `https://login.microsoftonline.com/common/oauth2/v2.0/authorize`
	// even though the issuer would be like `https://login.microsoftonline.com/d8cdd43f-fd94-4576-8deb-f3bfea72dc2e/v2.0`
	if s.Provider.tenant == CommonTenant ||
		s.Provider.tenant == OrganizationsTenant ||
		s.Provider.tenant == ConsumersTenant {
		return user, nil
	}
	idToken, ok := s.oauth().Tokens.Extra("id_token").(string)
	if !ok {
		return user, nil
	}
	idTokenVerifier := rp.NewIDTokenVerifier(s.Provider.issuer(), s.Provider.OAuthConfig().ClientID, rp.NewRemoteKeySet(s.Provider.HttpClient(), s.Provider.keysEndpoint()))
	s.oauth().Tokens.IDTokenClaims, err = rp.VerifyTokens[*oidc.IDTokenClaims](ctx, s.oauth().Tokens.AccessToken, idToken, idTokenVerifier)
	if err != nil {
		return nil, err
	}
	s.oauth().Tokens.IDToken = idToken
	return user, nil
}

func (s *Session) ExpiresAt() time.Time {
	if s.OAuthSession == nil {
		return time.Time{}
	}
	return s.OAuthSession.ExpiresAt()
}

// Tokens returns the [oidc.Tokens] of the underlying [oauth.Session].
func (s *Session) Tokens() *oidc.Tokens[*oidc.IDTokenClaims] {
	return s.oauth().Tokens
}

func (s *Session) oauth() *oauth.Session {
	if s.OAuthSession != nil {
		return s.OAuthSession
	}
	s.OAuthSession = &oauth.Session{
		Code:     s.Code,
		Provider: s.Provider.Provider,
	}
	return s.OAuthSession
}
