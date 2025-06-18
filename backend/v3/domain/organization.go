package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

//go:generate enumer -type OrgState -transform lower -trimprefix OrgState
type OrgState int

const (
	OrgStateActive OrgState = iota
	OrgStateInactive
)

type Organization struct {
	ID         string     `json:"id,omitempty" db:"id"`
	Name       string     `json:"name,omitempty" db:"name"`
	InstanceID string     `json:"instanceId,omitempty" db:"instance_id"`
	State      string     `json:"state,omitempty" db:"state"`
	CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt  time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}

// organizationColumns define all the columns of the instance table.
type organizationColumns interface {
	// IDColumn returns the column for the id field.
	IDColumn() database.Column
	// NameColumn returns the column for the name field.
	NameColumn() database.Column
	// InstanceIDColumn returns the column for the default org id field
	InstanceIDColumn() database.Column
	// StateColumn returns the column for the name field.
	StateColumn() database.Column
	// CreatedAtColumn returns the column for the created at field.
	CreatedAtColumn() database.Column
	// UpdatedAtColumn returns the column for the updated at field.
	UpdatedAtColumn() database.Column
	// DeletedAtColumn returns the column for the deleted at field.
	DeletedAtColumn() database.Column
}

// organizationConditions define all the conditions for the instance table.
type organizationConditions interface {
	// IDCondition returns an equal filter on the id field.
	IDCondition(instanceID string) database.Condition
	// NameCondition returns a filter on the name field.
	NameCondition(op database.TextOperation, name string) database.Condition
	// StateCondition returns a filter on the name field.
	StateCondition(state OrgState) database.Condition
}

// organizationChanges define all the changes for the instance table.
type organizationChanges interface {
	// SetName sets the name column.
	SetName(name string) database.Change
	// SetState sets the name column.
	SetState(state OrgState) database.Change
}

// OrganizationRepository is the interface for the instance repository.
type OrganizationRepository interface {
	organizationColumns
	organizationConditions
	organizationChanges

	Get(ctx context.Context, opts ...database.Condition) (*Organization, error)
	List(ctx context.Context, opts ...database.Condition) ([]*Organization, error)

	Create(ctx context.Context, instance *Organization) error
	Update(ctx context.Context, condition database.Condition, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, condition database.Condition) (int64, error)
}

type CreateOrganization struct {
	Name string `json:"name"`
}
