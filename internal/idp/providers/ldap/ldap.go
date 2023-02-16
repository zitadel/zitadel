package ldap

import (
	"context"

	"github.com/zitadel/zitadel/internal/idp"
)

const DefaultPort = "389"

var _ idp.Provider = (*Provider)(nil)

// Provider is the [idp.Provider] implementation for a generic LDAP provider
type Provider struct {
	name                string
	host                string
	port                string
	tls                 bool
	baseDN              string
	userObjectClass     string
	userUniqueAttribute string
	admin               string
	password            string
	loginUrl            string

	isLinkingAllowed  bool
	isCreationAllowed bool
	isAutoCreation    bool
	isAutoUpdate      bool

	idAttribute                string
	firstNameAttribute         string
	lastNameAttribute          string
	displayNameAttribute       string
	nickNameAttribute          string
	preferredUsernameAttribute string
	emailAttribute             string
	emailVerifiedAttribute     string
	phoneAttribute             string
	phoneVerifiedAttribute     string
	preferredLanguageAttribute string
	avatarURLAttribute         string
	profileAttribute           string
}

type ProviderOpts func(provider *Provider)

// WithLinkingAllowed allows end users to link the federated user to an existing one.
func WithLinkingAllowed() ProviderOpts {
	return func(p *Provider) {
		p.isLinkingAllowed = true
	}
}

// WithCreationAllowed allows end users to create a new user using the federated information.
func WithCreationAllowed() ProviderOpts {
	return func(p *Provider) {
		p.isCreationAllowed = true
	}
}

// WithAutoCreation enables that federated users are automatically created if not already existing.
func WithAutoCreation() ProviderOpts {
	return func(p *Provider) {
		p.isAutoCreation = true
	}
}

// WithAutoUpdate enables that information retrieved from the provider is automatically used to update
// the existing user on each authentication.
func WithAutoUpdate() ProviderOpts {
	return func(p *Provider) {
		p.isAutoUpdate = true
	}
}

// WithCustomPort configures a custom port used for the communication instead of :389 as per default
func WithCustomPort(port string) ProviderOpts {
	return func(p *Provider) {
		p.port = port
	}
}

// Insecure configures to communication insecure with the LDAP server without TLS
func Insecure() ProviderOpts {
	return func(p *Provider) {
		p.tls = false
	}
}

// WithCustomIDAttribute configures to map the LDAP attribute to the user, default is the uniqueUserAttribute
func WithCustomIDAttribute(name string) ProviderOpts {
	return func(p *Provider) {
		p.idAttribute = name
	}
}

// WithFirstNameAttribute configures to map the LDAP attribute to the user
func WithFirstNameAttribute(name string) ProviderOpts {
	return func(p *Provider) {
		p.firstNameAttribute = name
	}
}

// WithLastNameAttribute configures to map the LDAP attribute to the user
func WithLastNameAttribute(name string) ProviderOpts {
	return func(p *Provider) {
		p.lastNameAttribute = name
	}
}

// WithDisplayNameAttribute configures to map the LDAP attribute to the user
func WithDisplayNameAttribute(name string) ProviderOpts {
	return func(p *Provider) {
		p.displayNameAttribute = name
	}
}

// WithNickNameAttribute configures to map the LDAP attribute to the user
func WithNickNameAttribute(name string) ProviderOpts {
	return func(p *Provider) {
		p.nickNameAttribute = name
	}
}

// WithPreferredUsernameAttribute configures to map the LDAP attribute to the user
func WithPreferredUsernameAttribute(name string) ProviderOpts {
	return func(p *Provider) {
		p.preferredUsernameAttribute = name
	}
}

// WithEmailAttribute configures to map the LDAP attribute to the user
func WithEmailAttribute(name string) ProviderOpts {
	return func(p *Provider) {
		p.emailAttribute = name
	}
}

// WithEmailVerifiedAttribute configures to map the LDAP attribute to the user
func WithEmailVerifiedAttribute(name string) ProviderOpts {
	return func(p *Provider) {
		p.emailVerifiedAttribute = name
	}
}

// WithPhoneAttribute configures to map the LDAP attribute to the user
func WithPhoneAttribute(name string) ProviderOpts {
	return func(p *Provider) {
		p.phoneAttribute = name
	}
}

// WithPhoneVerifiedAttribute configures to map the LDAP attribute to the user
func WithPhoneVerifiedAttribute(name string) ProviderOpts {
	return func(p *Provider) {
		p.phoneVerifiedAttribute = name
	}
}

// WithPreferredLanguageAttribute configures to map the LDAP attribute to the user
func WithPreferredLanguageAttribute(name string) ProviderOpts {
	return func(p *Provider) {
		p.preferredLanguageAttribute = name
	}
}

// WithAvatarURLAttribute configures to map the LDAP attribute to the user
func WithAvatarURLAttribute(name string) ProviderOpts {
	return func(p *Provider) {
		p.avatarURLAttribute = name
	}
}

// WithProfileAttribute configures to map the LDAP attribute to the user
func WithProfileAttribute(name string) ProviderOpts {
	return func(p *Provider) {
		p.profileAttribute = name
	}
}

func New(
	name string,
	host string,
	baseDN string,
	userObjectClass string,
	userUniqueAttribute string,
	admin string,
	password string,
	loginUrl string,
	options ...ProviderOpts,
) *Provider {
	provider := &Provider{
		name:                name,
		host:                host,
		port:                DefaultPort,
		tls:                 true,
		baseDN:              baseDN,
		userObjectClass:     userObjectClass,
		userUniqueAttribute: userUniqueAttribute,
		admin:               admin,
		password:            password,
		loginUrl:            loginUrl,
		idAttribute:         userUniqueAttribute,
	}
	for _, option := range options {
		option(provider)
	}
	return provider
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) BeginAuth(ctx context.Context, state string, params ...any) (idp.Session, error) {
	return &Session{
		Provider: p,
		loginUrl: p.loginUrl + "?state=" + state,
	}, nil
}

func (p *Provider) IsLinkingAllowed() bool {
	return p.isLinkingAllowed
}

func (p *Provider) IsCreationAllowed() bool {
	return p.isCreationAllowed
}

func (p *Provider) IsAutoCreation() bool {
	return p.isAutoCreation
}

func (p *Provider) IsAutoUpdate() bool {
	return p.isAutoUpdate
}
