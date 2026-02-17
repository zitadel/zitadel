package repository

import (
	"context"
	"slices"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func SecretGeneratorSettings() domain.SecretGeneratorSettingsRepository {
	return new(secretGeneratorSetting)
}

type secretGeneratorSetting struct {
	setting[domain.SecretGeneratorSetting]
}

func (sgs secretGeneratorSetting) Ensure(ctx context.Context, client database.QueryExecutor, instanceID string, organizationID *string, changes ...database.Change) error {
	return sgs.ensure(ctx, client, instanceID, organizationID, domain.SettingStateActive, changes...)
}

// Get implements [domain.SecretGeneratorSettingsRepository].
func (sgs secretGeneratorSetting) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.SecretGeneratorSetting, error) {
	setting, err := sgs.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	setting.Attributes.Setting = setting.Setting
	return setting.Attributes, nil
}

// List implements [domain.SecretGeneratorSettingsRepository].
func (sgs secretGeneratorSetting) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.SecretGeneratorSetting, error) {
	settings, err := sgs.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.SecretGeneratorSetting, len(settings))
	for i, setting := range settings {
		result[i] = setting.Attributes
		result[i].Setting = setting.Setting
	}
	return result, nil
}

var (
	secretGeneratorSettingClientSecretPathPrefix             = []string{"clientSecret"}
	secretGeneratorSettingInitializeUserCodePathPrefix       = []string{"initializeUserCode"}
	secretGeneratorSettingEmailVerificationCodePathPrefix    = []string{"emailVerificationCode"}
	secretGeneratorSettingPhoneVerificationCodePathPrefix    = []string{"phoneVerificationCode"}
	secretGeneratorSettingPasswordVerificationCodePathPrefix = []string{"passwordVerificationCode"}
	secretGeneratorSettingPasswordlessInitCodePathPrefix     = []string{"passwordlessInitCode"}
	secretGeneratorSettingDomainVerificationPathPrefix       = []string{"domainVerification"}
	secretGeneratorSettingOTPSMSPathPrefix                   = []string{"otpSms"}
	secretGeneratorSettingOTPEmailPathPrefix                 = []string{"otpEmail"}
	secretGeneratorSettingPasswordResetCodePathPrefix        = []string{"passwordResetCode"}
	secretGeneratorSettingAppSecretPathPrefix                = []string{"appSecret"}
	secretGeneratorSettingInviteCodePathPrefix               = []string{"inviteCode"}
	secretGeneratorSettingSigningKeyPathPrefix               = []string{"signingKey"}

	secretGeneratorSettingExpiryPathSuffix              = []string{"expiry"}
	secretGeneratorSettingIncludeLowerLettersPathSuffix = []string{"includeLowerLetters"}
	secretGeneratorSettingIncludeUpperLettersPathSuffix = []string{"includeUpperLetters"}
	secretGeneratorSettingIncludeDigitsPathSuffix       = []string{"includeDigits"}
	secretGeneratorSettingIncludeSymbolsPathSuffix      = []string{"includeSymbols"}
	secretGeneratorSettingLengthPathSuffix              = []string{"length"}
)

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// Overwrite implements [domain.SecretGeneratorSettingsRepository].
func (sgs secretGeneratorSetting) Overwrite(value *domain.SecretGeneratorSetting) database.Changes {
	changes := sgs.setAttributes(&value.Setting)
	if value.ClientSecret != nil {
		changes = append(changes, sgs.SetClientSecretSecretGenerator(
			sgs.secretGeneratorAttributesToChanges(value.ClientSecret.SecretGeneratorAttrs, value.ClientSecret.Expiry)...),
		)
	}
	if value.InitializeUserCode != nil {
		changes = append(changes, sgs.SetInitializeUserCodeSecretGenerator(
			sgs.secretGeneratorAttributesToChanges(value.InitializeUserCode.SecretGeneratorAttrs, value.InitializeUserCode.Expiry)...),
		)
	}
	if value.EmailVerificationCode != nil {
		changes = append(changes, sgs.SetEmailVerificationCodeSecretGenerator(
			sgs.secretGeneratorAttributesToChanges(value.EmailVerificationCode.SecretGeneratorAttrs, value.EmailVerificationCode.Expiry)...),
		)
	}
	if value.PhoneVerificationCode != nil {
		changes = append(changes, sgs.SetPhoneVerificationCodeSecretGenerator(
			sgs.secretGeneratorAttributesToChanges(value.PhoneVerificationCode.SecretGeneratorAttrs, value.PhoneVerificationCode.Expiry)...),
		)
	}
	if value.PasswordlessInitCode != nil {
		changes = append(changes, sgs.SetPasswordlessInitCodeSecretGenerator(
			sgs.secretGeneratorAttributesToChanges(value.PasswordlessInitCode.SecretGeneratorAttrs, value.PasswordlessInitCode.Expiry)...),
		)
	}
	if value.DomainVerification != nil {
		changes = append(changes, sgs.SetDomainVerificationSecretGenerator(
			sgs.secretGeneratorAttributesToChanges(value.DomainVerification.SecretGeneratorAttrs, nil)...),
		)
	}
	if value.OTPSMS != nil {
		changes = append(changes, sgs.SetOTPSMSSecretGenerator(
			sgs.secretGeneratorAttributesToChanges(value.OTPSMS.SecretGeneratorAttrs, value.OTPSMS.Expiry)...),
		)
	}
	if value.OTPEmail != nil {
		changes = append(changes, sgs.SetOTPEmailSecretGenerator(
			sgs.secretGeneratorAttributesToChanges(value.OTPEmail.SecretGeneratorAttrs, value.OTPEmail.Expiry)...),
		)
	}
	return changes
}

