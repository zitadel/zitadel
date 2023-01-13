package azuread

import (
	"fmt"

	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
)

const (
	authURLTemplate  string = "https://login.microsoftonline.com/%s/oauth2/v2.0/authorize"
	tokenURLTemplate string = "https://login.microsoftonline.com/%s/oauth2/v2.0/token"
	userinfoURL      string = "https://graph.microsoft.com/oidc/userinfo"
)

type TenantType string

const (
	// CommonTenant allows users with both personal Microsoft accounts and work/school accounts from Azure Active
	// Directory to sign in to the application.
	CommonTenant TenantType = "common"

	// OrganizationsTenant allows only users with work/school accounts from Azure Active Directory to sign in to the application.
	OrganizationsTenant TenantType = "organizations"

	// ConsumersTenant allows only users with personal Microsoft accounts (MSA) to sign in to the application.
	ConsumersTenant TenantType = "consumers"
)

var _ idp.Provider = (*Provider)(nil)

type Provider struct {
	provider *oauth.Provider
	name     string
}

type ProviderOpts struct {
	Tenant        TenantType
	EmailVerified bool
}

func New(clientID, clientSecret, redirectURI string, opts ProviderOpts) (*Provider, error) {
	if opts.Tenant == "" {
		opts.Tenant = CommonTenant
	}
	config := newConfig(opts.Tenant, clientID, clientSecret, redirectURI, []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail})
	rp, err := oauth.New(
		config,
		userinfoURL,
		func() oauth.UserInfoMapper {
			return &User{isEmailVerified: opts.EmailVerified}
		},
	)
	if err != nil {
		return nil, err
	}
	provider := &Provider{
		provider: rp,
	}

	return provider, nil
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) BeginAuth(state string) (idp.Session, error) {
	return p.provider.BeginAuth(state)
}

func (p *Provider) FetchUser(session idp.Session) (idp.User, error) {
	return p.provider.FetchUser(session)
}

func newConfig(tenant TenantType, clientID, secret, callbackURL string, scopes []string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		RedirectURL:  callbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf(authURLTemplate, tenant),
			TokenURL: fmt.Sprintf(tokenURLTemplate, tenant),
		},
		Scopes: []string{oidc.ScopeOpenID},
	}
	if len(scopes) > 0 {
		c.Scopes = scopes
	}

	return c
}

type User struct {
	Sub               string `json:"sub"`
	FamilyName        string `json:"family_name"`
	GivenName         string `json:"given_name"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	Picture           string `json:"picture"`
	isEmailVerified   bool
}

func (u *User) GetID() string {
	return u.Sub
}

func (u *User) GetFirstName() string {
	return u.GivenName
}

func (u *User) GetLastName() string {
	return u.FamilyName
}

func (u *User) GetDisplayName() string {
	return u.Name
}

func (u *User) GetNickName() string {
	return ""
}

func (u *User) GetPreferredUsername() string {
	return u.PreferredUsername
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) IsEmailVerified() bool {
	return u.isEmailVerified
}

func (u *User) GetPhone() string {
	return "" //TODO: ?
}

func (u *User) IsPhoneVerified() bool {
	return false //TODO: ?
}

func (u *User) GetPreferredLanguange() language.Tag {
	return language.Und //TODO: ?
}

func (u *User) GetAvatarURL() string {
	return u.Picture
}

func (u *User) GetProfile() string {
	return "" //TODO: ?
}

func (u *User) RawData() any {
	return u
}
