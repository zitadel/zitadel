package domain

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

//go:generate enumer -type SettingType -transform snake -trimprefix SettingType -sql
type SettingType uint8

const (
	SettingTypeUnspecified SettingType = iota
	SettingTypeLogin
	SettingTypeLabel
	SettingTypePasswordComplexity
	SettingTypePasswordExpiry
	SettingTypeDomain
	SettingTypeLockout
	SettingTypeSecurity
	SettingTypeOrganization
)

//go:generate enumer -type OwnerType -transform snake -trimprefix OwnerType -sql
type OwnerType uint8

const (
	OwnerTypeUnspecified OwnerType = iota
	OwnerTypeSystem
	OwnerTypeInstance
	OwnerTypeOrganization
)

type Setting struct {
	ID         string          `json:"id,omitempty" db:"id"`
	InstanceID string          `json:"instanceId,omitempty" db:"instance_id"`
	OrgID      *string         `json:"orgId,omitempty" db:"org_id"`
	Type       SettingType     `json:"type,omitempty" db:"type"`
	OwnerType  OwnerType       `json:"ownerType,omitempty" db:"owner_type"`
	LabelState *LabelState     `json:"labelState,omitempty" db:"label_state"`
	Settings   json.RawMessage `json:"settings,omitempty" db:"settings"`
	CreatedAt  time.Time       `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt  *time.Time      `json:"updatedAt,omitzero" db:"updated_at"`
}

type PasswordlessType int32

const (
	PasswordlessTypeNotAllowed PasswordlessType = iota
	PasswordlessTypeAllowed
)

type MultiFactorType int32

const (
	MultiFactorTypeUnspecified MultiFactorType = iota
	MultiFactorTypeU2FWithPIN
)

type SecondFactorType int32

const (
	SecondFactorTypeUnspecified SecondFactorType = iota
	SecondFactorTypeTOTP
	SecondFactorTypeU2F
	SecondFactorTypeOTPEmail
	SecondFactorTypeOTPSMS
)

type LoginSettings struct {
	AllowUserNamePassword      *bool             `json:"allowUsernamePassword,omitempty"`
	AllowRegister              *bool             `json:"allowRegister,omitempty"`
	AllowExternalIDP           *bool             `json:"allowExternalIdp,omitempty"`
	ForceMFA                   *bool             `json:"forceMfa,omitempty"`
	ForceMFALocalOnly          *bool             `json:"forceMfaLocalOnly,omitempty"`
	HidePasswordReset          *bool             `json:"hidePasswordReset,omitempty"`
	IgnoreUnknownUsernames     *bool             `json:"ignoreUnknownUsernames,omitempty"`
	AllowDomainDiscovery       *bool             `json:"allowDomainDiscovery,omitempty"`
	DisableLoginWithEmail      *bool             `json:"disableLoginWithEmail,omitempty"`
	DisableLoginWithPhone      *bool             `json:"disableLoginWithPhone,omitempty"`
	PasswordlessType           *PasswordlessType `json:"passwordlessType,omitempty"`
	DefaultRedirectURI         string            `json:"defaultRedirectUri,omitempty"`
	PasswordCheckLifetime      time.Duration     `json:"passwordCheckLifetime,omitempty"`
	ExternalLoginCheckLifetime time.Duration     `json:"externalLoginCheckLifetime,omitempty"`
	MFAInitSkipLifetime        time.Duration     `json:"mfaInitSkipLifetime,omitempty"`
	SecondFactorCheckLifetime  time.Duration     `json:"secondFactorCheckLifetime,omitempty"`
	MultiFactorCheckLifetime   time.Duration     `json:"multiFactorCheckLifetime,omitempty"`

	MFAType           []MultiFactorType  `json:"mfaType"`
	SecondFactorTypes []SecondFactorType `json:"secondFactors"`
}

type LoginSetting struct {
	*Setting
	Settings LoginSettings
}

type LabelPolicyThemeMode int32

const (
	LabelPolicyThemeAuto LabelPolicyThemeMode = iota
	LabelPolicyThemeLight
	LabelPolicyThemeDark
)

//go:generate enumer -type LabelState -transform snake -trimprefix LabelState -sql
type LabelState int32

const (
	LabelStatePreview LabelState = iota + 1
	LabelStateActivated
)

type LabelSettings struct {
	PrimaryColor        string               `json:"primaryColor,omitempty"`
	BackgroundColor     string               `json:"backgroundColor,omitempty"`
	WarnColor           string               `json:"warnColor,omitempty"`
	FontColor           string               `json:"fontColor,omitempty"`
	PrimaryColorDark    string               `json:"primaryColorDark,omitempty"`
	BackgroundColorDark string               `json:"backgroundColorDark,omitempty"`
	WarnColorDark       string               `json:"warnColorDark,omitempty"`
	FontColorDark       string               `json:"fontColorDark,omitempty"`
	HideLoginNameSuffix bool                 `json:"hideLoginNameSuffix,omitempty"`
	ErrorMsgPopup       bool                 `json:"errorMsgPopup,omitempty"`
	DisableWatermark    bool                 `json:"disableMsgPopup,omitempty"`
	ThemeMode           LabelPolicyThemeMode `json:"themeMode,omitempty"`

	LabelPolicyLightLogoURL *string `json:"labelPolicyLightLogoURL,omitempty"`
	LabelPolicyDarkLogoURL  *string `json:"labelPolicyDarkLogoURL,omitempty"`

	LabelPolicyLightIconURL *url.URL `json:"labelPolicyLightIconURL,omitempty"`
	LabelPolicyDarkIconURL  *url.URL `json:"labelPolicyDarkIconURL,omitempty"`

	LabelPolicyFontURL *string `json:"labelPolicyLightFontURL,omitempty"`
}

type LabelSetting struct {
	*Setting
	Settings LabelSettings
}

type PasswordComplexitySettings struct {
	MinLength    uint64 `json:"minLength,omitempty"`
	HasLowercase bool   `json:"hasLowercase,omitempty"`
	HasUppercase bool   `json:"hasUppercase,omitempty"`
	HasNumber    bool   `json:"hasNumber,omitempty"`
	HasSymbol    bool   `json:"hasSymbol,omitempty"`
}

type PasswordComplexitySetting struct {
	*Setting
	Settings PasswordComplexitySettings
}

type PasswordExpirySettings struct {
	ExpireWarnDays uint64 `json:"expireWarnDays,omitempty"`
	MaxAgeDays     uint64 `json:"maxAgeDays,omitempty"`
}

type PasswordExpirySetting struct {
	*Setting
	Settings PasswordExpirySettings
}

type LockoutSettings struct {
	MaxPasswordAttempts uint64 `json:"maxPasswordAttempts,omitempty"`
	MaxOTPAttempts      uint64 `json:"maxOtpAttempts,omitempty"`
	ShowLockOutFailures bool   `json:"showLockOutFailures,omitempty"`
}

type LockoutSetting struct {
	*Setting
	Settings LockoutSettings
}

type DomainSettings struct {
	UserLoginMustBeDomain                  bool `json:"userLoginMustBeDomain,omitempty"`
	ValidateOrgDomains                     bool `json:"validateOrgDomains,omitempty"`
	SMTPSenderAddressMatchesInstanceDomain bool `json:"smtpSenderAddressMatchesInstanceDomain,omitempty"`
}

type DomainSetting struct {
	*Setting
	Settings DomainSettings
}

type SecuritySettings struct {
	Enabled               bool     `json:"enabled,omitempty"`
	EnableIframeEmbedding bool     `json:"enableIframe_embedding,omitempty"`
	AllowedOrigins        []string `json:"allowedOrigins,omitempty"`
	EnableImpersonation   bool     `json:"enableImpersonation,omitempty"`
}

type SecuritySetting struct {
	*Setting
	Settings SecuritySettings
}

type OrganizationSettings struct {
	OrganizationScopedUsernames    bool     `json:"organizationScopedUsernames,omitempty"`
	OldOrganizationScopedUsernames bool     `json:"oldOrganizationScopedUsernames,omitempty"`
	UsernameChanges                []string `json:"usernameChanges,omitempty"`
}

type OrganizationSetting struct {
	*Setting
	Settings OrganizationSettings
}

type settingsColumns interface {
	IDColumn() database.Column
	InstanceIDColumn() database.Column
	OrgIDColumn() database.Column
	TypeColumn() database.Column
	LabelStateColumn() database.Column
	SettingsColumn() database.Column
	CreatedAtColumn() database.Column
	UpdatedAtColumn() database.Column
}

type settingsConditions interface {
	InstanceIDCondition(id string) database.Condition
	OrgIDCondition(id *string) database.Condition
	IDCondition(id string) database.Condition
	TypeCondition(typ SettingType) database.Condition
	LabelStateCondition(typ LabelState) database.Condition
}

type Settings interface {
	GetSettings() []byte
}

type settingsChanges interface {
	SetType(state SettingType) database.Change
	SetSettings(settings string) database.Change
	SetUpdatedAt(updatedAt *time.Time) database.Change
}

type SettingsRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string, typ SettingType, opts ...database.QueryOption) (*Setting, error)
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*Setting, error)

	// CreateLogin(ctx context.Context, client database.QueryExecutor, setting *LoginSetting, changes ...database.Change) error
	// GetLogin(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*LoginSetting, error)
	// UpdateLogin(ctx context.Context, client database.QueryExecutor, setting *LoginSetting, changes ...database.Change) (int64, error)

	CreateLabel(ctx context.Context, client database.QueryExecutor, setting *LabelSetting) error
	GetLabel(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string, state LabelState) (*LabelSetting, error)
	UpdateLabel(ctx context.Context, client database.QueryExecutor, setting *LabelSetting, changes ...database.Change) (int64, error)
	ActivateLabelSetting(ctx context.Context, client database.QueryExecutor, setting *LabelSetting) error

	CreatePasswordComplexity(ctx context.Context, client database.QueryExecutor, setting *PasswordComplexitySetting) error
	GetPasswordComplexity(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*PasswordComplexitySetting, error)
	UpdatePasswordComplexity(ctx context.Context, client database.QueryExecutor, setting *PasswordComplexitySetting, changes ...database.Change) (int64, error)

	CreatePasswordExpiry(ctx context.Context, client database.QueryExecutor, setting *PasswordExpirySetting) error
	GetPasswordExpiry(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*PasswordExpirySetting, error)
	UpdatePasswordExpiry(ctx context.Context, client database.QueryExecutor, setting *PasswordExpirySetting, changes ...database.Change) (int64, error)

	CreateLockout(ctx context.Context, client database.QueryExecutor, setting *LockoutSetting) error
	GetLockout(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*LockoutSetting, error)
	UpdateLockout(ctx context.Context, client database.QueryExecutor, setting *LockoutSetting, changes ...database.Change) (int64, error)

	CreateSecurity(ctx context.Context, client database.QueryExecutor, setting *SecuritySetting) error
	GetSecurity(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*SecuritySetting, error)
	UpdateSecurity(ctx context.Context, client database.QueryExecutor, setting *SecuritySetting, changes ...database.Change) (int64, error)

	CreateDomain(ctx context.Context, client database.QueryExecutor, setting *DomainSetting) error
	GetDomain(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*DomainSetting, error)
	UpdateDomain(ctx context.Context, client database.QueryExecutor, setting *DomainSetting, changes ...database.Change) (int64, error)

	CreateOrg(ctx context.Context, client database.QueryExecutor, setting *OrganizationSetting) error
	GetOrg(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*OrganizationSetting, error)
	UpdateOrg(ctx context.Context, client database.QueryExecutor, setting *OrganizationSetting, changes ...database.Change) (int64, error)

	// Create is used for events reduction
	Create(ctx context.Context, client database.QueryExecutor, setting *Setting) error
	Delete(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string, typ SettingType) (int64, error)

	// DeleteSettingsForInstance is used when a Instance is deleted
	DeleteSettingsForInstance(ctx context.Context, client database.QueryExecutor, instanceID string) (int64, error)
	// DeleteSettingsForOrg is used ehwn an Organization is deleted
	DeleteSettingsForOrg(ctx context.Context, client database.QueryExecutor, orgID string) (int64, error)
}

type LoginRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*LoginSetting, error)
	Set(ctx context.Context, client database.QueryExecutor, setting *LoginSetting, changes ...database.Change) error
}

type LabelRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*LabelSetting, error)
	Set(ctx context.Context, client database.QueryExecutor, setting *LabelSetting, changes ...database.Change) error
	// ActivateLabelSetting(ctx context.Context, client database.QueryExecutor, setting *LabelSetting) error
}

type PasswordComplexityRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*PasswordComplexitySetting, error)
	Set(ctx context.Context, client database.QueryExecutor, setting *PasswordComplexitySetting, changes ...database.Change) error
}

type PasswordExpiryRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*PasswordExpirySetting, error)
	Set(ctx context.Context, client database.QueryExecutor, setting *PasswordExpirySetting, changes ...database.Change) error
}

type LockoutRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*LockoutSetting, error)
	Set(ctx context.Context, client database.QueryExecutor, setting *LockoutSetting, changes ...database.Change) error
}

type SecurityRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*SecuritySetting, error)
	Set(ctx context.Context, client database.QueryExecutor, setting *SecuritySetting, changes ...database.Change) error
}

type DomainRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*DomainSetting, error)
	Set(ctx context.Context, client database.QueryExecutor, setting *DomainSetting, changes ...database.Change) error
}

type OrganizationSettingRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*OrganizationSetting, error)
	Set(ctx context.Context, client database.QueryExecutor, setting *OrganizationSetting, changes ...database.Change) error
}
