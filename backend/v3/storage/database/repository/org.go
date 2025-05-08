package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

type org struct {
	repository
}

func OrgRepository(client database.QueryExecutor) domain.OrgRepository {
	return &org{
		repository: repository{
			client: client,
		},
	}
}

// Create implements [domain.OrgRepository].
func (o *org) Create(ctx context.Context, org *domain.Org) error {
	org.CreatedAt = time.Now()
	org.UpdatedAt = org.CreatedAt
	return nil
}

// Delete implements [domain.OrgRepository].
func (o *org) Delete(ctx context.Context, condition database.Condition) error {
	return nil
}

// Get implements [domain.OrgRepository].
func (o *org) Get(ctx context.Context, opts ...database.QueryOption) (*domain.Org, error) {
	panic("unimplemented")
}

// List implements [domain.OrgRepository].
func (o *org) List(ctx context.Context, opts ...database.QueryOption) ([]*domain.Org, error) {
	panic("unimplemented")
}

// Update implements [domain.OrgRepository].
func (o *org) Update(ctx context.Context, condition database.Condition, changes ...database.Change) error {
	panic("unimplemented")
}

func (o *org) Member() domain.MemberRepository {
	return &orgMember{o}
}

func (o *org) Domain() domain.DomainRepository {
	return &orgDomain{o}
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetName implements [domain.orgChanges].
func (o *org) SetName(name string) database.Change {
	return database.NewChange(o.NameColumn(), name)
}

// SetState implements [domain.orgChanges].
func (o *org) SetState(state domain.OrgState) database.Change {
	return database.NewChange(o.StateColumn(), state)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// IDCondition implements [domain.orgConditions].
func (o *org) IDCondition(orgID string) database.Condition {
	return database.NewTextCondition(o.IDColumn(), database.TextOperationEqual, orgID)
}

// InstanceIDCondition implements [domain.orgConditions].
func (o *org) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(o.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// NameCondition implements [domain.orgConditions].
func (o *org) NameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(o.NameColumn(), op, name)
}

// StateCondition implements [domain.orgConditions].
func (o *org) StateCondition(op database.NumberOperation, state domain.OrgState) database.Condition {
	return database.NewNumberCondition(o.StateColumn(), op, state)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// CreatedAtColumn implements [domain.orgColumns].
func (o *org) CreatedAtColumn() database.Column {
	return database.NewColumn("created_at")
}

// DeletedAtColumn implements [domain.orgColumns].
func (o *org) DeletedAtColumn() database.Column {
	return database.NewColumn("deleted_at")
}

// IDColumn implements [domain.orgColumns].
func (o *org) IDColumn() database.Column {
	return database.NewColumn("id")
}

// InstanceIDColumn implements [domain.orgColumns].
func (o *org) InstanceIDColumn() database.Column {
	return database.NewColumn("instance_id")
}

// NameColumn implements [domain.orgColumns].
func (o *org) NameColumn() database.Column {
	return database.NewColumn("name")
}

// StateColumn implements [domain.orgColumns].
func (o *org) StateColumn() database.Column {
	return database.NewColumn("state")
}

// UpdatedAtColumn implements [domain.orgColumns].
func (o *org) UpdatedAtColumn() database.Column {
	return database.NewColumn("updated_at")
}

var _ domain.OrgRepository = (*org)(nil)
