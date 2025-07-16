package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.OrganizationDomainRepository = (*orgDomain)(nil)

type orgDomain struct {
	repository
	*org
}

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

// Add implements [domain.OrganizationDomainRepository].
func (o *orgDomain) Add(ctx context.Context, domain *domain.AddOrganizationDomain) error {
	panic("unimplemented")
}

// Update implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method (*org).Update of orgDomain.org.
func (o *orgDomain) Update(ctx context.Context, condition database.Condition, changes ...database.Change) error {
	panic("unimplemented")
}

// Remove implements [domain.OrganizationDomainRepository].
func (o *orgDomain) Remove(ctx context.Context, condition database.Condition) error {
	panic("unimplemented")
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetPrimary implements [domain.OrganizationDomainRepository].
func (o *orgDomain) SetPrimary() database.Change {
	panic("unimplemented")
}

// SetVerificationType implements [domain.OrganizationDomainRepository].
func (o *orgDomain) SetVerificationType(verificationType domain.DomainVerificationType) database.Change {
	panic("unimplemented")
}

// SetVerified implements [domain.OrganizationDomainRepository].
func (o *orgDomain) SetVerified() database.Change {
	panic("unimplemented")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// DomainCondition implements [domain.OrganizationDomainRepository].
func (o *orgDomain) DomainCondition(op database.TextOperation, domain string) database.Condition {
	panic("unimplemented")
}

// InstanceIDCondition implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method (*org).InstanceIDCondition of orgDomain.org.
func (o *orgDomain) InstanceIDCondition(instanceID string) database.Condition {
	panic("unimplemented")
}

// IsPrimaryCondition implements [domain.OrganizationDomainRepository].
func (o *orgDomain) IsPrimaryCondition(isPrimary bool) database.Condition {
	panic("unimplemented")
}

// IsVerifiedCondition implements [domain.OrganizationDomainRepository].
func (o *orgDomain) IsVerifiedCondition(isVerified bool) database.Condition {
	panic("unimplemented")
}

// OrgIDCondition implements [domain.OrganizationDomainRepository].
func (o *orgDomain) OrgIDCondition(orgID string) database.Condition {
	panic("unimplemented")
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// CreatedAtColumn implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method (*org).CreatedAtColumn of orgDomain.org.
func (o *orgDomain) CreatedAtColumn() database.Column {
	panic("unimplemented")
}

// DomainColumn implements [domain.OrganizationDomainRepository].
func (o *orgDomain) DomainColumn() database.Column {
	panic("unimplemented")
}

// InstanceIDColumn implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method (*org).InstanceIDColumn of orgDomain.org.
func (o *orgDomain) InstanceIDColumn() database.Column {
	panic("unimplemented")
}

// IsPrimaryColumn implements [domain.OrganizationDomainRepository].
func (o *orgDomain) IsPrimaryColumn() database.Column {
	panic("unimplemented")
}

// IsVerifiedColumn implements [domain.OrganizationDomainRepository].
func (o *orgDomain) IsVerifiedColumn() database.Column {
	panic("unimplemented")
}

// OrgIDColumn implements [domain.OrganizationDomainRepository].
func (o *orgDomain) OrgIDColumn() database.Column {
	panic("unimplemented")
}

// UpdatedAtColumn implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method (*org).UpdatedAtColumn of orgDomain.org.
func (o *orgDomain) UpdatedAtColumn() database.Column {
	panic("unimplemented")
}

// VerificationTypeColumn implements [domain.OrganizationDomainRepository].
func (o *orgDomain) VerificationTypeColumn() database.Column {
	panic("unimplemented")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------
