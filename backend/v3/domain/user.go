package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

//go:generate enumer -type UserState -transform lower -trimprefix UserState -sql
type UserState uint8

const (
	UserStateInitial UserState = iota
	UserStateActive
	UserStateInactive
	UserStateLocked
	UserStateSuspended
)

// User represents a the polymorphic user in the system.
// It can be a human user or a machine user.
// Meaning that either Human or Machine is set, the other is nil.
type User struct {
	InstanceID          string          `json:"instanceId,omitempty" db:"instance_id"`
	OrgID               string          `json:"orgId,omitempty" db:"org_id"`
	ID                  string          `json:"id,omitempty" db:"id"`
	Username            string          `json:"username,omitempty" db:"username"`
	IsUsernameOrgUnique bool            `json:"usernameOrgUnique,omitempty" db:"username_org_unique"`
	State               UserState       `json:"state,omitempty" db:"state"`
	Metadata            []*UserMetadata `json:"metadata,omitempty" db:"-"` // metadata need to be handled separately

	*Machine `db:"machine"`
	*Human   `db:"human"`

	CreatedAt time.Time `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt,omitzero" db:"updated_at"`
}

// UserType defines the type of user
// It is used for List and Get methods of the [UserRepository] to filter users by type.
type UserType uint8

const (
	UserTypeHuman UserType = iota + 1
	UserTypeMachine
)

type UserRepository interface {
	Repository

	userColumns
	userConditions
	userChanges

	Human() HumanUserRepository
	Machine() MachineUserRepository

	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*User, error)
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*User, error)
	// DISCUSS(adlerhurst): Instead of having this method the Create methods could be on the sub-repositories?
	// The sub repos should then get the CreateMachineCommand or CreateHumanCommand as parameter instead of the User.
	// Passing the command instead of the object would generally simplify the domain logic for creation.
	Create(ctx context.Context, client database.QueryExecutor, user *User) error
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)

	LoadMetadata() UserRepository
}

type userColumns interface {
	PrimaryKeyColumns() []database.Column
	InstanceIDColumn() database.Column
	OrgIDColumn() database.Column
	IDColumn() database.Column
	UsernameColumn() database.Column
	UsernameOrgUniqueColumn() database.Column
	StateColumn() database.Column
	CreatedAtColumn() database.Column
	UpdatedAtColumn() database.Column
	TypeCondition(userType UserType) database.Condition
}

type userConditions interface {
	PrimaryKeyCondition(instanceID, userID string) database.Condition
	InstanceIDCondition(instanceID string) database.Condition
	OrgIDCondition(orgID string) database.Condition
	IDCondition(userID string) database.Condition
	UsernameCondition(op database.TextOperation, username string) database.Condition
	UsernameOrgUniqueCondition(condition bool) database.Condition
	StateCondition(state UserState) database.Condition
	CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition
	UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition

	ExistsMetadata(cond database.Condition) database.Condition
}

type userChanges interface {
	SetUsername(username string) database.Change
	SetUsernameOrgUnique(usernameOrgUnique bool) database.Change
	SetState(state UserState) database.Change
	SetUpdatedAt(updatedAt time.Time) database.Change
}
