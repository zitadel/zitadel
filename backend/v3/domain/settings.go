package domain

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	db_json "github.com/zitadel/zitadel/backend/v3/storage/database/json"
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
	OrgID      *string         `json:"organizationId,omitempty" db:"organization_id"`
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
	AllowUserNamePassword      bool             `json:"allowUsernamePassword,omitempty"`
	AllowRegister              bool             `json:"allowRegister,omitempty"`
	AllowExternalIDP           bool             `json:"allowExternalIdp,omitempty"`
	ForceMFA                   bool             `json:"forceMfa,omitempty"`
	ForceMFALocalOnly          bool             `json:"forceMFALocalOnly,omitempty"`
	HidePasswordReset          bool             `json:"hidePasswordReset,omitempty"`
	IgnoreUnknownUsernames     bool             `json:"ignoreUnknownUsernames,omitempty"`
	AllowDomainDiscovery       bool             `json:"allowDomainDiscovery,omitempty"`
	DisableLoginWithEmail      bool             `json:"disableLoginWithEmail,omitempty"`
	DisableLoginWithPhone      bool             `json:"disableLoginWithPhone,omitempty"`
	PasswordlessType           PasswordlessType `json:"passwordlessType,omitempty"`
	DefaultRedirectURI         string           `json:"defaultRedirectUri,omitempty"`
	PasswordCheckLifetime      time.Duration    `json:"passwordCheckLifetime,omitempty"`
	ExternalLoginCheckLifetime time.Duration    `json:"externalLoginCheckLifetime,omitempty"`
	MFAInitSkipLifetime        time.Duration    `json:"mfaInitSkipLifetime,omitempty"`
	SecondFactorCheckLifetime  time.Duration    `json:"secondFactorCheckLifetime,omitempty"`
	MultiFactorCheckLifetime   time.Duration    `json:"multiFactorCheckLifetime,omitempty"`

	MFAType           []MultiFactorType  `json:"mfaType"`
	SecondFactorTypes []SecondFactorType `json:"secondFactors"`
}
type loginSettingsJSONFieldsChanges interface {
	SetAllowUserNamePasswordField(value bool) db_json.JsonUpdate
	SetAllowRegisterField(value bool) db_json.JsonUpdate
	SetAllowExternalIDPField(value bool) db_json.JsonUpdate
	SetForceMFAField(value bool) db_json.JsonUpdate
	SetForceMFALocalOnlyField(value bool) db_json.JsonUpdate
	SetHidePasswordResetField(value bool) db_json.JsonUpdate
	SetIgnoreUnknownUsernamesField(value bool) db_json.JsonUpdate
	SetAllowDomainDiscoveryField(value bool) db_json.JsonUpdate
	SetDisableLoginWithEmailField(value bool) db_json.JsonUpdate
	SetDisableLoginWithPhoneField(value bool) db_json.JsonUpdate
	SetPasswordlessTypeField(value PasswordlessType) db_json.JsonUpdate
	SetDefaultRedirectURIField(value string) db_json.JsonUpdate
	SetPasswordCheckLifetimeField(value time.Duration) db_json.JsonUpdate
	SetExternalLoginCheckLifetimeField(value time.Duration) db_json.JsonUpdate
	SetMFAInitSkipLifetimeField(value time.Duration) db_json.JsonUpdate
	SetSecondFactorCheckLifetimeField(value time.Duration) db_json.JsonUpdate
	SetMultiFactorCheckLifetimeField(value time.Duration) db_json.JsonUpdate
	AddMFAType(value MultiFactorType) db_json.JsonUpdate
	RemoveMFAType(value MultiFactorType) db_json.JsonUpdate
	SetSecondFactorTypesField(value []SecondFactorType) db_json.JsonUpdate
}

type loginSettingsJsonChanges interface {
	loginSettingsJSONFieldsChanges
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

	LabelPolicyLightLogoURL *url.URL `json:"labelPolicyLightLogoURL,omitempty"`
	LabelPolicyDarkLogoURL  *url.URL `json:"labelPolicyDarkLogoURL,omitempty"`

	LabelPolicyLightIconURL *url.URL `json:"labelPolicyLightIconURL,omitempty"`
	LabelPolicyDarkIconURL  *url.URL `json:"labelPolicyDarkIconURL,omitempty"`

	LabelPolicyFontURL *url.URL `json:"labelPolicyLightFontURL,omitempty"`
}

type labelSettingsJSONFieldsChanges interface {
	SetPrimaryColorField(value string) db_json.JsonUpdate
	SetBackgroundColorField(value string) db_json.JsonUpdate
	SetWarnColorField(value string) db_json.JsonUpdate
	SetFontColorField(value string) db_json.JsonUpdate
	SetPrimaryColorDarkField(value string) db_json.JsonUpdate
	SetBackgroundColorDarkField(value string) db_json.JsonUpdate
	SetWarnColorDarkField(value string) db_json.JsonUpdate
	SetFontColorDarkField(value string) db_json.JsonUpdate
	SetHideLoginNameSuffixField(value bool) db_json.JsonUpdate
	SetErrorMsgPopupField(value bool) db_json.JsonUpdate
	SetDisableWatermarkField(value bool) db_json.JsonUpdate
	SetThemeModeField(value LabelPolicyThemeMode) db_json.JsonUpdate
	SetLabelPolicyLightLogoURL(value *url.URL) db_json.JsonUpdate
	SetLabelPolicyDarkLogoURL(value *url.URL) db_json.JsonUpdate
	SetLabelPolicyLightIconURL(value *url.URL) db_json.JsonUpdate
	SetLabelPolicyDarkIconURL(value *url.URL) db_json.JsonUpdate
	SetLabelPolicyFontURL(value *url.URL) db_json.JsonUpdate
}

