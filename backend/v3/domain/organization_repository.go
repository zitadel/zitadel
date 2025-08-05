package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// OrganizationRepository is the interface for the instance repository.
type OrganizationRepository interface {
	organizationColumns
	organizationConditions
	organizationChanges

	Get(ctx context.Context, opts ...database.QueryOption) (*Organization, error)
	List(ctx context.Context, opts ...database.QueryOption) ([]*Organization, error)

	Create(ctx context.Context, instance *Organization) error
	Update(ctx context.Context, condition database.Condition, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, condition database.Condition) (int64, error)

	// Domains returns the domain sub repository for the organization.
	// If shouldLoad is true, the domains will be loaded from the database and written to the [Organization].Domains field.
	// If shouldLoad is set to true once, the Domains field will be set even if shouldLoad is false in the future.
	Domains(shouldLoad bool) OrganizationDomainRepository
}

// organizationColumns define all the columns of the instance table.
type organizationColumns interface {
	// IDColumn returns the column for the id field.
	// `qualified` indicates if the column should be qualified with the table name.
	IDColumn(qualified bool) database.Column
	// NameColumn returns the column for the name field.
	// `qualified` indicates if the column should be qualified with the table name.
	NameColumn(qualified bool) database.Column
	// InstanceIDColumn returns the column for the default org id field
	// `qualified` indicates if the column should be qualified with the table name.
	InstanceIDColumn(qualified bool) database.Column
	// StateColumn returns the column for the name field.
	// `qualified` indicates if the column should be qualified with the table name.
	StateColumn(qualified bool) database.Column
	// CreatedAtColumn returns the column for the created at field.
	// `qualified` indicates if the column should be qualified with the table name.
	CreatedAtColumn(qualified bool) database.Column
	// UpdatedAtColumn returns the column for the updated at field.
	// `qualified` indicates if the column should be qualified with the table name.
	UpdatedAtColumn(qualified bool) database.Column
}

// organizationConditions define all the conditions for the instance table.
type organizationConditions interface {
	// IDCondition returns an equal filter on the id field.
	IDCondition(id string) database.Condition
	// NameCondition returns a filter on the name field.
	NameCondition(op database.TextOperation, name string) database.Condition
	// InstanceIDCondition returns a filter on the instance id field.
	InstanceIDCondition(instanceID string) database.Condition
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
