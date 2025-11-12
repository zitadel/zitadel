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
		s.InstanceIDColumn(),
		s.SessionIDColumn(),
		s.KeyColumn(),
	}
}

func (s sessionMetadata) InstanceIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "instance_id")
}

func (s sessionMetadata) SessionIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "session_id")
}

func (s sessionMetadata) KeyColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "key")
}

func (s sessionMetadata) CreatedAtColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "created_at")
}

func (s sessionMetadata) UpdatedAtColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "updated_at")
}

func (s sessionMetadata) ValueColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "value")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// PrimaryKeyCondition implements [domain.SessionMetadataRepository].
func (s sessionMetadata) PrimaryKeyCondition(instanceID string, orgID string, key string) database.Condition {
	return database.And(
		s.InstanceIDCondition(instanceID),
		s.SessionIDCondition(orgID),
		s.KeyCondition(database.TextOperationEqual, key),
	)
}

// InstanceIDCondition implements [domain.SessionMetadataRepository].
func (s sessionMetadata) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(s.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// SessionIDCondition implements [domain.SessionMetadataRepository].
func (s sessionMetadata) SessionIDCondition(orgID string) database.Condition {
	return database.NewTextCondition(s.SessionIDColumn(), database.TextOperationEqual, orgID)
}

// KeyCondition implements [domain.SessionMetadataRepository].
func (s sessionMetadata) KeyCondition(op database.TextOperation, key string) database.Condition {
	return database.NewTextCondition(s.KeyColumn(), op, key)
}

// ValueCondition implements [domain.SessionMetadataRepository].
func (s sessionMetadata) ValueCondition(op database.BytesOperation, value []byte) database.Condition {
	return database.NewBytesCondition[[]byte](database.SHA256Column(s.ValueColumn()), op, database.SHA256Value(value))
}
