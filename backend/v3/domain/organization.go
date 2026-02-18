package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

//go:generate enumer -type OrgState -transform lower -trimprefix OrgState -sql
type OrgState uint8

// Must be in the same order and quantity as zitadel/org/v2/org.proto
const (
	OrgStateUnspecified OrgState = iota
	OrgStateActive
	OrgStateInactive

	// TODO(IAM-Marco): This should be removed in next versions of Zitadel I believe. We are hard deleting,
	// so not sure when this state would be used. It is kept here just for compatibility with the gRPC model
	OrgStateRemoved
)

type Organization struct {
	ID         string    `json:"id,omitempty" db:"id"`
	Name       string    `json:"name,omitempty" db:"name"`
	InstanceID string    `json:"instanceId,omitempty" db:"instance_id"`
	State      OrgState  `json:"state,omitempty" db:"state"`
	CreatedAt  time.Time `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt,omitzero" db:"updated_at"`

	Domains  []*OrganizationDomain   `json:"domains,omitempty" db:"-"`  // domains need to be handled separately
	Metadata []*OrganizationMetadata `json:"metadata,omitempty" db:"-"` // metadata need to be handled separately
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
	// PrimaryKeyCondition returns a filter on the primary key fields.
	PrimaryKeyCondition(instanceID, organizationID string) database.Condition
	// IDCondition returns an equal filter on the id field.
	IDCondition(orgID string) database.Condition
	// NameCondition returns a filter on the name field.
	NameCondition(op database.TextOperation, name string) database.Condition
	// InstanceIDCondition returns a filter on the instance id field.
	InstanceIDCondition(instanceID string) database.Condition
	// StateCondition returns a filter on the name field.
	StateCondition(state OrgState) database.Condition
	// ExistsDomain returns a filter on the organizations domains.
	ExistsDomain(cond database.Condition) database.Condition
	// ExistsMetadata returns a filter on the organizations metadata.
	ExistsMetadata(cond database.Condition) database.Condition
}

// organizationChanges define all the changes for the instance table.
type organizationChanges interface {
	// SetName sets the name column.
	SetName(name string) database.Change
	// SetState sets the name column.
	SetState(state OrgState) database.Change
	// SetUpdatedAt sets the updated at column.
	SetUpdatedAt(updatedAt time.Time) database.Change
}

//go:generate mockgen -typed -package domainmock -destination ./mock/org.mock.go . OrganizationRepository

// OrganizationRepository is the interface for the instance repository.
type OrganizationRepository interface {
	Repository

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
	// LoadMetadata loads the metadata of the given organizations.
	// If it is called the [Organization].Metadata field will be set on future calls to Get or List.
	LoadMetadata() OrganizationRepository
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

func (o *Organization) PrimaryDomain() string {
	for _, d := range o.Domains {
		if d.IsPrimary {
			return d.Domain
		}
	}

	return ""
}