type labelSettingsJsonChanges interface {
	labelSettingsJSONFieldsChanges
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

type passwordComplexityJSONFieldsChanges interface {
	SetMinLengthField(value uint64) db_json.JsonUpdate
	SetHasLowercaseField(value bool) db_json.JsonUpdate
	SetHasUppercaseField(value bool) db_json.JsonUpdate
	SetHasNumberField(value bool) db_json.JsonUpdate
	SetHasSymbolField(value bool) db_json.JsonUpdate
}

type passwordComplexitySettingsJsonChanges interface {
	passwordComplexityJSONFieldsChanges
}

type PasswordComplexitySetting struct {
	*Setting
	Settings PasswordComplexitySettings
}

type PasswordExpirySettings struct {
	ExpireWarnDays uint64 `json:"expireWarnDays,omitempty"`
	MaxAgeDays     uint64 `json:"maxAgeDays,omitempty"`
}

type passwordExpiryJsonUpdates interface {
	SetExpireWarnDays(value uint64) db_json.JsonUpdate
	SetMaxAgeDays(value uint64) db_json.JsonUpdate
}

type passwordExpirySettingsJsonChanges interface {
	passwordExpiryJsonUpdates
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

type lockoutJsonUpdates interface {
	SetMaxPasswordAttempts(value uint64) db_json.JsonUpdate
	SetMaxOTPAttempts(value uint64) db_json.JsonUpdate
	SetShowLockOutFailures(value bool) db_json.JsonUpdate
}

type lockoutSettingsJsonChanges interface {
	lockoutJsonUpdates
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

type domainJsonUpdates interface {
	SetUserLoginMustBeDomain(value bool) db_json.JsonUpdate
	SetValidateOrgDomains(value bool) db_json.JsonUpdate
	SetSMTPSenderAddressMatchesInstanceDomain(value bool) db_json.JsonUpdate
}

type domainSettingsJsonChanges interface {
	domainJsonUpdates
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

type securityJsonUpdates interface {
	SetEnabled(value bool) db_json.JsonUpdate
	SetEnableIframeEmbedding(value bool) db_json.JsonUpdate
	SetAllowedOrigins(value []string) db_json.JsonUpdate
	SetEnableImpersonation(value bool) db_json.JsonUpdate
}

type securitySettingsJsonChanges interface {
	securityJsonUpdates
}

type SecuritySetting struct {
	*Setting
	Settings SecuritySettings
}

type OrganizationSettings struct {
	OrganizationScopedUsernames bool `json:"organizationScopedUsernames,omitempty"`
}

type organizationJsonUpdates interface {
	SetOrganizationScopedUsernames(value bool) db_json.JsonUpdate
}

type organizationSettingsJsonChanges interface {
	organizationJsonUpdates
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
	OwnerTypeColumn() database.Column
	LabelStateColumn() database.Column
	SettingsColumn() database.Column
	CreatedAtColumn() database.Column
	UpdatedAtColumn() database.Column
	SetLabelSettings(changes ...db_json.JsonUpdate) database.Change
}

type settingsConditions interface {
	InstanceIDCondition(id string) database.Condition
	OrgIDCondition(id *string) database.Condition
	IDCondition(id string) database.Condition
	TypeCondition(typ SettingType) database.Condition
	OwnerTypeCondition(typ OwnerType) database.Condition
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

	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*Setting, error)
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*Setting, error)

	// Create is used for events reduction
	Create(ctx context.Context, client database.QueryExecutor, setting *Setting) error
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)

	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
}

type LoginRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges
	loginSettingsJsonChanges

	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*LoginSetting, error)
	Set(ctx context.Context, client database.QueryExecutor, setting *LoginSetting) error
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
}

type LabelRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	labelSettingsJsonChanges

	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*LabelSetting, error)
	Set(ctx context.Context, client database.QueryExecutor, setting *LabelSetting) error
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	Reset(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
	ActivateLabelSetting(ctx context.Context, client database.QueryExecutor, setting *LabelSetting) error
	ActivateLabelSettingEvent(ctx context.Context, client database.QueryExecutor, condition database.Condition, UpdateAt time.Time) (int64, error)
}

type PasswordComplexityRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	passwordComplexitySettingsJsonChanges

	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*PasswordComplexitySetting, error)
	Set(ctx context.Context, client database.QueryExecutor, setting *PasswordComplexitySetting) error
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	Reset(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}

type PasswordExpiryRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	passwordExpirySettingsJsonChanges

	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*PasswordExpirySetting, error)
	Set(ctx context.Context, client database.QueryExecutor, setting *PasswordExpirySetting) error
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	Reset(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}

type LockoutRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	lockoutSettingsJsonChanges

	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*LockoutSetting, error)
	Set(ctx context.Context, client database.QueryExecutor, setting *LockoutSetting) error
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	Reset(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}

type SecurityRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	securitySettingsJsonChanges

	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*SecuritySetting, error)
	Set(ctx context.Context, client database.QueryExecutor, setting *SecuritySetting) error
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	Reset(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
	SetEvent(ctx context.Context, client database.QueryExecutor, setting *SecuritySetting, changes ...database.Change) (int64, error)
}

type DomainRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	domainSettingsJsonChanges

	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*DomainSetting, error)
	Set(ctx context.Context, client database.QueryExecutor, setting *DomainSetting) error
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	Reset(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}

type OrganizationSettingRepository interface {
	settingsColumns
	settingsConditions
	settingsChanges

	organizationSettingsJsonChanges

	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*OrganizationSetting, error)
	Set(ctx context.Context, client database.QueryExecutor, setting *OrganizationSetting) error
	SetEvent(ctx context.Context, client database.QueryExecutor, setting *OrganizationSetting) (int64, error)
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	Reset(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}
