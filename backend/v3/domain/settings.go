package domain

import (
	"context"
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
	SettingTypeBranding
	SettingTypePasswordComplexity
	SettingTypePasswordExpiry
	SettingTypeDomain
	SettingTypeLockout
	SettingTypeSecurity
	SettingTypeOrganization
	SettingTypeNotification
	SettingTypeLegalAndSupport
)

//go:generate enumer -type OwnerType -transform snake -trimprefix OwnerType -sql
type OwnerType uint8

const (
	OwnerTypeUnspecified OwnerType = iota
	OwnerTypeSystem
	OwnerTypeInstance
	OwnerTypeOrganization
)

//go:generate enumer -type SettingState -transform snake -trimprefix SettingState -sql
type SettingState int32

const (
	SettingStateUnspecified SettingState = iota
	SettingStateActivated
	SettingStatePreview
)

type Settings struct {
	ID             string       `json:"id,omitempty" db:"id"`
	InstanceID     string       `json:"instanceId,omitempty" db:"instance_id"`
	OrganizationID *string      `json:"organizationId,omitempty" db:"organization_id"`
	Type           SettingType  `json:"type,omitempty" db:"type"`
	OwnerType      OwnerType    `json:"ownerType,omitempty" db:"owner_type"`
	State          SettingState `json:"state,omitempty" db:"state"`
	Settings       []byte       `json:"settings,omitempty" db:"settings"`
	CreatedAt      time.Time    `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt      *time.Time   `json:"updatedAt,omitzero" db:"updated_at"`
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

type loginSettingsJSONChanges interface {
	SetSettingFields(value *LoginSettings) database.Change

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
	AddSecondFactorTypesField(value SecondFactorType) db_json.JsonUpdate
	RemoveSecondFactorTypesField(value SecondFactorType) db_json.JsonUpdate
}

type LoginSettings struct {
	Settings
	AllowUserNamePassword      *bool             `json:"allowUsernamePassword,omitempty"`
	AllowRegister              *bool             `json:"allowRegister,omitempty"`
	AllowExternalIDP           *bool             `json:"allowExternalIdp,omitempty"`
	ForceMFA                   *bool             `json:"forceMfa,omitempty"`
	ForceMFALocalOnly          *bool             `json:"forceMFALocalOnly,omitempty"`
	HidePasswordReset          *bool             `json:"hidePasswordReset,omitempty"`
	IgnoreUnknownUsernames     *bool             `json:"ignoreUnknownUsernames,omitempty"`
	AllowDomainDiscovery       *bool             `json:"allowDomainDiscovery,omitempty"`
	DisableLoginWithEmail      *bool             `json:"disableLoginWithEmail,omitempty"`
	DisableLoginWithPhone      *bool             `json:"disableLoginWithPhone,omitempty"`
	PasswordlessType           *PasswordlessType `json:"passwordlessType,omitempty"`
	DefaultRedirectURI         *string           `json:"defaultRedirectUri,omitempty"`
	PasswordCheckLifetime      *time.Duration    `json:"passwordCheckLifetime,omitempty"`
	ExternalLoginCheckLifetime *time.Duration    `json:"externalLoginCheckLifetime,omitempty"`
	MFAInitSkipLifetime        *time.Duration    `json:"mfaInitSkipLifetime,omitempty"`
	SecondFactorCheckLifetime  *time.Duration    `json:"secondFactorCheckLifetime,omitempty"`
	MultiFactorCheckLifetime   *time.Duration    `json:"multiFactorCheckLifetime,omitempty"`

	MFAType           []MultiFactorType  `json:"mfaType"`
	SecondFactorTypes []SecondFactorType `json:"secondFactors"`
}

type BrandingPolicyThemeMode int32

const (
	BrandingPolicyThemeAuto BrandingPolicyThemeMode = iota
	BrandingPolicyThemeLight
	BrandingPolicyThemeDark
)

type brandingSettingsJSONChanges interface {
	SetSettingFields(value *BrandingSettings) database.Change

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
	SetThemeModeField(value BrandingPolicyThemeMode) db_json.JsonUpdate
	SetLightLogoURL(value *url.URL) db_json.JsonUpdate
	SetDarkLogoURL(value *url.URL) db_json.JsonUpdate
	SetLightIconURL(value *url.URL) db_json.JsonUpdate
	SetDarkIconURL(value *url.URL) db_json.JsonUpdate
	SetFontURL(value *url.URL) db_json.JsonUpdate
}

type BrandingSettings struct {
	Settings
	PrimaryColor        *string                  `json:"primaryColor,omitempty"`
	BackgroundColor     *string                  `json:"backgroundColor,omitempty"`
	WarnColor           *string                  `json:"warnColor,omitempty"`
	FontColor           *string                  `json:"fontColor,omitempty"`
	PrimaryColorDark    *string                  `json:"primaryColorDark,omitempty"`
	BackgroundColorDark *string                  `json:"backgroundColorDark,omitempty"`
	WarnColorDark       *string                  `json:"warnColorDark,omitempty"`
	FontColorDark       *string                  `json:"fontColorDark,omitempty"`
	HideLoginNameSuffix *bool                    `json:"hideLoginNameSuffix,omitempty"`
	ErrorMsgPopup       *bool                    `json:"errorMsgPopup,omitempty"`
	DisableWatermark    *bool                    `json:"disableMsgPopup,omitempty"`
	ThemeMode           *BrandingPolicyThemeMode `json:"themeMode,omitempty"`

	LightLogoURL *url.URL `json:"lightLogoUrl,omitempty"`
	DarkLogoURL  *url.URL `json:"darkLogoUrl,omitempty"`

	LightIconURL *url.URL `json:"lightIconUrl,omitempty"`
	DarkIconURL  *url.URL `json:"darkIconUrl,omitempty"`

	FontURL *url.URL `json:"fontUrl,omitempty"`
}

type passwordComplexitySettingsJSONChanges interface {
	SetSettingFields(value *PasswordComplexitySettings) database.Change

	SetMinLengthField(value uint64) db_json.JsonUpdate
	SetHasLowercaseField(value bool) db_json.JsonUpdate
	SetHasUppercaseField(value bool) db_json.JsonUpdate
	SetHasNumberField(value bool) db_json.JsonUpdate
	SetHasSymbolField(value bool) db_json.JsonUpdate
}

type PasswordComplexitySettings struct {
	Settings
	MinLength    *uint64 `json:"minLength,omitempty"`
	HasLowercase *bool   `json:"hasLowercase,omitempty"`
	HasUppercase *bool   `json:"hasUppercase,omitempty"`
	HasNumber    *bool   `json:"hasNumber,omitempty"`
	HasSymbol    *bool   `json:"hasSymbol,omitempty"`
}

type passwordExpirySettingsJSONChanges interface {
	SetSettingFields(value *PasswordExpirySettings) database.Change

	SetExpireWarnDays(value uint64) db_json.JsonUpdate
	SetMaxAgeDays(value uint64) db_json.JsonUpdate
}

type PasswordExpirySettings struct {
	Settings
	ExpireWarnDays *uint64 `json:"expireWarnDays,omitempty"`
	MaxAgeDays     *uint64 `json:"maxAgeDays,omitempty"`
}

type lockoutSettingsJSONChanges interface {
	SetSettingFields(value *LockoutSettings) database.Change

	SetMaxPasswordAttempts(value uint64) db_json.JsonUpdate
	SetMaxOTPAttempts(value uint64) db_json.JsonUpdate
	SetShowLockOutFailures(value bool) db_json.JsonUpdate
}

type LockoutSettings struct {
	Settings
	MaxPasswordAttempts *uint64 `json:"maxPasswordAttempts,omitempty"`
	MaxOTPAttempts      *uint64 `json:"maxOtpAttempts,omitempty"`
	ShowLockOutFailures *bool   `json:"showLockOutFailures,omitempty"`
}

type domainSettingsJSONChanges interface {
	SetSettingFields(value *DomainSettings) database.Change

	SetUserLoginMustBeDomain(value bool) db_json.JsonUpdate
	SetValidateOrgDomains(value bool) db_json.JsonUpdate
	SetSMTPSenderAddressMatchesInstanceDomain(value bool) db_json.JsonUpdate
}

type DomainSettings struct {
	Settings
	UserLoginMustBeDomain                  *bool `json:"userLoginMustBeDomain,omitempty"`
	ValidateOrgDomains                     *bool `json:"validateOrgDomains,omitempty"`
	SMTPSenderAddressMatchesInstanceDomain *bool `json:"smtpSenderAddressMatchesInstanceDomain,omitempty"`
}

type securitySettingsJSONChanges interface {
	SetSettingFields(value *SecuritySettings) database.Change

	SetEnabled(value bool) db_json.JsonUpdate
	SetEnableIframeEmbedding(value bool) db_json.JsonUpdate
	AddAllowedOrigins(value string) db_json.JsonUpdate
	RemoveAllowedOrigins(value string) db_json.JsonUpdate
	SetEnableImpersonation(value bool) db_json.JsonUpdate
}

type SecuritySettings struct {
	Settings
	Enabled               *bool    `json:"enabled,omitempty"`
	EnableIframeEmbedding *bool    `json:"enableIframe_embedding,omitempty"`
	AllowedOrigins        []string `json:"allowedOrigins,omitempty"`
	EnableImpersonation   *bool    `json:"enableImpersonation,omitempty"`
}

type organizationSettingsJSONChanges interface {
	SetSettingFields(value *OrganizationSettings) database.Change

	SetOrganizationScopedUsernames(value bool) db_json.JsonUpdate
}

type OrganizationSettings struct {
	Settings
	OrganizationScopedUsernames *bool `json:"organizationScopedUsernames,omitempty"`
}

type settingsColumns interface {
	IDColumn() database.Column
	InstanceIDColumn() database.Column
	OrganizationIDColumn() database.Column
	TypeColumn() database.Column
	OwnerTypeColumn() database.Column
	StateColumn() database.Column
	SettingsColumn() database.Column
	CreatedAtColumn() database.Column
	UpdatedAtColumn() database.Column
}

type settingsConditions interface {
	InstanceIDCondition(id string) database.Condition
	OrganizationIDCondition(id *string) database.Condition
	IDCondition(id string) database.Condition
	TypeCondition(typ SettingType) database.Condition
	OwnerTypeCondition(typ OwnerType) database.Condition
	StateCondition(typ SettingState) database.Condition
}

type settingsChanges interface {
	SetSettings(settings string) database.Change
	SetUpdatedAt(updatedAt *time.Time) database.Change
}

type setting interface {
	ToJsonChanges() []db_json.JsonUpdate
	GetSettings() []byte
}

type settingsRepository[T setting] interface {
	settingsColumns
	settingsConditions
	settingsChanges

	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*setting, error)
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*setting, error)

	Set(ctx context.Context, client database.QueryExecutor, setting *setting) error
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}

type LoginSettingsRepository interface {
	settingsRepository[LoginSettings]
	loginSettingsJSONChanges
}

type BrandingSettingsRepository interface {
	settingsRepository[BrandingSettings]
	brandingSettingsJSONChanges
	Activate(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}

type PasswordComplexitySettingsRepository interface {
	settingsRepository[PasswordComplexitySettings]
	passwordComplexitySettingsJSONChanges
}

type PasswordExpirySettingsRepository interface {
	settingsRepository[PasswordExpirySettings]
	passwordExpirySettingsJSONChanges
}

type LockoutSettingsRepository interface {
	settingsRepository[LockoutSettings]
	lockoutSettingsJSONChanges
}

type SecuritySettingsRepository interface {
	settingsRepository[SecuritySettings]
	securitySettingsJSONChanges
}

type DomainSettingsRepository interface {
	settingsRepository[DomainSettings]
	domainSettingsJSONChanges
}

type OrganizationSettingsRepository interface {
	settingsRepository[OrganizationSettings]
	organizationSettingsJSONChanges
}
