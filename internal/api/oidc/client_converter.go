package oidc

import (
	"slices"
	"strings"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

type Client struct {
	client            *query.OIDCClient
	defaultLoginURL   string
	defaultLoginURLV2 string
	allowedScopes     []string
}

func ClientFromBusiness(client *query.OIDCClient, defaultLoginURL, defaultLoginURLV2 string) op.Client {
	allowedScopes := make([]string, len(client.ProjectRoleKeys))
	for i, roleKey := range client.ProjectRoleKeys {
		allowedScopes[i] = ScopeProjectRolePrefix + roleKey
	}

	return &Client{
		client:            client,
		defaultLoginURL:   defaultLoginURL,
		defaultLoginURLV2: defaultLoginURLV2,
		allowedScopes:     allowedScopes,
	}
}

func (c *Client) ApplicationType() op.ApplicationType {
	return op.ApplicationType(c.client.ApplicationType)
}

func (c *Client) AuthMethod() oidc.AuthMethod {
	return authMethodToOIDC(c.client.AuthMethodType)
}

func (c *Client) GetID() string {
	return c.client.ClientID
}

func (c *Client) LoginURL(id string) string {
	if strings.HasPrefix(id, command.IDPrefixV2) {
		return c.defaultLoginURLV2 + id
	}
	return c.defaultLoginURL + id
}

func (c *Client) RedirectURIs() []string {
	return c.client.RedirectURIs
}

func (c *Client) PostLogoutRedirectURIs() []string {
	return c.client.PostLogoutRedirectURIs
}

func (c *Client) ResponseTypes() []oidc.ResponseType {
	return responseTypesToOIDC(c.client.ResponseTypes)
}

func (c *Client) GrantTypes() []oidc.GrantType {
	return grantTypesToOIDC(c.client.GrantTypes)
}

func (c *Client) DevMode() bool {
	return c.client.IsDevMode
}

func (c *Client) RestrictAdditionalIdTokenScopes() func(scopes []string) []string {
	return func(scopes []string) []string {
		if c.client.IDTokenRoleAssertion {
			return scopes
		}
		return removeScopeWithPrefix(scopes, ScopeProjectRolePrefix)
	}
}

func (c *Client) RestrictAdditionalAccessTokenScopes() func(scopes []string) []string {
	return func(scopes []string) []string {
		if c.client.AccessTokenRoleAssertion {
			return scopes
		}
		return removeScopeWithPrefix(scopes, ScopeProjectRolePrefix)
	}
}

func (c *Client) AccessTokenLifetime() time.Duration {
	return c.client.Settings.AccessTokenLifetime
}

func (c *Client) IDTokenLifetime() time.Duration {
	return c.client.Settings.IdTokenLifetime
}

func (c *Client) AccessTokenType() op.AccessTokenType {
	return accessTokenTypeToOIDC(c.client.AccessTokenType)
}

func (c *Client) IsScopeAllowed(scope string) bool {
	if strings.HasPrefix(scope, domain.OrgDomainPrimaryScope) {
		return true
	}
	if strings.HasPrefix(scope, domain.OrgIDScope) {
		return true
	}
	if strings.HasPrefix(scope, domain.ProjectIDScope) {
		return true
	}
	if strings.HasPrefix(scope, domain.SelectIDPScope) {
		return true
	}
	if scope == ScopeUserMetaData {
		return true
	}
	if scope == ScopeResourceOwner {
		return true
	}
	if scope == ScopeProjectsRoles {
		return true
	}
	return slices.Contains(c.allowedScopes, scope)
}

func (c *Client) ClockSkew() time.Duration {
	return c.client.ClockSkew
}

func (c *Client) IDTokenUserinfoClaimsAssertion() bool {
	return c.client.IDTokenUserinfoAssertion
}

func (c *Client) RedirectURIGlobs() []string {
	if c.DevMode() {
		return c.RedirectURIs()
	}
	return nil
}

func (c *Client) PostLogoutRedirectURIGlobs() []string {
	if c.DevMode() {
		return c.PostLogoutRedirectURIs()
	}
	return nil
}

func accessTokenTypeToOIDC(tokenType domain.OIDCTokenType) op.AccessTokenType {
	switch tokenType {
	case domain.OIDCTokenTypeBearer:
		return op.AccessTokenTypeBearer
	case domain.OIDCTokenTypeJWT:
		return op.AccessTokenTypeJWT
	default:
		return op.AccessTokenTypeBearer
	}
}

func authMethodToOIDC(authType domain.OIDCAuthMethodType) oidc.AuthMethod {
	switch authType {
	case domain.OIDCAuthMethodTypeBasic:
		return oidc.AuthMethodBasic
	case domain.OIDCAuthMethodTypePost:
		return oidc.AuthMethodPost
	case domain.OIDCAuthMethodTypeNone:
		return oidc.AuthMethodNone
	case domain.OIDCAuthMethodTypePrivateKeyJWT:
		return oidc.AuthMethodPrivateKeyJWT
	default:
		return oidc.AuthMethodBasic
	}
}

func responseTypesToOIDC(responseTypes []domain.OIDCResponseType) []oidc.ResponseType {
	oidcTypes := make([]oidc.ResponseType, len(responseTypes))
	for i, t := range responseTypes {
		oidcTypes[i] = responseTypeToOIDC(t)
	}
	return oidcTypes
}

func responseTypeToOIDC(responseType domain.OIDCResponseType) oidc.ResponseType {
	switch responseType {
	case domain.OIDCResponseTypeCode:
		return oidc.ResponseTypeCode
	case domain.OIDCResponseTypeIDTokenToken:
		return oidc.ResponseTypeIDToken
	case domain.OIDCResponseTypeIDToken:
		return oidc.ResponseTypeIDTokenOnly
	default:
		return oidc.ResponseTypeCode
	}
}

func grantTypesToOIDC(grantTypes []domain.OIDCGrantType) []oidc.GrantType {
	oidcTypes := make([]oidc.GrantType, len(grantTypes))
	for i, t := range grantTypes {
		oidcTypes[i] = grantTypeToOIDC(t)
	}
	return oidcTypes
}

func grantTypeToOIDC(grantType domain.OIDCGrantType) oidc.GrantType {
	switch grantType {
	case domain.OIDCGrantTypeAuthorizationCode:
		return oidc.GrantTypeCode
	case domain.OIDCGrantTypeImplicit:
		return oidc.GrantTypeImplicit
	case domain.OIDCGrantTypeRefreshToken:
		return oidc.GrantTypeRefreshToken
	case domain.OIDCGrantTypeDeviceCode:
		return oidc.GrantTypeDeviceCode
	case domain.OIDCGrantTypeTokenExchange:
		return oidc.GrantTypeTokenExchange
	default:
		return oidc.GrantTypeCode
	}
}

func removeScopeWithPrefix(scopes []string, scopePrefix ...string) []string {
	newScopeList := make([]string, 0)
	for _, scope := range scopes {
		hasPrefix := false
		for _, prefix := range scopePrefix {
			if strings.HasPrefix(scope, prefix) {
				hasPrefix = true
				continue
			}
		}
		if !hasPrefix {
			newScopeList = append(newScopeList, scope)
		}
	}
	return newScopeList
}

func clientIDFromCredentials(cc *op.ClientCredentials) (clientID string, assertion bool, err error) {
	if cc.ClientAssertion != "" {
		claims := new(oidc.JWTTokenRequest)
		if _, err := oidc.ParseToken(cc.ClientAssertion, claims); err != nil {
			return "", false, oidc.ErrInvalidClient().WithParent(err)
		}
		return claims.Issuer, true, nil
	}
	return cc.ClientID, false, nil
}
