package repository

import (
	"context"
	"time"

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

const queryOrganizationDomainStmt = `SELECT instance_id, org_id, domain, is_verified, is_primary, validation_type, created_at, updated_at ` +
	`FROM zitadel.org_domains`

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
	var (
		builder              database.StatementBuilder
		createdAt, updatedAt any = database.DefaultInstruction, database.DefaultInstruction
	)
	if !domain.CreatedAt.IsZero() {
		createdAt = domain.CreatedAt
	}
	if !domain.UpdatedAt.IsZero() {
		updatedAt = domain.UpdatedAt
	}

	builder.WriteString(`INSERT INTO zitadel.org_domains (instance_id, org_id, domain, is_verified, is_primary, validation_type, created_at, updated_at) VALUES (`)
	builder.WriteArgs(domain.InstanceID, domain.OrgID, domain.Domain, domain.IsVerified, domain.IsPrimary, domain.ValidationType, createdAt, updatedAt)
	builder.WriteString(`) RETURNING created_at, updated_at`)

	return o.client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&domain.CreatedAt, &domain.UpdatedAt)
}

// Update implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).Update of orgDomain.org.
func (o *orgDomain) Update(ctx context.Context, condition database.Condition, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}

	var builder database.StatementBuilder

	builder.WriteString(`UPDATE zitadel.org_domains SET `)
	database.Changes(changes).Write(&builder)
	writeCondition(&builder, condition)

	return o.client.Exec(ctx, builder.String(), builder.Args()...)
}

// Remove implements [domain.OrganizationDomainRepository].
func (o *orgDomain) Remove(ctx context.Context, condition database.Condition) (int64, error) {
	var builder database.StatementBuilder

	builder.WriteString(`DELETE FROM zitadel.org_domains `)
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

// SetValidationType implements [domain.OrganizationDomainRepository].
func (o orgDomain) SetValidationType(verificationType domain.DomainValidationType) database.Change {
	return database.NewChange(o.ValidationTypeColumn(), verificationType)
}

// SetVerified implements [domain.OrganizationDomainRepository].
func (o orgDomain) SetVerified() database.Change {
	return database.NewChange(o.IsVerifiedColumn(), true)
}

// SetUpdatedAt implements [domain.OrganizationDomainRepository].
func (o orgDomain) SetUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(o.UpdatedAtColumn(), updatedAt)
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
	return database.NewColumn("org_domains", "created_at")
}

// DomainColumn implements [domain.OrganizationDomainRepository].
func (orgDomain) DomainColumn() database.Column {
	return database.NewColumn("org_domains", "domain")
}

// InstanceIDColumn implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).InstanceIDColumn of orgDomain.org.
func (orgDomain) InstanceIDColumn() database.Column {
	return database.NewColumn("org_domains", "instance_id")
}

// IsPrimaryColumn implements [domain.OrganizationDomainRepository].
func (orgDomain) IsPrimaryColumn() database.Column {
	return database.NewColumn("org_domains", "is_primary")
}

// IsVerifiedColumn implements [domain.OrganizationDomainRepository].
func (orgDomain) IsVerifiedColumn() database.Column {
	return database.NewColumn("org_domains", "is_verified")
}

// OrgIDColumn implements [domain.OrganizationDomainRepository].
func (orgDomain) OrgIDColumn() database.Column {
	return database.NewColumn("org_domains", "org_id")
}

// UpdatedAtColumn implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).UpdatedAtColumn of orgDomain.org.
func (orgDomain) UpdatedAtColumn() database.Column {
	return database.NewColumn("org_domains", "updated_at")
}

// ValidationTypeColumn implements [domain.OrganizationDomainRepository].
func (orgDomain) ValidationTypeColumn() database.Column {
	return database.NewColumn("org_domains", "validation_type")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------

func scanOrganizationDomain(ctx context.Context, client database.Querier, builder *database.StatementBuilder) (*domain.OrganizationDomain, error) {
	rows, err := client.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	domain := &domain.OrganizationDomain{}
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(domain); err != nil {
		return nil, err
	}
	return domain, nil
}

func scanOrganizationDomains(ctx context.Context, client database.Querier, builder *database.StatementBuilder) ([]*domain.OrganizationDomain, error) {
	rows, err := client.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var domains []*domain.OrganizationDomain
	if err := rows.(database.CollectableRows).Collect(&domains); err != nil {
		return nil, err
	}
	return domains, nil
}
