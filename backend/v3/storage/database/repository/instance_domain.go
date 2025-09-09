package repository

import (
	"context"
	"time"

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

const queryInstanceDomainStmt = `SELECT instance_domains.instance_id, instance_domains.domain, instance_domains.is_primary, instance_domains.created_at, instance_domains.updated_at ` +
	`FROM zitadel.instance_domains`

// Get implements [domain.InstanceDomainRepository].
// Subtle: this method shadows the method ([domain.InstanceRepository]).Get of instanceDomain.instance.
func (i *instanceDomain) Get(ctx context.Context, opts ...database.QueryOption) (*domain.InstanceDomain, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	var builder database.StatementBuilder
	builder.WriteString(queryInstanceDomainStmt)
	options.Write(&builder)

	return scanInstanceDomain(ctx, i.client, &builder)
}

// List implements [domain.InstanceDomainRepository].
// Subtle: this method shadows the method ([domain.InstanceRepository]).List of instanceDomain.instance.
func (i *instanceDomain) List(ctx context.Context, opts ...database.QueryOption) ([]*domain.InstanceDomain, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	var builder database.StatementBuilder
	builder.WriteString(queryInstanceDomainStmt)
	options.Write(&builder)

	return scanInstanceDomains(ctx, i.client, &builder)
}

// Add implements [domain.InstanceDomainRepository].
func (i *instanceDomain) Add(ctx context.Context, domain *domain.AddInstanceDomain) error {
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

	builder.WriteString(`INSERT INTO zitadel.instance_domains (instance_id, domain, is_primary, is_generated, type, created_at, updated_at) VALUES (`)
	builder.WriteArgs(domain.InstanceID, domain.Domain, domain.IsPrimary, domain.IsGenerated, domain.Type, createdAt, updatedAt)
	builder.WriteString(`) RETURNING created_at, updated_at`)

	return i.client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&domain.CreatedAt, &domain.UpdatedAt)
}

// Update implements [domain.InstanceDomainRepository].
// Subtle: this method shadows the method ([domain.InstanceRepository]).Update of instanceDomain.instance.
func (i *instanceDomain) Update(ctx context.Context, condition database.Condition, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}
	var builder database.StatementBuilder

	builder.WriteString(`UPDATE zitadel.instance_domains SET `)
	database.Changes(changes).Write(&builder)

	writeCondition(&builder, condition)

	return i.client.Exec(ctx, builder.String(), builder.Args()...)
}

// Remove implements [domain.InstanceDomainRepository].
func (i *instanceDomain) Remove(ctx context.Context, condition database.Condition) (int64, error) {
	var builder database.StatementBuilder

	builder.WriteString(`DELETE FROM zitadel.instance_domains WHERE `)
	condition.Write(&builder)

	return i.client.Exec(ctx, builder.String(), builder.Args()...)
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

// CreatedAtColumn implements [domain.InstanceDomainRepository].
// Subtle: this method shadows the method ([domain.InstanceRepository]).CreatedAtColumn of instanceDomain.instance.
func (instanceDomain) CreatedAtColumn() database.Column {
	return database.NewColumn("instance_domains", "created_at")
}

// DomainColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) DomainColumn() database.Column {
	return database.NewColumn("instance_domains", "domain")
}

// InstanceIDColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) InstanceIDColumn() database.Column {
	return database.NewColumn("instance_domains", "instance_id")
}

// IsPrimaryColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) IsPrimaryColumn() database.Column {
	return database.NewColumn("instance_domains", "is_primary")
}

// UpdatedAtColumn implements [domain.InstanceDomainRepository].
// Subtle: this method shadows the method ([domain.InstanceRepository]).UpdatedAtColumn of instanceDomain.instance.
func (instanceDomain) UpdatedAtColumn() database.Column {
	return database.NewColumn("instance_domains", "updated_at")
}

// IsGeneratedColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) IsGeneratedColumn() database.Column {
	return database.NewColumn("instance_domains", "is_generated")
}

// TypeColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) TypeColumn() database.Column {
	return database.NewColumn("instance_domains", "type")
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
