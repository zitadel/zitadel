package domain

import (
	"context"
	"time"

	v4 "github.com/zitadel/zitadel/backend/v3/storage/database/repository/stmt/v4"
)

type userColumns interface {
	// TODO: move v4.columns to domain
	InstanceIDColumn() v4.Column
	OrgIDColumn() v4.Column
	IDColumn() v4.Column
	usernameColumn() v4.Column
	CreatedAtColumn() v4.Column
	UpdatedAtColumn() v4.Column
	DeletedAtColumn() v4.Column
}

type userConditions interface {
	InstanceIDCondition(instanceID string) v4.Condition
	OrgIDCondition(orgID string) v4.Condition
	IDCondition(userID string) v4.Condition
	UsernameCondition(op v4.TextOperator, username string) v4.Condition
	CreatedAtCondition(op v4.NumberOperator, createdAt time.Time) v4.Condition
	UpdatedAtCondition(op v4.NumberOperator, updatedAt time.Time) v4.Condition
	DeletedCondition(isDeleted bool) v4.Condition
	DeletedAtCondition(op v4.NumberOperator, deletedAt time.Time) v4.Condition
}

type userChanges interface {
	SetUsername(username string) v4.Change
}

type UserRepository interface {
	userColumns
	userConditions
	userChanges
	// TODO: move condition to domain
	Get(ctx context.Context, opts v4.QueryOption) (*User, error)
	List(ctx context.Context, opts v4.QueryOption) ([]*User, error)
	Delete(ctx context.Context, condition v4.Condition) error

	Human() HumanRepository
	Machine() MachineRepository
}

type humanColumns interface {
	userColumns
	FirstNameColumn() v4.Column
	LastNameColumn() v4.Column
	EmailAddressColumn() v4.Column
	EmailVerifiedAtColumn() v4.Column
	PhoneNumberColumn() v4.Column
	PhoneVerifiedAtColumn() v4.Column
}

type humanConditions interface {
	userConditions
	FirstNameCondition(op v4.TextOperator, firstName string) v4.Condition
	LastNameCondition(op v4.TextOperator, lastName string) v4.Condition
	EmailAddressCondition(op v4.TextOperator, email string) v4.Condition
	EmailAddressVerifiedCondition(isVerified bool) v4.Condition
	EmailVerifiedAtCondition(op v4.TextOperator, emailVerifiedAt string) v4.Condition
	PhoneNumberCondition(op v4.TextOperator, phoneNumber string) v4.Condition
	PhoneNumberVerifiedCondition(isVerified bool) v4.Condition
	PhoneVerifiedAtCondition(op v4.TextOperator, phoneVerifiedAt string) v4.Condition
}

type humanChanges interface {
	userChanges
	SetFirstName(firstName string) v4.Change
	SetLastName(lastName string) v4.Change

	SetEmail(address string, verified *time.Time) v4.Change
	SetEmailAddress(email string) v4.Change
	SetEmailVerifiedAt(emailVerifiedAt time.Time) v4.Change

	SetPhone(number string, verifiedAt *time.Time) v4.Change
	SetPhoneNumber(phoneNumber string) v4.Change
	SetPhoneVerifiedAt(phoneVerifiedAt time.Time) v4.Change
}

type HumanRepository interface {
	humanColumns
	humanConditions
	humanChanges

	GetEmail(ctx context.Context, condition v4.Condition) (*Email, error)
	// TODO: replace any with add email update columns
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, condition v4.Condition, changes ...v4.Change) error
}

type machineColumns interface {
	userColumns
	DescriptionColumn() v4.Column
}

type machineConditions interface {
	userConditions
	DescriptionCondition(op v4.TextOperator, description string) v4.Condition
}

type machineChanges interface {
	userChanges
	SetDescription(description string) v4.Change
}

type MachineRepository interface {
	machineColumns
	machineConditions
	machineChanges

	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, condition v4.Condition, changes ...v4.Change) error
}

// type UserRepository interface {
// 	// Get(ctx context.Context, clauses ...UserClause) (*User, error)
// 	// Search(ctx context.Context, clauses ...UserClause) ([]*User, error)

// 	UserQuery[UserOperation]
// 	Human() HumanQuery
// 	Machine() MachineQuery
// }

// type UserQuery[Op UserOperation] interface {
// 	ByID(id string) UserQuery[Op]
// 	Username(username string) UserQuery[Op]
// 	Exec() Op
// }

// type HumanQuery interface {
// 	UserQuery[HumanOperation]
// 	Email(op TextOperation, email string) HumanQuery
// 	HumanOperation
// }

// type MachineQuery interface {
// 	UserQuery[MachineOperation]
// 	MachineOperation
// }

// type UserClause interface {
// 	Field() UserField
// 	Operation() Operation
// 	Args() []any
// }

// type UserField uint8

// const (
// 	// Fields used for all users
// 	UserFieldInstanceID UserField = iota + 1
// 	UserFieldOrgID
// 	UserFieldID
// 	UserFieldUsername

// 	// Fields used for human users
// 	UserHumanFieldEmail
// 	UserHumanFieldEmailVerified

// 	// Fields used for machine users
// 	UserMachineFieldDescription
// )

// type userByIDClause struct {
// 	id string
// }

// func (c *userByIDClause) Field() UserField {
// 	return UserFieldID
// }

// func (c *userByIDClause) Operation() Operation {
// 	return TextOperationEqual
// }

// func (c *userByIDClause) Args() []any {
// 	return []any{c.id}
// }

// type UserOperation interface {
// 	Delete(ctx context.Context) error
// 	SetUsername(ctx context.Context, username string) error
// }

// type HumanOperation interface {
// 	UserOperation
// 	SetEmail(ctx context.Context, email string) error
// 	SetEmailVerified(ctx context.Context, email string) error
// 	GetEmail(ctx context.Context) (*Email, error)
// }

// type MachineOperation interface {
// 	UserOperation
// 	SetDescription(ctx context.Context, description string) error
// }

type User struct {
	v4.User
}

type Email struct {
	v4.Email
	IsVerified bool
}

// type userTraits interface {
// 	isUserTraits()
// }

// type Human struct {
// 	Email *Email `json:"email"`
// }

// func (*Human) isUserTraits() {}

// type Machine struct {
// 	Description string `json:"description"`
// }

// func (*Machine) isUserTraits() {}

// type Email struct {
// 	Address    string `json:"address"`
// 	IsVerified bool   `json:"isVerified"`
// }
