package azuread

import (
	"fmt"

	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
)

const (
	authURLTemplate  string = "https://login.microsoftonline.com/%s/oauth2/v2.0/authorize"
	tokenURLTemplate string = "https://login.microsoftonline.com/%s/oauth2/v2.0/token"
	userinfoURL      string = "https://graph.microsoft.com/v1.0/me"
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

// Provider is the [idp.Provider] implementation for AzureAD (V2 Endpoints)
type Provider struct {
	*oauth.Provider
	tenant        TenantType
	emailVerified bool
	options       []oauth.ProviderOpts
}

type ProviderOptions func(*Provider)

// WithTenant allows to set a [TenantType] (can also be a Tenant ID)
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

// WithOAuthOptions allows to specify [oauth.ProviderOpts] like [oauth.WithLinkingAllowed]
func WithOAuthOptions(opts ...oauth.ProviderOpts) ProviderOptions {
	return func(p *Provider) {
		p.options = append(p.options, opts...)
	}
}

// New creates an AzureAD provider using the [oauth.Provider] (OAuth 2.0 generic provider).
// By default, it uses the [CommonTenant] and unverified emails.
func New(name, clientID, clientSecret, redirectURI string, scopes []string, opts ...ProviderOptions) (*Provider, error) {
	provider := &Provider{
		tenant:  CommonTenant,
		options: make([]oauth.ProviderOpts, 0),
	}
	for _, opt := range opts {
		opt(provider)
	}
	config := newConfig(provider.tenant, clientID, clientSecret, redirectURI, scopes)
	rp, err := oauth.New(
		config,
		name,
		userinfoURL,
		func() idp.User {
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

// User represents the structure return on the userinfo endpoint and implements the [idp.User] interface
//
// AzureAD does not return an `email_verified` claim.
// The verification can be automatically activated on the provider ([WithEmailVerified])
type User struct {
	ID                string               `json:"id"`
	BusinessPhones    []domain.PhoneNumber `json:"businessPhones"`
	DisplayName       string               `json:"displayName"`
	FirstName         string               `json:"givenName"`
	JobTitle          string               `json:"jobTitle"`
	Email             domain.EmailAddress  `json:"mail"`
	MobilePhone       domain.PhoneNumber   `json:"mobilePhone"`
	OfficeLocation    string               `json:"officeLocation"`
	PreferredLanguage string               `json:"preferredLanguage"`
	LastName          string               `json:"surname"`
	UserPrincipalName string               `json:"userPrincipalName"`
	isEmailVerified   bool
}

// GetID is an implementation of the [idp.User] interface.
func (u *User) GetID() string {
	return u.ID
}

// GetFirstName is an implementation of the [idp.User] interface.
func (u *User) GetFirstName() string {
	return u.FirstName
}

// GetLastName is an implementation of the [idp.User] interface.
func (u *User) GetLastName() string {
	return u.LastName
}

// GetDisplayName is an implementation of the [idp.User] interface.
func (u *User) GetDisplayName() string {
	return u.DisplayName
}

// GetNickname is an implementation of the [idp.User] interface.
// It returns an empty string because AzureAD does not provide the user's nickname.
func (u *User) GetNickname() string {
	return ""
}

// GetPreferredUsername is an implementation of the [idp.User] interface.
func (u *User) GetPreferredUsername() string {
	return u.UserPrincipalName
}

// GetEmail is an implementation of the [idp.User] interface.
func (u *User) GetEmail() domain.EmailAddress {
	if u.Email == "" {
		// if the user used a social login on Azure as well, the email will be empty
		// but is used as username
		return domain.EmailAddress(u.UserPrincipalName)
	}
	return u.Email
}

// IsEmailVerified is an implementation of the [idp.User] interface
// returning the value specified in the creation of the [Provider].
// Default is false because AzureAD does not return an `email_verified` claim.
// The verification can be automatically activated on the provider ([WithEmailVerified]).
func (u *User) IsEmailVerified() bool {
	return u.isEmailVerified
}

// GetPhone is an implementation of the [idp.User] interface.
// It returns an empty string because AzureAD does not provide the user's phone.
func (u *User) GetPhone() domain.PhoneNumber {
	return ""
}

// IsPhoneVerified is an implementation of the [idp.User] interface.
// It returns false because AzureAD does not provide the user's phone.
func (u *User) IsPhoneVerified() bool {
	return false
}

// GetPreferredLanguage is an implementation of the [idp.User] interface.
func (u *User) GetPreferredLanguage() language.Tag {
	return language.Make(u.PreferredLanguage)
}

// GetProfile is an implementation of the [idp.User] interface.
// It returns an empty string because AzureAD does not provide the user's profile page.
func (u *User) GetProfile() string {
	return ""
}

// GetAvatarURL is an implementation of the [idp.User] interface.
func (u *User) GetAvatarURL() string {
	return ""
}
