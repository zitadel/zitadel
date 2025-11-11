package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.InstanceDomainRepository = (*instanceDomain)(nil)

type instanceDomain struct{}

func InstanceDomainRepository() domain.InstanceDomainRepository {
	return new(instanceDomain)
}

func (instanceDomain) qualifiedTableName() string {
	return "zitadel.instance_domains"
}

func (instanceDomain) unqualifiedTableName() string {
	return "instance_domains"
}

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

// Get implements [domain.InstanceDomainRepository].
func (i instanceDomain) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.InstanceDomain, error) {
	builder, err := i.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return scanInstanceDomain(ctx, client, builder)
}

// List implements [domain.InstanceDomainRepository].
func (i instanceDomain) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.InstanceDomain, error) {
	builder, err := i.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return scanInstanceDomains(ctx, client, builder)
}

// Add implements [domain.InstanceDomainRepository].
func (i instanceDomain) Add(ctx context.Context, client database.QueryExecutor, domain *domain.AddInstanceDomain) error {
	builder := database.NewStatementBuilder(`INSERT INTO `)
	builder.WriteString(i.qualifiedTableName())
	builder.WriteString(` (instance_id, domain, is_primary, is_generated, type, created_at, updated_at) VALUES (`)
	builder.WriteArgs(domain.InstanceID, domain.Domain, domain.IsPrimary, domain.IsGenerated, domain.Type, defaultTimestamp(domain.CreatedAt), defaultTimestamp(domain.UpdatedAt))
	builder.WriteString(`) RETURNING created_at, updated_at`)

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&domain.CreatedAt, &domain.UpdatedAt)
}

// Update implements [domain.InstanceDomainRepository].
func (i instanceDomain) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if err := checkRestrictingColumns(condition, i.InstanceIDColumn()); err != nil {
		return 0, err
	}
	return updateSub(ctx, client, i, condition, changes...)
}

// Remove implements [domain.InstanceDomainRepository].
func (i instanceDomain) Remove(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if err := checkRestrictingColumns(condition, i.InstanceIDColumn()); err != nil {
		return 0, err
	}
	return delete(ctx, client, i, condition)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetPrimary implements [domain.InstanceDomainRepository].
func (i instanceDomain) SetPrimary() database.Change {
	return database.NewChange(i.IsPrimaryColumn(), true)
}

// SetUpdatedAt implements [domain.OrganizationDomainRepository].
func (i instanceDomain) SetUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(i.UpdatedAtColumn(), updatedAt)
}

// SetType implements [domain.InstanceDomainRepository].
func (i instanceDomain) SetType(typ domain.DomainType) database.Change {
	return database.NewChange(i.TypeColumn(), typ)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (i instanceDomain) PrimaryKeyCondition(domain string) database.Condition {
	return i.DomainCondition(database.TextOperationEqual, domain)
}

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

// TypeCondition implements [domain.InstanceDomainRepository].
func (i instanceDomain) TypeCondition(typ domain.DomainType) database.Condition {
	return database.NewTextCondition(i.TypeColumn(), database.TextOperationEqual, typ.String())
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// PrimaryKeyColumns implements [domain.Repository].
func (i instanceDomain) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		i.DomainColumn(),
	}
}

// CreatedAtColumn implements [domain.InstanceDomainRepository].
func (i instanceDomain) CreatedAtColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "created_at")
}

// DomainColumn implements [domain.InstanceDomainRepository].
func (i instanceDomain) DomainColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "domain")
}

// InstanceIDColumn implements [domain.InstanceDomainRepository].
func (i instanceDomain) InstanceIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "instance_id")
}

// IsPrimaryColumn implements [domain.InstanceDomainRepository].
func (i instanceDomain) IsPrimaryColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "is_primary")
}

// UpdatedAtColumn implements [domain.InstanceDomainRepository].
func (i instanceDomain) UpdatedAtColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "updated_at")
}

// IsGeneratedColumn implements [domain.InstanceDomainRepository].
func (i instanceDomain) IsGeneratedColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "is_generated")
}

// TypeColumn implements [domain.InstanceDomainRepository].
func (i instanceDomain) TypeColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "type")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------

func scanInstanceDomains(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.InstanceDomain, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var domains []*domain.InstanceDomain
	if err := rows.(database.CollectableRows).Collect(&domains); err != nil {
		return nil, err
	}

	return domains, nil
}

func scanInstanceDomain(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.InstanceDomain, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	domain := new(domain.InstanceDomain)
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(domain); err != nil {
		return nil, err
	}

	return domain, nil
}

// -------------------------------------------------------------
// helpers
// -------------------------------------------------------------

const queryInstanceDomainStmt = `SELECT instance_domains.instance_id, instance_domains.domain, instance_domains.is_primary, instance_domains.created_at, instance_domains.updated_at ` +
	`FROM `

func (i instanceDomain) prepareQuery(opts []database.QueryOption) (*database.StatementBuilder, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	builder := database.NewStatementBuilder(queryInstanceDomainStmt + i.qualifiedTableName())
	options.Write(builder)

	return builder, nil
}
