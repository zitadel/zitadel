package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// userColumns define all the columns of the user table.
type userColumns interface {
	// InstanceIDColumn returns the column for the instance id field.
	InstanceIDColumn() database.Column
	// OrgIDColumn returns the column for the org id field.
	OrgIDColumn() database.Column
	// IDColumn returns the column for the id field.
	IDColumn() database.Column
	// UsernameColumn returns the column for the username field.
	UsernameColumn() database.Column
	// CreatedAtColumn returns the column for the created at field.
	CreatedAtColumn() database.Column
	// UpdatedAtColumn returns the column for the updated at field.
	UpdatedAtColumn() database.Column
	// DeletedAtColumn returns the column for the deleted at field.
	DeletedAtColumn() database.Column
}

// userConditions define all the conditions for the user table.
type userConditions interface {
	// InstanceIDCondition returns an equal filter on the instance id field.
	InstanceIDCondition(instanceID string) database.Condition
	// OrgIDCondition returns an equal filter on the org id field.
	OrgIDCondition(orgID string) database.Condition
	// IDCondition returns an equal filter on the id field.
	IDCondition(userID string) database.Condition
	// UsernameCondition returns a filter on the username field.
	UsernameCondition(op database.TextOperation, username string) database.Condition
	// CreatedAtCondition returns a filter on the created at field.
	CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition
	// UpdatedAtCondition returns a filter on the updated at field.
	UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition
	// DeletedAtCondition filters for deleted users is isDeleted is set to true otherwise only not deleted users must be filtered.
	DeletedCondition(isDeleted bool) database.Condition
	// DeletedAtCondition filters for deleted users based on the given parameters.
	DeletedAtCondition(op database.NumberOperation, deletedAt time.Time) database.Condition
}

// userChanges define all the changes for the user table.
type userChanges interface {
	// SetUsername sets the username column.
	SetUsername(username string) database.Change
}

// UserRepository is the interface for the user repository.
type UserRepository interface {
	userColumns
	userConditions
	userChanges
	// Get returns a user based on the given condition.
	Get(ctx context.Context, opts ...database.QueryOption) (*User, error)
	// List returns a list of users based on the given condition.
	List(ctx context.Context, opts ...database.QueryOption) ([]*User, error)
	// Create creates a new user.
	Create(ctx context.Context, user *User) error
	// Delete removes users based on the given condition.
	Delete(ctx context.Context, condition database.Condition) error
	// Human returns the [HumanRepository].
	Human() HumanRepository
	// Machine returns the [MachineRepository].
	Machine() MachineRepository
}

// humanColumns define all the columns of the human table which inherits the user table.
type humanColumns interface {
	userColumns
	// FirstNameColumn returns the column for the first name field.
	FirstNameColumn() database.Column
	// LastNameColumn returns the column for the last name field.
	LastNameColumn() database.Column
	// EmailAddressColumn returns the column for the email address field.
	EmailAddressColumn() database.Column
	// EmailVerifiedAtColumn returns the column for the email verified at field.
	EmailVerifiedAtColumn() database.Column
	// PhoneNumberColumn returns the column for the phone number field.
	PhoneNumberColumn() database.Column
	// PhoneVerifiedAtColumn returns the column for the phone verified at field.
	PhoneVerifiedAtColumn() database.Column
}

// humanConditions define all the conditions for the human table which inherits the user table.
type humanConditions interface {
	userConditions
	// FirstNameCondition returns a filter on the first name field.
	FirstNameCondition(op database.TextOperation, firstName string) database.Condition
	// LastNameCondition returns a filter on the last name field.
	LastNameCondition(op database.TextOperation, lastName string) database.Condition
	// EmailAddressCondition returns a filter on the email address field.
	EmailAddressCondition(op database.TextOperation, email string) database.Condition
	// EmailVerifiedCondition returns a filter that checks if the email is verified or not.
	EmailVerifiedCondition(isVerified bool) database.Condition
	// EmailVerifiedAtCondition returns a filter on the email verified at field.
	EmailVerifiedAtCondition(op database.NumberOperation, emailVerifiedAt time.Time) database.Condition

	// PhoneNumberCondition returns a filter on the phone number field.
	PhoneNumberCondition(op database.TextOperation, phoneNumber string) database.Condition
	// PhoneVerifiedCondition returns a filter that checks if the phone is verified or not.
	PhoneVerifiedCondition(isVerified bool) database.Condition
	// PhoneVerifiedAtCondition returns a filter on the phone verified at field.
	PhoneVerifiedAtCondition(op database.NumberOperation, phoneVerifiedAt time.Time) database.Condition
}

