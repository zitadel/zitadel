package google

import (
	"fmt"

	openid "github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oidc"
	"github.com/zitadel/zitadel/internal/zerrors"
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

// WithHostedDomain restricts sign-in to users from the specified Google Workspace domain
// by adding the `hd` parameter to the auth URL. Leave empty to allow any Google account.
func WithHostedDomain(hostedDomain string) oidc.ProviderOpts {
	return oidc.WithAuthURLParam("hd", hostedDomain)
}

// WithEnforceHostedDomain adds server-side validation of the `hd` claim returned in the Google
// ID token. When set, FetchUser will reject any Google account whose `hd` claim does not exactly
// match the configured domain, preventing users from bypassing the hosted-domain restriction by
// switching accounts after the auth redirect.
//
// Requires WithHostedDomain to be set with the same domain.
func WithEnforceHostedDomain(domain string) oidc.ProviderOpts {
	return oidc.WithUserValidator(func(info *openid.UserInfo) error {
		hd, _ := info.Claims["hd"].(string)
		if hd != domain {
			return zerrors.ThrowPermissionDenied(nil, "GOOGLE-HD-01", fmt.Sprintf("hd claim %q does not match required domain %q", hd, domain))
		}
		return nil
	})
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
