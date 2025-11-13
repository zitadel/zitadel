package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

//go:generate enumer -type  AuthorizationState -transform lower -trimprefix AuthorizationState -sql
type AuthorizationState uint8

const (
	AuthorizationStateActive AuthorizationState = iota
	AuthorizationStateInactive
)

type Authorization struct {
	ID         string             `json:"id,omitempty" db:"id"`
	UserID     string             `json:"userId" db:"user_id"`
	ProjectID  string             `json:"projectId" db:"project_id"`
	GrantID    string             `json:"grantId" db:"grant_id"`
	InstanceID string             `json:"instanceId,omitempty" db:"instance_id"`
	State      AuthorizationState `json:"state,omitempty" db:"state"`
	CreatedAt  time.Time          `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt  time.Time          `json:"updatedAt,omitempty" db:"updated_at"`
	Roles      []string           `json:"roles,omitempty" db:"roles"`
}

// authorizationColumns defines all the columns of the authorizations table.
type authorizationColumns interface {
	// PrimaryKeyColumns returns the columns for the primary key fields
	PrimaryKeyColumns() []database.Column
	// IDColumn returns the column for the id field.
	IDColumn() database.Column
	// UserIDColumn returns the column for the user_id field.
	UserIDColumn() database.Column
	// GrantIDColumn returns the column for the project grant_id field.
	GrantIDColumn() database.Column
	// ProjectIDColumn returns the column for the project_id field.
	ProjectIDColumn() database.Column
	// InstanceIDColumn returns the column for the instance_id field.
	InstanceIDColumn() database.Column
	// StateColumn returns the column for the state field.
	StateColumn() database.Column
	// CreatedAtColumn returns the column for the created_at field.
	CreatedAtColumn() database.Column
	// UpdatedAtColumn returns the column for the updated_at field.
	UpdatedAtColumn() database.Column
}

// authorizationConditions defines all the conditions for the authorizations table.
type authorizationConditions interface {
	// PrimaryKeyCondition returns a filter on the primary key fields.
	PrimaryKeyCondition(instanceID, id string) database.Condition
	// IDCondition returns an equal filter on the id field.
	IDCondition(orgID string) database.Condition
	// InstanceIDCondition returns a filter on the instance_id field.
	InstanceIDCondition(instanceID string) database.Condition
	// ProjectIDCondition returns a filter on the project_id field.
	ProjectIDCondition(projectID string) database.Condition
	// GrantIDCondition returns a filter on the grant_id field.
	GrantIDCondition(grantID string) database.Condition
	// UserIDCondition returns a filter on the user_id field.
	UserIDCondition(userID string) database.Condition
	// RolesCondition returns a filter on the roles field.
	RolesCondition(op database.TextOperation, roles string) database.Condition
	// StateCondition returns a filter on the name field.
	StateCondition(state AuthorizationState) database.Condition
}

// authorizationChanges defines all the changes for the authorizations table.
type authorizationChanges interface {
	// SetUpdatedAt sets the updated at column.
	// Only use this when reducing events,
	// during regular updates the DB sets this column automatically.
	SetUpdatedAt(updatedAt time.Time) database.Change
	// SetState sets the state column.
	SetState(state AuthorizationState) database.Change
}

//go:generate mockgen -typed -package domainmock -destination ./mock/authorization.mock.go . AuthorizationRepository

// AuthorizationRepository is the interface for the authorization repository.
type AuthorizationRepository interface {
	Repository

	authorizationConditions
	authorizationChanges
	authorizationColumns

	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*Authorization, error)
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*Authorization, error)

	Create(ctx context.Context, client database.QueryExecutor, authorization *Authorization) error
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, roles []string, changes ...database.Change) (int64, error)
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}
