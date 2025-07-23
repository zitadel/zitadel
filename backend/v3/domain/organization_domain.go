package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type OrganizationDomain struct {
	InstanceID     string               `json:"instanceId,omitempty" db:"instance_id"`
	OrgID          string               `json:"orgId,omitempty" db:"org_id"`
	Domain         string               `json:"domain,omitempty" db:"domain"`
	IsVerified     bool                 `json:"isVerified,omitempty" db:"is_verified"`
	IsPrimary      bool                 `json:"isPrimary,omitempty" db:"is_primary"`
	ValidationType DomainValidationType `json:"validationType,omitempty" db:"validation_type"`

	CreatedAt time.Time `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt,omitzero" db:"updated_at"`
}

type AddOrganizationDomain struct {
	InstanceID     string               `json:"instanceId,omitempty" db:"instance_id"`
	OrgID          string               `json:"orgId,omitempty" db:"org_id"`
	Domain         string               `json:"domain,omitempty" db:"domain"`
	IsVerified     bool                 `json:"isVerified,omitempty" db:"is_verified"`
	IsPrimary      bool                 `json:"isPrimary,omitempty" db:"is_primary"`
	ValidationType DomainValidationType `json:"validationType,omitempty" db:"validation_type"`

	// CreatedAt is the time when the domain was added.
	// It is set by the repository and should not be set by the caller.
	CreatedAt time.Time `json:"createdAt,omitzero" db:"created_at"`
	// UpdatedAt is the time when the domain was added.
	// It is set by the repository and should not be set by the caller.
	UpdatedAt time.Time `json:"updatedAt,omitzero" db:"updated_at"`
}

type organizationDomainColumns interface {
	domainColumns
	// OrgIDColumn returns the column for the org id field.
	OrgIDColumn(qualified bool) database.Column
}

type organizationDomainConditions interface {
	domainConditions
	// OrgIDCondition returns a filter on the org id field.
	OrgIDCondition(orgID string) database.Condition
}

type organizationDomainChanges interface {
	domainChanges
}

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
