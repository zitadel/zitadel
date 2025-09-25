package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

//go:generate go tool enumer -type OrgState -transform lower -trimprefix OrgState -sql
type OrgState uint8

const (
	OrgStateActive OrgState = iota
	OrgStateInactive
)

type Organization struct {
	ID         string    `json:"id,omitempty" db:"id"`
	Name       string    `json:"name,omitempty" db:"name"`
	InstanceID string    `json:"instanceId,omitempty" db:"instance_id"`
	State      OrgState  `json:"state,omitempty" db:"state"`
	CreatedAt  time.Time `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt,omitzero" db:"updated_at"`

	Domains []*OrganizationDomain `json:"domains,omitempty" db:"-"` // domains need to be handled separately
}

// organizationColumns define all the columns of the instance table.
type organizationColumns interface {
	// IDColumn returns the column for the id field.
	IDColumn() database.Column
	// NameColumn returns the column for the name field.
	NameColumn() database.Column
	// InstanceIDColumn returns the column for the instance id field
	InstanceIDColumn() database.Column
	// StateColumn returns the column for the name field.
	StateColumn() database.Column
	// CreatedAtColumn returns the column for the created at field.
	CreatedAtColumn() database.Column
	// UpdatedAtColumn returns the column for the updated at field.
	UpdatedAtColumn() database.Column
}

// organizationConditions define all the conditions for the instance table.
type organizationConditions interface {
	// IDCondition returns an equal filter on the id field.
	IDCondition(instanceID string) database.Condition
	// NameCondition returns a filter on the name field.
	NameCondition(op database.TextOperation, name string) database.Condition
	// InstanceIDCondition returns a filter on the instance id field.
	InstanceIDCondition(instanceID string) database.Condition
	// StateCondition returns a filter on the name field.
	StateCondition(state OrgState) database.Condition
	// ExistsDomain returns a filter on the organizations domains.
	ExistsDomain(cond database.Condition) database.Condition
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

	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*Organization, error)
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*Organization, error)

	Create(ctx context.Context, client database.QueryExecutor, org *Organization) error
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)

	// LoadDomains loads the domains of the given organizations.
	// If it is called the [Organization].Domains field will be set on future calls to Get or List.
	LoadDomains() OrganizationRepository
}

type CreateOrganization struct {
	Name string `json:"name"`
}

// MemberRepository is a sub repository of the org repository and maybe the instance repository.
type MemberRepository interface {
	AddMember(ctx context.Context, orgID, userID string, roles []string) error
	SetMemberRoles(ctx context.Context, orgID, userID string, roles []string) error
	RemoveMember(ctx context.Context, orgID, userID string) error
}
