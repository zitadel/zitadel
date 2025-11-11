package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Session struct {
	InstanceID string            `json:"instance_id,omitempty" db:"instance_id"`
	ID         string            `json:"id,omitempty" db:"id"`
	Token      string            `json:"token,omitempty" db:"token"`
	Lifetime   time.Duration     `json:"lifetime,omitempty" db:"lifetime"`
	Expiration time.Time         `json:"expiration,omitzero" db:"expiration"`
	UserID     string            `json:"user_id,omitempty" db:"user_id"`
	CreatedAt  time.Time         `json:"created_at,omitzero" db:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at,omitzero" db:"updated_at"`
	Factors    SessionFactors    `json:"factors,omitempty"`
	Challenges SessionChallenges `json:"challenges,omitempty"`
	Metadata   []SessionMetadata `json:"metadata,omitempty"`
	UserAgent  *SessionUserAgent `json:"user_agent,omitempty"`
}

//go:generate mockgen -typed -package domainmock -destination ./mock/session.mock.go . SessionRepository

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
	// Update one or more existing sessions.
	// The condition must include at least the instanceID of the session to update.
	Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error)
	// Delete removes sessions based on the given condition.
	Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) error
}

// sessionColumns define all the columns of the session table.
type sessionColumns interface {
	// InstanceIDColumn returns the column for the instance id field.
	InstanceIDColumn() database.Column
	// IDColumn returns the column for the id field.
	IDColumn() database.Column
	// TokenColumn returns the column for the token field.
	TokenColumn() database.Column
	// LifetimeColumn returns the column for the lifetime field.
	LifetimeColumn() database.Column
	// ExpirationColumn returns the column for the expiration field.
	ExpirationColumn() database.Column
	// UserIDColumn returns the column for the user id field.
	UserIDColumn() database.Column
	// CreatedAtColumn returns the column for the created at field.
	CreatedAtColumn() database.Column
	// UpdatedAtColumn returns the column for the updated at field.
	UpdatedAtColumn() database.Column

	// FactorColumns returns the columns for the factors fields.
	FactorColumns() sessionFactorColumns
	// MetadataColumns returns the columns for the metadata fields.
	MetadataColumns() sessionMetadataColumns
	// UserAgentColumns returns the columns for the user agent fields.
	UserAgentColumns() sessionUserAgentColumns
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

	// FactorConditions returns the conditions for the factors fields.
	FactorConditions() sessionFactorConditions
	// MetadataConditions returns the conditions for the metadata fields.
	MetadataConditions() sessionMetadataConditions
	// UserAgentConditions returns the conditions for the user agent fields.
	UserAgentConditions() sessionUserAgentConditions
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

	// SetChallenge adds or updates the challenge of the corresponding type.
	SetChallenge(challenge SessionChallenge) database.Change
	// SetFactor adds or updates the factor of the corresponding type.
	SetFactor(factor SessionFactor) database.Change
	// ClearFactor resets the factor of the corresponding type.
	ClearFactor(factor SessionFactorType) database.Change
	// SetMetadata adds or updates the metadata of the session.
	SetMetadata(metadata []SessionMetadata) database.Change
	// SetUserAgent adds or updates the user agent of the session.
	SetUserAgent(userAgent SessionUserAgent) database.Change
}
