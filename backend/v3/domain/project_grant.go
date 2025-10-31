package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

//go:generate enumer -type ProjectGrantState -transform lower -trimprefix ProjectGrantState -sql
type ProjectGrantState uint8

const (
	ProjectGrantStateActive ProjectGrantState = iota
	ProjectGrantStateInactive
)

type ProjectGrant struct {
	InstanceID             string              `json:"instanceId,omitempty" db:"instance_id"`
	ID                     string              `json:"id,omitempty" db:"id"`
	ProjectID              string              `json:"projectId,omitempty" db:"project_id"`
	GrantedOrganizationID  string              `json:"grantedOrganizationId,omitempty" db:"granted_organization_id"`
	GrantingOrganizationID string              `json:"grantingOrganizationId,omitempty" db:"granting_organization_id"`
	CreatedAt              time.Time           `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt              time.Time           `json:"updatedAt,omitzero" db:"updated_at"`
	State                  ProjectGrantState   `json:"state,omitempty" db:"state"`
	Roles                  []*ProjectGrantRole `json:"roles,omitempty" db:"-"` // roles need to be handled separately
}

type projectGrantColumns interface {
	// PrimaryKeyColumns returns the columns for the primary key fields
	PrimaryKeyColumns() []database.Column
	// InstanceIDColumn returns the column for the instance id field
	InstanceIDColumn() database.Column
	// IDColumn returns the column for the id field
	IDColumn() database.Column
	// ProjectIDColumn returns the column for the project id field
	ProjectIDColumn() database.Column
	// GrantingOrganizationIDColumn returns the column for the granting organization id field
	GrantingOrganizationIDColumn() database.Column
	// GrantedOrganizationIDColumn returns the column for the granted organization id field
	GrantedOrganizationIDColumn() database.Column
	// CreatedAtColumn returns the column for the created at field.
	CreatedAtColumn() database.Column
	// UpdatedAtColumn returns the column for the updated at field.
	UpdatedAtColumn() database.Column
}

type projectGrantConditions interface {
	// PrimaryKeyCondition returns a filter on the primary key fields.
	PrimaryKeyCondition(instanceID, id string) database.Condition
	// InstanceIDCondition returns a filter on the instance id field.
	InstanceIDCondition(instanceID string) database.Condition
	// IDCondition returns a filter on the id field.
	IDCondition(id string) database.Condition
	// ProjectIDCondition returns a filter on the project id field.
	ProjectIDCondition(projectID string) database.Condition
	// GrantingOrganizationIDCondition returns a filter on the granting organization id field.
	GrantingOrganizationIDCondition(grantingOrgID string) database.Condition
	// GrantedOrganizationIDCondition returns a filter on the granted organization id field.
	GrantedOrganizationIDCondition(grantedOrgID string) database.Condition
	// StateCondition returns a filter on the state field.
	StateCondition(state ProjectGrantState) database.Condition
}

type projectGrantChanges interface {
	// SetUpdatedAt sets the updated at column.
	// Only use this when reducing events,
	// during regular updates the DB sets this column automatically.
	SetUpdatedAt(updatedAt time.Time) database.Change
	// SetState sets the state column.
	SetState(state ProjectGrantState) database.Change
}

// ProjectGrantRepository manages project grants.
//
//go:generate mockgen -typed -package domainmock -destination ./mock/project_grant.mock.go . ProjectGrantRepository
type ProjectGrantRepository interface {
	Repository

	projectGrantColumns
	projectGrantConditions
	projectGrantChanges

	// Get a single project grant. An error is returned if not exactly one project is found.
	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*ProjectGrant, error)
	// List projects grant. An empty list is returned if no project grants are found.
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*ProjectGrant, error)
	// Create a new project grant.
	Create(ctx context.Context, client database.QueryExecutor, project *ProjectGrant) error
	// Update an existing project grant.
	// The condition must include the instanceID and ID of the project grant to update.
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	// Delete an existing project grant.
	// The condition must include the instanceID and ID of the project grant to delete.
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)

	// Role returns the sub-repository for project grant roles.
	Role() ProjectGrantRoleRepository
}

type ProjectGrantRole struct {
	InstanceID   string    `json:"instanceId,omitempty" db:"instance_id"`
	GrantID      string    `json:"grantId,omitempty" db:"grant_id"`
	ProjectOrgID string    `json:"projectOrgId,omitempty" db:"project_org_id"`
	ProjectID    string    `json:"projectId,omitempty" db:"project_id"`
	Key          string    `json:"key,omitempty" db:"key"`
	CreatedAt    time.Time `json:"createdAt,omitzero" db:"created_at"`
}

type projectGrantRoleColumns interface {
	// PrimaryKeyColumns returns the columns for the primary key fields
	PrimaryKeyColumns() []database.Column
	// InstanceIDColumn returns the column for the instance id field
	InstanceIDColumn() database.Column
	// GrantIDColumn returns the column for the grant id field
	GrantIDColumn() database.Column
	// KeyColumn returns the column for the key field
	KeyColumn() database.Column
	// CreatedAtColumn returns the column for the created at field.
	CreatedAtColumn() database.Column
	// ProjectIDColumn returns the column for the project id field
	ProjectIDColumn() database.Column
	// ProjectOrgIDColumn returns the column for the project organization id field
	ProjectOrgIDColumn() database.Column
}

type projectGrantRoleConditions interface {
	// PrimaryKeyCondition returns a filter on the primary key fields.
	PrimaryKeyCondition(instanceID, grantID, key string) database.Condition
	// InstanceIDCondition returns an equal filter on the instance id field.
	InstanceIDCondition(instanceID string) database.Condition
	// GrantIDCondition returns an equal filter on the grant id field.
	GrantIDCondition(grantID string) database.Condition
	// KeyCondition returns an equal filter on the key field.
	KeyCondition(key string) database.Condition
}

type projectGrantRoleChanges interface{}

type ProjectGrantRoleRepository interface {
	projectGrantRoleColumns
	projectGrantRoleConditions
	projectGrantRoleChanges

	// List project grant roles. An empty list is returned if no project roles are found.
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*ProjectGrantRole, error)
	// Add a new project grant role.
	Add(ctx context.Context, client database.QueryExecutor, role *ProjectGrantRole) error
	// Remove an existing project grant role.
	// The condition must include the instanceID, grantID and key of the project grant role to delete.
	Remove(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}
