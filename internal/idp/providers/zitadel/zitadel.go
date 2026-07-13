package zitadel

import (
	"net/http"

	"github.com/zitadel/oidc/v3/pkg/client/rp"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

const (
	name = "ZITADEL"
)

var _ idp.Provider = (*Provider)(nil)

// Provider is the [idp.Provider] implementation for ZITADEL
type Provider struct {
	*oidc.Provider
}

func New(issuer, clientID, clientSecret, redirectURI string, scopes []string, httpClient *http.Client, options ...oidc.ProviderOpts) (*Provider, error) {
	options = append(options, oidc.WithRelyingPartyOption(rp.WithPKCE(nil)))
	provider, err := oidc.New(name, issuer, clientID, clientSecret, redirectURI, scopes, oidc.DefaultMapper, httpClient, options...)
	if err != nil {
		return nil, err
	}
	return &Provider{
		Provider: provider,
	}, nil
}
