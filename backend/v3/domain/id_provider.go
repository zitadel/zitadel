package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
)

//go:generate enumer -type IDPType -transform lower -trimprefix IDPType
type IDPType uint8

const (
	IDPTypeOIDC IDPType = iota + 1
	IDPTypeJWT
	IDPTypeOAuth
	IDPTypeSAML
	IDPTypeLDAP
	IDPTypeGitHub
	IDPTypeGitHubEnterprise
	IDPTypeGitLab
	IDPTypeGitLabSelfHosted
	IDPTypeAzure
	IDPTypeGoogle
	IDPTypeApple
)

//go:generate enumer -type IDPState -transform lower -trimprefix IDPState -sql
type IDPState uint8

const (
	IDPStateActive IDPState = iota
	IDPStateInactive
)

//go:generate enumer -type IDPAutoLinkingOption -transform lower -trimprefix IDPAutoLinkingOption
type IDPAutoLinkingOption uint8

const (
	IDPAutoLinkingOptionUserName IDPAutoLinkingOption = iota + 1
	IDPAutoLinkingOptionEmail
)

type OIDCMappingField int8

const (
	OIDCMappingFieldUnspecified OIDCMappingField = iota
	OIDCMappingFieldPreferredLoginName
	OIDCMappingFieldEmail
	// count is for validation purposes
	//nolint: unused
	oidcMappingFieldCount
)

type IdentityProvider struct {
	InstanceID        string          `json:"instanceId,omitempty" db:"instance_id"`
	OrgID             *string         `json:"orgId,omitempty" db:"org_id"`
	ID                string          `json:"id,omitempty" db:"id"`
	State             IDPState        `json:"state,omitempty" db:"state"`
	Name              string          `json:"name,omitempty" db:"name"`
	Type              *int16          `json:"type,omitempty" db:"type"`
	AllowCreation     bool            `json:"allowCreation,omitempty" db:"allow_creation"`
	AutoRegister      bool            `json:"autoRegister,omitempty" db:"auto_register"`
	AllowAutoCreation bool            `json:"allowAutoCreation,omitempty" db:"allow_auto_creation"`
	AllowAutoUpdate   bool            `json:"allowAutoUpdate,omitempty" db:"allow_auto_update"`
	AllowLinking      bool            `json:"allowLinking,omitempty" db:"allow_linking"`
	AllowAutoLinking  *int16          `json:"allowAutoLinking,omitempty" db:"allow_auto_linking"`
	StylingType       *int16          `json:"stylingType,omitempty" db:"styling_type"`
	Payload           json.RawMessage `json:"payload,omitempty" db:"payload"`
	CreatedAt         time.Time       `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt         time.Time       `json:"updatedAt,omitzero" db:"updated_at"`
}

type OIDC struct {
	ClientID              string              `json:"clientId,omitempty"`
	ClientSecret          *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Issuer                string              `json:"issuer,omitempty"`
	AuthorizationEndpoint string              `json:"authorizationEndpoint,omitempty"`
	TokenEndpoint         string              `json:"tokenEndpoint,omitempty"`
	Scopes                []string            `json:"scopes,omitempty"`
	IDPDisplayNameMapping OIDCMappingField    `json:"IDPDisplayNameMapping,omitempty"`
	UserNameMapping       OIDCMappingField    `json:"usernameMapping,omitempty"`
	IsIDTokenMapping      bool                `json:"idTokenMapping,omitempty"`
	UsePKCE               bool                `json:"usePKCE,omitempty"`
}

type IDPOIDC struct {
	*IdentityProvider
	OIDC
}

type JWT struct {
	IDPConfigID  string `json:"idpConfigId"`
	JWTEndpoint  string `json:"jwtEndpoint,omitempty"`
	Issuer       string `json:"issuer,omitempty"`
	KeysEndpoint string `json:"keysEndpoint,omitempty"`
	HeaderName   string `json:"headerName,omitempty"`
}

type IDPJWT struct {
	*IdentityProvider
	JWT
}

type OAuth struct {
	ClientID              string              `json:"clientId,omitempty"`
	ClientSecret          *crypto.CryptoValue `json:"clientSecret,omitempty"`
	AuthorizationEndpoint string              `json:"authorizationEndpoint,omitempty"`
	TokenEndpoint         string              `json:"tokenEndpoint,omitempty"`
	UserEndpoint          string              `json:"userEndpoint,omitempty"`
	Scopes                []string            `json:"scopes,omitempty"`
	IDAttribute           string              `json:"idAttribute,omitempty"`
	UsePKCE               bool                `json:"usePKCE,omitempty"`
}

type IDPOAuth struct {
	*IdentityProvider
	OAuth
}

//go:generate enumer -type AzureTenantType -transform lower -trimprefix AzureTenantType -sql
type AzureTenantType uint8

const (
	AzureTenantTypeCommon AzureTenantType = iota
	AzureTenantTypeOrganizations
	AzureTenantTypeConsumers
)

type Azure struct {
	ClientID        string              `json:"client_id,omitempty"`
	ClientSecret    *crypto.CryptoValue `json:"client_secret,omitempty"`
	Scopes          []string            `json:"scopes,omitempty"`
	Tenant          AzureTenantType     `json:"tenant,omitempty"`
	IsEmailVerified bool                `json:"isEmailVerified,omitempty"`
}

type IDPAzureAD struct {
	*IdentityProvider
	Azure
}

type Google struct {
	ClientID     string              `json:"clientId"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret"`
	Scopes       []string            `json:"scopes,omitempty"`
}

