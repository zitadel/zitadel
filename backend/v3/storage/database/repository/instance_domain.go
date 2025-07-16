package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.InstanceDomainRepository = (*instanceDomain)(nil)

type instanceDomain struct {
	repository
	*instance
}

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

// Add implements [domain.InstanceDomainRepository].
func (i *instanceDomain) Add(ctx context.Context, domain *domain.AddInstanceDomain) error {
	panic("unimplemented")
}

// Remove implements [domain.InstanceDomainRepository].
func (i *instanceDomain) Remove(ctx context.Context, condition database.Condition) error {
	panic("unimplemented")
}

// Update implements [domain.InstanceDomainRepository].
// Subtle: this method shadows the method (instance).Update of instanceDomain.instance.
func (i *instanceDomain) Update(ctx context.Context, condition database.Condition, changes ...database.Change) error {
	panic("unimplemented")
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetVerificationType implements [domain.InstanceDomainRepository].
func (i *instanceDomain) SetVerificationType(verificationType domain.DomainVerificationType) database.Change {
	panic("unimplemented")
}

// SetPrimary implements [domain.InstanceDomainRepository].
func (i *instanceDomain) SetPrimary() database.Change {
	panic("unimplemented")
}

// SetVerified implements [domain.InstanceDomainRepository].
func (i *instanceDomain) SetVerified() database.Change {
	panic("unimplemented")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// DomainCondition implements [domain.InstanceDomainRepository].
func (i *instanceDomain) DomainCondition(op database.TextOperation, domain string) database.Condition {
	panic("unimplemented")
}

// InstanceIDCondition implements [domain.InstanceDomainRepository].
func (i *instanceDomain) InstanceIDCondition(instanceID string) database.Condition {
	panic("unimplemented")
}

// IsPrimaryCondition implements [domain.InstanceDomainRepository].
func (i *instanceDomain) IsPrimaryCondition(isPrimary bool) database.Condition {
	panic("unimplemented")
}

// IsVerifiedCondition implements [domain.InstanceDomainRepository].
func (i *instanceDomain) IsVerifiedCondition(isVerified bool) database.Condition {
	panic("unimplemented")
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// CreatedAtColumn implements [domain.InstanceDomainRepository].
// Subtle: this method shadows the method (instance).CreatedAtColumn of instanceDomain.instance.
func (i *instanceDomain) CreatedAtColumn() database.Column {
	panic("unimplemented")
}

// DomainColumn implements [domain.InstanceDomainRepository].
func (i *instanceDomain) DomainColumn() database.Column {
	panic("unimplemented")
}

// InstanceIDColumn implements [domain.InstanceDomainRepository].
func (i *instanceDomain) InstanceIDColumn() database.Column {
	panic("unimplemented")
}

// IsPrimaryColumn implements [domain.InstanceDomainRepository].
func (i *instanceDomain) IsPrimaryColumn() database.Column {
	panic("unimplemented")
}

// IsVerifiedColumn implements [domain.InstanceDomainRepository].
func (i *instanceDomain) IsVerifiedColumn() database.Column {
	panic("unimplemented")
}

// UpdatedAtColumn implements [domain.InstanceDomainRepository].
// Subtle: this method shadows the method (instance).UpdatedAtColumn of instanceDomain.instance.
func (i *instanceDomain) UpdatedAtColumn() database.Column {
	panic("unimplemented")
}

// VerificationTypeColumn implements [domain.InstanceDomainRepository].
func (i *instanceDomain) VerificationTypeColumn() database.Column {
	panic("unimplemented")
}

// IsGeneratedColumn implements [domain.InstanceDomainRepository].
func (i *instanceDomain) IsGeneratedColumn() database.Column {
	return database.NewColumn("is_generated")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------
