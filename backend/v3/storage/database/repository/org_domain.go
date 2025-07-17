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

const queryOrganizationDomainStmt = `SELECT instance_id, org_id, domain, is_verified, is_primary, verification_type, created_at, updated_at ` +
	`FROM zitadel.organization_domains`

// Get implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).Get of orgDomain.org.
func (o *orgDomain) Get(ctx context.Context, opts ...database.QueryOption) (*domain.OrganizationDomain, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	var builder database.StatementBuilder
	builder.WriteString(queryOrganizationDomainStmt)
	options.Write(&builder)

	return scanOrganizationDomain(ctx, o.client, &builder)
}

// List implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).List of orgDomain.org.
func (o *orgDomain) List(ctx context.Context, opts ...database.QueryOption) ([]*domain.OrganizationDomain, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	var builder database.StatementBuilder
	builder.WriteString(queryOrganizationDomainStmt)
	options.Write(&builder)

	return scanOrganizationDomains(ctx, o.client, &builder)
}

// Add implements [domain.OrganizationDomainRepository].
func (o *orgDomain) Add(ctx context.Context, domain *domain.AddOrganizationDomain) error {
	var builder database.StatementBuilder

	builder.WriteString(`INSERT INTO zitadel.organization_domains (instance_id, org_id, domain, is_verified, is_primary, verification_type) ` +
		`VALUES ($1, $2, $3, $4, $5, $6)` +
		` RETURNING created_at, updated_at`)

	builder.AppendArgs(domain.InstanceID, domain.OrgID, domain.Domain, domain.IsVerified, domain.IsPrimary, domain.VerificationType)

	return o.client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&domain.CreatedAt, &domain.UpdatedAt)
}

// Update implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).Update of orgDomain.org.
func (o *orgDomain) Update(ctx context.Context, condition database.Condition, changes ...database.Change) (int64, error) {
	var builder database.StatementBuilder

	builder.WriteString(`UPDATE zitadel.organization_domains SET `)
	database.Changes(changes).Write(&builder)

	writeCondition(&builder, condition)

	return o.client.Exec(ctx, builder.String(), builder.Args()...)
}

// Remove implements [domain.OrganizationDomainRepository].
func (o *orgDomain) Remove(ctx context.Context, condition database.Condition) (int64, error) {
	var builder database.StatementBuilder

	builder.WriteString(`DELETE FROM zitadel.organization_domains `)
	writeCondition(&builder, condition)

	return o.client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetPrimary implements [domain.OrganizationDomainRepository].
func (o orgDomain) SetPrimary() database.Change {
	return database.NewChange(o.IsPrimaryColumn(), true)
}

// SetVerificationType implements [domain.OrganizationDomainRepository].
func (o orgDomain) SetVerificationType(verificationType domain.DomainVerificationType) database.Change {
	return database.NewChange(o.VerificationTypeColumn(), verificationType)
}

// SetVerified implements [domain.OrganizationDomainRepository].
func (o orgDomain) SetVerified() database.Change {
	return database.NewChange(o.IsVerifiedColumn(), true)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// DomainCondition implements [domain.OrganizationDomainRepository].
func (o orgDomain) DomainCondition(op database.TextOperation, domain string) database.Condition {
	return database.NewTextCondition(o.DomainColumn(), op, domain)
}

// InstanceIDCondition implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).InstanceIDCondition of orgDomain.org.
func (o orgDomain) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(o.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// IsPrimaryCondition implements [domain.OrganizationDomainRepository].
func (o orgDomain) IsPrimaryCondition(isPrimary bool) database.Condition {
	return database.NewBooleanCondition(o.IsPrimaryColumn(), isPrimary)
}

// IsVerifiedCondition implements [domain.OrganizationDomainRepository].
func (o orgDomain) IsVerifiedCondition(isVerified bool) database.Condition {
	return database.NewBooleanCondition(o.IsVerifiedColumn(), isVerified)
}

// OrgIDCondition implements [domain.OrganizationDomainRepository].
func (o orgDomain) OrgIDCondition(orgID string) database.Condition {
	return database.NewTextCondition(o.OrgIDColumn(), database.TextOperationEqual, orgID)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// CreatedAtColumn implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).CreatedAtColumn of orgDomain.org.
func (orgDomain) CreatedAtColumn() database.Column {
	return database.NewColumn("created_at")
}

// DomainColumn implements [domain.OrganizationDomainRepository].
func (orgDomain) DomainColumn() database.Column {
	return database.NewColumn("domain")
}

// InstanceIDColumn implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).InstanceIDColumn of orgDomain.org.
func (orgDomain) InstanceIDColumn() database.Column {
	return database.NewColumn("instance_id")
}

// IsPrimaryColumn implements [domain.OrganizationDomainRepository].
func (orgDomain) IsPrimaryColumn() database.Column {
	return database.NewColumn("is_primary")
}

// IsVerifiedColumn implements [domain.OrganizationDomainRepository].
func (orgDomain) IsVerifiedColumn() database.Column {
	return database.NewColumn("is_verified")
}

// OrgIDColumn implements [domain.OrganizationDomainRepository].
func (orgDomain) OrgIDColumn() database.Column {
	return database.NewColumn("org_id")
}

// UpdatedAtColumn implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).UpdatedAtColumn of orgDomain.org.
func (orgDomain) UpdatedAtColumn() database.Column {
	return database.NewColumn("updated_at")
}

// VerificationTypeColumn implements [domain.OrganizationDomainRepository].
func (orgDomain) VerificationTypeColumn() database.Column {
	return database.NewColumn("verification_type")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------

func scanOrganizationDomain(ctx context.Context, client database.Querier, builder *database.StatementBuilder) (*domain.OrganizationDomain, error) {
	rows, err := client.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	
	organizationDomain := &domain.OrganizationDomain{}
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(organizationDomain); err != nil {
		return nil, err
	}
	return organizationDomain, nil
}

func scanOrganizationDomains(ctx context.Context, client database.Querier, builder *database.StatementBuilder) ([]*domain.OrganizationDomain, error) {
	rows, err := client.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var organizationDomains []*domain.OrganizationDomain
	if err := rows.(database.CollectableRows).Collect(&organizationDomains); err != nil {
		return nil, err
	}
	return organizationDomains, nil
}
