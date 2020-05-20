package oidc

import (
	"time"

	"github.com/caos/oidc/pkg/op"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/project/model"
)

type Client struct {
	*model.ApplicationView
	defaultLoginURL string
	tokenLifetime   time.Duration
}

func ClientFromBusiness(app *model.ApplicationView, defaultLoginURL string, tokenLifetime time.Duration) (op.Client, error) {
	if !app.IsOIDC {
		return nil, errors.ThrowInvalidArgument(nil, "OIDC-d5bhD", "client is not a proper oidc application")
	}
	return &Client{ApplicationView: app, defaultLoginURL: defaultLoginURL, tokenLifetime: tokenLifetime}, nil
}

func (c *Client) ApplicationType() op.ApplicationType {
	return op.ApplicationType(c.OIDCApplicationType)
}

func (c *Client) GetAuthMethod() op.AuthMethod {
	return authMethodToOIDC(c.OIDCAuthMethodType)
}

func (c *Client) GetID() string {
	return c.OIDCClientID
}

func (c *Client) LoginURL(id string) string {
	return c.defaultLoginURL + id //TODO: still needed
}

func (c *Client) RedirectURIs() []string {
	return c.OIDCRedirectUris
}

func (c *Client) PostLogoutRedirectURIs() []string {
	return c.OIDCPostLogoutRedirectUris
}

func (c *Client) AccessTokenLifetime() time.Duration {
	return c.tokenLifetime //TODO: impl from real client
}

func (c *Client) IDTokenLifetime() time.Duration {
	return c.tokenLifetime //TODO: impl from real client
}

func (c *Client) AccessTokenType() op.AccessTokenType {
	return op.AccessTokenTypeBearer //TODO: impl from real client
}

func authMethodToOIDC(authType model.OIDCAuthMethodType) op.AuthMethod {
	switch authType {
	case model.OIDCAUTHMETHODTYPE_BASIC:
		return op.AuthMethodBasic
	case model.OIDCAUTHMETHODTYPE_POST:
		return op.AuthMethodPost
	case model.OIDCAUTHMETHODTYPE_NONE:
		return op.AuthMethodNone
	default:
		return op.AuthMethodBasic
	}
}
