package domain

import (
	"context"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
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
	SettingTypeSecretGenerator
)

//go:generate enumer -type SettingState -transform snake -trimprefix SettingState -sql
type SettingState int32

const (
	SettingStateActive SettingState = iota
	SettingStatePreview
)

type SettingTypes interface {
	Type() SettingType
	Base() Setting

	LoginSetting | BrandingSetting | PasswordComplexitySetting |
		PasswordExpirySetting | DomainSetting | LockoutSetting |
		SecuritySetting | OrganizationSetting | NotificationSetting |
		LegalAndSupportSetting | SecretGeneratorSetting
}

type Setting struct {
	ID             string    `json:"id,omitempty" db:"id"`
	InstanceID     string    `json:"instanceId,omitempty" db:"instance_id"`
	OrganizationID *string   `json:"organizationId,omitempty" db:"organization_id"`
	CreatedAt      time.Time `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt,omitzero" db:"updated_at"`
}

type settingColumns interface {
	InstanceIDColumn() database.Column
	OrganizationIDColumn() database.Column
	TypeColumn() database.Column
	StateColumn() database.Column
	CreatedAtColumn() database.Column
	UpdatedAtColumn() database.Column
	PrimaryKeyColumns() []database.Column
}

type settingConditions interface {
	InstanceIDCondition(id string) database.Condition
	OrganizationIDCondition(id *string) database.Condition
	UniqueCondition(instanceID string, orgID *string, typ SettingType, state SettingState) database.Condition
}

type settingChanges[T SettingTypes] interface {
	SetUpdatedAt(updatedAt time.Time) database.Change
	// TODO(adlerhurst): use this function to replace [SettingsRepository.Set]
	// Overwrite is used to overwrite the setting
	// If the setting does not exist it is created or updated with the given value if present.
	// Overwrite(value *T) database.Changes
}

type SettingsRepository[T SettingTypes] interface {
	settingColumns
	settingConditions
	settingChanges[T]

	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*T, error)
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*T, error)

	// Sets the setting with the given value, if the setting does not exist it is created.
	Set(ctx context.Context, client database.QueryExecutor, setting *T) error
	// Ensure is used to ensure the setting with the given changes.
	// If the setting already exists it is updated with the given changes, otherwise it is created with the given changes.
	Ensure(ctx context.Context, client database.QueryExecutor, instanceID string, organizationID *string, changes ...database.Change) error
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
	SecondFactorTypeRecoveryCodes
)

type loginSettingChanges interface {
	SetAllowUsernamePassword(value bool) database.Change
	SetAllowRegister(value bool) database.Change
	SetAllowExternalIDP(value bool) database.Change
	SetForceMultiFactor(value bool) database.Change
	SetForceMultiFactorLocalOnly(value bool) database.Change
	SetHidePasswordReset(value bool) database.Change
	SetIgnoreUnknownUsernames(value bool) database.Change
	SetAllowDomainDiscovery(value bool) database.Change
	SetDisableLoginWithEmail(value bool) database.Change
	SetDisableLoginWithPhone(value bool) database.Change
	SetPasswordlessType(value PasswordlessType) database.Change
	SetDefaultRedirectURI(value string) database.Change
	SetPasswordCheckLifetime(value time.Duration) database.Change
	SetExternalLoginCheckLifetime(value time.Duration) database.Change
	SetMultiFactorInitSkipLifetime(value time.Duration) database.Change
	SetSecondFactorCheckLifetime(value time.Duration) database.Change
	SetMultiFactorCheckLifetime(value time.Duration) database.Change

	AddMultiFactorType(value MultiFactorType) database.Change
	RemoveMultiFactorType(value MultiFactorType) database.Change
	AddSecondFactorType(value SecondFactorType) database.Change
	RemoveSecondFactorType(value SecondFactorType) database.Change
}

