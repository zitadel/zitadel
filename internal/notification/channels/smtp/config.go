package smtp

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/denylist"
)

type Config struct {
	SMTP           SMTP
	Tls            bool
	From           string
	FromName       string
	ReplyToAddress string
}

type SMTP struct {
	Host        string
	PlainAuth   *PlainAuthConfig
	XOAuth2Auth *XOAuth2AuthConfig
	// Restrictions limits where the SMTP host and the XOAuth2 token endpoint may point to.
	// It is set for configs owned by an organization, and nil for instance level configs.
	Restrictions *Restrictions
}

// Restrictions keeps tenants from using their SMTP config to reach internal services.
type Restrictions struct {
	// DenyList is checked against the SMTP host before it is dialed.
	// The resolved IPs are checked again right before the connection is made.
	DenyList []denylist.AddressChecker
	// HTTPClient sends the XOAuth2 token request.
	// It must enforce DenyList at the transport, like the client from the HTTPClient config.
	HTTPClient *http.Client
}

type ConfigHTTP struct {
	Endpoint string
}