// humanChanges define all the changes for the human table which inherits the user table.
type humanChanges interface {
	userChanges
	// SetFirstName sets the first name field of the human.
	SetFirstName(firstName string) database.Change
	// SetLastName sets the last name field of the human.
	SetLastName(lastName string) database.Change

	// SetEmail sets the email address and verified field of the email
	// if verifiedAt is nil the email is not verified
	SetEmail(address string, verifiedAt *time.Time) database.Change
	// SetEmailAddress sets the email address field of the email
	SetEmailAddress(email string) database.Change
	// SetEmailVerifiedAt sets the verified column of the email
	// if at is zero the statement uses the database timestamp
	SetEmailVerifiedAt(at time.Time) database.Change

	// SetPhone sets the phone number and verified field
	// if verifiedAt is nil the phone is not verified
	SetPhone(number string, verifiedAt *time.Time) database.Change
	// SetPhoneNumber sets the phone number field
	SetPhoneNumber(phoneNumber string) database.Change
	// SetPhoneVerifiedAt sets the verified field of the phone
	// if at is zero the statement uses the database timestamp
	SetPhoneVerifiedAt(at time.Time) database.Change
}

// HumanRepository is the interface for the human repository it inherits the user repository.
type HumanRepository interface {
	humanColumns
	humanConditions
	humanChanges

	// Get returns an email based on the given condition.
	GetEmail(ctx context.Context, condition database.Condition) (*Email, error)
	// Update updates human users based on the given condition and changes.
	Update(ctx context.Context, condition database.Condition, changes ...database.Change) error
}

// machineColumns define all the columns of the machine table which inherits the user table.
type machineColumns interface {
	userColumns
	// DescriptionColumn returns the column for the description field.
	DescriptionColumn() database.Column
}

// machineConditions define all the conditions for the machine table which inherits the user table.
type machineConditions interface {
	userConditions
	// DescriptionCondition returns a filter on the description field.
	DescriptionCondition(op database.TextOperation, description string) database.Condition
}

// machineChanges define all the changes for the machine table which inherits the user table.
type machineChanges interface {
	userChanges
	// SetDescription sets the description field of the machine.
	SetDescription(description string) database.Change
}

// MachineRepository is the interface for the machine repository it inherits the user repository.
type MachineRepository interface {
	// Update updates machine users based on the given condition and changes.
	Update(ctx context.Context, condition database.Condition, changes ...database.Change) error

	machineColumns
	machineConditions
	machineChanges
}

// UserTraits is implemented by [Human] and [Machine].
type UserTraits interface {
	Type() UserType
}

type UserType string

const (
	UserTypeHuman   UserType = "human"
	UserTypeMachine UserType = "machine"
)

type User struct {
	InstanceID string
	OrgID      string
	ID         string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	Username string

	Traits UserTraits
}

type Human struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     *Email `json:"email,omitempty"`
	Phone     *Phone `json:"phone,omitempty"`
}

// Type implements [UserTraits].
func (h *Human) Type() UserType {
	return UserTypeHuman
}

var _ UserTraits = (*Human)(nil)

type Email struct {
	Address    string    `json:"address"`
	VerifiedAt time.Time `json:"verifiedAt"`
}

type Phone struct {
	Number     string    `json:"number"`
	VerifiedAt time.Time `json:"verifiedAt"`
}

type Machine struct {
	Description string `json:"description"`
}

// Type implements [UserTraits].
func (m *Machine) Type() UserType {
	return UserTypeMachine
}

var _ UserTraits = (*Machine)(nil)
