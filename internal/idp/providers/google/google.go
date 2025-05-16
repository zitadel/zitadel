package google

import (
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
func New(clientID, clientSecret, redirectURI string, scopes []string, opts ...oidc.ProviderOpts) (*Provider, error) {
	rp, err := oidc.New(name, issuer, clientID, clientSecret, redirectURI, scopes, userMapper, append(opts, oidc.WithSelectAccount())...)
	if err != nil {
		return nil, err
	}
	return &Provider{
		Provider: rp,
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
