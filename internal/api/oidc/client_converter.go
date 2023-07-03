package oidc

import (
	"strings"
	"time"

	"github.com/zitadel/oidc/v2/pkg/oidc"
	"github.com/zitadel/oidc/v2/pkg/op"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
)

type Client struct {
	app                        *query.App
	defaultLoginURL            string
	defaultLoginURLV2          string
	defaultAccessTokenLifetime time.Duration
	defaultIdTokenLifetime     time.Duration
	allowedScopes              []string
}

func ClientFromBusiness(app *query.App, defaultLoginURL, defaultLoginURLV2 string, defaultAccessTokenLifetime, defaultIdTokenLifetime time.Duration, allowedScopes []string) (op.Client, error) {
	if app.OIDCConfig == nil {
		return nil, errors.ThrowInvalidArgument(nil, "OIDC-d5bhD", "client is not a proper oidc application")
	}
	return &Client{
			app:                        app,
			defaultLoginURL:            defaultLoginURL,
			defaultLoginURLV2:          defaultLoginURLV2,
			defaultAccessTokenLifetime: defaultAccessTokenLifetime,
			defaultIdTokenLifetime:     defaultIdTokenLifetime,
			allowedScopes:              allowedScopes},
		nil
}

func (c *Client) ApplicationType() op.ApplicationType {
	return op.ApplicationType(c.app.OIDCConfig.AppType)
}

func (c *Client) AuthMethod() oidc.AuthMethod {
	return authMethodToOIDC(c.app.OIDCConfig.AuthMethodType)
}

func (c *Client) GetID() string {
	return c.app.OIDCConfig.ClientID
}

func (c *Client) LoginURL(id string) string {
	if strings.HasPrefix(id, command.IDPrefixV2) {
		return c.defaultLoginURLV2 + id
	}
	return c.defaultLoginURL + id
}

func (c *Client) RedirectURIs() []string {
	return c.app.OIDCConfig.RedirectURIs
}

func (c *Client) PostLogoutRedirectURIs() []string {
	return c.app.OIDCConfig.PostLogoutRedirectURIs
}

func (c *Client) ResponseTypes() []oidc.ResponseType {
	return responseTypesToOIDC(c.app.OIDCConfig.ResponseTypes)
}

func (c *Client) GrantTypes() []oidc.GrantType {
	return grantTypesToOIDC(c.app.OIDCConfig.GrantTypes)
}

func (c *Client) DevMode() bool {
	return c.app.OIDCConfig.IsDevMode
}

func (c *Client) RestrictAdditionalIdTokenScopes() func(scopes []string) []string {
	return func(scopes []string) []string {
		if c.app.OIDCConfig.AssertIDTokenRole {
			return scopes
		}
		return removeScopeWithPrefix(scopes, ScopeProjectRolePrefix)
	}
}

func (c *Client) RestrictAdditionalAccessTokenScopes() func(scopes []string) []string {
	return func(scopes []string) []string {
		if c.app.OIDCConfig.AssertAccessTokenRole {
			return scopes
		}
		return removeScopeWithPrefix(scopes, ScopeProjectRolePrefix)
	}
}

func (c *Client) AccessTokenLifetime() time.Duration {
	return c.defaultAccessTokenLifetime //PLANNED: impl from real client
}

func (c *Client) IDTokenLifetime() time.Duration {
	return c.defaultIdTokenLifetime //PLANNED: impl from real client
}

func (c *Client) AccessTokenType() op.AccessTokenType {
	return accessTokenTypeToOIDC(c.app.OIDCConfig.AccessTokenType)
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
	for _, allowedScope := range c.allowedScopes {
		if scope == allowedScope {
			return true
		}
	}
	return false
}

func (c *Client) ClockSkew() time.Duration {
	return c.app.OIDCConfig.ClockSkew
}

func (c *Client) IDTokenUserinfoClaimsAssertion() bool {
	return c.app.OIDCConfig.AssertIDTokenUserinfo
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
