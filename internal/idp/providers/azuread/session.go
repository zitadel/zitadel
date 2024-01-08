package azuread

import (
	"context"
	"net/http"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
)

// Session extends the [oauth.Session] to extend it with the [idp.SessionSupportsMigration] functionality
type Session struct {
	*Provider
	Code string

	OAuthSession *oauth.Session
}

func (s *Session) GetAuth(ctx context.Context) (content string, redirect bool) {
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

// FetchUser implements the [idp.Session] interface.
// It will execute an OAuth 2.0 code exchange if needed to retrieve the access token,
// call the specified userEndpoint and map the received information into an [idp.User].
func (s *Session) FetchUser(ctx context.Context) (user idp.User, err error) {
	user, err = s.oauth().FetchUser(ctx)
	if err != nil {
		return nil, err
	}
	// since azure will sign the
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

func (s *Session) Tokens() *oidc.Tokens[*oidc.IDTokenClaims] {
	return s.oauth().Tokens
}
