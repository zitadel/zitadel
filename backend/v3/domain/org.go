package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/cache"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type OrgState uint8

const (
	OrgStateActive OrgState = iota + 1
	OrgStateInactive
)

type Org struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	State OrgState `json:"state"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type orgCacheIndex uint8

const (
	orgCacheIndexUndefined orgCacheIndex = iota
	orgCacheIndexID
)

// Keys implements [cache.Entry].
func (o *Org) Keys(index orgCacheIndex) (key []string) {
	if index == orgCacheIndexID {
		return []string{o.ID}
	}
	return nil
}

var _ cache.Entry[orgCacheIndex, string] = (*Org)(nil)

type orgColumns interface {
	// InstanceIDColumn returns the column for the instance id field.
	InstanceIDColumn() database.Column
	// IDColumn returns the column for the id field.
	IDColumn() database.Column
	// NameColumn returns the column for the name field.
	NameColumn() database.Column
	// StateColumn returns the column for the state field.
	StateColumn() database.Column
	// CreatedAtColumn returns the column for the created at field.
	CreatedAtColumn() database.Column
	// UpdatedAtColumn returns the column for the updated at field.
	UpdatedAtColumn() database.Column
	// DeletedAtColumn returns the column for the deleted at field.
	DeletedAtColumn() database.Column
}

type orgConditions interface {
	// InstanceIDCondition returns an equal filter on the instance id field.
	InstanceIDCondition(instanceID string) database.Condition
	// IDCondition returns an equal filter on the id field.
	IDCondition(orgID string) database.Condition
	// NameCondition returns a filter on the name field.
	NameCondition(op database.TextOperation, name string) database.Condition
	// StateCondition returns a filter on the state field.
	StateCondition(op database.NumberOperation, state OrgState) database.Condition
}

type orgChanges interface {
	// SetName sets the name column.
	SetName(name string) database.Change
	// SetState sets the state column.
	SetState(state OrgState) database.Change
}

type OrgRepository interface {
	orgColumns
	orgConditions
	orgChanges

	// Member returns the admin repository.
	Member() MemberRepository
	// Domain returns the domain repository.
	Domain() DomainRepository

	// Get returns an org based on the given condition.
	Get(ctx context.Context, opts ...database.QueryOption) (*Org, error)
	// List returns a list of orgs based on the given condition.
	List(ctx context.Context, opts ...database.QueryOption) ([]*Org, error)
	// Create creates a new org.
	Create(ctx context.Context, org *Org) error
	// Delete removes orgs based on the given condition.
	Delete(ctx context.Context, condition database.Condition) error
	// Update executes the given changes based on the given condition.
	Update(ctx context.Context, condition database.Condition, changes ...database.Change) error
}

type OrgOperation interface {
	MemberRepository
	DomainRepository
	Update(ctx context.Context, org *Org) error
	Delete(ctx context.Context) error
}

type MemberRepository interface {
	AddMember(ctx context.Context, orgID, userID string, roles []string) error
	SetMemberRoles(ctx context.Context, orgID, userID string, roles []string) error
	RemoveMember(ctx context.Context, orgID, userID string) error
}

type DomainRepository interface {
	AddDomain(ctx context.Context, domain string) error
	SetDomainVerified(ctx context.Context, domain string) error
	RemoveDomain(ctx context.Context, domain string) error
}
