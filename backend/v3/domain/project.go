package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

//go:generate enumer -type ProjectState -transform lower -trimprefix ProjectState -sql
type ProjectState uint8

const (
	ProjectStateActive ProjectState = iota
	ProjectStateInactive
)

type Project struct {
	InstanceID               string       `json:"instanceId,omitempty" db:"instance_id"`
	OrganizationID           string       `json:"organizationId,omitempty" db:"organization_id"`
	ID                       string       `json:"id,omitempty" db:"id"`
	CreatedAt                time.Time    `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt                time.Time    `json:"updatedAt,omitzero" db:"updated_at"`
	Name                     string       `json:"name,omitempty" db:"name"`
	State                    ProjectState `json:"state,omitempty" db:"state"`
	ShouldAssertRole         bool         `json:"shouldAssertRole,omitempty" db:"should_assert_role"`
	IsAuthorizationRequired  bool         `json:"isAuthorizationRequired,omitempty" db:"is_authorization_required"`
	IsProjectAccessRequired  bool         `json:"isProjectAccessRequired,omitempty" db:"is_project_access_required"`
	UsedLabelingSettingOwner int16        `json:"usedLabelingSettingOwner,omitempty" db:"used_labeling_setting_owner"`
}

type projectColumns interface {
	// PrimaryKeyColumns returns the columns for the primary key fields
	PrimaryKeyColumns() []database.Column
	// InstanceIDColumn returns the column for the instance id field
	InstanceIDColumn() database.Column
	// OrganizationIDColumn returns the column for the organization id field
	OrganizationIDColumn() database.Column
	// IDColumn returns the column for the id field.
	IDColumn() database.Column
	// CreatedAtColumn returns the column for the created at field.
	CreatedAtColumn() database.Column
	// UpdatedAtColumn returns the column for the updated at field.
	UpdatedAtColumn() database.Column
	// NameColumn returns the column for the name field.
	NameColumn() database.Column
	// StateColumn returns the column for the state field.
	StateColumn() database.Column
	// ShouldAssertRoleColumn returns the column for the should assert role field.
	ShouldAssertRoleColumn() database.Column
	// IsAuthorizationRequiredColumn returns the column for the is authorization required field.
	IsAuthorizationRequiredColumn() database.Column
	// IsProjectAccessRequiredColumn returns the column for the is project access required field.
	IsProjectAccessRequiredColumn() database.Column
	// UsedLabelingSettingOwnerColumn returns the column for the used labeling setting owner field.
	UsedLabelingSettingOwnerColumn() database.Column
}

type projectConditions interface {
	// PrimaryKeyCondition returns a filter on the primary key fields.
	PrimaryKeyCondition(instanceID, projectID string) database.Condition
	// InstanceIDCondition returns a filter on the instance id field.
	InstanceIDCondition(instanceID string) database.Condition
	// OrganizationIDCondition returns a filter on the organization id field.
	OrganizationIDCondition(organizationID string) database.Condition
	// IDCondition returns an equal filter on the id field.
	IDCondition(projectID string) database.Condition
	// NameCondition returns a filter on the name field.
	NameCondition(op database.TextOperation, name string) database.Condition
	// StateCondition returns a filter on the state field.
	StateCondition(state ProjectState) database.Condition
}

type projectChanges interface {
	// SetUpdatedAt sets the updated at column.
	// Only use this when reducing events,
	// during regular updates the DB sets this column automatically.
	SetUpdatedAt(updatedAt time.Time) database.Change
	// SetName sets the name column.
	SetName(name string) database.Change
	// SetState sets the state column.
	SetState(state ProjectState) database.Change
	// SetShouldAssertRole sets the should assert role column.
	SetShouldAssertRole(shouldAssertRole bool) database.Change
	// SetIsAuthorizationRequired sets the is authorization required column.
	SetIsAuthorizationRequired(isAuthorizationRequired bool) database.Change
	// SetIsProjectAccessRequired sets the is project access required column.
	SetIsProjectAccessRequired(isProjectAccessRequired bool) database.Change
	// SetUsedLabelingSettingOwner sets the used labeling setting owner column.
	SetUsedLabelingSettingOwner(usedLabelingSettingOwner int16) database.Change
}

// ProjectRepository manages projects and project roles.
type ProjectRepository interface {
	projectColumns
	projectConditions
	projectChanges

	// Get a single project. An error is returned if not exactly one project is found.
	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*Project, error)
	// List projects. An empty list is returned if no projects are found.
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*Project, error)
	// Create a new project.
	Create(ctx context.Context, client database.QueryExecutor, project *Project) error
	// Update an existing project.
	// The condition must include the instanceID and ID of the project to update.
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	// Delete an existing project.
	// The condition must include the instanceID and ID of the project to delete.
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)

	// Role returns the sub-repository for project roles.
	Role() ProjectRoleRepository
}

type ProjectRole struct {
	InstanceID     string    `json:"instanceId,omitempty" db:"instance_id"`
	OrganizationID string    `json:"organizationId,omitempty" db:"organization_id"`
	ProjectID      string    `json:"projectId,omitempty" db:"project_id"`
	CreatedAt      time.Time `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt,omitzero" db:"updated_at"`
	Key            string    `json:"key,omitempty" db:"key"`
	DisplayName    string    `json:"displayName,omitempty" db:"display_name"`
	RoleGroup      *string   `json:"roleGroup,omitempty" db:"role_group"`
}

