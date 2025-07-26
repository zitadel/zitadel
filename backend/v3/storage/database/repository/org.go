package repository

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

var _ domain.OrganizationRepository = (*org)(nil)

type org struct {
	repository
	shouldLoadDomains bool
	domainRepo        domain.OrganizationDomainRepository
}

func OrganizationRepository(client database.QueryExecutor) domain.OrganizationRepository {
	return &org{
		repository: repository{
			client: client,
		},
	}
}

const queryOrganizationStmt = `SELECT organizations.id, organizations.name, organizations.instance_id, organizations.state, organizations.created_at, organizations.updated_at` +
	` , CASE WHEN count(org_domains.domain) > 0 THEN jsonb_agg(json_build_object('domain', org_domains.domain, 'isVerified', org_domains.is_verified, 'isPrimary', org_domains.is_primary, 'validationType', org_domains.validation_type, 'createdAt', org_domains.created_at, 'updatedAt', org_domains.updated_at)) ELSE NULL::JSONB END domains` +
	` FROM zitadel.organizations`

// Get implements [domain.OrganizationRepository].
func (o *org) Get(ctx context.Context, opts ...database.QueryOption) (*domain.Organization, error) {
	opts = append(opts,
		o.joinDomains(),
		database.WithGroupBy(o.InstanceIDColumn(true), o.IDColumn(true)),
	)

	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	var builder database.StatementBuilder
	builder.WriteString(queryOrganizationStmt)
	options.Write(&builder)

	return scanOrganization(ctx, o.client, &builder)
}

// List implements [domain.OrganizationRepository].
func (o *org) List(ctx context.Context, opts ...database.QueryOption) ([]*domain.Organization, error) {
	opts = append(opts,
		o.joinDomains(),
		database.WithGroupBy(o.InstanceIDColumn(true), o.IDColumn(true)),
	)

	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	var builder database.StatementBuilder
	builder.WriteString(queryOrganizationStmt)
	options.Write(&builder)

	return scanOrganizations(ctx, o.client, &builder)
}

func (o *org) joinDomains() database.QueryOption {
	columns := make([]database.Condition, 0, 3)
	columns = append(columns,
		database.NewColumnCondition(o.InstanceIDColumn(true), o.Domains(false).InstanceIDColumn(true)),
		database.NewColumnCondition(o.IDColumn(true), o.Domains(false).OrgIDColumn(true)),
	)

	// If domains should not be joined, we make sure to return null for the domain columns
	// the query optimizer of the dialect should optimize this away if no domains are requested
	if !o.shouldLoadDomains {
		columns = append(columns, database.IsNull(o.domainRepo.OrgIDColumn(true)))
	}

	return database.WithLeftJoin(
		"zitadel.org_domains",
		database.And(columns...),
	)
}

const createOrganizationStmt = `INSERT INTO zitadel.organizations (id, name, instance_id, state)` +
	` VALUES ($1, $2, $3, $4)` +
	` RETURNING created_at, updated_at`

// Create implements [domain.OrganizationRepository].
func (o *org) Create(ctx context.Context, organization *domain.Organization) error {
	builder := database.StatementBuilder{}
	builder.AppendArgs(organization.ID, organization.Name, organization.InstanceID, organization.State)
	builder.WriteString(createOrganizationStmt)

	return o.client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&organization.CreatedAt, &organization.UpdatedAt)
}

