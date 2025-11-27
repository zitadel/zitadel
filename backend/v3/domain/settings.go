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
	SettingTypeLogin SettingType = iota
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

//go:generate enumer -type SettingState -transform snake -trimprefix SettingState -sql
type SettingState int32

const (
	SettingStateActive SettingState = iota
	SettingStatePreview
)

type Settings struct {
	ID             string       `json:"id,omitempty" db:"id"`
	InstanceID     string       `json:"instanceId,omitempty" db:"instance_id"`
	OrganizationID *string      `json:"organizationId,omitempty" db:"organization_id"`
	Type           SettingType  `json:"type,omitempty" db:"type"`
	State          SettingState `json:"state,omitempty" db:"state"`
	Settings       []byte       `json:"settings,omitempty" db:"settings"`
	CreatedAt      time.Time    `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt      time.Time    `json:"updatedAt,omitzero" db:"updated_at"`
}

type settingsColumns interface {
	IDColumn() database.Column
	InstanceIDColumn() database.Column
	OrganizationIDColumn() database.Column
	TypeColumn() database.Column
	StateColumn() database.Column
	SettingsColumn() database.Column
	CreatedAtColumn() database.Column
	UpdatedAtColumn() database.Column
	PrimaryKeyColumns() []database.Column
	UniqueColumns() []database.Column
}

type settingsConditions interface {
	InstanceIDCondition(id string) database.Condition
	OrganizationIDCondition(id *string) database.Condition
	IDCondition(id string) database.Condition
	PrimaryKeyCondition(instanceID, id string) database.Condition
	UniqueCondition(instanceID string, orgID *string, typ SettingType, state SettingState) database.Condition
}

type settingsChanges interface {
	SetSettings(settings string) database.Change
	SetUpdatedAt(updatedAt *time.Time) database.Change
}

type settingsRepository[T any] interface {
	settingsColumns
	settingsConditions
	settingsChanges

	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*T, error)
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*T, error)

	Set(ctx context.Context, client database.QueryExecutor, settings *T) error
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
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
	SetSettingFields(value LoginSettingsAttributes) database.Change

	SetAllowUserNamePassword(value bool) db_json.JsonUpdate
	SetAllowRegister(value bool) db_json.JsonUpdate
	SetAllowExternalIDP(value bool) db_json.JsonUpdate
	SetForceMFA(value bool) db_json.JsonUpdate
	SetForceMFALocalOnly(value bool) db_json.JsonUpdate
	SetHidePasswordReset(value bool) db_json.JsonUpdate
	SetIgnoreUnknownUsernames(value bool) db_json.JsonUpdate
	SetAllowDomainDiscovery(value bool) db_json.JsonUpdate
	SetDisableLoginWithEmail(value bool) db_json.JsonUpdate
	SetDisableLoginWithPhone(value bool) db_json.JsonUpdate
	SetPasswordlessType(value PasswordlessType) db_json.JsonUpdate
	SetDefaultRedirectURI(value string) db_json.JsonUpdate
	SetPasswordCheckLifetime(value time.Duration) db_json.JsonUpdate
	SetExternalLoginCheckLifetime(value time.Duration) db_json.JsonUpdate
	SetMFAInitSkipLifetime(value time.Duration) db_json.JsonUpdate
	SetSecondFactorCheckLifetime(value time.Duration) db_json.JsonUpdate
	SetMultiFactorCheckLifetime(value time.Duration) db_json.JsonUpdate

	AddMFAType(value MultiFactorType) database.Change
	RemoveMFAType(value MultiFactorType) database.Change
	AddSecondFactorTypes(value SecondFactorType) database.Change
	RemoveSecondFactorTypes(value SecondFactorType) database.Change
}

type LoginSettings struct {
	Settings
	LoginSettingsAttributes
}

type LoginSettingsAttributes struct {
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

	MFAType           []MultiFactorType  `json:"mfaType,omitempty"`
	SecondFactorTypes []SecondFactorType `json:"secondFactors,omitempty"`
}

type LoginSettingsRepository interface {
	settingsRepository[LoginSettings]
	loginSettingsJSONChanges

	SetColumns(ctx context.Context, client database.QueryExecutor, settings *Settings, changes ...database.Change) error
}

type BrandingPolicyThemeMode int32

const (
	BrandingPolicyThemeAuto BrandingPolicyThemeMode = iota
	BrandingPolicyThemeLight
	BrandingPolicyThemeDark
)

type brandingSettingsJSONChanges interface {
	SetSettingFields(value BrandingSettingsAttributes) database.Change

	SetPrimaryColorLight(value string) db_json.JsonUpdate
	SetBackgroundColorLight(value string) db_json.JsonUpdate
	SetWarnColorLight(value string) db_json.JsonUpdate
	SetFontColorLight(value string) db_json.JsonUpdate
	SetLogoURLLight(value url.URL) db_json.JsonUpdate
	SetIconURLLight(value url.URL) db_json.JsonUpdate

	SetPrimaryColorDark(value string) db_json.JsonUpdate
	SetBackgroundColorDark(value string) db_json.JsonUpdate
	SetWarnColorDark(value string) db_json.JsonUpdate
	SetFontColorDark(value string) db_json.JsonUpdate
	SetLogoURLDark(value url.URL) db_json.JsonUpdate
	SetIconURLDark(value url.URL) db_json.JsonUpdate

	SetHideLoginNameSuffix(value bool) db_json.JsonUpdate
	SetErrorMsgPopup(value bool) db_json.JsonUpdate
	SetDisableWatermark(value bool) db_json.JsonUpdate
	SetThemeMode(value BrandingPolicyThemeMode) db_json.JsonUpdate
	SetFontURL(value url.URL) db_json.JsonUpdate
}

type BrandingSettings struct {
	Settings
	BrandingSettingsAttributes
}

type BrandingSettingsAttributes struct {
	PrimaryColorLight    *string  `json:"primaryColorLight,omitempty"`
	BackgroundColorLight *string  `json:"backgroundColorLight,omitempty"`
	WarnColorLight       *string  `json:"warnColorLight,omitempty"`
	FontColorLight       *string  `json:"fontColorLight,omitempty"`
	LogoURLLight         *url.URL `json:"logoUrlLight,omitempty"`
	IconURLLight         *url.URL `json:"iconUrlLight,omitempty"`

	PrimaryColorDark    *string  `json:"primaryColorDark,omitempty"`
	BackgroundColorDark *string  `json:"backgroundColorDark,omitempty"`
	WarnColorDark       *string  `json:"warnColorDark,omitempty"`
	FontColorDark       *string  `json:"fontColorDark,omitempty"`
	LogoURLDark         *url.URL `json:"logoUrlDark,omitempty"`
	IconURLDark         *url.URL `json:"iconUrlDark,omitempty"`

	HideLoginNameSuffix *bool                    `json:"hideLoginNameSuffix,omitempty"`
	ErrorMsgPopup       *bool                    `json:"errorMsgPopup,omitempty"`
	DisableWatermark    *bool                    `json:"disableMsgPopup,omitempty"`
	ThemeMode           *BrandingPolicyThemeMode `json:"themeMode,omitempty"`
	FontURL             *url.URL                 `json:"fontUrl,omitempty"`
}

type BrandingSettingsRepository interface {
	settingsRepository[BrandingSettings]
	brandingSettingsJSONChanges

	StateCondition(typ SettingState) database.Condition
	Activate(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
	ActivateAt(ctx context.Context, client database.QueryExecutor, condition database.Condition, updatedAt time.Time) (int64, error)
	SetColumns(ctx context.Context, client database.QueryExecutor, settings *Settings, changes ...database.Change) error
}

type passwordComplexitySettingsJSONChanges interface {
	SetSettingFields(value PasswordComplexitySettingsAttributes) database.Change

	SetMinLength(value uint64) db_json.JsonUpdate
	SetHasLowercase(value bool) db_json.JsonUpdate
	SetHasUppercase(value bool) db_json.JsonUpdate
	SetHasNumber(value bool) db_json.JsonUpdate
	SetHasSymbol(value bool) db_json.JsonUpdate
}

type PasswordComplexitySettings struct {
	Settings
	PasswordComplexitySettingsAttributes
}

type PasswordComplexitySettingsAttributes struct {
	MinLength    *uint64 `json:"minLength,omitempty"`
	HasLowercase *bool   `json:"hasLowercase,omitempty"`
	HasUppercase *bool   `json:"hasUppercase,omitempty"`
	HasNumber    *bool   `json:"hasNumber,omitempty"`
	HasSymbol    *bool   `json:"hasSymbol,omitempty"`
}

type PasswordComplexitySettingsRepository interface {
	settingsRepository[PasswordComplexitySettings]
	passwordComplexitySettingsJSONChanges
}

type passwordExpirySettingsJSONChanges interface {
	SetSettingFields(value PasswordExpirySettingsAttributes) database.Change

	SetExpireWarnDays(value uint64) db_json.JsonUpdate
	SetMaxAgeDays(value uint64) db_json.JsonUpdate
}

type PasswordExpirySettings struct {
	Settings
	PasswordExpirySettingsAttributes
}

type PasswordExpirySettingsAttributes struct {
	ExpireWarnDays *uint64 `json:"expireWarnDays,omitempty"`
	MaxAgeDays     *uint64 `json:"maxAgeDays,omitempty"`
}

type PasswordExpirySettingsRepository interface {
	settingsRepository[PasswordExpirySettings]
	passwordExpirySettingsJSONChanges
}

type lockoutSettingsJSONChanges interface {
	SetSettingFields(value LockoutSettingsAttributes) database.Change

	SetMaxPasswordAttempts(value uint64) db_json.JsonUpdate
	SetMaxOTPAttempts(value uint64) db_json.JsonUpdate
	SetShowLockOutFailures(value bool) db_json.JsonUpdate
}

type LockoutSettings struct {
	Settings
	LockoutSettingsAttributes
}

type LockoutSettingsAttributes struct {
	MaxPasswordAttempts *uint64 `json:"maxPasswordAttempts,omitempty"`
	MaxOTPAttempts      *uint64 `json:"maxOtpAttempts,omitempty"`
	ShowLockOutFailures *bool   `json:"showLockOutFailures,omitempty"`
}

type LockoutSettingsRepository interface {
	settingsRepository[LockoutSettings]
	lockoutSettingsJSONChanges
}

type securitySettingsJSONChanges interface {
	SetSettingFields(value SecuritySettingsAttributes) database.Change

	SetEnableIframeEmbedding(value bool) db_json.JsonUpdate
	SetAllowedOrigins(values []string) db_json.JsonUpdate
	SetEnableImpersonation(value bool) db_json.JsonUpdate
}

type SecuritySettings struct {
	Settings
	SecuritySettingsAttributes
}

type SecuritySettingsAttributes struct {
	EnableIframeEmbedding *bool    `json:"enableIframeEmbedding,omitempty"`
	AllowedOrigins        []string `json:"allowedOrigins,omitempty"`
	EnableImpersonation   *bool    `json:"enableImpersonation,omitempty"`
}

type SecuritySettingsRepository interface {
	settingsRepository[SecuritySettings]
	securitySettingsJSONChanges
}

type domainSettingsJSONChanges interface {
	SetSettingFields(value DomainSettingsAttributes) database.Change

	SetLoginNameIncludesDomain(value bool) db_json.JsonUpdate
	SetRequireOrgDomainVerification(value bool) db_json.JsonUpdate
	SetSMTPSenderAddressMatchesInstanceDomain(value bool) db_json.JsonUpdate
}

type DomainSettings struct {
	Settings
	DomainSettingsAttributes
}
type DomainSettingsAttributes struct {
	LoginNameIncludesDomain                *bool `json:"loginNameIncludesDomain,omitempty"`
	RequireOrgDomainVerification           *bool `json:"requireOrgDomainVerification,omitempty"`
	SMTPSenderAddressMatchesInstanceDomain *bool `json:"smtpSenderAddressMatchesInstanceDomain,omitempty"`
}

type DomainSettingsRepository interface {
	settingsRepository[DomainSettings]
	domainSettingsJSONChanges
}

type organizationSettingsJSONChanges interface {
	SetSettingFields(value OrganizationSettingsAttributes) database.Change

	SetOrganizationScopedUsernames(value bool) db_json.JsonUpdate
}

type OrganizationSettings struct {
	Settings
	OrganizationSettingsAttributes
}
type OrganizationSettingsAttributes struct {
	OrganizationScopedUsernames *bool `json:"organizationScopedUsernames,omitempty"`
}

type OrganizationSettingsRepository interface {
	settingsRepository[OrganizationSettings]
	organizationSettingsJSONChanges
}

type notificationSettingsJSONChanges interface {
	SetSettingFields(value NotificationSettingsAttributes) database.Change

	SetPasswordChange(value bool) db_json.JsonUpdate
}

type NotificationSettings struct {
	Settings
	NotificationSettingsAttributes
}

type NotificationSettingsAttributes struct {
	PasswordChange *bool `json:"passwordChange,omitempty"`
}

type NotificationSettingsRepository interface {
	settingsRepository[NotificationSettings]
	notificationSettingsJSONChanges
}

type legalAndSupportSettingsJSONChanges interface {
	SetSettingFields(value LegalAndSupportSettingsAttributes) database.Change

	SetTOSLink(value string) db_json.JsonUpdate
	SetPrivacyPolicyLink(value string) db_json.JsonUpdate
	SetHelpLink(value string) db_json.JsonUpdate
	SetSupportEmail(value string) db_json.JsonUpdate
	SetDocsLink(value string) db_json.JsonUpdate
	SetCustomLink(value string) db_json.JsonUpdate
	SetCustomLinkText(value string) db_json.JsonUpdate
}

type LegalAndSupportSettings struct {
	Settings
	LegalAndSupportSettingsAttributes
}

type LegalAndSupportSettingsAttributes struct {
	TOSLink           *string `json:"tosLink,omitempty"`
	PrivacyPolicyLink *string `json:"privacyPolicyLink,omitempty"`
	HelpLink          *string `json:"helpLink,omitempty"`
	SupportEmail      *string `json:"supportEmail,omitempty"`
	DocsLink          *string `json:"docsLink,omitempty"`
	CustomLink        *string `json:"customLink,omitempty"`
	CustomLinkText    *string `json:"customLinkText,omitempty"`
}

type LegalAndSupportSettingsRepository interface {
	settingsRepository[LegalAndSupportSettings]
	legalAndSupportSettingsJSONChanges
}
