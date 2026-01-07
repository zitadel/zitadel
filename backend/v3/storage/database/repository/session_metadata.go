package repository

import "github.com/zitadel/zitadel/backend/v3/storage/database"

type sessionMetadata struct {
}

func (s sessionMetadata) qualifiedTableName() string {
	return "zitadel.session_metadata"
}

func (s sessionMetadata) unqualifiedTableName() string {
	return "session_metadata"
}

func (s sessionMetadata) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		s.instanceIDColumn(),
		s.sessionIDColumn(),
		s.keyColumn(),
	}
}

func (s sessionMetadata) instanceIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "instance_id")
}

func (s sessionMetadata) sessionIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "session_id")
}

func (s sessionMetadata) keyColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "key")
}

func (s sessionMetadata) valueColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "value")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// PrimaryKeyCondition implements [domain.SessionMetadataRepository].
func (s sessionMetadata) PrimaryKeyCondition(instanceID string, sessionID string, key string) database.Condition {
	return database.And(
		s.InstanceIDCondition(instanceID),
		s.SessionIDCondition(sessionID),
		s.KeyCondition(database.TextOperationEqual, key),
	)
}

// InstanceIDCondition implements [domain.SessionMetadataRepository].
func (s sessionMetadata) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(s.instanceIDColumn(), database.TextOperationEqual, instanceID)
}

// SessionIDCondition implements [domain.SessionMetadataRepository].
func (s sessionMetadata) SessionIDCondition(sessionID string) database.Condition {
	return database.NewTextCondition(s.sessionIDColumn(), database.TextOperationEqual, sessionID)
}

// KeyCondition implements [domain.SessionMetadataRepository].
func (s sessionMetadata) KeyCondition(op database.TextOperation, key string) database.Condition {
	return database.NewTextCondition(s.keyColumn(), op, key)
}

// ValueCondition implements [domain.SessionMetadataRepository].
func (s sessionMetadata) ValueCondition(op database.BytesOperation, value []byte) database.Condition {
	return database.NewBytesCondition[[]byte](database.SHA256Column(s.valueColumn()), op, database.SHA256Value(value))
}
