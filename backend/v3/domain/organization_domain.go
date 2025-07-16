package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type OrganizationDomain struct {
	InstanceID       string                 `json:"instanceId,omitempty" db:"instance_id"`
	OrgID            string                 `json:"orgId,omitempty" db:"org_id"`
	Domain           string                 `json:"domain,omitempty" db:"domain"`
	IsVerified       bool                   `json:"isVerified,omitempty" db:"is_verified"`
	IsPrimary        bool                   `json:"isPrimary,omitempty" db:"is_primary"`
	VerificationType DomainVerificationType `json:"verificationType,omitempty" db:"verification_type"`

	CreatedAt string `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt string `json:"updatedAt,omitempty" db:"updated_at"`
}

type AddOrganizationDomain struct {
	InstanceID       string                 `json:"instanceId,omitempty" db:"instance_id"`
	OrgID            string                 `json:"orgId,omitempty" db:"org_id"`
	Domain           string                 `json:"domain,omitempty" db:"domain"`
	IsVerified       bool                   `json:"isVerified,omitempty" db:"is_verified"`
	IsPrimary        bool                   `json:"isPrimary,omitempty" db:"is_primary"`
	VerificationType DomainVerificationType `json:"verificationType,omitempty" db:"verification_type"`
}

type organizationDomainColumns interface {
	domainColumns
	// OrgIDColumn returns the column for the org id field.
	OrgIDColumn() database.Column
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

	// Add adds a new domain to the organization.
	Add(ctx context.Context, domain *AddOrganizationDomain) error
	// Update updates an existing domain in the organization.
	Update(ctx context.Context, condition database.Condition, changes ...database.Change) error
	// Remove removes a domain from the organization.
	Remove(ctx context.Context, condition database.Condition) error
}
