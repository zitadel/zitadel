package google

import (
	"net/http"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	openid "github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

const (
	issuer = "https://accounts.google.com"
	name   = "Google"
)

var _ idp.Provider = (*Provider)(nil)

// Provider is the [idp.Provider] implementation for Google
type Provider struct {
	*oidc.Provider
}

// New creates a Google provider using the [oidc.Provider] (OIDC generic provider)
func New(clientID, clientSecret, redirectURI string, scopes []string, httpClient *http.Client, opts ...oidc.ProviderOpts) (*Provider, error) {
	// PKCE is used by default
	opts = append(opts, oidc.WithSelectAccount(), oidc.WithRelyingPartyOption(rp.WithPKCE(nil)))
	provider, err := oidc.New(name, issuer, clientID, clientSecret, redirectURI, scopes, userMapper, httpClient, opts...)
	if err != nil {
		return nil, err
	}
	return &Provider{
		Provider: provider,
	}, nil
}

var userMapper = func(info *openid.UserInfo) idp.User {
	return &User{oidc.DefaultMapper(info)}
}

func InitUser() idp.User {
	return &User{oidc.InitUser()}
}

// User is a representation of the authenticated Google and implements the [idp.User] interface
// by wrapping an [idp.User] (implemented by [oidc.User]). It overwrites the [GetPreferredUsername] to use the `email` claim.
type User struct {
	idp.User
}

// GetPreferredUsername implements the [idp.User] interface.
// It returns the email, because Google does not return a username.
func (u *User) GetPreferredUsername() string {
	return string(u.GetEmail())
}
