package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.OrganizationDomainRepository = (*orgDomain)(nil)

type orgDomain struct{}

func OrganizationDomainRepository() domain.OrganizationDomainRepository {
	return new(orgDomain)
}

func (o orgDomain) qualifiedTableName() string {
	return "zitadel." + o.unqualifiedTableName()
}

func (orgDomain) unqualifiedTableName() string {
	return "org_domains"
}

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

// Get implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).Get of orgDomain.org.
func (o orgDomain) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.OrganizationDomain, error) {
	builder, err := o.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return get[domain.OrganizationDomain](ctx, client, builder)
}

// List implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).List of orgDomain.org.
func (o orgDomain) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.OrganizationDomain, error) {
	builder, err := o.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return list[domain.OrganizationDomain](ctx, client, builder)
}

// Add implements [domain.OrganizationDomainRepository].
func (o orgDomain) Add(ctx context.Context, client database.QueryExecutor, domain *domain.AddOrganizationDomain) error {
	builder := database.NewStatementBuilder(`INSERT INTO `)
	builder.WriteString(o.qualifiedTableName())
	builder.WriteString(` (instance_id, organization_id, domain, is_verified, is_primary, validation_type, created_at, updated_at) VALUES (`)
	builder.WriteArgs(domain.InstanceID, domain.OrgID, domain.Domain, domain.IsVerified, domain.IsPrimary, domain.ValidationType, defaultTimestamp(domain.CreatedAt), defaultTimestamp(domain.UpdatedAt))
	builder.WriteString(`) RETURNING created_at, updated_at`)

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&domain.CreatedAt, &domain.UpdatedAt)
}

// Update implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).Update of orgDomain.org.
func (o orgDomain) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return update(ctx, client, o, condition, changes...)
}

// Remove implements [domain.OrganizationDomainRepository].
func (o orgDomain) Remove(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if err := checkRestrictingColumns(condition, o.InstanceIDColumn(), o.OrgIDColumn()); err != nil {
		return 0, err
	}
	return delete(ctx, client, o, condition)
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

func (o orgDomain) PrimaryKeyCondition(instanceID, orgID, domain string) database.Condition {
	return database.And(
		o.InstanceIDCondition(instanceID),
		o.OrgIDCondition(orgID),
		o.DomainCondition(database.TextOperationEqual, domain),
	)
}

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

// PrimaryKeyColumns implements [domain.Repository].
func (o orgDomain) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		o.InstanceIDColumn(),
		o.OrgIDColumn(),
		o.DomainColumn(),
	}
}

// CreatedAtColumn implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).CreatedAtColumn of orgDomain.org.
func (o orgDomain) CreatedAtColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "created_at")
}

// DomainColumn implements [domain.OrganizationDomainRepository].
func (o orgDomain) DomainColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "domain")
}

// InstanceIDColumn implements [domain.OrganizationDomainRepository].
func (o orgDomain) InstanceIDColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "instance_id")
}

// IsPrimaryColumn implements [domain.OrganizationDomainRepository].
func (o orgDomain) IsPrimaryColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "is_primary")
}

// IsVerifiedColumn implements [domain.OrganizationDomainRepository].
func (o orgDomain) IsVerifiedColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "is_verified")
}

// OrgIDColumn implements [domain.OrganizationDomainRepository].
func (o orgDomain) OrgIDColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "organization_id")
}

// UpdatedAtColumn implements [domain.OrganizationDomainRepository].
func (o orgDomain) UpdatedAtColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "updated_at")
}

// ValidationTypeColumn implements [domain.OrganizationDomainRepository].
func (o orgDomain) ValidationTypeColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "validation_type")
}

// -------------------------------------------------------------
// helpers
// -------------------------------------------------------------

const queryOrganizationDomainStmt = `SELECT instance_id, organization_id, domain, is_verified, is_primary, validation_type, created_at, updated_at ` +
	`FROM `

func (o orgDomain) prepareQuery(opts []database.QueryOption) (*database.StatementBuilder, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if err := checkRestrictingColumns(options.Condition, o.InstanceIDColumn()); err != nil {
		return nil, err
	}
	builder := database.NewStatementBuilder(queryOrganizationDomainStmt + o.qualifiedTableName())
	options.Write(builder)

	return builder, nil
}
