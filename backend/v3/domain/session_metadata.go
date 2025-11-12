package domain

import "github.com/zitadel/zitadel/backend/v3/storage/database"

type SessionMetadata struct {
	Metadata
	SessionID string `json:"sessionId,omitempty" db:"session_id"`
}

type sessionMetadataColumns interface {
	MetadataColumns
	// SessionIDColumn returns the column for the session id field.
	SessionIDColumn() database.Column
}

type sessionMetadataConditions interface {
	MetadataConditions
	// PrimaryKeyCondition returns a filter on the primary key fields.
	PrimaryKeyCondition(instanceID, sessionID, key string) database.Condition
	// SessionIDCondition returns an equal filter on the session id field.
	SessionIDCondition(sessionID string) database.Condition
}
