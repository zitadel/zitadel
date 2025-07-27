package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

//go:generate enumer -type OrgState -transform lower -trimprefix OrgState
type OrgState uint8

const (
	OrgStateActive OrgState = iota
	OrgStateInactive
)

type Organization struct {
	ID         string    `json:"id,omitempty" db:"id"`
	Name       string    `json:"name,omitempty" db:"name"`
	InstanceID string    `json:"instanceId,omitempty" db:"instance_id"`
	State      string    `json:"state,omitempty" db:"state"`
	CreatedAt  time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

// OrgIdentifierCondition is used to help specify a single Organization,
// it will either be used as the organization ID or organization name,
// as organizations can be identified either using (instnaceID + ID) OR (instanceID + name)
type OrgIdentifierCondition interface {
	database.Condition
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
	IDCondition(instanceID string) OrgIdentifierCondition
	// NameCondition returns a filter on the name field.
	NameCondition(name string) OrgIdentifierCondition
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

// OrganizationRepository is the interface for the instance repository.
type OrganizationRepository interface {
	organizationColumns
	organizationConditions
	organizationChanges

	Get(ctx context.Context, id OrgIdentifierCondition, instance_id string, opts ...database.Condition) (*Organization, error)
	List(ctx context.Context, conditions ...database.Condition) ([]*Organization, error)

	Create(ctx context.Context, instance *Organization) error
	Update(ctx context.Context, id OrgIdentifierCondition, instance_id string, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, id OrgIdentifierCondition, instance_id string) (int64, error)
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

// DomainRepository is a sub repository of the org repository and maybe the instance repository.
type DomainRepository interface {
	AddDomain(ctx context.Context, domain string) error
	SetDomainVerified(ctx context.Context, domain string) error
	RemoveDomain(ctx context.Context, domain string) error
}
