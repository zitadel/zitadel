package repository

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userMetadataConditions struct {
	user
}

// MetadataKeyCondition implements [domain.UserMetadataConditions].
func (u userMetadataConditions) MetadataKeyCondition(op database.TextOperation, key string) database.Condition {
	return database.NewTextCondition(u.metadataKeyColumn(), op, key)
}

// MetadataValueCondition implements [domain.UserMetadataConditions].
func (u userMetadataConditions) MetadataValueCondition(op database.BytesOperation, value []byte) database.Condition {
	return database.NewBytesCondition[[]byte](database.SHA256Column(u.metadataValueColumn()), op, database.SHA256Value(value))
}

var _ domain.UserMetadataConditions = (*userMetadataConditions)(nil)
