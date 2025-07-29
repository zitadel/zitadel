package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type OrganizationDomainRepository interface {
	organizationDomainColumns
	organizationDomainConditions
	organizationDomainChanges

	// Get returns a single domain based on the criteria.
	// If no domain is found, it returns an error of type [database.ErrNotFound].
	// If multiple domains are found, it returns an error of type [database.ErrMultipleRows].
	Get(ctx context.Context, opts ...database.QueryOption) (*OrganizationDomain, error)
	// List returns a list of domains based on the criteria.
	// If no domains are found, it returns an empty slice.
	List(ctx context.Context, opts ...database.QueryOption) ([]*OrganizationDomain, error)

	// Add adds a new domain to the organization.
	Add(ctx context.Context, domain *AddOrganizationDomain) error
	// Update updates an existing domain in the organization.
	Update(ctx context.Context, condition database.Condition, changes ...database.Change) (int64, error)
	// Remove removes a domain from the organization.
	Remove(ctx context.Context, condition database.Condition) (int64, error)
}

type organizationDomainColumns interface {
	domainColumns
	// OrgIDColumn returns the column for the org id field.
	// `qualified` indicates if the column should be qualified with the table name.
	OrgIDColumn(qualified bool) database.Column
	// IsVerifiedColumn returns the column for the is verified field.
	// `qualified` indicates if the column should be qualified with the table name.
	IsVerifiedColumn(qualified bool) database.Column
	// ValidationTypeColumn returns the column for the verification type field.
	// `qualified` indicates if the column should be qualified with the table name.
	ValidationTypeColumn(qualified bool) database.Column
}

type organizationDomainConditions interface {
	domainConditions
	// OrgIDCondition returns a filter on the org id field.
	OrgIDCondition(orgID string) database.Condition
	// IsVerifiedCondition returns a filter on the is verified field.
	IsVerifiedCondition(isVerified bool) database.Condition
}

type organizationDomainChanges interface {
	domainChanges
	// SetVerified sets the is verified column to true.
	SetVerified() database.Change
	// SetValidationType sets the verification type column.
	// If the domain is already verified, this is a no-op.
	SetValidationType(verificationType DomainValidationType) database.Change
}
