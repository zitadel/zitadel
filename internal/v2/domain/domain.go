package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// Domain represents a unified domain entry for both organization and instance domains
type Domain struct {
	ID             string
	InstanceID     string
	OrganizationID *string // nil for instance domains
	Domain         string
	IsVerified     bool
	IsPrimary      bool
	ValidationType domain.OrgDomainValidationType
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

// IsInstanceDomain returns true if this is an instance domain (org_id is nil)
func (d *Domain) IsInstanceDomain() bool {
	return d.OrganizationID == nil
}

// IsOrganizationDomain returns true if this is an organization domain (org_id is not nil)
func (d *Domain) IsOrganizationDomain() bool {
	return d.OrganizationID != nil
}

// DomainSearchCriteria defines criteria for searching domains
type DomainSearchCriteria struct {
	ID             *string
	InstanceID     *string
	OrganizationID *string
	Domain         *string
	IsVerified     *bool
	IsPrimary      *bool
}

// DomainPagination defines pagination options for domain queries
type DomainPagination struct {
	Limit  int
	Offset int
	// SortBy can be "created_at", "updated_at", or "domain"
	SortBy string
	// SortOrder can be "ASC" or "DESC"
	SortOrder string
}

// InstanceDomainRepository defines the repository interface for instance domains
type InstanceDomainRepository interface {
	// Add creates a new instance domain (always verified)
	Add(ctx context.Context, instanceID, domain string) (*Domain, error)
	
	// SetPrimary sets a domain as primary for an instance
	SetPrimary(ctx context.Context, instanceID, domain string) error
	
	// Remove removes an instance domain
	Remove(ctx context.Context, instanceID, domain string) error
	
	// Get returns a single instance domain matching criteria
	Get(ctx context.Context, criteria DomainSearchCriteria) (*Domain, error)
	
	// List returns instance domains matching criteria with pagination
	List(ctx context.Context, criteria DomainSearchCriteria, pagination DomainPagination) ([]*Domain, int64, error)
}

// OrganizationDomainRepository defines the repository interface for organization domains
type OrganizationDomainRepository interface {
	// Add creates a new organization domain
	Add(ctx context.Context, instanceID, organizationID, domain string, validationType domain.OrgDomainValidationType) (*Domain, error)
	
	// SetVerified marks a domain as verified
	SetVerified(ctx context.Context, instanceID, organizationID, domain string) error
	
	// SetPrimary sets a domain as primary for an organization
	SetPrimary(ctx context.Context, instanceID, organizationID, domain string) error
	
	// Remove removes an organization domain
	Remove(ctx context.Context, instanceID, organizationID, domain string) error
	
	// Get returns a single organization domain matching criteria
	Get(ctx context.Context, criteria DomainSearchCriteria) (*Domain, error)
	
	// List returns organization domains matching criteria with pagination
	List(ctx context.Context, criteria DomainSearchCriteria, pagination DomainPagination) ([]*Domain, int64, error)
}

// DomainErrors
var (
	ErrDomainNotFound      = zerrors.ThrowNotFound(nil, "DOMAIN-3n8sd", "Errors.Domain.NotFound")
	ErrDomainAlreadyExists = zerrors.ThrowAlreadyExists(nil, "DOMAIN-2n8sd", "Errors.Domain.AlreadyExists")
	ErrMultipleDomainsFound = zerrors.ThrowInvalidArgument(nil, "DOMAIN-4n8sd", "Errors.Domain.MultipleFound")
)