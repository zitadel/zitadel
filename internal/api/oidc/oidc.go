package oidc

import (
	"github.com/caos/oidc/pkg/op"
)

type Config struct {
	Config          *op.Config
	DefaultLoginURL string
	TokenLifetime   string
	//UserAgentCookieConfig *auth.UserAgentCookieConfig
	Endpoints *EndpointConfig
}

type EndpointConfig struct {
	Auth       *Endpoint
	Token      *Endpoint
	Userinfo   *Endpoint
	EndSession *Endpoint
	Keys       *Endpoint
}

type Endpoint struct {
	Path string
	URL  string
}
