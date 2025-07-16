package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type InstanceDomain struct {
	InstanceID       string                 `json:"instanceId,omitempty" db:"instance_id"`
	Domain           string                 `json:"domain,omitempty" db:"domain"`
	IsVerified       bool                   `json:"isVerified,omitempty" db:"is_verified"`
	IsPrimary        bool                   `json:"isPrimary,omitempty" db:"is_primary"`
	VerificationType DomainVerificationType `json:"verificationType,omitempty" db:"verification_type"`

	CreatedAt string `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt string `json:"updatedAt,omitempty" db:"updated_at"`
}

type AddInstanceDomain struct {
	InstanceID       string                 `json:"instanceId,omitempty" db:"instance_id"`
	Domain           string                 `json:"domain,omitempty" db:"domain"`
	IsVerified       bool                   `json:"isVerified,omitempty" db:"is_verified"`
	IsPrimary        bool                   `json:"isPrimary,omitempty" db:"is_primary"`
	VerificationType DomainVerificationType `json:"verificationType,omitempty" db:"verification_type"`
}

type instanceDomainColumns interface {
	domainColumns
	// IsGeneratedColumn returns the column for the is generated field.
	IsGeneratedColumn() database.Column
}

type instanceDomainConditions interface {
	domainConditions
}

type instanceDomainChanges interface {
	domainChanges
}

type InstanceDomainRepository interface {
	instanceDomainColumns
	instanceDomainConditions
	instanceDomainChanges

	// Add adds a new domain to the instance.
	Add(ctx context.Context, domain *AddInstanceDomain) error
	// Update updates an existing domain in the instance.
	Update(ctx context.Context, condition database.Condition, changes ...database.Change) error
	// Remove removes a domain from the instance.
	Remove(ctx context.Context, condition database.Condition) error
}