type IDPGoogle struct {
	*IdentityProvider
	Google
}

type Github struct {
	ClientID     string              `json:"clientId"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret"`
	Scopes       []string            `json:"scopes,omitempty"`
}

type IDPGithub struct {
	*IdentityProvider
	Github
}

type GithubEnterprise struct {
	ClientID              string              `json:"clientId,omitempty"`
	ClientSecret          *crypto.CryptoValue `json:"clientSecret,omitempty"`
	AuthorizationEndpoint string              `json:"authorizationEndpoint,omitempty"`
	TokenEndpoint         string              `json:"tokenEndpoint,omitempty"`
	UserEndpoint          string              `json:"userEndpoint,omitempty"`
	Scopes                []string            `json:"scopes,omitempty"`
}

type IDPGithubEnterprise struct {
	*IdentityProvider
	GithubEnterprise
}

type Gitlab struct {
	ClientID     string              `json:"clientId,omitempty"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Scopes       []string            `json:"scopes,omitempty"`
}

type IDPGitlab struct {
	*IdentityProvider
	Gitlab
}

type GitlabSelfHosting struct {
	Issuer       string              `json:"issuer"`
	ClientID     string              `json:"clientId,omitempty"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Scopes       []string            `json:"scopes,omitempty"`
}

type IDPGitlabSelfHosting struct {
	*IdentityProvider
	GitlabSelfHosting
}

type LDAP struct {
	Servers           []string            `json:"servers"`
	StartTLS          bool                `json:"startTLS"`
	BaseDN            string              `json:"baseDN"`
	BindDN            string              `json:"bindDN"`
	BindPassword      *crypto.CryptoValue `json:"bindPassword"`
	UserBase          string              `json:"userBase"`
	UserObjectClasses []string            `json:"userObjectClasses"`
	UserFilters       []string            `json:"userFilters"`
	Timeout           time.Duration       `json:"timeout"`
	RootCA            []byte              `json:"rootCA"`

	LDAPAttributes
}

type LDAPAttributes struct {
	IDAttribute                string `json:"idAttribute,omitempty"`
	FirstNameAttribute         string `json:"firstNameAttribute,omitempty"`
	LastNameAttribute          string `json:"lastNameAttribute,omitempty"`
	DisplayNameAttribute       string `json:"displayNameAttribute,omitempty"`
	NickNameAttribute          string `json:"nickNameAttribute,omitempty"`
	PreferredUsernameAttribute string `json:"preferredUsernameAttribute,omitempty"`
	EmailAttribute             string `json:"emailAttribute,omitempty"`
	EmailVerifiedAttribute     string `json:"emailVerifiedAttribute,omitempty"`
	PhoneAttribute             string `json:"phoneAttribute,omitempty"`
	PhoneVerifiedAttribute     string `json:"phoneVerifiedAttribute,omitempty"`
	PreferredLanguageAttribute string `json:"preferredLanguageAttribute,omitempty"`
	AvatarURLAttribute         string `json:"avatarURLAttribute,omitempty"`
	ProfileAttribute           string `json:"profileAttribute,omitempty"`
}

type IDPLDAP struct {
	*IdentityProvider
	LDAP
}

type Apple struct {
	ClientID   string              `json:"clientId"`
	TeamID     string              `json:"teamId"`
	KeyID      string              `json:"keyId"`
	PrivateKey *crypto.CryptoValue `json:"privateKey"`
	Scopes     []string            `json:"scopes,omitempty"`
}

type IDPApple struct {
	*IdentityProvider
	Apple
}

type SAML struct {
	Metadata                      []byte                   `json:"metadata,omitempty"`
	Key                           *crypto.CryptoValue      `json:"key,omitempty"`
	Certificate                   []byte                   `json:"certificate,omitempty"`
	Binding                       string                   `json:"binding,omitempty"`
	WithSignedRequest             bool                     `json:"withSignedRequest,omitempty"`
	NameIDFormat                  *domain.SAMLNameIDFormat `json:"nameIDFormat,omitempty"`
	TransientMappingAttributeName string                   `json:"transientMappingAttributeName,omitempty"`
	FederatedLogoutEnabled        bool                     `json:"federatedLogoutEnabled,omitempty"`
	SignatureAlgorithm            string                   `json:"signatureAlgorithm,omitempty"`
}

