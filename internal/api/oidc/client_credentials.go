package oidc

import (
	"time"

	"github.com/zitadel/oidc/v2/pkg/oidc"
	"github.com/zitadel/oidc/v2/pkg/op"
)

type clientCredentialsRequest struct {
	sub    string
	aud    []string
	scopes []string
}

func (c *clientCredentialsRequest) GetSubject() string {
	return c.sub
}

func (c *clientCredentialsRequest) GetAudience() []string {
	return c.aud
}

func (c *clientCredentialsRequest) GetScopes() []string {
	return c.scopes
}

type clientCredentialsClient struct {
	id string
}

func (c *clientCredentialsClient) AccessTokenType() op.AccessTokenType {
	return op.AccessTokenTypeBearer
}

func (c *clientCredentialsClient) GetID() string {
	return c.id
}

func (c *clientCredentialsClient) RedirectURIs() []string {
	return nil
}

func (c *clientCredentialsClient) PostLogoutRedirectURIs() []string {
	return nil
}

func (c *clientCredentialsClient) ApplicationType() op.ApplicationType {
	return op.ApplicationTypeWeb
}

func (c *clientCredentialsClient) AuthMethod() oidc.AuthMethod {
	return oidc.AuthMethodBasic
}

func (c *clientCredentialsClient) ResponseTypes() []oidc.ResponseType {
	return nil
}

func (c *clientCredentialsClient) GrantTypes() []oidc.GrantType {
	return []oidc.GrantType{
		oidc.GrantTypeClientCredentials,
	}
}

func (c *clientCredentialsClient) LoginURL(_ string) string {
	return ""
}

func (c *clientCredentialsClient) AccessTokenLifetime() time.Duration {
	return 0
}

func (c *clientCredentialsClient) IDTokenLifetime() time.Duration {
	return 0
}

func (c *clientCredentialsClient) DevMode() bool {
	return false
}

func (c *clientCredentialsClient) RestrictAdditionalIdTokenScopes() func(scopes []string) []string {
	return nil
}

func (c *clientCredentialsClient) RestrictAdditionalAccessTokenScopes() func(scopes []string) []string {
	return func(scopes []string) []string {
		return scopes
	}
}

func (c *clientCredentialsClient) IsScopeAllowed(scope string) bool {
	return true
}

func (c *clientCredentialsClient) IDTokenUserinfoClaimsAssertion() bool {
	return false
}

func (c *clientCredentialsClient) ClockSkew() time.Duration {
	return 0
}
