package domain

import (
	"context"
	"encoding/json"
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
	SettingTypeBranding
	SettingTypeDomain
	SettingTypeLegalAndSupport
	SettingTypeLockout
	SettingTypeGeneral
	SettingTypeSecurity
	SettingTypeOrganization
)

type Setting struct {
	ID         string          `json:"id,omitempty" db:"id"`
	InstanceID string          `json:"instanceId,omitempty" db:"instance_id"`
	OrgID      *string         `json:"orgId,omitempty" db:"org_id"`
	Type       SettingType     `json:"type,omitempty" db:"type"`
	Settings   json.RawMessage `json:"settings,omitempty" db:"settings"`
	CreatedAt  time.Time       `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt  time.Time       `json:"updatedAt,omitzero" db:"updated_at"`
}

type PasswordlessType int32

const (
	PasswordlessTypeNotAllowed PasswordlessType = iota
	PasswordlessTypeAllowed

	passwordlessCount
)

type MultiFactorType int32

const (
	MultiFactorTypeUnspecified MultiFactorType = iota
	MultiFactorTypeU2FWithPIN

	multiFactorCount
)

type SecondFactorType int32

const (
	SecondFactorTypeUnspecified SecondFactorType = iota
	SecondFactorTypeTOTP
	SecondFactorTypeU2F
	SecondFactorTypeOTPEmail
	SecondFactorTypeOTPSMS

	secondFactorCount
)

type LoginSettings struct {
	IsDefault                  bool             `json:"isDefault,omitempty"`
	AllowUserNamePassword      bool             `json:"allowUsernamePassword,omitempty"`
	AllowRegister              bool             `json:"allowRegister,omitempty"`
	AllowExternalSetting       bool             `json:"allowExternalIdp,omitempty"`
	ForceMFA                   bool             `json:"forceMFA,omitempty"`
	ForceMFALocalOnly          bool             `json:"forceMFALocalOnly,omitempty"`
	HidePasswordReset          bool             `json:"hidePasswordReset,omitempty"`
	IgnoreUnknownUsernames     bool             `json:"ignoreUnknownUsernames,omitempty"`
	AllowDomainDiscovery       bool             `json:"allowDomainDiscovery,omitempty"`
	DisableLoginWithEmail      bool             `json:"disableLoginWithEmail,omitempty"`
	DisableLoginWithPhone      bool             `json:"disableLoginWithPhone,omitempty"`
	PasswordlessType           PasswordlessType `json:"passwordlessType,omitempty"`
	DefaultRedirectURI         string           `json:"defaultRedirectURI,omitempty"`
	PasswordCheckLifetime      time.Duration    `json:"passwordCheckLifetime,omitempty"`
	ExternalLoginCheckLifetime time.Duration    `json:"externalLoginCheckLifetime,omitempty"`
	MFAInitSkipLifetime        time.Duration    `json:"mfaInitSkipLifetime,omitempty"`
	SecondFactorCheckLifetime  time.Duration    `json:"secondFactorCheckLifetime,omitempty"`
	MultiFactorCheckLifetime   time.Duration    `json:"multiFactorCheckLifetime,omitempty"`

	MFAType           []MultiFactorType  `json:"mfaType"`
	SecondFactorTypes []SecondFactorType `json:"second_factors"`
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

type LabelSettings struct {
	IsDefault           bool                 `json:"isDefault,omitempty"`
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

	LabelPolicyLightIconURL *string `json:"labelPolicyLightIconURL,omitempty"`
	LabelPolicyDarkIconURL  *string `json:"labelPolicyDarkIconURL,omitempty"`

	LabelPolicyFontURL *string `json:"labelPolicyLightFontURL,omitempty"`
}

type LabelSetting struct {
	*Setting
	Settings LabelSettings
}

type PasswordComplexitySettings struct {
	IsDefault    bool   `json:"isDefault,omitempty"`
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
	IsDefault      bool   `json:"isDefault,omitempty"`
	ExpireWarnDays uint64 `json:"expireWarnDays,omitempty"`
	MaxAgeDays     uint64 `json:"maxAgeDays,omitempty"`
}

type PasswordExpirySetting struct {
	*Setting
	Settings PasswordExpirySettings
}

type LockoutSettings struct {
	IsDefault           bool   `json:"isDefault,omitempty"`
	MaxPasswordAttempts uint64 `json:"maxPasswordAttempts,omitempty"`
	MaxOTPAttempts      uint64 `json:"maxOTPAttempts,omitempty"`
	ShowLockOutFailures bool   `json:"showLockOutFailures,omitempty"`
}

type LockoutSetting struct {
	*Setting
	Settings LockoutSettings
}

type DomainSettings struct {
	IsDefault                              bool `json:"isDefault,omitempty"`
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
	EnableIframeEmbedding bool     `json:"enable_iframe_embedding,omitempty"`
	AllowedOrigins        []string `json:"allowedOrigins,omitempty"`
	EnableImpersonation   bool     `json:"enable_impersonation,omitempty"`
}

type SecuritySetting struct {
	*Setting
	Settings SecuritySettings
}

type OrgSettings struct {
	OrganizationScopedUsernames    bool     `json:"organizationScopedUsernames,omitempty"`
	oldOrganizationScopedUsernames bool     `json:"oldOrganizationScopedUsernames,omitempty"`
	usernameChanges                []string `json:"usernameChanges,omitempty"`
}

type OrgSetting struct {
	*Setting
	Settings OrgSettings
}

func (s *OrgSetting) getSettings() any {
	return s.Settings
}

type settingsColumns interface {
	IDColumn() database.Column
	InstanceIDColumn() database.Column
	OrgIDColumn() database.Column
	TypeColumn() database.Column
	SettingsColumn() database.Column
	CreatedAtColumn() database.Column
	UpdatedAtColumn() database.Column
}

type settingsConditions interface {
	InstanceIDCondition(id string) database.Condition
	OrgIDCondition(id *string) database.Condition
	IDCondition(id string) database.Condition
	TypeCondition(typ SettingType) database.Condition
}

type Settings interface {
	GetSettings() []byte
}

type settingsChanges interface {
	SetType(state SettingType) database.Change
	SetSettings(settings string) database.Change
}

type SettingsRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	Get(ctx context.Context, instanceID string, orgID *string, typ SettingType) (*Setting, error)
	List(ctx context.Context, conditions ...database.Condition) ([]*Setting, error)

	GetLogin(ctx context.Context, instanceID string, orgID *string) (*LoginSetting, error)
	UpdateLogin(ctx context.Context, setting *LoginSetting) (int64, error)
	DeleteLogin(ctx context.Context, instanceID string, orgID *string) (int64, error)

	GetLabel(ctx context.Context, instanceID string, orgID *string) (*LabelSetting, error)
	UpdateLabel(ctx context.Context, setting *LabelSetting) (int64, error)
	// DeleteLabel(ctx context.Context, instanceID string, orgID *string) (int64, error)

	GetPasswordComplexity(ctx context.Context, instanceID string, orgID *string) (*PasswordComplexitySetting, error)
	UpdatePasswordComplexity(ctx context.Context, setting *PasswordComplexitySetting) (int64, error)

	GetPasswordExpiry(ctx context.Context, instanceID string, orgID *string) (*PasswordExpirySetting, error)
	UpdatePasswordExpiry(ctx context.Context, setting *PasswordExpirySetting) (int64, error)

	GetLockout(ctx context.Context, instanceID string, orgID *string) (*LockoutSetting, error)
	UpdateLockout(ctx context.Context, setting *LockoutSetting) (int64, error)

	GetSecurity(ctx context.Context, instanceID string, orgID *string) (*SecuritySetting, error)
	UpdateSecurity(ctx context.Context, setting *SecuritySetting) (int64, error)

	GetDomain(ctx context.Context, instanceID string, orgID *string) (*DomainSetting, error)
	UpdateDomain(ctx context.Context, setting *DomainSetting) (int64, error)

	GetOrg(ctx context.Context, instanceID string, orgID *string) (*OrgSetting, error)
	UpdateOrg(ctx context.Context, setting *OrgSetting) (int64, error)

	Create(ctx context.Context, setting *Setting) error
	// Update(ctx context.Context, id string, instanceID string, orgID *string, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, instanceID string, orgID *string, typ SettingType) (int64, error)

	DeleteSettingsForInstance(ctx context.Context, instanceID string) (int64, error)
	DeleteSettingsForOrg(ctx context.Context, orgID *string) (int64, error)
}
