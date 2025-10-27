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

type SetPasswordVerification interface{}

type SetPasswordVerificationCurrentPassword struct {
	CurrentPassword []byte
}

type SetPasswordVerificationVerificationCode struct {
	VerificationCode string
}

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
}

type humanConditions interface {
	userConditions
	FirstNameCondition(op database.TextOperation, name string) database.Condition
	LastNameCondition(op database.TextOperation, name string) database.Condition
	NicknameCondition(op database.TextOperation, name string) database.Condition
	DisplayNameCondition(op database.TextOperation, name string) database.Condition
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
	// SetPassword sets the password hash, if verifiedAt is Zero, NOW() is used
	SetPassword(password []byte, verifiedAt time.Time) database.Change
}

type HumanUserRepository interface {
	humanColumns
	humanConditions
	humanChanges

	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)

	// Security() HumanSecurityRepository
}

//go:generate enumer -type ContactType -transform lower -trimprefix ContactType -sql
// type ContactType uint8

// const (
// 	ContactTypeUnspecified ContactType = iota
// 	ContactTypeEmail
// 	ContactTypePhone
// )

// // human contact type
// type HumanContact struct {
// 	// InstanceID      string      `json:"instanceId,omitempty" db:"instance_id"`
// 	// OrgID           string      `json:"orgId,omitempty" db:"org_id"`
// 	// UserId          string      `json:"userId,omitempty" db:"user_id"`
// 	Type            *ContactType `json:"type,omitempty" db:"type"`
// 	Value           *string      `json:"value,omitempty" db:"value"`
// 	IsVerified      *bool        `json:"isVerified,omitempty" db:"is_verified"`
// 	UnverifiedValue *string      `json:"unverifiedValue,omitempty" db:"unverified_value"`
// }

// // human security
// type HumanSecurity struct {
// 	// InstanceID string `json:"instanceId,omitempty" db:"instance_id"`
// 	// OrgID      string `json:"orgId,omitempty" db:"org_id"`
// 	// UserId     string `json:"userId,omitempty" db:"user_id"`

// 	PasswordChangeRequired bool       `json:"passwordChangeRequired,omitempty" db:"password_change_required"`
// 	PasswordChange         *time.Time `json:"passwordChange,omitempty" db:"password_change"`
// 	MFAInitSkipped         bool       `json:"mfaInitSkipped,omitempty" db:"mfa_init_skipped"`
// }

// type HumanSecurityRepository interface {
// 	humanSecurityColumns
// 	humanSecurityConditions
// 	humanSecurityChanges

// 	Create(ctx context.Context, client database.QueryExecutor, user *Human) error
// 	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*HumanSecurity, error)
// 	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
// }

// type humanSecurityColumns interface {
// 	InstanceIDColumn() database.Column
// 	OrgIDColumn() database.Column
// 	UserIDColumn() database.Column
// 	PasswordChangeRequiredColumn() database.Column
// 	PasswordChangedColumn() database.Column
// 	MFAInitSkippedColumn() database.Column
// }

// type humanSecurityConditions interface {
// 	InstanceIDCondition(instanceID string) database.Condition
// 	OrgIDCondition(orgID string) database.Condition
// 	UserIDCondition(userID string) database.Condition
// 	PasswordChangeRequiredCondition(required bool) database.Condition
// 	PasswordChangeCondition(op database.NumberOperation, time time.Time) database.Condition
// 	MFAInitSkippedCondition(skipped bool) database.Condition
// }

// type humanSecurityChanges interface {
// 	SetPasswordChangeRequired(required bool) database.Change
// 	SetPasswordChanged(time time.Time) database.Change
// 	SetMFAInitSkipped(skipped bool) database.Change
// }

// type humanContactColumns interface {
// 	InstanceIDColumn() database.Column
// 	OrgIDColumn() database.Column
// 	UserIDColumn() database.Column
// 	TypeCondition() database.Column
// 	CurrentValueColumn() database.Column
// 	VerifiedColumn() database.Column
// 	UnverifiedValueColumn() database.Column
// }

// type humanContactConditions interface {
// 	InstanceIDCondition(instanceID string) database.Condition
// 	OrgIDCondition(orgID string) database.Condition
// 	UserIDCondition(userID string) database.Condition
// 	TypeCondition(typ ContactType) database.Condition
// 	CurrentValueCondition(value string) database.Condition
// 	VerifiedCondition(verified bool) database.Condition
// 	UnverifiedValueCondition(value string) database.Condition
// }

// type humanContactChanges interface {
// 	SetInstanceID(instanceID string) database.Change
// 	SetOrgID(orgID string) database.Change
// 	SetUserID(userID string) database.Change
// 	SetType(typ ContactType) database.Change
// 	SetCurrentValue(value string) database.Change
// 	SetVerified(verified bool) database.Change
// 	SetUnverifiedValue(value string) database.Change
// }
