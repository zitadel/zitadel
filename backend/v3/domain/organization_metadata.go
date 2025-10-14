package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type OrganizationMetadata struct {
	Metadata
	OrganizationID string `json:"orgId,omitempty" db:"organization_id"`
}

type organizationMetadataColumns interface {
	MetadataColumns
	// OrganizationIDColumn returns the column for the org id field.
	OrganizationIDColumn() database.Column
}

type organizationMetadataConditions interface {
	MetadataConditions
	// PrimaryKeyCondition returns the condition for the primary key fields.
	PrimaryKeyCondition(instanceID, orgID, key string) database.Condition
	// OrganizationIDCondition returns a filter on the org id field.
	OrganizationIDCondition(orgID string) database.Condition
}

type OrganizationMetadataRepository interface {
	organizationMetadataColumns
	organizationMetadataConditions

	// Get returns one metadata based on the criteria.
	// If none is found, it returns an error of type [database.ErrNotFound].
	// If multiple were found, it returns an error of type [database.ErrMultipleRows].
	Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*OrganizationMetadata, error)
	// List returns a list of metadata based on the criteria.
	// If none are found, it returns an empty slice.
	List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*OrganizationMetadata, error)

	// Set sets the given metadata for the organization.
	// If a metadata with the same key already exists, it will be overwritten, otherwise created.
	Set(ctx context.Context, client database.QueryExecutor, metadata ...*OrganizationMetadata) error
	// Remove removes a metadata from the organization.
	Remove(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error)
}
