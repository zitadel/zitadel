package domain

import "github.com/zitadel/zitadel/backend/v3/storage/database"

type SessionMetadata struct {
	Metadata
	SessionID string `json:"sessionId,omitempty" db:"session_id"`
}

type SessionMetadataColumns interface {
	MetadataColumns
	// SessionIDColumn returns the column for the session id field.
	SessionIDColumn() database.Column
}

type SessionMetadataConditions interface {
	MetadataConditions
	// PrimaryKeyCondition returns a filter on the primary key fields.
	PrimaryKeyCondition(instanceID, sessionID, key string) database.Condition
	// SessionIDCondition returns an equal filter on the session id field.
	SessionIDCondition(sessionID string) database.Condition
}