func (sgs secretGeneratorSetting) secretGeneratorAttributesToChanges(attributes domain.SecretGeneratorAttrs, expiry *time.Duration) []database.Change {
	changes := make([]database.Change, 0, 6)
	changes = append(changes,
		sgs.SetLength(attributes.Length),
		sgs.SetIncludeLowerLetters(attributes.IncludeLowerLetters),
		sgs.SetIncludeUpperLetters(attributes.IncludeUpperLetters),
		sgs.SetIncludeDigits(attributes.IncludeDigits),
		sgs.SetIncludeSymbols(attributes.IncludeSymbols),
	)
	if expiry != nil {
		changes = slices.Grow(changes, 1)
		changes = append(changes, sgs.SetExpiry(*expiry))
	}
	return changes
}

func (sgs secretGeneratorSetting) SetClientSecretSecretGenerator(changes ...database.Change) database.Change {
	result := make(database.Changes, len(changes))
	for i, change := range changes {
		constructor, ok := change.(secretGeneratorSettingAttributesChange)
		if !ok {
			// TODO(adlerhurst): log error about wrong change type
			continue
		}
		result[i] = &attributeChange{
			change: constructor(secretGeneratorSettingClientSecretPathPrefix),
		}
	}
	return result
}

func (sgs secretGeneratorSetting) SetInitializeUserCodeSecretGenerator(changes ...database.Change) database.Change {
	result := make(database.Changes, len(changes))
	for i, change := range changes {
		constructor, ok := change.(secretGeneratorSettingAttributesChange)
		if !ok {
			// TODO(adlerhurst): log error about wrong change type
			continue
		}
		result[i] = &attributeChange{
			change: constructor(secretGeneratorSettingInitializeUserCodePathPrefix),
		}
	}
	return result
}

func (sgs secretGeneratorSetting) SetEmailVerificationCodeSecretGenerator(changes ...database.Change) database.Change {
	result := make(database.Changes, len(changes))
	for i, change := range changes {
		constructor, ok := change.(secretGeneratorSettingAttributesChange)
		if !ok {
			// TODO(adlerhurst): log error about wrong change type
			continue
		}
		result[i] = &attributeChange{
			change: constructor(secretGeneratorSettingEmailVerificationCodePathPrefix),
		}
	}
	return result
}

func (sgs secretGeneratorSetting) SetPhoneVerificationCodeSecretGenerator(changes ...database.Change) database.Change {
	result := make(database.Changes, len(changes))
	for i, change := range changes {
		constructor, ok := change.(secretGeneratorSettingAttributesChange)
		if !ok {
			// TODO(adlerhurst): log error about wrong change type
			continue
		}
		result[i] = &attributeChange{
			change: constructor(secretGeneratorSettingPhoneVerificationCodePathPrefix),
		}
	}
	return result
}