// Update implements [domain.OrganizationRepository].
func (o *org) Update(ctx context.Context, id domain.OrgIdentifierCondition, instanceID string, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.organizations SET `)

	instanceIDCondition := o.InstanceIDCondition(instanceID)

	conditions := []database.Condition{id, instanceIDCondition}
	database.Changes(changes).Write(&builder)
	writeCondition(&builder, database.And(conditions...))

	stmt := builder.String()

	rowsAffected, err := o.client.Exec(ctx, stmt, builder.Args()...)
	return rowsAffected, err
}

// Delete implements [domain.OrganizationRepository].
func (o *org) Delete(ctx context.Context, id domain.OrgIdentifierCondition, instanceID string) (int64, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(`DELETE FROM zitadel.organizations`)

	instanceIDCondition := o.InstanceIDCondition(instanceID)

	conditions := []database.Condition{id, instanceIDCondition}
	writeCondition(&builder, database.And(conditions...))

	return o.client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetName implements [domain.organizationChanges].
func (o org) SetName(name string) database.Change {
	return database.NewChange(o.NameColumn(false), name)
}

// SetState implements [domain.organizationChanges].
func (o org) SetState(state domain.OrgState) database.Change {
	return database.NewChange(o.StateColumn(false), state)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// IDCondition implements [domain.organizationConditions].
func (o org) IDCondition(id string) domain.OrgIdentifierCondition {
	return database.NewTextCondition(o.IDColumn(true), database.TextOperationEqual, id)
}

// NameCondition implements [domain.organizationConditions].
func (o org) NameCondition(name string) domain.OrgIdentifierCondition {
	return database.NewTextCondition(o.NameColumn(true), database.TextOperationEqual, name)
}

// InstanceIDCondition implements [domain.organizationConditions].
func (o org) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(o.InstanceIDColumn(true), database.TextOperationEqual, instanceID)
}

// StateCondition implements [domain.organizationConditions].
func (o org) StateCondition(state domain.OrgState) database.Condition {
	return database.NewTextCondition(o.StateColumn(true), database.TextOperationEqual, state)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// IDColumn implements [domain.organizationColumns].
func (org) IDColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("organizations.id")
	}
	return database.NewColumn("id")
}

// NameColumn implements [domain.organizationColumns].
func (org) NameColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("organizations.name")
	}
	return database.NewColumn("name")
}

// InstanceIDColumn implements [domain.organizationColumns].
func (org) InstanceIDColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("organizations.instance_id")
	}
	return database.NewColumn("instance_id")
}

// StateColumn implements [domain.organizationColumns].
func (org) StateColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("organizations.state")
	}
	return database.NewColumn("state")
}

// CreatedAtColumn implements [domain.organizationColumns].
func (org) CreatedAtColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("organizations.created_at")
	}
	return database.NewColumn("created_at")
}

// UpdatedAtColumn implements [domain.organizationColumns].
func (org) UpdatedAtColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("organizations.updated_at")
	}
	return database.NewColumn("updated_at")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------

type rawOrganization struct {
	*domain.Organization
	RawDomains json.RawMessage `json:"domains,omitempty" db:"domains"`
}

func scanOrganization(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.Organization, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var org rawOrganization
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(&org); err != nil {
		return nil, err
	}
	if len(org.RawDomains) > 0 {
		if err := json.Unmarshal(org.RawDomains, &org.Domains); err != nil {
			return nil, err
		}
	}

	return org.Organization, nil
}

func scanOrganizations(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.Organization, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var rawOrgs []*rawOrganization
	if err := rows.(database.CollectableRows).Collect(&rawOrgs); err != nil {
		return nil, err
	}

	organizations := make([]*domain.Organization, len(rawOrgs))
	for i, org := range rawOrgs {
		if len(org.RawDomains) > 0 {
			if err := json.Unmarshal(org.RawDomains, &org.Domains); err != nil {
				return nil, err
			}
		}
		organizations[i] = org.Organization
	}
	return organizations, nil
}

// -------------------------------------------------------------
// sub repositories
// -------------------------------------------------------------

// Domains implements [domain.OrganizationRepository].
func (o *org) Domains(shouldLoad bool) domain.OrganizationDomainRepository {
	if !o.shouldLoadDomains {
		o.shouldLoadDomains = shouldLoad
	}

	if o.domainRepo != nil {
		return o.domainRepo
	}

	o.domainRepo = &orgDomain{
		repository: o.repository,
		org:        o,
	}

	return o.domainRepo
}
