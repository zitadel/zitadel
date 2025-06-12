package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var State []string = []string{
	"ACTIVE",
	"INACTIVE",
}

type Organization struct {
	ID         string     `json:"id,omitempty" db:"id"`
	Name       string     `json:"name,omitempty" db:"name"`
	InstanceID string     `json:"instance_id,omitempty" db:"instance_id"`
	State      string     `json:"state,omitempty" db:"state"`
	CreatedAt  time.Time  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
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
}

// organizationChanges define all the changes for the instance table.
type organizationChanges interface {
	// SetName sets the name column.
	SetName(name string) database.Change
}

// OrganizationRepository is the interface for the instance repository.
type OrganizationRepository interface {
	organizationColumns
	organizationConditions
	organizationChanges

	Get(ctx context.Context, opts ...database.Condition) (*Organization, error)
	List(ctx context.Context, opts ...database.Condition) ([]Organization, error)

	Create(ctx context.Context, instance *Organization) error
	Update(ctx context.Context, condition database.Condition, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, condition database.Condition) error
}

type CreateOrganization struct {
	Name string `json:"name"`
}