type IDPSAML struct {
	*IdentityProvider
	SAML
}

// IDPIdentifierCondition is used to help specify a single identity_provider,
// it will either be used as the  identity_provider ID or identity_provider name,
// as identity_provider can be identified either using (instanceID + OrgID + ID) OR (instanceID + OrgID + name)
type IDPIdentifierCondition interface {
	database.Condition
}

type idProviderColumns interface {
	InstanceIDColumn() database.Column
	OrgIDColumn() database.Column
	IDColumn() database.Column
	StateColumn() database.Column
	NameColumn() database.Column
	TypeColumn() database.Column
	AllowCreationColumn() database.Column
	AutoRegisterColumn() database.Column
	AllowAutoCreationColumn() database.Column
	AllowAutoUpdateColumn() database.Column
	AllowLinkingColumn() database.Column
	AllowAutoLinkingColumn() database.Column
	StylingTypeColumn() database.Column
	PayloadColumn() database.Column
	CreatedAtColumn() database.Column
	UpdatedAtColumn() database.Column
}

type idProviderConditions interface {
	InstanceIDCondition(id string) database.Condition
	OrgIDCondition(id *string) database.Condition
	IDCondition(id string) IDPIdentifierCondition
	StateCondition(state IDPState) database.Condition
	NameCondition(name string) IDPIdentifierCondition
	TypeCondition(typee IDPType) database.Condition
	AutoRegisterCondition(allow bool) database.Condition
	AllowCreationCondition(allow bool) database.Condition
	AllowAutoCreationCondition(allow bool) database.Condition
	AllowAutoUpdateCondition(allow bool) database.Condition
	AllowLinkingCondition(allow bool) database.Condition
	AllowAutoLinkingCondition(linkingType IDPAutoLinkingOption) database.Condition
	StylingTypeCondition(style int16) database.Condition
	PayloadCondition(payload string) database.Condition
}

type idProviderChanges interface {
	SetName(name string) database.Change
	SetState(state IDPState) database.Change
	SetAllowCreation(allow bool) database.Change
	SetAutoRegister(allow bool) database.Change
	SetAllowAutoCreation(allow bool) database.Change
	SetAllowAutoUpdate(allow bool) database.Change
	SetAllowLinking(allow bool) database.Change
	SetAutoAllowLinking(allow bool) database.Change
	SetStylingType(stylingType int16) database.Change
	SetPayload(payload string) database.Change
}

type IDProviderRepository interface {
	idProviderColumns
	idProviderConditions
	idProviderChanges

	Get(ctx context.Context, id IDPIdentifierCondition, instanceID string, orgID *string) (*IdentityProvider, error)
	List(ctx context.Context, conditions ...database.Condition) ([]*IdentityProvider, error)

	Create(ctx context.Context, idp *IdentityProvider) error
	Update(ctx context.Context, id IDPIdentifierCondition, instanceID string, orgID *string, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, id IDPIdentifierCondition, instanceID string, orgID *string) (int64, error)

	GetOIDC(ctx context.Context, id IDPIdentifierCondition, instanceID string, orgID *string) (*IDPOIDC, error)
	GetJWT(ctx context.Context, id IDPIdentifierCondition, instanceID string, orgID *string) (*IDPJWT, error)

	GetOAuth(ctx context.Context, id IDPIdentifierCondition, instanceID string, orgID *string) (*IDPOAuth, error)

	GetAzureAD(ctx context.Context, id IDPIdentifierCondition, instanceID string, orgID *string) (*IDPAzureAD, error)
	GetGoogle(ctx context.Context, id IDPIdentifierCondition, instanceID string, orgID *string) (*IDPGoogle, error)
	GetGithub(ctx context.Context, id IDPIdentifierCondition, instanceID string, orgID *string) (*IDPGithub, error)
	GetGithubEnterprise(ctx context.Context, id IDPIdentifierCondition, instanceID string, orgID *string) (*IDPGithubEnterprise, error)
	GetGitlab(ctx context.Context, id IDPIdentifierCondition, instanceID string, orgID *string) (*IDPGitlab, error)
	GetGitlabSelfHosting(ctx context.Context, id IDPIdentifierCondition, instanceID string, orgID *string) (*IDPGitlabSelfHosting, error)
	GetLDAP(ctx context.Context, id IDPIdentifierCondition, instanceID string, orgID *string) (*IDPLDAP, error)
	GetApple(ctx context.Context, id IDPIdentifierCondition, instanceID string, orgID *string) (*IDPApple, error)
	GetSAML(ctx context.Context, id IDPIdentifierCondition, instanceID string, orgID *string) (*IDPSAML, error)
}
