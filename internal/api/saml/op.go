package saml

import (
	"github.com/caos/oidc/pkg/op"
	"github.com/caos/zitadel/internal/api/http/middleware"
)

type OPHandlerConfig struct {
	OPConfig              *op.Config
	StorageConfig         StorageConfig
	UserAgentCookieConfig *middleware.UserAgentCookieConfig
	Cache                 *middleware.CacheConfig
	Endpoints             *EndpointConfig
}

type EndpointConfig struct {
	Metadata        *Endpoint
	ArtifactResolve *Endpoint
	RedirectSSO     *Endpoint
}

type Endpoint struct {
	Path string
	URL  string
}
