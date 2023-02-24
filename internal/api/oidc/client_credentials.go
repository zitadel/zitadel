package oidc

import (
	"time"

	"github.com/zitadel/oidc/v2/pkg/oidc"
	"github.com/zitadel/oidc/v2/pkg/op"
)

type clientCredentialsRequest struct {
	sub      string
	audience []string
	scopes   []string
}

// GetSubject returns the subject for token to be created because of the client credentials request
// the subject will be the id of the service user
func (c *clientCredentialsRequest) GetSubject() string {
	return c.sub
}

// GetAudience returns the audience for token to be created because of the client credentials request
func (c *clientCredentialsRequest) GetAudience() []string {
	return c.audience
}

func (c *clientCredentialsRequest) GetScopes() []string {
	return c.scopes
}

type clientCredentialsClient struct {
	id        string
	tokenType op.AccessTokenType
}

// AccessTokenType returns the AccessTokenType for the token to be created because of the client credentials request
// machine users currently only have opaque tokens ([op.AccessTokenTypeBearer])
func (c *clientCredentialsClient) AccessTokenType() op.AccessTokenType {
	return c.tokenType
}

// GetID returns the client_id (username of the machine user) for the token to be created because of the client credentials request
func (c *clientCredentialsClient) GetID() string {
	return c.id
}

// RedirectURIs returns nil as there are no redirect uris
func (c *clientCredentialsClient) RedirectURIs() []string {
	return nil
}

// PostLogoutRedirectURIs returns nil as there are no logout redirect uris
func (c *clientCredentialsClient) PostLogoutRedirectURIs() []string {
	return nil
}

// ApplicationType returns [op.ApplicationTypeWeb] as the machine users is a confidential client
func (c *clientCredentialsClient) ApplicationType() op.ApplicationType {
	return op.ApplicationTypeWeb
}

// AuthMethod returns the allowed auth method type for machine user.
// It returns Basic Auth
func (c *clientCredentialsClient) AuthMethod() oidc.AuthMethod {
	return oidc.AuthMethodBasic
}

// ResponseTypes returns nil as the types are only required for an authorization request
func (c *clientCredentialsClient) ResponseTypes() []oidc.ResponseType {
	return nil
}

// GrantTypes returns the grant types supported by the machine users, which is currently only client credentials ([oidc.GrantTypeClientCredentials])
func (c *clientCredentialsClient) GrantTypes() []oidc.GrantType {
	return []oidc.GrantType{
		oidc.GrantTypeClientCredentials,
	}
}

// LoginURL returns an empty string as there is no login UI involved
func (c *clientCredentialsClient) LoginURL(_ string) string {
	return ""
}

// IDTokenLifetime returns 0 as there is no id_token issued
func (c *clientCredentialsClient) IDTokenLifetime() time.Duration {
	return 0
}

// DevMode returns false as there is no dev mode
func (c *clientCredentialsClient) DevMode() bool {
	return false
}

// RestrictAdditionalIdTokenScopes returns nil as no id_token is issued
func (c *clientCredentialsClient) RestrictAdditionalIdTokenScopes() func(scopes []string) []string {
	return nil
}

// RestrictAdditionalAccessTokenScopes returns the scope allowed for the token to be created because of the client credentials request
// currently it allows all scopes to be used in the access token
func (c *clientCredentialsClient) RestrictAdditionalAccessTokenScopes() func(scopes []string) []string {
	return func(scopes []string) []string {
		return scopes
	}
}

// IsScopeAllowed returns null false as the check is executed during the auth request validation
func (c *clientCredentialsClient) IsScopeAllowed(scope string) bool {
	return false
}

// IDTokenUserinfoClaimsAssertion returns null false as no id_token is issued
func (c *clientCredentialsClient) IDTokenUserinfoClaimsAssertion() bool {
	return false
}

// ClockSkew enable handling clock skew of the token validation. The duration (0-5s) will be added to exp claim and subtracted from iats,
// auth_time and nbf of the token to be created because of the client credentials request.
// It returns 0 as clock skew is not implemented on machine users.
func (c *clientCredentialsClient) ClockSkew() time.Duration {
	return 0
}
