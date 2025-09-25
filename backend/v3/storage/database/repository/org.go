package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

var _ domain.OrganizationRepository = (*org)(nil)

type org struct {
	shouldLoadDomains  bool
	domainRepo         orgDomain
	shouldLoadMetadata bool
	metadataRepo       orgMetadata
}

func (o org) unqualifiedTableName() string {
	return "organizations"
}

func OrganizationRepository() domain.OrganizationRepository {
	return new(org)
}

const queryOrganizationStmt = `SELECT organizations.id, organizations.name, organizations.instance_id, organizations.state, organizations.created_at, organizations.updated_at` +
	` , jsonb_agg(json_build_object('instanceId', org_domains.instance_id, 'orgId', org_domains.org_id, 'domain', org_domains.domain, 'isVerified', org_domains.is_verified, 'isPrimary', org_domains.is_primary, 'validationType', org_domains.validation_type, 'createdAt', org_domains.created_at, 'updatedAt', org_domains.updated_at)) FILTER (WHERE org_domains.org_id IS NOT NULL) AS domains` +
	` , jsonb_agg(json_build_object('instanceId', org_metadata.instance_id, 'orgId', org_metadata.org_id, 'key', org_metadata.key, 'value', encode(org_metadata.value, 'base64'), 'createdAt', org_metadata.created_at, 'updatedAt', org_metadata.updated_at)) FILTER (WHERE org_metadata.org_id IS NOT NULL) AS metadata` +
	` FROM zitadel.organizations`

// Get implements [domain.OrganizationRepository].
func (o org) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Organization, error) {
	opts = append(opts,
		o.joinDomains(),
		o.joinMetadata(),
		database.WithGroupBy(o.InstanceIDColumn(), o.IDColumn()),
	)

	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	if !options.Condition.IsRestrictingColumn(o.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(o.InstanceIDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(queryOrganizationStmt)
	options.Write(&builder)

	return scanOrganization(ctx, client, &builder)
}

// List implements [domain.OrganizationRepository].
func (o org) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Organization, error) {
	opts = append(opts,
		o.joinDomains(),
		o.joinMetadata(),
		database.WithGroupBy(o.InstanceIDColumn(), o.IDColumn()),
	)

	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	if !options.Condition.IsRestrictingColumn(o.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(o.InstanceIDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(queryOrganizationStmt)
	options.Write(&builder)

	return scanOrganizations(ctx, client, &builder)
}

const createOrganizationStmt = `INSERT INTO zitadel.organizations (id, name, instance_id, state)` +
	` VALUES ($1, $2, $3, $4)` +
	` RETURNING created_at, updated_at`

// Create implements [domain.OrganizationRepository].
func (o org) Create(ctx context.Context, client database.QueryExecutor, organization *domain.Organization) error {
	builder := database.StatementBuilder{}
	builder.AppendArgs(organization.ID, organization.Name, organization.InstanceID, organization.State)
	builder.WriteString(createOrganizationStmt)

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&organization.CreatedAt, &organization.UpdatedAt)
}

// Update implements [domain.OrganizationRepository].
func (o org) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}
	if !condition.IsRestrictingColumn(o.InstanceIDColumn()) {
		return 0, database.NewMissingConditionError(o.InstanceIDColumn())
	}
	if !database.Changes(changes).IsOnColumn(o.UpdatedAtColumn()) {
		changes = append(changes, database.NewChange(o.UpdatedAtColumn(), database.NullInstruction))
	}

	var builder database.StatementBuilder
	builder.WriteString(`UPDATE zitadel.organizations SET `)
	database.Changes(changes).Write(&builder)
	writeCondition(&builder, condition)

	stmt := builder.String()

	rowsAffected, err := client.Exec(ctx, stmt, builder.Args()...)
	return rowsAffected, err
}

