package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/database"
)

// Domain represents a unified domain that can belong to either an instance or an organization
type Domain struct {
	ID             string
	InstanceID     string
	OrganizationID *string                          // nil for instance domains
	Domain         string
	IsVerified     bool
	IsPrimary      bool
	ValidationType *domain.OrgDomainValidationType // nil for instance domains
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

// IsInstanceDomain returns true if this is an instance domain (OrganizationID is nil)
func (d *Domain) IsInstanceDomain() bool {
	return d.OrganizationID == nil
}

// IsOrganizationDomain returns true if this is an organization domain (OrganizationID is not nil)
func (d *Domain) IsOrganizationDomain() bool {
	return d.OrganizationID != nil
}

// DomainSearchCriteria defines the search criteria for domains
type DomainSearchCriteria struct {
	ID             *string
	Domain         *string
	InstanceID     *string
	OrganizationID *string
	IsVerified     *bool
	IsPrimary      *bool
}

// DomainPagination defines pagination options for domain queries
type DomainPagination struct {
	Limit  uint32
	Offset uint32
	SortBy DomainSortField
	Order  database.SortOrder
}

// DomainSortField defines the fields available for sorting domain results
type DomainSortField int

const (
	DomainSortFieldCreatedAt DomainSortField = iota
	DomainSortFieldUpdatedAt
	DomainSortFieldDomain
)

// DomainList represents a paginated list of domains
type DomainList struct {
	Domains    []*Domain
	TotalCount uint64
}

// InstanceDomainRepository defines the repository interface for instance domain operations
type InstanceDomainRepository interface {
	// Add creates a new instance domain (always verified)
	Add(ctx context.Context, instanceID, domain string) (*Domain, error)
	
	// SetPrimary sets the primary domain for an instance
	SetPrimary(ctx context.Context, instanceID, domain string) error
	
	// Remove soft deletes an instance domain
	Remove(ctx context.Context, instanceID, domain string) error
	
	// Get returns a single instance domain matching the criteria
	// Returns error if multiple domains are found
	Get(ctx context.Context, criteria DomainSearchCriteria) (*Domain, error)
	
	// List returns a list of instance domains matching the criteria with pagination
	List(ctx context.Context, criteria DomainSearchCriteria, pagination DomainPagination) (*DomainList, error)
}

// OrganizationDomainRepository defines the repository interface for organization domain operations
type OrganizationDomainRepository interface {
	// Add creates a new organization domain
	Add(ctx context.Context, instanceID, organizationID, domain string, validationType domain.OrgDomainValidationType) (*Domain, error)
	
	// SetVerified marks an organization domain as verified
	SetVerified(ctx context.Context, instanceID, organizationID, domain string) error
	
	// SetPrimary sets the primary domain for an organization
	SetPrimary(ctx context.Context, instanceID, organizationID, domain string) error
	
	// Remove soft deletes an organization domain
	Remove(ctx context.Context, instanceID, organizationID, domain string) error
	
	// Get returns a single organization domain matching the criteria
	// Returns error if multiple domains are found
	Get(ctx context.Context, criteria DomainSearchCriteria) (*Domain, error)
	
	// List returns a list of organization domains matching the criteria with pagination
	List(ctx context.Context, criteria DomainSearchCriteria, pagination DomainPagination) (*DomainList, error)
}