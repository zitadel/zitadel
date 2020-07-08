package oidc

import (
	"time"

	"github.com/caos/oidc/pkg/op"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/project/model"
)

type Client struct {
	*model.ApplicationView
	defaultLoginURL            string
	defaultAccessTokenLifetime time.Duration
	defaultIdTokenLifetime     time.Duration
}

func ClientFromBusiness(app *model.ApplicationView, defaultLoginURL string, defaultAccessTokenLifetime, defaultIdTokenLifetime time.Duration) (op.Client, error) {
	if !app.IsOIDC {
		return nil, errors.ThrowInvalidArgument(nil, "OIDC-d5bhD", "client is not a proper oidc application")
	}
	return &Client{ApplicationView: app, defaultLoginURL: defaultLoginURL, defaultAccessTokenLifetime: defaultAccessTokenLifetime, defaultIdTokenLifetime: defaultIdTokenLifetime}, nil
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
	return c.defaultLoginURL + id
}

func (c *Client) RedirectURIs() []string {
	return c.OIDCRedirectUris
}

func (c *Client) PostLogoutRedirectURIs() []string {
	return c.OIDCPostLogoutRedirectUris
}

func (c *Client) AccessTokenLifetime() time.Duration {
	return c.defaultAccessTokenLifetime //PLANNED: impl from real client
}

func (c *Client) IDTokenLifetime() time.Duration {
	return c.defaultIdTokenLifetime //PLANNED: impl from real client
}

func (c *Client) AccessTokenType() op.AccessTokenType {
	return op.AccessTokenTypeBearer //PLANNED: impl from real client
}

func authMethodToOIDC(authType model.OIDCAuthMethodType) op.AuthMethod {
	switch authType {
	case model.OIDCAuthMethodTypeBasic:
		return op.AuthMethodBasic
	case model.OIDCAuthMethodTypePost:
		return op.AuthMethodPost
	case model.OIDCAuthMethodTypeNone:
		return op.AuthMethodNone
	default:
		return op.AuthMethodBasic
	}
}
