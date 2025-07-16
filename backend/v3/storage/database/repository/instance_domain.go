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
	var builder database.StatementBuilder

	builder.WriteString(`INSERT INTO zitadel.instance_domains (instance_id, domain, is_verified, is_primary, verification_type) ` +
		`VALUES ($1, $2, $3, $4, $5)` +
		` RETURNING created_at, updated_at`)

	builder.AppendArgs(domain.InstanceID, domain.Domain, domain.IsVerified, domain.IsPrimary, domain.VerificationType)

	return i.client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&domain.CreatedAt, &domain.UpdatedAt)
}

// Remove implements [domain.InstanceDomainRepository].
func (i *instanceDomain) Remove(ctx context.Context, condition database.Condition) (int64, error) {
	var builder database.StatementBuilder

	builder.WriteString(`DELETE FROM zitadel.instance_domains WHERE `)
	writeCondition(&builder, condition)

	return i.client.Exec(ctx, builder.String(), builder.Args()...)
}

// Update implements [domain.InstanceDomainRepository].
// Subtle: this method shadows the method (instance).Update of instanceDomain.instance.
func (i *instanceDomain) Update(ctx context.Context, condition database.Condition, changes ...database.Change) (int64, error) {
	var builder database.StatementBuilder

	builder.WriteString(`UPDATE zitadel.instance_domains SET `)
	database.Changes(changes).Write(&builder)

	writeCondition(&builder, condition)

	return i.client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetVerificationType implements [domain.InstanceDomainRepository].
func (i instanceDomain) SetVerificationType(verificationType domain.DomainVerificationType) database.Change {
	return database.NewChange(i.VerificationTypeColumn(), verificationType)
}

// SetPrimary implements [domain.InstanceDomainRepository].
func (i instanceDomain) SetPrimary() database.Change {
	return database.NewChange(i.IsPrimaryColumn(), true)
}

// SetVerified implements [domain.InstanceDomainRepository].
func (i instanceDomain) SetVerified() database.Change {
	return database.NewChange(i.IsVerifiedColumn(), true)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// DomainCondition implements [domain.InstanceDomainRepository].
func (i instanceDomain) DomainCondition(op database.TextOperation, domain string) database.Condition {
	return database.NewTextCondition(i.DomainColumn(), op, domain)
}

// InstanceIDCondition implements [domain.InstanceDomainRepository].
func (i instanceDomain) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(i.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// IsPrimaryCondition implements [domain.InstanceDomainRepository].
func (i instanceDomain) IsPrimaryCondition(isPrimary bool) database.Condition {
	return database.NewBooleanCondition(i.IsPrimaryColumn(), isPrimary)
}

// IsVerifiedCondition implements [domain.InstanceDomainRepository].
func (i instanceDomain) IsVerifiedCondition(isVerified bool) database.Condition {
	return database.NewBooleanCondition(i.IsVerifiedColumn(), isVerified)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// CreatedAtColumn implements [domain.InstanceDomainRepository].
// Subtle: this method shadows the method (instance).CreatedAtColumn of instanceDomain.instance.
func (instanceDomain) CreatedAtColumn() database.Column {
	return database.NewColumn("created_at")
}

// DomainColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) DomainColumn() database.Column {
	return database.NewColumn("domain")
}

// InstanceIDColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) InstanceIDColumn() database.Column {
	return database.NewColumn("instance_id")
}

// IsPrimaryColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) IsPrimaryColumn() database.Column {
	return database.NewColumn("is_primary")
}

// IsVerifiedColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) IsVerifiedColumn() database.Column {
	return database.NewColumn("is_verified")
}

// UpdatedAtColumn implements [domain.InstanceDomainRepository].
// Subtle: this method shadows the method (instance).UpdatedAtColumn of instanceDomain.instance.
func (instanceDomain) UpdatedAtColumn() database.Column {
	return database.NewColumn("updated_at")
}

// VerificationTypeColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) VerificationTypeColumn() database.Column {
	return database.NewColumn("verification_type")
}

// IsGeneratedColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) IsGeneratedColumn() database.Column {
	return database.NewColumn("is_generated")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------
