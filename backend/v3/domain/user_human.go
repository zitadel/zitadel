package domain

import (
	"context"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Gender uint8

const (
	GenderUnspecified Gender = iota
	GenderFemale
	GenderMale
	GenderDiverse
)

type Human struct {
	FirstName         string        `json:"firstName,omitempty" db:"first_name"`
	LastName          string        `json:"lastName,omitempty" db:"last_name"`
	Nickname          string        `json:"nickName,omitempty" db:"nick_name"`
	DisplayName       string        `json:"displayName,omitempty" db:"display_name"`
	PreferredLanguage *language.Tag `json:"preferredLanguage,omitempty" db:"preferred_language"`
	Gender            *Gender       `json:"gender,omitempty" db:"gender"`
	AvatarKey         *string       `json:"avatarKey,omitempty" db:"avatar_key"`
	Avatar            []byte        `json:"avatar,omitempty" db:"avatar"`
}

//go:generate mockgen -typed -package domainmock -destination ./mock/user_human.mock.go . HumanUserRepository
type HumanUserRepository interface {
	humanColumns
	humanConditions
	humanChanges

	// Update updates a human user.
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)

	// SetPassword sets the password based on the given verification type.
	SetPassword(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification VerificationType) (int64, error)
	// GetPasswordVerification retrieves the password verification based on the given condition.
	GetPasswordVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*Verification, error)

	// SetEmail sets the email based on the given verification type.
	// * [VerificationTypeInit] to initialize email verification, previously verified email remains verified
	// * [VerificationTypeVerified] to mark email as verified, a verification must exist
	// * [VerificationTypeUpdate] to update email verification, a verification must exist (e.g. resend code)
	// * [VerificationTypeSkipped] to skip email verification, existing verification is removed (e.g. admin set email)
	SetEmail(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification VerificationType) (int64, error)
	// GetEmailVerification retrieves the email verification based on the given condition.
	GetEmailVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*Verification, error)

	// SetPhone sets the phone based on the given verification type.
	// * [VerificationTypeInit] to initialize phone verification, previously verified phone remains verified
	// * [VerificationTypeVerified] to mark phone as verified, a verification must exist
	// * [VerificationTypeUpdate] to update phone verification, a verification must exist (e.g. resend code)
	// * [VerificationTypeSkipped] to skip phone verification, existing verification is removed (e.g. admin set phone)
	SetPhone(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification VerificationType) (int64, error)
	// GetPhoneVerification retrieves the phone verification based on the given condition.
	GetPhoneVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*Verification, error)

	// SetTOTP sets the TOTP based on the given verification type.
	SetTOTP(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification VerificationType) (int64, error)
	// GetTOTPVerification retrieves the TOTP verification based on the given condition.
	GetTOTPVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*Verification, error)

	// SetEmailOTPCheck sets the email OTP check based on the given check type.
	// * [CheckTypeInit] to initialize email OTP check, it overwrites the existing check
	// * [CheckTypeSuccessful] to mark email OTP check as successful, updating last successful check time
	// * [CheckTypeFailed] to mark email OTP check as failed, increasing failed attempts count
	SetEmailOTPCheck(ctx context.Context, client database.QueryExecutor, condition database.Condition, check CheckType) (int64, error)
	// GetEmailOTPCheck retrieves the email OTP check based on the given condition.
	GetEmailOTPCheck(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*Check, error)
	// SetSMSOTPCheck sets the SMS OTP check based on the given check type.
	// * [CheckTypeInit] to initialize SMS OTP check, it overwrites the existing check
	// * [CheckTypeSuccessful] to mark SMS OTP check as successful, updating last successful check time
	// * [CheckTypeFailed] to mark SMS OTP check as failed, increasing failed attempts count
	SetSMSOTPCheck(ctx context.Context, client database.QueryExecutor, condition database.Condition, check CheckType) (int64, error)
	// GetSMSOTPCheck retrieves the SMS OTP check based on the given condition.
	GetSMSOTPCheck(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*Check, error)
	// SetTOTPCheck sets the TOTP check based on the given check type.
	// * [CheckTypeInit] is not allowed for TOTP check. because the secret must already be set.
	// * [CheckTypeSuccessful] to mark TOTP check as successful, updating last successful check time
	// * [CheckTypeFailed] to mark TOTP check as failed, increasing failed attempts count
	SetTOTPCheck(ctx context.Context, client database.QueryExecutor, condition database.Condition, check CheckType) (int64, error)
	// GetTOTPCheck retrieves the TOTP check based on the given condition.
	GetTOTPCheck(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*Check, error)

	// LoadIdentityProviderLinks enables fetching of identity provider links when getting or listing human users
	LoadIdentityProviderLinks() HumanUserRepository
}

type humanColumns interface {
	userColumns
	FirstNameColumn() database.Column
	LastNameColumn() database.Column
	DisplayNameColumn() database.Column
	NicknameColumn() database.Column
	PreferredLanguageColumn() database.Column
	GenderColumn() database.Column
	AvatarKeyColumn() database.Column

	PasswordColumn() database.Column
	PasswordVerifiedAtColumn() database.Column
	FailedPasswordAttemptsColumn() database.Column

	EmailColumn() database.Column
	PhoneColumn() database.Column
}

type humanConditions interface {
	userConditions
	FirstNameCondition(op database.TextOperation, name string) database.Condition
	LastNameCondition(op database.TextOperation, name string) database.Condition
	NicknameCondition(op database.TextOperation, name string) database.Condition
	DisplayNameCondition(op database.TextOperation, name string) database.Condition

	EmailCondition(op database.TextOperation, email string) database.Condition
	PhoneCondition(op database.TextOperation, phone string) database.Condition
}

type humanChanges interface {
	userChanges
	SetFirstName(name string) database.Change
	SetLastName(name string) database.Change
	SetNickname(name string) database.Change
	SetDisplayName(name string) database.Change
	// SetPreferredLanguage sets the preferred language,
	// nil and [language.Und] are treated as unset
	SetPreferredLanguage(language *language.Tag) database.Change
	// SetGender sets the gender,
	// nil and [GenderUnspecified] are treated as unset
	SetGender(gender *Gender) database.Change
	// SetAvatarKey sets the avatar key,
	// nil removes the avatar key
	SetAvatarKey(key *string) database.Change

	SetPasswordChangeRequired(required bool) database.Change

	// SetMFAInitSkippedAt sets the time when MFA initialization was skipped,
	SetMFAInitSkippedAt(skippedAt *time.Time) database.Change
	// SetEmailOTPEnabledAt sets the email OTP enabled at time,
	// If [time.Time.IsZero] is treated as NOW()
	SetEmailOTPEnabledAt(enabledAt time.Time) database.Change
	// SetSMSOTPEnabledAt sets the SMS OTP enabled at time,
	// If [time.Time.IsZero] is treated as NOW()
	SetSMSOTPEnabledAt(enabledAt time.Time) database.Change

	// RemoveTOTP removes the TOTP settings.
	RemoveTOTP() database.Change
	// RemovePhone removes the phone number and its verification.
	RemovePhone() database.Change
}