type projectRoleColumns interface {
	// PrimaryKeyColumns returns the columns for the primary key fields
	PrimaryKeyColumns() []database.Column
	// InstanceIDColumn returns the column for the instance id field
	InstanceIDColumn() database.Column
	// OrganizationIDColumn returns the column for the organization id field
	OrganizationIDColumn() database.Column
	// ProjectIDColumn returns the column for the project id field
	ProjectIDColumn() database.Column
	// CreatedAtColumn returns the column for the created at field.
	CreatedAtColumn() database.Column
	// UpdatedAtColumn returns the column for the updated at field.
	UpdatedAtColumn() database.Column
	// KeyColumn returns the column for the key field.
	KeyColumn() database.Column
	// DisplayNameColumn returns the column for the display name field.
	DisplayNameColumn() database.Column
	// RoleGroupColumn returns the column for the role group field.
	RoleGroupColumn() database.Column
}

type projectRoleConditions interface {
	// PrimaryKeyCondition returns a filter on the primary key fields.
	PrimaryKeyCondition(instanceID, projectID, key string) database.Condition
	// InstanceIDCondition returns an equal filter on the instance id field.
	InstanceIDCondition(instanceID string) database.Condition
	// ProjectIDCondition returns an equal filter on the project id field.
	ProjectIDCondition(projectID string) database.Condition
	// KeyCondition returns an equal filter on the key field.
	KeyCondition(key string) database.Condition
	// DisplayNameCondition returns a filter on the display name field.
	DisplayNameCondition(op database.TextOperation, displayName string) database.Condition
	// RoleGroupCondition returns a filter on the role group field.
	RoleGroupCondition(op database.TextOperation, roleGroup string) database.Condition
}

type projectRoleChanges interface {
	// SetUpdatedAt sets the updated at column.
	// Only use this when reducing events,
	// during regular updates the DB sets this column automatically.
	SetUpdatedAt(updatedAt time.Time) database.Change
	// SetDisplayName sets the display name column.
	SetDisplayName(displayName string) database.Change
	// SetRoleGroup sets the role group column.
	SetRoleGroup(roleGroup string) database.Change
}

type ProjectRoleRepository interface {
	projectRoleColumns
	projectRoleConditions
	projectRoleChanges

	// Get a single project role. An error is returned if not exactly one project role is found.
	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*ProjectRole, error)
	// List project roles. An empty list is returned if no project roles are found.
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*ProjectRole, error)
	// Create a new project role.
	Create(ctx context.Context, client database.QueryExecutor, role *ProjectRole) error
	// Update an existing project role.
	// The condition must include the instanceID, projectID and key of the project role to update.
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	// Delete an existing project role.
	// The condition must include the instanceID, projectID and key of the project role to delete.
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}
