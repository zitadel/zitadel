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

type SetPasswordVerification interface {
	isSetPasswordVerification()
}

type SetPasswordVerificationCurrentPasswordChecked struct {
	VerifiedAt time.Time
}

func (s SetPasswordVerificationCurrentPasswordChecked) isSetPasswordVerification() {}

type SetPasswordVerificationVerificationCode struct {
	VerifiedAt time.Time
}

func (s SetPasswordVerificationVerificationCode) isSetPasswordVerification() {}

// TODO(adlerhurst): is there a code present in that case?
type SetPasswordVerificationChangeRequired struct {
	VerifiedAt time.Time
}

func (s SetPasswordVerificationChangeRequired) isSetPasswordVerification() {}

type Human struct {
	// HumanEmailContact HumanContact  `db:"email"`
	// HumanPhoneContact *HumanContact `db:"phone"`

	// HumanSecurity

	FirstName         string        `json:"firstName,omitempty" db:"first_name"`
	LastName          string        `json:"lastName,omitempty" db:"last_name"`
	Nickname          string        `json:"nickName,omitempty" db:"nick_name"`
	DisplayName       string        `json:"displayName,omitempty" db:"display_name"`
	PreferredLanguage *language.Tag `json:"preferredLanguage,omitempty" db:"preferred_language"`
	Gender            *Gender       `json:"gender,omitempty" db:"gender"`
	AvatarKey         *string       `json:"avatarKey,omitempty" db:"avatar_key"`
	Avatar            []byte        `json:"avatar,omitempty" db:"avatar"`
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
	// nil and [language.Und] are treated as unset
	SetPreferredLanguage(language *language.Tag) database.Change
	// nil and [GenderUnspecified] are treated as unset
	SetGender(gender *Gender) database.Change
	SetAvatarKey(key *string) database.Change

	SetPasswordChangeRequired(required bool) database.Change
	IncrementFailedPasswordAttempts() database.Change
	ResetFailedPasswordAttempts() database.Change

	RemovePhone() database.Change
}

type HumanUserRepository interface {
	humanColumns
	humanConditions
	humanChanges

	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)

	// SetPassword sets the password based on the given verification type.
	SetPassword(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification VerificationType) (int64, error)
	GetPasswordVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*Verification, error)
	IncrementPasswordVerificationAttempts(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)

	SetEmail(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification VerificationType) (int64, error)
	GetEmailVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*Verification, error)
	IncrementEmailVerificationAttempts(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)

	SetPhone(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification VerificationType) (int64, error)
	GetPhoneVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*Verification, error)
	IncrementPhoneVerificationAttempts(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}
