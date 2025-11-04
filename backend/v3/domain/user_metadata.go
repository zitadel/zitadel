package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type UserMetadata struct {
	Metadata
	UserID string `json:"userId,omitempty" db:"user_id"`
}

type userMetadataColumns interface {
	MetadataColumns
	// UserIDColumn returns the column for the user id field.
	UserIDColumn() database.Column
}

type userMetadataConditions interface {
	MetadataConditions
	// PrimaryKeyCondition returns the condition for the primary key fields.
	PrimaryKeyCondition(instanceID, userID, key string) database.Condition
	// UserIDCondition returns a filter on the user id field.
	UserIDCondition(userID string) database.Condition
}

type UserMetadataRepository interface {
	userMetadataColumns
	userMetadataConditions

	// Get returns one metadata based on the criteria.
	// If none is found, it returns an error of type [database.ErrNotFound].
	// If multiple were found, it returns an error of type [database.ErrMultipleRows].
	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*UserMetadata, error)
	// List returns a list of metadata based on the criteria.
	// If none are found, it returns an empty slice.
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*UserMetadata, error)

	// Set sets the given metadata for the user.
	// If a metadata with the same key already exists, it will be overwritten, otherwise created.
	Set(ctx context.Context, client database.QueryExecutor, metadata ...*UserMetadata) error
	// Remove removes a metadata from the user.
	Remove(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}
