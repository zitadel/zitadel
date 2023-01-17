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

// TenantType are the well known tenant types to scope the users that can authenticate. TenantType is not an
// exclusive list of Azure Tenants which can be used. A consumer can also use their own Tenant ID to scope
// authentication to their specific Tenant either through the Tenant ID or the friendly domain name.
//
// see also https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-v2-protocols#endpoints
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

// Provider is the idp.Provider implementation for AzureAD (V2 Endpoints)
type Provider struct {
	*oauth.Provider
	tenant        TenantType
	emailVerified bool
	options       []oauth.ProviderOpts
}

type ProviderOptions func(*Provider)

// WithTenant allows to set a TenantType (can also be a tenantID)
// default is CommonTenant
func WithTenant(tenantType TenantType) ProviderOptions {
	return func(p *Provider) {
		p.tenant = tenantType
	}
}

// WithEmailVerified allows to set every email received as verified
func WithEmailVerified() ProviderOptions {
	return func(p *Provider) {
		p.emailVerified = true
	}
}

// WithOAuthOptions allows to specify oauth.ProviderOpts like oauth.WithLinkingAllowed()
func WithOAuthOptions(opts ...oauth.ProviderOpts) ProviderOptions {
	return func(p *Provider) {
		p.options = append(p.options, opts...)
	}
}

// New creates an AzureAD provider using the oauth.Provider (OAuth 2.0 generic provider)
// By default it uses the CommonTenant and unverified emails
func New(name, clientID, clientSecret, redirectURI string, opts ...ProviderOptions) (*Provider, error) {
	provider := &Provider{
		tenant:  CommonTenant,
		options: make([]oauth.ProviderOpts, 0),
	}
	for _, opt := range opts {
		opt(provider)
	}
	config := newConfig(provider.tenant, clientID, clientSecret, redirectURI, []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail})
	rp, err := oauth.New(
		config,
		name,
		userinfoURL,
		func() oauth.UserInfoMapper {
			return &User{isEmailVerified: provider.emailVerified}
		},
		provider.options...,
	)
	if err != nil {
		return nil, err
	}
	provider.Provider = rp
	return provider, nil
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

// User represents the structure return on the userinfo endpoint
//
// AzureAD does not return an `email_verified` claim.
// The verification can be automatically activated on the provider (WithEmailVerified())
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
	// AzureAD does not provide the user's nickname
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
	// AzureAD does not provide the user's phone
	return ""
}

func (u *User) IsPhoneVerified() bool {
	// AzureAD does not provide the user's phone
	return false
}

func (u *User) GetPreferredLanguage() language.Tag {
	// AzureAD does not provide the user's language
	return language.Und
}

func (u *User) GetProfile() string {
	// AzureAD does not provide the user's profile page
	return ""
}

func (u *User) GetAvatarURL() string {
	return u.Picture
}

func (u *User) RawData() any {
	return u
}
