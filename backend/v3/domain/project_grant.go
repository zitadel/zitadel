package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type ProjectGrant struct {
	InstanceID             string              `json:"instanceId,omitempty" db:"instance_id"`
	ID                     string              `json:"id,omitempty" db:"id"`
	ProjectID              string              `json:"projectId,omitempty" db:"project_id"`
	GrantedOrganizationID  string              `json:"grantedOrganizationId,omitempty" db:"granted_organization_id"`
	GrantingOrganizationID string              `json:"grantingOrganizationId,omitempty" db:"granting_organization_id"`
	CreatedAt              time.Time           `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt              time.Time           `json:"updatedAt,omitzero" db:"updated_at"`
	Roles                  []*ProjectGrantRole `json:"roles,omitempty" db:"-"` // roles need to be handled separately
}

type ProjectGrantRole struct {
	InstanceID string    `json:"instanceId,omitempty" db:"instance_id"`
	GrantID    string    `json:"grantId,omitempty" db:"grant_id"`
	ProjectID  string    `json:"projectId,omitempty" db:"project_id"`
	RoleKey    string    `json:"roleKey,omitempty" db:"role_key"`
	CreatedAt  time.Time `json:"createdAt,omitzero" db:"created_at"`
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
}

type projectGrantChanges interface{}

type ProjectGrantRepository interface {
	projectGrantColumns
	projectGrantConditions
	projectGrantChanges

	// Get a single project. An error is returned if not exactly one project is found.
	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*ProjectGrant, error)
	// List projects. An empty list is returned if no projects are found.
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*ProjectGrant, error)
	// Create a new project.
	Create(ctx context.Context, client database.QueryExecutor, project *ProjectGrant) error
	// Update an existing project.
	// The condition must include the instanceID and ID of the project grant to update.
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	// Delete an existing project.
	// The condition must include the instanceID and ID of the project grant to delete.
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}
