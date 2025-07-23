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
	var builder database.StatementBuilder

	builder.WriteString(`INSERT INTO zitadel.org_domains (instance_id, org_id, domain, is_verified, is_primary, validation_type) VALUES (`)
	builder.WriteArgs(domain.InstanceID, domain.OrgID, domain.Domain, domain.IsVerified, domain.IsPrimary, domain.ValidationType)
	builder.WriteString(`) RETURNING created_at, updated_at`)

	return o.client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&domain.CreatedAt, &domain.UpdatedAt)
}

// Update implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).Update of orgDomain.org.
func (o *orgDomain) Update(ctx context.Context, condition database.Condition, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.NoChangesError
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
	return database.NewChange(o.IsPrimaryColumn(false), true)
}

// SetValidationType implements [domain.OrganizationDomainRepository].
func (o orgDomain) SetValidationType(verificationType domain.DomainValidationType) database.Change {
	return database.NewChange(o.ValidationTypeColumn(false), verificationType)
}

// SetVerified implements [domain.OrganizationDomainRepository].
func (o orgDomain) SetVerified() database.Change {
	return database.NewChange(o.IsVerifiedColumn(false), true)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// DomainCondition implements [domain.OrganizationDomainRepository].
func (o orgDomain) DomainCondition(op database.TextOperation, domain string) database.Condition {
	return database.NewTextCondition(o.DomainColumn(true), op, domain)
}

// InstanceIDCondition implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).InstanceIDCondition of orgDomain.org.
func (o orgDomain) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(o.InstanceIDColumn(true), database.TextOperationEqual, instanceID)
}

// IsPrimaryCondition implements [domain.OrganizationDomainRepository].
func (o orgDomain) IsPrimaryCondition(isPrimary bool) database.Condition {
	return database.NewBooleanCondition(o.IsPrimaryColumn(true), isPrimary)
}

// IsVerifiedCondition implements [domain.OrganizationDomainRepository].
func (o orgDomain) IsVerifiedCondition(isVerified bool) database.Condition {
	return database.NewBooleanCondition(o.IsVerifiedColumn(true), isVerified)
}

// OrgIDCondition implements [domain.OrganizationDomainRepository].
func (o orgDomain) OrgIDCondition(orgID string) database.Condition {
	return database.NewTextCondition(o.OrgIDColumn(true), database.TextOperationEqual, orgID)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// CreatedAtColumn implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).CreatedAtColumn of orgDomain.org.
func (orgDomain) CreatedAtColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("org_domains.created_at")
	}
	return database.NewColumn("created_at")
}

// DomainColumn implements [domain.OrganizationDomainRepository].
func (orgDomain) DomainColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("org_domains.domain")
	}
	return database.NewColumn("domain")
}

// InstanceIDColumn implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).InstanceIDColumn of orgDomain.org.
func (orgDomain) InstanceIDColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("org_domains.instance_id")
	}
	return database.NewColumn("instance_id")
}

// IsPrimaryColumn implements [domain.OrganizationDomainRepository].
func (orgDomain) IsPrimaryColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("org_domains.is_primary")
	}
	return database.NewColumn("is_primary")
}

// IsVerifiedColumn implements [domain.OrganizationDomainRepository].
func (orgDomain) IsVerifiedColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("org_domains.is_verified")
	}
	return database.NewColumn("is_verified")
}

// OrgIDColumn implements [domain.OrganizationDomainRepository].
func (orgDomain) OrgIDColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("org_domains.org_id")
	}
	return database.NewColumn("org_id")
}

// UpdatedAtColumn implements [domain.OrganizationDomainRepository].
// Subtle: this method shadows the method ([domain.OrganizationRepository]).UpdatedAtColumn of orgDomain.org.
func (orgDomain) UpdatedAtColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("org_domains.updated_at")
	}
	return database.NewColumn("updated_at")
}

// ValidationTypeColumn implements [domain.OrganizationDomainRepository].
func (orgDomain) ValidationTypeColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("org_domains.validation_type")
	}
	return database.NewColumn("validation_type")
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
