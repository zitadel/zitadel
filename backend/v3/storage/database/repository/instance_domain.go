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

const queryInstanceDomainStmt = `SELECT instance_domains.instance_id, instance_domains.domain, instance_domains.is_verified, instance_domains.is_primary, instance_domains.validation_type, instance_domains.created_at, instance_domains.updated_at ` +
	`FROM zitadel.instance_domains id`

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
	var builder database.StatementBuilder

	builder.WriteString(`INSERT INTO zitadel.instance_domains (instance_id, domain, is_verified, is_primary, validation_type) ` +
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
// Subtle: this method shadows the method ([domain.InstanceRepository]).Update of instanceDomain.instance.
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

// SetValidationType implements [domain.InstanceDomainRepository].
func (i instanceDomain) SetValidationType(verificationType domain.DomainValidationType) database.Change {
	return database.NewChange(i.ValidationTypeColumn(false), verificationType)
}

// SetPrimary implements [domain.InstanceDomainRepository].
func (i instanceDomain) SetPrimary() database.Change {
	return database.NewChange(i.IsPrimaryColumn(false), true)
}

// SetVerified implements [domain.InstanceDomainRepository].
func (i instanceDomain) SetVerified() database.Change {
	return database.NewChange(i.IsVerifiedColumn(false), true)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// DomainCondition implements [domain.InstanceDomainRepository].
func (i instanceDomain) DomainCondition(op database.TextOperation, domain string) database.Condition {
	return database.NewTextCondition(i.DomainColumn(true), op, domain)
}

// InstanceIDCondition implements [domain.InstanceDomainRepository].
func (i instanceDomain) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(i.InstanceIDColumn(true), database.TextOperationEqual, instanceID)
}

// IsPrimaryCondition implements [domain.InstanceDomainRepository].
func (i instanceDomain) IsPrimaryCondition(isPrimary bool) database.Condition {
	return database.NewBooleanCondition(i.IsPrimaryColumn(true), isPrimary)
}

// IsVerifiedCondition implements [domain.InstanceDomainRepository].
func (i instanceDomain) IsVerifiedCondition(isVerified bool) database.Condition {
	return database.NewBooleanCondition(i.IsVerifiedColumn(true), isVerified)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// CreatedAtColumn implements [domain.InstanceDomainRepository].
// Subtle: this method shadows the method ([domain.InstanceRepository]).CreatedAtColumn of instanceDomain.instance.
func (instanceDomain) CreatedAtColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instance_domains.created_at")
	}
	return database.NewColumn("created_at")
}

// DomainColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) DomainColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instance_domains.domain")
	}
	return database.NewColumn("domain")
}

// InstanceIDColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) InstanceIDColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instance_domains.instance_id")
	}
	return database.NewColumn("instance_id")
}

// IsPrimaryColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) IsPrimaryColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instance_domains.is_primary")
	}
	return database.NewColumn("is_primary")
}

// IsVerifiedColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) IsVerifiedColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instance_domains.is_verified")
	}
	return database.NewColumn("is_verified")
}

// UpdatedAtColumn implements [domain.InstanceDomainRepository].
// Subtle: this method shadows the method ([domain.InstanceRepository]).UpdatedAtColumn of instanceDomain.instance.
func (instanceDomain) UpdatedAtColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instance_domains.updated_at")
	}
	return database.NewColumn("updated_at")
}

// ValidationTypeColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) ValidationTypeColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instance_domains.validation_type")
	}
	return database.NewColumn("validation_type")
}

// IsGeneratedColumn implements [domain.InstanceDomainRepository].
func (instanceDomain) IsGeneratedColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instance_domains.is_generated")
	}
	return database.NewColumn("is_generated")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------

func scanInstanceDomains(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.InstanceDomain, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var instanceDomains []*domain.InstanceDomain
	if err := rows.(database.CollectableRows).Collect(&instanceDomains); err != nil {
		return nil, err
	}

	return instanceDomains, nil
}

func scanInstanceDomain(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.InstanceDomain, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	instanceDomain := &domain.InstanceDomain{}
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(instanceDomain); err != nil {
		return nil, err
	}

	return instanceDomain, nil
}