type LoginSetting struct {
	Setting `json:"-"`

	AllowUsernamePassword       bool               `json:"allowUsernamePassword"`
	AllowRegister               bool               `json:"allowRegister"`
	AllowExternalIDP            bool               `json:"allowExternalIdp"`
	ForceMultiFactor            bool               `json:"forceMultiFactor"`
	ForceMultiFactorLocalOnly   bool               `json:"forceMultiFactorLocalOnly"`
	HidePasswordReset           bool               `json:"hidePasswordReset"`
	IgnoreUnknownUsernames      bool               `json:"ignoreUnknownUsernames"`
	AllowDomainDiscovery        bool               `json:"allowDomainDiscovery"`
	DisableLoginWithEmail       bool               `json:"disableLoginWithEmail"`
	DisableLoginWithPhone       bool               `json:"disableLoginWithPhone"`
	PasswordlessType            PasswordlessType   `json:"passwordlessType"`
	DefaultRedirectURI          string             `json:"defaultRedirectUri"`
	PasswordCheckLifetime       time.Duration      `json:"passwordCheckLifetime"`
	ExternalLoginCheckLifetime  time.Duration      `json:"externalLoginCheckLifetime"`
	MultiFactorInitSkipLifetime time.Duration      `json:"multiFactorInitSkipLifetime"`
	SecondFactorCheckLifetime   time.Duration      `json:"secondFactorCheckLifetime"`
	MultiFactorCheckLifetime    time.Duration      `json:"multiFactorCheckLifetime"`
	MultiFactorTypes            []MultiFactorType  `json:"multiFactorTypes,omitempty"`
	SecondFactorTypes           []SecondFactorType `json:"secondFactorTypes,omitempty"`
}

func (LoginSetting) Type() SettingType {
	return SettingTypeLogin
}

func (s LoginSetting) Base() Setting {
	return s.Setting
}

type LoginSettingsRepository interface {
	SettingsRepository[LoginSetting]
	loginSettingChanges
}

type BrandingPolicyThemeMode int32

const (
	BrandingPolicyThemeAuto BrandingPolicyThemeMode = iota
	BrandingPolicyThemeLight
	BrandingPolicyThemeDark
)

type brandingSettingChanges interface {
	SetPrimaryColorLight(value string) database.Change
	SetBackgroundColorLight(value string) database.Change
	SetWarnColorLight(value string) database.Change
	SetFontColorLight(value string) database.Change
	SetLogoURLLight(value *url.URL) database.Change
	SetIconURLLight(value *url.URL) database.Change

	SetPrimaryColorDark(value string) database.Change
	SetBackgroundColorDark(value string) database.Change
	SetWarnColorDark(value string) database.Change
	SetFontColorDark(value string) database.Change
	SetLogoURLDark(value *url.URL) database.Change
	SetIconURLDark(value *url.URL) database.Change

	SetHideLoginNameSuffix(value bool) database.Change
	SetErrorMessagePopup(value bool) database.Change
	SetDisableWatermark(value bool) database.Change
	SetThemeMode(value BrandingPolicyThemeMode) database.Change
	SetFontURL(value *url.URL) database.Change

	// TODO(adlerhurst): move Activate() database.Change to changes instead of having separate operations for it
	// TODO(adlerhurst): move ActivateAt(t time.Time) database.Change to changes instead of having separate operations for it
}

type brandingSettingConditions interface {
	StateCondition(state SettingState) database.Condition
}

type BrandingSetting struct {
	Setting `json:"-"`

	State SettingState `json:"-"`

	PrimaryColorLight    string   `json:"primaryColorLight,omitempty"`
	BackgroundColorLight string   `json:"backgroundColorLight,omitempty"`
	WarnColorLight       string   `json:"warnColorLight,omitempty"`
	FontColorLight       string   `json:"fontColorLight,omitempty"`
	LogoURLLight         *url.URL `json:"logoUrlLight,omitzero"`
	IconURLLight         *url.URL `json:"iconUrlLight,omitzero"`

	PrimaryColorDark    string   `json:"primaryColorDark,omitempty"`
	BackgroundColorDark string   `json:"backgroundColorDark,omitempty"`
	WarnColorDark       string   `json:"warnColorDark,omitempty"`
	FontColorDark       string   `json:"fontColorDark,omitempty"`
	LogoURLDark         *url.URL `json:"logoUrlDark,omitzero"`
	IconURLDark         *url.URL `json:"iconUrlDark,omitzero"`

	HideLoginNameSuffix bool                    `json:"hideLoginNameSuffix,omitempty"`
	ErrorMessagePopup   bool                    `json:"errorMessagePopup,omitempty"`
	DisableWatermark    bool                    `json:"disableWatermark,omitempty"`
	ThemeMode           BrandingPolicyThemeMode `json:"themeMode,omitempty"`
	FontURL             *url.URL                `json:"fontUrl,omitzero"`
}

