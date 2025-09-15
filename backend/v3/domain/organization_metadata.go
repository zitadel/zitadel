package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type OrganizationMetadata struct {
	Metadata
	OrgID string `json:"orgId,omitempty" db:"org_id"`
}

type organizationMetadataColumns interface {
	MetadataColumns
	// OrgIDColumn returns the column for the org id field.
	OrgIDColumn() database.Column
}

type organizationMetadataConditions interface {
	MetadataConditions
	// OrgIDCondition returns a filter on the org id field.
	OrgIDCondition(orgID string) database.Condition
}

type OrganizationMetadataRepository interface {
	organizationMetadataColumns
	organizationMetadataConditions

	// Get returns one metadata based on the criteria.
	// If none is found, it returns an error of type [database.ErrNotFound].
	// If multiple were found, it returns an error of type [database.ErrMultipleRows].
	Get(ctx context.Context, opts ...database.QueryOption) (*OrganizationMetadata, error)
	// List returns a list of metadata based on the criteria.
	// If none are found, it returns an empty slice.
	List(ctx context.Context, opts ...database.QueryOption) ([]*OrganizationMetadata, error)

	// Set sets the given metadata for the organization.
	// If a metadata with the same key already exists, it will be overwritten, otherwise created.
	Set(ctx context.Context, metadata ...*OrganizationMetadata) error
	// Remove removes a metadata from the organization.
	Remove(ctx context.Context, condition database.Condition) (int64, error)
}
