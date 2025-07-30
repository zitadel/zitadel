package domain

import (
	"time"
)

type OrganizationDomain struct {
	InstanceID     string                `json:"instanceId,omitempty" db:"instance_id"`
	OrgID          string                `json:"orgId,omitempty" db:"org_id"`
	Domain         string                `json:"domain,omitempty" db:"domain"`
	IsVerified     bool                  `json:"isVerified,omitempty" db:"is_verified"`
	IsPrimary      bool                  `json:"isPrimary,omitempty" db:"is_primary"`
	ValidationType *DomainValidationType `json:"validationType,omitempty" db:"validation_type"`

	CreatedAt time.Time `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt,omitzero" db:"updated_at"`
}

type AddOrganizationDomain struct {
	InstanceID     string                `json:"instanceId,omitempty" db:"instance_id"`
	OrgID          string                `json:"orgId,omitempty" db:"org_id"`
	Domain         string                `json:"domain,omitempty" db:"domain"`
	IsVerified     bool                  `json:"isVerified,omitempty" db:"is_verified"`
	IsPrimary      bool                  `json:"isPrimary,omitempty" db:"is_primary"`
	ValidationType *DomainValidationType `json:"validationType,omitempty" db:"validation_type"`

	// CreatedAt is the time when the domain was added.
	// It is set by the repository and should not be set by the caller.
	CreatedAt time.Time `json:"createdAt,omitzero" db:"created_at"`
	// UpdatedAt is the time when the domain was added.
	// It is set by the repository and should not be set by the caller.
	UpdatedAt time.Time `json:"updatedAt,omitzero" db:"updated_at"`
}