func (sgs secretGeneratorSetting) SetPasswordVerificationCodeSecretGenerator(changes ...database.Change) database.Change {
	result := make(database.Changes, len(changes))
	for i, change := range changes {
		constructor, ok := change.(secretGeneratorSettingAttributesChange)
		if !ok {
			// TODO(adlerhurst): log error about wrong change type
			continue
		}
		result[i] = &attributeChange{
			change: constructor(secretGeneratorSettingPasswordVerificationCodePathPrefix),
		}
	}
	return result
}

func (sgs secretGeneratorSetting) SetPasswordlessInitCodeSecretGenerator(changes ...database.Change) database.Change {
	result := make(database.Changes, len(changes))
	for i, change := range changes {
		constructor, ok := change.(secretGeneratorSettingAttributesChange)
		if !ok {
			// TODO(adlerhurst): log error about wrong change type
			continue
		}
		result[i] = &attributeChange{
			change: constructor(secretGeneratorSettingPasswordlessInitCodePathPrefix),
		}
	}
	return result
}

func (sgs secretGeneratorSetting) SetDomainVerificationSecretGenerator(changes ...database.Change) database.Change {
	result := make(database.Changes, len(changes))
	for i, change := range changes {
		constructor, ok := change.(secretGeneratorSettingAttributesChange)
		if !ok {
			// TODO(adlerhurst): log error about wrong change type
			continue
		}
		result[i] = &attributeChange{
			change: constructor(secretGeneratorSettingDomainVerificationPathPrefix),
		}
	}
	return result
}

func (sgs secretGeneratorSetting) SetOTPSMSSecretGenerator(changes ...database.Change) database.Change {
	result := make(database.Changes, len(changes))
	for i, change := range changes {
		constructor, ok := change.(secretGeneratorSettingAttributesChange)
		if !ok {
			// TODO(adlerhurst): log error about wrong change type
			continue
		}
		result[i] = &attributeChange{
			change: constructor(secretGeneratorSettingOTPSMSPathPrefix),
		}
	}
	return result
}

func (sgs secretGeneratorSetting) SetOTPEmailSecretGenerator(changes ...database.Change) database.Change {
	result := make(database.Changes, len(changes))
	for i, change := range changes {
		constructor, ok := change.(secretGeneratorSettingAttributesChange)
		if !ok {
			// TODO(adlerhurst): log error about wrong change type
			continue
		}
		result[i] = &attributeChange{
			change: constructor(secretGeneratorSettingOTPEmailPathPrefix),
		}
	}
	return result
}

func (sgs secretGeneratorSetting) SetPasswordResetCodeSecretGenerator(changes ...database.Change) database.Change {
	result := make(database.Changes, len(changes))
	for i, change := range changes {
		constructor, ok := change.(secretGeneratorSettingAttributesChange)
		if !ok {
			// TODO(adlerhurst): log error about wrong change type
			continue
		}
		result[i] = &attributeChange{
			change: constructor(secretGeneratorSettingPasswordResetCodePathPrefix),
		}
	}
	return result
}
func (sgs secretGeneratorSetting) SetAppSecretSecretGenerator(changes ...database.Change) database.Change {
	result := make(database.Changes, len(changes))
	for i, change := range changes {
		constructor, ok := change.(secretGeneratorSettingAttributesChange)
		if !ok {
			// TODO(adlerhurst): log error about wrong change type
			continue
		}
		result[i] = &attributeChange{
			change: constructor(secretGeneratorSettingAppSecretPathPrefix),
		}
	}
	return result
}
func (sgs secretGeneratorSetting) SetInviteCodeSecretGenerator(changes ...database.Change) database.Change {
	result := make(database.Changes, len(changes))
	for i, change := range changes {
		constructor, ok := change.(secretGeneratorSettingAttributesChange)
		if !ok {
			// TODO(adlerhurst): log error about wrong change type
			continue
		}
		result[i] = &attributeChange{
			change: constructor(secretGeneratorSettingInviteCodePathPrefix),
		}
	}
	return result
}
func (sgs secretGeneratorSetting) SetSigningKeySecretGenerator(changes ...database.Change) database.Change {
	result := make(database.Changes, len(changes))
	for i, change := range changes {
		constructor, ok := change.(secretGeneratorSettingAttributesChange)
		if !ok {
			// TODO(adlerhurst): log error about wrong change type
			continue
		}
		result[i] = &attributeChange{
			change: constructor(secretGeneratorSettingSigningKeyPathPrefix),
		}
	}
	return result
}

