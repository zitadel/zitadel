package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

//go:generate enumer -type UserState -transform lower -trimprefix UserState -sql
type UserState uint8

const (
	UserStateInital UserState = iota
	UserStateActive
	UserStateInactive
	UserStateLocked
	UserStateSuspended
)

// user
type User struct {
	InstanceID        string    `json:"instanceId,omitempty" db:"instance_id"`
	OrgID             string    `json:"orgId,omitempty" db:"org_id"`
	ID                string    `json:"id,omitempty" db:"id"`
	Username          string    `json:"username,omitempty" db:"username"`
	UsernameOrgUnique bool      `json:"usernameOrgUnique,omitempty" db:"username_org_unique"`
	State             UserState `json:"state,omitempty" db:"state"`

	CreatedAt time.Time `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt,omitzero" db:"updated_at"`
}

type userColumns interface {
	InstanceIDColumn() database.Column
	OrgIDColumn() database.Column
	IDColumn() database.Column
	UsernameColumn() database.Column
	UsernameOrgUniqueColumn() database.Column
	StateColumn() database.Column
	CreatedAtColumn() database.Column
	UpdatedAtColumn() database.Column
}

type userConditions interface {
	InstanceIDCondition(instanceID string) database.Condition
	OrgIDCondition(orgID string) database.Condition
	IDCondition(userID string) database.Condition
	UsernameCondition(op database.TextOperation, username string) database.Condition
	UsernameOrgUniqueCondition(condition bool) database.Condition
	StateCondition(state UserState) database.Condition
	CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition
	UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition
}

type userChanges interface {
	// SetInstanceID(instanceID string) database.Condition
	// SetOrgID(orgID string) database.Condition
	// SetID(userID string) database.Condition
	SetUsername(username string) database.Change
	SetUsernameOrgUnique(op bool) database.Change
	SetState(state UserState) database.Change
	// SetCreatedAt(op database.NumberOperation, createdAt time.Time) database.Condition
	SetUpdatedAt(updatedAt time.Time) database.Change
}

// machine user
type Machine struct {
	User
	Name        string  `json:"name,omitempty" db:"name"`
	Description *string `json:"description,omitempty" db:"description"`
}

type machineColumns interface {
	userColumns
	NameColumn() database.Column
	DescriptionColumn() database.Column
}

type machineConditions interface {
	userConditions
	NameCondition(op database.TextOperation, name string) database.Condition
	DescriptionCondition(op database.TextOperation, description string) database.Condition
}

type machineChanges interface {
	userChanges
	SetName(name string) database.Change
	SetDescription(description string) database.Change
}

type MachineRepository interface {
	machineColumns
	machineConditions
	machineChanges
}

// human user
type Human struct {
	User
	HumanEmailContact HumanContact  `db:"email"`
	HumanPhoneContact *HumanContact `db:"phone"`

	HumanSecurity

	FirstName         string  `json:"firstName,omitempty" db:"first_name"`
	LastName          string  `json:"lastName,omitempty" db:"last_name"`
	NickName          string  `json:"nickName,omitempty" db:"nick_name"`
	DisplayName       string  `json:"displayName,omitempty" db:"display_name"`
	PreferredLanguage string  `json:"preferredLanguage,omitempty" db:"preferred_language"`
	Gender            uint8   `json:"gender,omitempty" db:"gender"`
	AvatarKey         *string `json:"avataryKey,omitempty" db:"avatar_key"`
}

type humanColumns interface {
	userColumns
	FirstNameColumn() database.Column
	LastNameColumn() database.Column
	DisplayNameColumn() database.Column
	PreferredLanguageColumn() database.Column
	GenderColumn() database.Column
	AvatarKeyColumn() database.Column
}

type humanConditions interface {
	userConditions
	FirstNameCondition(op database.TextOperation, name string) database.Condition
	LastNameCondition(op database.TextOperation, name string) database.Condition
	NickNameCondition(op database.TextOperation, name string) database.Condition
	DisplayNameCondition(op database.TextOperation, name string) database.Condition
	PreferredLanguageCondition(language string) database.Condition
	GenderCondition(gender uint8) database.Condition
	AvatarKeyCondition(key string) database.Condition
}

type humanChanges interface {
	userChanges
	SetFirstName(name string) database.Change
	SetLastName(name string) database.Change
	SetNickName(name string) database.Change
	SetDisplayName(name string) database.Change
	SetPreferredLanguage(language string) database.Change
	SetGender(gender uint8) database.Change
	SetAvatarKey(key string) database.Change
}

type HumanRepository interface {
	humanColumns
	humanConditions
	humanChanges
}

//go:generate enumer -type ContactType -transform lower -trimprefix ContactType -sql
type ContactType uint8

const (
	ContactTypeUnspecified ContactType = iota
	ContactTypeEmail
	ContactTypePhone
)

// human contact type
type HumanContact struct {
	// InstanceID      string      `json:"instanceId,omitempty" db:"instance_id"`
	// OrgID           string      `json:"orgId,omitempty" db:"org_id"`
	// UserId          string      `json:"userId,omitempty" db:"user_id"`
	Type            *ContactType `json:"type,omitempty" db:"type"`
	Value           *string      `json:"value,omitempty" db:"value"`
	IsVerified      *bool        `json:"isVerified,omitempty" db:"is_verified"`
	UnverifiedValue *string      `json:"unverifiedValue,omitempty" db:"unverified_value"`
}

// human security
type HumanSecurity struct {
	// InstanceID string `json:"instanceId,omitempty" db:"instance_id"`
	// OrgID      string `json:"orgId,omitempty" db:"org_id"`
	// UserId     string `json:"userId,omitempty" db:"user_id"`

	PasswordChangeRequired bool       `json:"passwordChangeRequired,omitempty" db:"password_change_required"`
	PasswordChange         *time.Time `json:"passwordChange,omitempty" db:"password_change"`
	MFAInitSkipped         bool       `json:"mfaInitSkipped,omitempty" db:"mfa_init_skipped"`
}

type humanContactColumns interface {
	InstanceIDCondition() database.Column
	OrgIDCondition() database.Column
	UserIDCondition() database.Column
	TypeCondition() database.Column
	CurrentValueCondition() database.Column
	VerifiedCondition() database.Column
	UnverifiedValueCondition() database.Column
}

type humanContactConditions interface {
	InstanceIDCondition(instanceID string) database.Condition
	OrgIDCondition(orgID string) database.Condition
	UserIDCondition(userID string) database.Condition
	TypeCondition(typ ContactType) database.Condition
	CurrentValueCondition(value string) database.Condition
	VerifiedCondition(verified bool) database.Condition
	UnverifiedValueCondition(value string) database.Condition
}

type humanContactChanges interface {
	SetInstanceID(instanceID string) database.Change
	SetOrgID(orgID string) database.Change
	SetUserID(userID string) database.Change
	SetType(typ ContactType) database.Change
	SetCurrentValue(value string) database.Change
	SetVerified(verified bool) database.Change
	SetUnverifiedValue(value string) database.Change
}

type UserRepository interface {
	Human() HumanRepository
	Machine() MachineRepository

	CreateHuman(ctx context.Context, client database.QueryExecutor, user *Human) (*Human, error)
	UpdateHuman(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	GetHuman(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*Human, error)
	ListHuman(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*Human, error)
	DeleteHuman(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)

	CreateMachine(ctx context.Context, client database.QueryExecutor, user *Machine) (*Machine, error)
	UpdateMachine(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	GetMachine(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*Machine, error)
	ListMachine(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*Machine, error)
	DeleteMachine(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}
