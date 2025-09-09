package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type InstanceDomain struct {
	InstanceID string `json:"instanceId,omitempty" db:"instance_id"`
	Domain     string `json:"domain,omitempty" db:"domain"`
	// IsPrimary indicates if the domain is the primary domain of the instance.
	// It is only set for custom domains.
	IsPrimary *bool `json:"isPrimary,omitempty" db:"is_primary"`
	// IsGenerated indicates if the domain is a generated domain.
	// It is only set for custom domains.
	IsGenerated *bool      `json:"isGenerated,omitempty" db:"is_generated"`
	Type        DomainType `json:"type,omitempty" db:"type"`

	CreatedAt time.Time `json:"createdAt,omitzero" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt,omitzero" db:"updated_at"`
}

type AddInstanceDomain struct {
	InstanceID  string     `json:"instanceId,omitempty" db:"instance_id"`
	Domain      string     `json:"domain,omitempty" db:"domain"`
	IsPrimary   *bool      `json:"isPrimary,omitempty" db:"is_primary"`
	IsGenerated *bool      `json:"isGenerated,omitempty" db:"is_generated"`
	Type        DomainType `json:"type,omitempty" db:"type"`

	// CreatedAt is the time when the domain was added.
	// It is set by the repository and should not be set by the caller.
	CreatedAt time.Time `json:"createdAt,omitzero" db:"created_at"`
	// UpdatedAt is the time when the domain was last updated.
	// It is set by the repository and should not be set by the caller.
	UpdatedAt time.Time `json:"updatedAt,omitzero" db:"updated_at"`
}

type instanceDomainColumns interface {
	domainColumns
	// IsGeneratedColumn returns the column for the is generated field.
	IsGeneratedColumn() database.Column
	// TypeColumn returns the column for the type field.
	TypeColumn() database.Column
}

type instanceDomainConditions interface {
	domainConditions
	// TypeCondition returns a filter for the type field.
	TypeCondition(typ DomainType) database.Condition
}

type instanceDomainChanges interface {
	domainChanges
	// SetType sets the type column.
	SetType(typ DomainType) database.Change
}

type InstanceDomainRepository interface {
	instanceDomainColumns
	instanceDomainConditions
	instanceDomainChanges

	// Get returns a single domain based on the criteria.
	// If no domain is found, it returns an error of type [database.ErrNotFound].
	// If multiple domains are found, it returns an error of type [database.ErrMultipleRows].
	Get(ctx context.Context, opts ...database.QueryOption) (*InstanceDomain, error)
	// List returns a list of domains based on the criteria.
	// If no domains are found, it returns an empty slice.
	List(ctx context.Context, opts ...database.QueryOption) ([]*InstanceDomain, error)

	// Add adds a new domain to the instance.
	Add(ctx context.Context, domain *AddInstanceDomain) error
	// Update updates an existing domain in the instance.
	Update(ctx context.Context, condition database.Condition, changes ...database.Change) (int64, error)
	// Remove removes a domain from the instance.
	Remove(ctx context.Context, condition database.Condition) (int64, error)
}