// SetExpiry implements [domain.SecretGeneratorSettingsRepository].
func (sgs secretGeneratorSetting) SetExpiry(expiry time.Duration) database.Change {
	return secretGeneratorSettingAttributesChange(func(path []string) func(column database.Column) database.Change {
		return func(column database.Column) database.Change {
			return database.SetJSONValue(column, append(path, secretGeneratorSettingExpiryPathSuffix...), expiry)
		}
	})
}

// SetIncludeDigits implements [domain.SecretGeneratorSettingsRepository].
func (sgs secretGeneratorSetting) SetIncludeDigits(includeDigits bool) database.Change {
	return secretGeneratorSettingAttributesChange(func(path []string) func(column database.Column) database.Change {
		return func(column database.Column) database.Change {
			return database.SetJSONValue(column, append(path, secretGeneratorSettingIncludeDigitsPathSuffix...), includeDigits)
		}
	})
}

// SetIncludeLowerLetters implements [domain.SecretGeneratorSettingsRepository].
func (sgs secretGeneratorSetting) SetIncludeLowerLetters(includeLowerLetters bool) database.Change {
	return secretGeneratorSettingAttributesChange(func(path []string) func(column database.Column) database.Change {
		return func(column database.Column) database.Change {
			return database.SetJSONValue(column, append(path, secretGeneratorSettingIncludeLowerLettersPathSuffix...), includeLowerLetters)
		}
	})
}

// SetIncludeSymbols implements [domain.SecretGeneratorSettingsRepository].
func (sgs secretGeneratorSetting) SetIncludeSymbols(includeSymbols bool) database.Change {
	return secretGeneratorSettingAttributesChange(func(path []string) func(column database.Column) database.Change {
		return func(column database.Column) database.Change {
			return database.SetJSONValue(column, append(path, secretGeneratorSettingIncludeSymbolsPathSuffix...), includeSymbols)
		}
	})
}

// SetIncludeUpperLetters implements [domain.SecretGeneratorSettingsRepository].
func (sgs secretGeneratorSetting) SetIncludeUpperLetters(includeUpperLetters bool) database.Change {
	return secretGeneratorSettingAttributesChange(func(path []string) func(column database.Column) database.Change {
		return func(column database.Column) database.Change {
			return database.SetJSONValue(column, append(path, secretGeneratorSettingIncludeUpperLettersPathSuffix...), includeUpperLetters)
		}
	})
}

// SetLength implements [domain.SecretGeneratorSettingsRepository].
func (sgs secretGeneratorSetting) SetLength(length uint) database.Change {
	return secretGeneratorSettingAttributesChange(func(path []string) func(column database.Column) database.Change {
		return func(column database.Column) database.Change {
			return database.SetJSONValue(column, append(path, secretGeneratorSettingLengthPathSuffix...), length)
		}
	})
}

type secretGeneratorSettingAttributesChange func(path []string) func(column database.Column) database.Change

// IsOnColumn implements [database.Change].
func (s secretGeneratorSettingAttributesChange) IsOnColumn(col database.Column) bool {
	return col.Equals(settingAttributesColumn)
}

// Matches implements [database.Change].
func (s secretGeneratorSettingAttributesChange) Matches(x any) bool {
	toMatch, ok := x.(secretGeneratorSettingAttributesChange)
	if !ok {
		return false
	}
	return s(nil)(settingAttributesColumn).Matches(toMatch(nil)(settingAttributesColumn))
}

// String implements [database.Change].
func (s secretGeneratorSettingAttributesChange) String() string {
	return "repository.secretGeneratorSettingAttributesChange"
}

// Write implements [database.Change].
// It must never be called without setting a type through the [domain.secretGeneratorSettingChanges].`Set...SecretGenerator` methods.
func (s secretGeneratorSettingAttributesChange) Write(builder *database.StatementBuilder) {
	panic("unimplemented")
}

// WriteArg implements [database.Change].
func (s secretGeneratorSettingAttributesChange) WriteArg(builder *database.StatementBuilder) {
	s.Write(builder)
}

var _ database.Change = (secretGeneratorSettingAttributesChange)(nil)

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

var _ domain.SecretGeneratorSettingsRepository = (*secretGeneratorSetting)(nil)
