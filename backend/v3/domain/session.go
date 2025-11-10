package domain

import (
	"time"

	"golang.org/x/net/context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Session struct {
	InstanceID         string
	ID                 string
	Token              string
	UserAgent          *SessionUserAgent
	Lifetime           time.Duration
	Expiration         time.Time
	UserID             string
	IdentityProviderID string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Metadata           []SessionMetadata
	Factors            SessionFactors
}

type SessionRepository interface {
	Repository

	sessionColumns
	sessionConditions
	sessionChanges

	// Get returns a session based on the given condition.
	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*Session, error)
	// List returns a list of sessions based on the given condition.
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*Session, error)
	// Create creates a new session.
	Create(ctx context.Context, client database.QueryExecutor, user *Session) error
	// Update an existing session.
	// The condition must include the instanceID and ID of the session to update.
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	// Delete removes sessions based on the given condition.
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) error

	// SetChallenge adds or updates the challenge of the session matching the condition.
	SetChallenge(ctx context.Context, client database.QueryExecutor, condition database.Condition, challenge SessionChallengeType) error
	// SetFactor adds or updates the factor of the session matching the condition.
	SetFactor(ctx context.Context, client database.QueryExecutor, condition database.Condition, factor SessionFactorType) error
	// ClearFactor resets the factor of the session matching the condition.
	ClearFactor(ctx context.Context, client database.QueryExecutor, condition database.Condition) error
	// SetUserAgent adds or updates the user agent of the session matching the condition.
	SetUserAgent(ctx context.Context, client database.QueryExecutor, condition database.Condition, userAgent SessionUserAgent) error
	// SetMetadata adds or updates the metadata of the session matching the condition.
	SetMetadata(ctx context.Context, client database.QueryExecutor, condition database.Condition, metadata []SessionMetadata) error
}

// sessionColumns define all the columns of the session table.
type sessionColumns interface {
	// InstanceIDColumn returns the column for the instance id field.
	InstanceIDColumn() database.Column
	// IDColumn returns the column for the id field.
	IDColumn() database.Column
	// TokenColumn returns the column for the token field.
	TokenColumn() database.Column
	// UserAgentIDColumn returns the column for the user agent id field.
	UserAgentIDColumn() database.Column
	// LifetimeColumn returns the column for the lifetime field.
	LifetimeColumn() database.Column
	// ExpirationColumn returns the column for the expiration field.
	ExpirationColumn() database.Column
	// UserIDColumn returns the column for the expiration field.
	UserIDColumn() database.Column
	// IdentityProviderIDColumn returns the column for the identity provider ID field.
	IdentityProviderIDColumn() database.Column
	// CreatedAtColumn returns the column for the created at field.
	CreatedAtColumn() database.Column
	// UpdatedAtColumn returns the column for the updated at field.
	UpdatedAtColumn() database.Column
}

// sessionConditions define all the conditions for the session table.
type sessionConditions interface {
	// PrimaryKeyCondition returns a filter on the primary key fields.
	PrimaryKeyCondition(instanceID, sessionID string) database.Condition
	// InstanceIDCondition returns an equal filter on the instance id field.
	InstanceIDCondition(instanceID string) database.Condition
	// IDCondition returns an equal filter on the id field.
	IDCondition(sessionID string) database.Condition
	// UserAgentIDCondition returns an equal filter on the user agent ID field.
	UserAgentIDCondition(userAgentID string) database.Condition
	// UserIDCondition returns an equal filter on the user id field.
	UserIDCondition(userID string) database.Condition
	// CreatedAtCondition returns a filter on the created at field.
	CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition
	// UpdatedAtCondition returns a filter on the updated at field.
	UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition

	FactorConditions() sessionFactorConditions
}

type sessionChanges interface {
	// SetUpdatedAt sets the updated at column.
	// Only use this when reducing events,
	// during regular updates the DB sets this column automatically.
	SetUpdatedAt(updatedAt time.Time) database.Change
	// SetToken sets the token field of the session.
	SetToken(token string) database.Change
	// SetLifetime sets the lifetime field of the session and will update the computed expiration field.
	SetLifetime(lifetime time.Duration) database.Change
}