func (BrandingSetting) Type() SettingType {
	return SettingTypeBranding
}

func (s BrandingSetting) Base() Setting {
	return s.Setting
}

type BrandingSettingsRepository interface {
	SettingsRepository[BrandingSetting]
	brandingSettingChanges
	brandingSettingConditions

	Activate(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
	ActivateAt(ctx context.Context, client database.QueryExecutor, condition database.Condition, t time.Time) (int64, error)
	// TODO(adlerhurst): EnsureWithState(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string, state SettingState, changes ...database.Change) error
}

type passwordComplexitySettingChanges interface {
	SetMinLength(value uint64) database.Change
	SetHasLowercase(value bool) database.Change
	SetHasUppercase(value bool) database.Change
	SetHasNumber(value bool) database.Change
	SetHasSymbol(value bool) database.Change
}

type PasswordComplexitySetting struct {
	Setting `json:"-"`

	MinLength    uint64 `json:"minLength,omitempty"`
	HasLowercase bool   `json:"hasLowercase,omitempty"`
	HasUppercase bool   `json:"hasUppercase,omitempty"`
	HasNumber    bool   `json:"hasNumber,omitempty"`
	HasSymbol    bool   `json:"hasSymbol,omitempty"`
}

func (PasswordComplexitySetting) Type() SettingType {
	return SettingTypePasswordComplexity
}

func (s PasswordComplexitySetting) Base() Setting {
	return s.Setting
}

type PasswordComplexitySettingsRepository interface {
	SettingsRepository[PasswordComplexitySetting]
	passwordComplexitySettingChanges
}

type passwordExpirySettingChanges interface {
	SetExpireWarnDays(value uint64) database.Change
	SetMaxAgeDays(value uint64) database.Change
}

type PasswordExpirySetting struct {
	Setting `json:"-"`

	ExpireWarnDays uint64 `json:"expireWarnDays,omitempty"`
	MaxAgeDays     uint64 `json:"maxAgeDays,omitempty"`
}

func (PasswordExpirySetting) Type() SettingType {
	return SettingTypePasswordExpiry
}

func (s PasswordExpirySetting) Base() Setting {
	return s.Setting
}

type PasswordExpirySettingsRepository interface {
	SettingsRepository[PasswordExpirySetting]
	passwordExpirySettingChanges
}

type lockoutSettingChanges interface {
	SetMaxPasswordAttempts(value uint64) database.Change
	SetMaxOTPAttempts(value uint64) database.Change
	SetShowLockOutFailures(value bool) database.Change
}

type LockoutSetting struct {
	Setting `json:"-"`

	MaxPasswordAttempts uint64 `json:"maxPasswordAttempts,omitempty"`
	MaxOTPAttempts      uint64 `json:"maxOtpAttempts,omitempty"`
	ShowLockOutFailures bool   `json:"showLockOutFailures,omitempty"`
}

func (LockoutSetting) Type() SettingType {
	return SettingTypeLockout
}

func (s LockoutSetting) Base() Setting {
	return s.Setting
}

type LockoutSettingsRepository interface {
	SettingsRepository[LockoutSetting]
	lockoutSettingChanges
}

type securitySettingChanges interface {
	SetEnableIframeEmbedding(value bool) database.Change
	SetAllowedOrigins(values []string) database.Change
	SetEnableImpersonation(value bool) database.Change
}

type SecuritySetting struct {
	Setting `json:"-"`

	EnableIframeEmbedding bool     `json:"enableIframeEmbedding,omitempty"`
	AllowedOrigins        []string `json:"allowedOrigins,omitempty"`
	EnableImpersonation   bool     `json:"enableImpersonation,omitempty"`
}

func (SecuritySetting) Type() SettingType {
	return SettingTypeSecurity
}

func (s SecuritySetting) Base() Setting {
	return s.Setting
}

type SecuritySettingsRepository interface {
	SettingsRepository[SecuritySetting]
	securitySettingChanges
}

type domainSettingChanges interface {
	SetLoginNameIncludesDomain(value bool) database.Change
	SetRequireOrgDomainVerification(value bool) database.Change
	SetSMTPSenderAddressMatchesInstanceDomain(value bool) database.Change
}

type DomainSetting struct {
	Setting `json:"-"`

	LoginNameIncludesDomain                bool `json:"loginNameIncludesDomain,omitempty"`
	RequireOrgDomainVerification           bool `json:"requireOrgDomainVerification,omitempty"`
	SMTPSenderAddressMatchesInstanceDomain bool `json:"smtpSenderAddressMatchesInstanceDomain,omitempty"`
}

func (DomainSetting) Type() SettingType {
	return SettingTypeDomain
}

func (s DomainSetting) Base() Setting {
	return s.Setting
}

type DomainSettingsRepository interface {
	SettingsRepository[DomainSetting]
	domainSettingChanges
}

type organizationSettingChanges interface {
	SetOrganizationScopedUsernames(value bool) database.Change
}

type OrganizationSetting struct {
	Setting `json:"-"`

	OrganizationScopedUsernames bool `json:"organizationScopedUsernames,omitempty"`
}

func (OrganizationSetting) Type() SettingType {
	return SettingTypeOrganization
}

func (s OrganizationSetting) Base() Setting {
	return s.Setting
}

type OrganizationSettingsRepository interface {
	SettingsRepository[OrganizationSetting]
	organizationSettingChanges
}

type notificationSettingChanges interface {
	SetPasswordChange(value bool) database.Change
}

type NotificationSetting struct {
	Setting `json:"-"`

	PasswordChange bool `json:"passwordChange,omitempty"`
}

func (NotificationSetting) Type() SettingType {
	return SettingTypeNotification
}

func (s NotificationSetting) Base() Setting {
	return s.Setting
}

type NotificationSettingsRepository interface {
	SettingsRepository[NotificationSetting]
	notificationSettingChanges
}

type legalAndSupportSettingChanges interface {
	SetTOSLink(value string) database.Change
	SetPrivacyPolicyLink(value string) database.Change
	SetHelpLink(value string) database.Change
	SetSupportEmail(value string) database.Change
	SetDocsLink(value string) database.Change
	SetCustomLink(value string) database.Change
	SetCustomLinkText(value string) database.Change
}

type LegalAndSupportSetting struct {
	Setting `json:"-"`

	TOSLink           string `json:"tosLink,omitempty"`
	PrivacyPolicyLink string `json:"privacyPolicyLink,omitempty"`
	HelpLink          string `json:"helpLink,omitempty"`
	SupportEmail      string `json:"supportEmail,omitempty"`
	DocsLink          string `json:"docsLink,omitempty"`
	CustomLink        string `json:"customLink,omitempty"`
	CustomLinkText    string `json:"customLinkText,omitempty"`
}

func (LegalAndSupportSetting) Type() SettingType {
	return SettingTypeLegalAndSupport
}

func (s LegalAndSupportSetting) Base() Setting {
	return s.Setting
}

type LegalAndSupportSettingsRepository interface {
	SettingsRepository[LegalAndSupportSetting]
	legalAndSupportSettingChanges
}

type SecretGeneratorSetting struct {
	Setting `json:"-"`

	ClientSecret             *ClientSecretAttributes             `json:"clientSecret,omitempty"`
	InitializeUserCode       *InitializeUserCodeAttributes       `json:"initializeUserCode,omitempty"`
	EmailVerificationCode    *EmailVerificationCodeAttributes    `json:"emailVerificationCode,omitempty"`
	PhoneVerificationCode    *PhoneVerificationCodeAttributes    `json:"phoneVerificationCode,omitempty"`
	PasswordVerificationCode *PasswordVerificationCodeAttributes `json:"passwordVerificationCode,omitempty"`
	PasswordlessInitCode     *PasswordlessInitCodeAttributes     `json:"passwordlessInitCode,omitempty"`
	DomainVerification       *DomainVerificationAttributes       `json:"domainVerification,omitempty"`
	OTPSMS                   *OTPSMSAttributes                   `json:"otpSms,omitempty"`
	OTPEmail                 *OTPEmailAttributes                 `json:"otpEmail,omitempty"`
	AppSecret                *AppSecretAttributes                `json:"appSecret,omitempty"`
	InviteCode               *InviteCodeAttributes               `json:"inviteCode,omitempty"`
	SigningKey               *SigningKeyAttributes               `json:"signingKey,omitempty"`
}

func (SecretGeneratorSetting) Type() SettingType {
	return SettingTypeSecretGenerator
}

func (s SecretGeneratorSetting) Base() Setting {
	return s.Setting
}

//go:generate mockgen -typed -package domainmock -destination ./mock/secret_generator_setting.mock.go . SecretGeneratorSettingsRepository
type SecretGeneratorSettingsRepository interface {
	SettingsRepository[SecretGeneratorSetting]
	secretGeneratorSettingChanges
}

type secretGeneratorSettingChanges interface {
	SetClientSecretSecretGenerator(changes ...database.Change) database.Change
	SetInitializeUserCodeSecretGenerator(changes ...database.Change) database.Change
	SetEmailVerificationCodeSecretGenerator(changes ...database.Change) database.Change
	SetPhoneVerificationCodeSecretGenerator(changes ...database.Change) database.Change
	SetPasswordVerificationCodeSecretGenerator(changes ...database.Change) database.Change
	SetPasswordlessInitCodeSecretGenerator(changes ...database.Change) database.Change
	SetDomainVerificationSecretGenerator(changes ...database.Change) database.Change
	SetOTPSMSSecretGenerator(changes ...database.Change) database.Change
	SetOTPEmailSecretGenerator(changes ...database.Change) database.Change
	SetAppSecretSecretGenerator(changes ...database.Change) database.Change
	SetInviteCodeSecretGenerator(changes ...database.Change) database.Change
	SetSigningKeySecretGenerator(changes ...database.Change) database.Change

	SetLength(length uint) database.Change
	SetExpiry(expiry time.Duration) database.Change
	SetIncludeLowerLetters(includeLowerLetters bool) database.Change
	SetIncludeUpperLetters(includeUpperLetters bool) database.Change
	SetIncludeDigits(includeDigits bool) database.Change
	SetIncludeSymbols(includeSymbols bool) database.Change
}

//go:generate enumer -type=SecretGeneratorType -transform=snake -trimprefix=SecretGeneratorType
type SecretGeneratorType uint8

const (
	SecretGeneratorTypeUnspecified SecretGeneratorType = iota
	SecretGeneratorTypeClientSecret
	SecretGeneratorTypeInitializeUserCode
	SecretGeneratorTypeEmailVerificationCode
	SecretGeneratorTypePhoneVerificationCode
	SecretGeneratorTypePasswordVerificationCode
	SecretGeneratorTypePasswordlessInitCode
	SecretGeneratorTypeDomainVerification
	SecretGeneratorTypeOTPSMS
	SecretGeneratorTypeOTPEmail
	SecretGeneratorTypeAppSecret
	SecretGeneratorTypeInviteCode
	SecretGeneratorTypeSigningKey
)

type AppSecretAttributes struct {
	SecretGeneratorAttrsWithExpiry
}

type InviteCodeAttributes struct {
	SecretGeneratorAttrsWithExpiry
}

type SigningKeyAttributes struct {
	SecretGeneratorAttrsWithExpiry
}

type ClientSecretAttributes struct {
	SecretGeneratorAttrsWithExpiry
}

type InitializeUserCodeAttributes struct {
	SecretGeneratorAttrsWithExpiry
}

type EmailVerificationCodeAttributes struct {
	SecretGeneratorAttrsWithExpiry
}

type PhoneVerificationCodeAttributes struct {
	SecretGeneratorAttrsWithExpiry
}

type PasswordVerificationCodeAttributes struct {
	SecretGeneratorAttrsWithExpiry
}

type PasswordlessInitCodeAttributes struct {
	SecretGeneratorAttrsWithExpiry
}

type DomainVerificationAttributes struct {
	SecretGeneratorAttrs
}

type OTPSMSAttributes struct {
	SecretGeneratorAttrsWithExpiry
}

type OTPEmailAttributes struct {
	SecretGeneratorAttrsWithExpiry
}

type SecretGeneratorAttrs struct {
	Length              uint `json:"length,omitempty"`
	IncludeLowerLetters bool `json:"includeLowerLetters,omitempty"`
	IncludeUpperLetters bool `json:"includeUpperLetters,omitempty"`
	IncludeDigits       bool `json:"includeDigits,omitempty"`
	IncludeSymbols      bool `json:"includeSymbols,omitempty"`
}

type SecretGeneratorAttrsWithExpiry struct {
	Expiry *time.Duration `json:"expiry,omitempty"`
	SecretGeneratorAttrs
}