// Delete implements [domain.OrganizationRepository].
func (o org) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if !condition.IsRestrictingColumn(o.InstanceIDColumn()) {
		return 0, database.NewMissingConditionError(o.InstanceIDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(`DELETE FROM zitadel.organizations`)
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetName implements [domain.organizationChanges].
func (o org) SetName(name string) database.Change {
	return database.NewChange(o.NameColumn(), name)
}

// SetState implements [domain.organizationChanges].
func (o org) SetState(state domain.OrgState) database.Change {
	return database.NewChange(o.StateColumn(), state)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// IDCondition implements [domain.organizationConditions].
func (o org) IDCondition(id string) database.Condition {
	return database.NewTextCondition(o.IDColumn(), database.TextOperationEqual, id)
}

// NameCondition implements [domain.organizationConditions].
func (o org) NameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(o.NameColumn(), op, name)
}

// InstanceIDCondition implements [domain.organizationConditions].
func (o org) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(o.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// StateCondition implements [domain.organizationConditions].
func (o org) StateCondition(state domain.OrgState) database.Condition {
	return database.NewTextCondition(o.StateColumn(), database.TextOperationEqual, state.String())
}

// ExistsDomain creates a correlated [database.Exists] condition on org_domains.
// Use this filter to make sure the organization returned contains a specific domain.
// Example usage:
//
//	domainRepo := orgRepo.Domains(true) // ensure domains are loaded/aggregated
//	org, _ := orgRepo.Get(ctx,
//	    database.WithCondition(
//	        database.And(
//	            orgRepo.InstanceIDCondition(instanceID),
//	            orgRepo.ExistsDomain(domainRepo.DomainCondition(database.TextOperationEqual, "example.com")),
//	        ),
//	    ),
//	)
func (o org) ExistsDomain(cond database.Condition) database.Condition {
	return database.Exists(
		o.domainRepo.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(o.InstanceIDColumn(), o.domainRepo.InstanceIDColumn()),
			database.NewColumnCondition(o.IDColumn(), o.domainRepo.OrgIDColumn()),
			cond,
		),
	)
}

// ExistsMetadata creates a correlated [database.Exists] condition on org_metadata.
// Use this when you want to filter organizations by a metadata condition but still return all metadata
// of the organization in the aggregated result.
// Example usage:
//
//	metadataRepo := orgRepo.Metadata(true) // ensure metadata are loaded/aggregated
//	org, _ := orgRepo.Get(ctx,
//	    database.WithCondition(
//	        database.And(
//	            orgRepo.InstanceIDCondition(instanceID),
//	            orgRepo.MetadataExists(metadataRepo.KeyCondition(database.TextOperationEqual, "urn:zitadel:org:custom:my-key")),
//	        ),
//	    ),
//	)
func (o org) ExistsMetadata(cond database.Condition) database.Condition {
	return database.Exists(
		o.metadataRepo.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(o.InstanceIDColumn(), o.metadataRepo.InstanceIDColumn()),
			database.NewColumnCondition(o.IDColumn(), o.metadataRepo.OrgIDColumn()),
			cond,
		),
	)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// IDColumn implements [domain.organizationColumns].
func (o org) IDColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "id")
}

// NameColumn implements [domain.organizationColumns].
func (o org) NameColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "name")
}

// InstanceIDColumn implements [domain.organizationColumns].
func (o org) InstanceIDColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "instance_id")
}

// StateColumn implements [domain.organizationColumns].
func (o org) StateColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "state")
}

// CreatedAtColumn implements [domain.organizationColumns].
func (o org) CreatedAtColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "created_at")
}

// UpdatedAtColumn implements [domain.organizationColumns].
func (o org) UpdatedAtColumn() database.Column {
	return database.NewColumn(o.unqualifiedTableName(), "updated_at")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------

type rawOrg struct {
	*domain.Organization
	Domains  JSONArray[domain.OrganizationDomain]   `json:"domains,omitempty" db:"domains"`
	Metadata JSONArray[domain.OrganizationMetadata] `json:"metadata,omitempty" db:"metadata"`
}

func scanOrganization(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.Organization, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var org rawOrg
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(&org); err != nil {
		return nil, err
	}
	org.Organization.Domains = org.Domains
	org.Organization.Metadata = org.Metadata

	return org.Organization, nil
}

func scanOrganizations(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.Organization, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var orgs []*rawOrg
	if err := rows.(database.CollectableRows).Collect(&orgs); err != nil {
		return nil, err
	}

	result := make([]*domain.Organization, len(orgs))
	for i, org := range orgs {
		result[i] = org.Organization
		result[i].Domains = org.Domains
		result[i].Metadata = org.Metadata
	}

	return result, nil
}

// -------------------------------------------------------------
// sub repositories
// -------------------------------------------------------------

func (o org) LoadDomains() domain.OrganizationRepository {
	return &org{
		shouldLoadDomains:  true,
		shouldLoadMetadata: o.shouldLoadMetadata,
	}
}

func (o org) joinDomains() database.QueryOption {
	columns := make([]database.Condition, 0, 3)
	columns = append(columns,
		database.NewColumnCondition(o.InstanceIDColumn(), o.domainRepo.InstanceIDColumn()),
		database.NewColumnCondition(o.IDColumn(), o.domainRepo.OrgIDColumn()),
	)

	// If domains should not be joined, we make sure to return null for the domain columns
	// the query optimizer of the dialect should optimize this away if no domains are requested
	if !o.shouldLoadDomains {
		columns = append(columns, database.IsNull(o.domainRepo.OrgIDColumn()))
	}

	return database.WithLeftJoin(
		o.domainRepo.qualifiedTableName(),
		database.And(columns...),
	)
}

func (o org) LoadMetadata() domain.OrganizationRepository {
	return &org{
		shouldLoadDomains:  o.shouldLoadDomains,
		shouldLoadMetadata: true,
	}
}

func (o org) joinMetadata() database.QueryOption {
	columns := make([]database.Condition, 0, 3)
	columns = append(columns,
		database.NewColumnCondition(o.InstanceIDColumn(), o.metadataRepo.InstanceIDColumn()),
		database.NewColumnCondition(o.IDColumn(), o.metadataRepo.OrgIDColumn()),
	)

	// If metadata should not be joined, we make sure to return null for the metadata columns
	// the query optimizer of the dialect should optimize this away if no metadata are requested
	if !o.shouldLoadMetadata {
		columns = append(columns, database.IsNull(o.metadataRepo.OrgIDColumn()))
	}

	return database.WithLeftJoin(
		o.metadataRepo.qualifiedTableName(),
		database.And(columns...),
	)
}
