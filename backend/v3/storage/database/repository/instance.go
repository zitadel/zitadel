package repository

import (
	"context"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.InstanceRepository = (*instance)(nil)

type instance struct {
	shouldLoadDomains bool
	domainRepo        instanceDomain
}

func InstanceRepository() domain.InstanceRepository {
	return new(instance)
}

func (i instance) qualifiedTableName() string {
	return "zitadel." + i.unqualifiedTableName()
}

func (instance) unqualifiedTableName() string {
	return "instances"
}

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

// Get implements [domain.InstanceRepository].
func (i instance) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Instance, error) {
	opts = append(opts,
		i.joinDomains(),
		database.WithGroupBy(i.IDColumn()),
	)

	builder, err := i.prepareQuery(opts)
	if err != nil {
		return nil, err
	}

	return scanInstance(ctx, client, builder)
}

// List implements [domain.InstanceRepository].
func (i instance) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Instance, error) {
	opts = append(opts,
		i.joinDomains(),
		database.WithGroupBy(i.IDColumn()),
	)

	builder, err := i.prepareQuery(opts)
	if err != nil {
		return nil, err
	}

	return scanInstances(ctx, client, builder)
}

// Create implements [domain.InstanceRepository].
func (i instance) Create(ctx context.Context, client database.QueryExecutor, instance *domain.Instance) error {
	builder := database.NewStatementBuilder(`INSERT INTO `)
	builder.WriteString(i.qualifiedTableName())
	builder.WriteString(` (id, name, default_organization_id, iam_project_id, console_client_id, console_application_id, default_language, created_at, updated_at) VALUES (`)
	builder.WriteArgs(instance.ID, instance.Name, instance.DefaultOrgID, instance.IAMProjectID, instance.ConsoleClientID, instance.ConsoleAppID, instance.DefaultLanguage, defaultTimestamp(instance.CreatedAt), defaultTimestamp(instance.UpdatedAt))
	builder.WriteString(`) RETURNING created_at, updated_at`)

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&instance.CreatedAt, &instance.UpdatedAt)
}

// Update implements [domain.InstanceRepository].
func (i instance) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return update(ctx, client, i, condition, changes...)
}

// Delete implements [domain.InstanceRepository].
func (i instance) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return delete(ctx, client, i, condition)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetName implements [domain.instanceChanges].
func (i instance) SetName(name string) database.Change {
	return database.NewChange(i.NameColumn(), name)
}

// SetUpdatedAt implements [domain.instanceChanges].
func (i instance) SetUpdatedAt(time time.Time) database.Change {
	return database.NewChange(i.UpdatedAtColumn(), time)
}

func (i instance) SetIAMProject(id string) database.Change {
	return database.NewChange(i.IAMProjectIDColumn(), id)
}
func (i instance) SetDefaultOrg(id string) database.Change {
	return database.NewChange(i.DefaultOrgIDColumn(), id)
}
func (i instance) SetDefaultLanguage(lang language.Tag) database.Change {
	return database.NewChange(i.DefaultLanguageColumn(), lang.String())
}
func (i instance) SetConsoleClientID(id string) database.Change {
	return database.NewChange(i.ConsoleClientIDColumn(), id)
}
func (i instance) SetConsoleAppID(id string) database.Change {
	return database.NewChange(i.ConsoleAppIDColumn(), id)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (i instance) PrimaryKeyCondition(instanceID string) database.Condition {
	return i.IDCondition(instanceID)
}

// IDCondition implements [domain.instanceConditions].
func (i instance) IDCondition(id string) database.Condition {
	return database.NewTextCondition(i.IDColumn(), database.TextOperationEqual, id)
}

// NameCondition implements [domain.instanceConditions].
func (i instance) NameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(i.NameColumn(), op, name)
}

// ExistsDomain creates a correlated [database.Exists] condition on instance_domains.
// Use this filter to make sure the Instance returned contains a specific domain.
// of the instance in the aggregated result.
// Example usage:
//
//	domainRepo := instanceRepo.Domains(true) // ensure domains are loaded/aggregated
//	instance, _ := instanceRepo.Get(ctx,
//	    database.WithCondition(
//	        database.And(
//	            instanceRepo.InstanceIDCondition(instanceID),
//	            instanceRepo.ExistsDomain(domainRepo.DomainCondition(database.TextOperationEqual, "example.com")),
//	        ),
//	    ),
//	)
func (i instance) ExistsDomain(cond database.Condition) database.Condition {
	return database.Exists(
		i.domainRepo.qualifiedTableName(),
		database.And(
			database.NewColumnCondition(i.IDColumn(), i.domainRepo.InstanceIDColumn()),
			cond,
		),
	)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// PrimaryKeyColumns implements [domain.Repository].
func (i instance) PrimaryKeyColumns() []database.Column {
	return []database.Column{i.IDColumn()}
}

// IDColumn implements [domain.instanceColumns].
func (i instance) IDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "id")
}

// NameColumn implements [domain.instanceColumns].
func (i instance) NameColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "name")
}

// CreatedAtColumn implements [domain.instanceColumns].
func (i instance) CreatedAtColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "created_at")
}

// DefaultOrgIdColumn implements [domain.instanceColumns].
func (i instance) DefaultOrgIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "default_organization_id")
}

// IAMProjectIDColumn implements [domain.instanceColumns].
func (i instance) IAMProjectIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "iam_project_id")
}

// ConsoleClientIDColumn implements [domain.instanceColumns].
func (i instance) ConsoleClientIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "console_client_id")
}

// ConsoleAppIDColumn implements [domain.instanceColumns].
func (i instance) ConsoleAppIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "console_application_id")
}

// DefaultLanguageColumn implements [domain.instanceColumns].
func (i instance) DefaultLanguageColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "default_language")
}

// UpdatedAtColumn implements [domain.instanceColumns].
func (i instance) UpdatedAtColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "updated_at")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------

type rawInstance struct {
	*domain.Instance
	Domains JSONArray[domain.InstanceDomain] `json:"domains,omitempty" db:"domains"`
}

func scanInstance(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.Instance, error) {
	instance, err := get[rawInstance](ctx, querier, builder)
	if err != nil {
		return nil, err
	}

	instance.Instance.Domains = instance.Domains
	return instance.Instance, nil
}

func scanInstances(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.Instance, error) {
	instances, err := list[rawInstance](ctx, querier, builder)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Instance, len(instances))
	for i, inst := range instances {
		result[i] = inst.Instance
		result[i].Domains = inst.Domains
	}
	return result, nil
}

// -------------------------------------------------------------
// sub repositories
// -------------------------------------------------------------

func (i *instance) LoadDomains() domain.InstanceRepository {
	return &instance{
		shouldLoadDomains: true,
	}
}

func (i *instance) joinDomains() database.QueryOption {
	columns := make([]database.Condition, 0, 2)
	columns = append(columns, database.NewColumnCondition(i.IDColumn(), i.domainRepo.InstanceIDColumn()))

	// If domains should not be joined, we make sure to return null for the domain columns
	// the query optimizer of the dialect should optimize this away if no domains are requested
	if !i.shouldLoadDomains {
		columns = append(columns, database.IsNull(i.domainRepo.InstanceIDColumn()))
	}

	return database.WithLeftJoin(
		i.domainRepo.qualifiedTableName(),
		database.And(columns...),
	)
}

// -------------------------------------------------------------
// helpers
// -------------------------------------------------------------

const (
	queryInstanceStmt = `SELECT instances.id, instances.name, instances.default_organization_id, instances.iam_project_id, instances.console_client_id, instances.console_application_id, instances.default_language, instances.created_at, instances.updated_at` +
		` , jsonb_agg(json_build_object('domain', instance_domains.domain, 'isPrimary', instance_domains.is_primary, 'isGenerated', instance_domains.is_generated, 'createdAt', instance_domains.created_at, 'updatedAt', instance_domains.updated_at)) FILTER (WHERE instance_domains.instance_id IS NOT NULL) AS domains` +
		` FROM `
)

func (i instance) prepareQuery(opts []database.QueryOption) (*database.StatementBuilder, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	builder := database.NewStatementBuilder(queryInstanceStmt + i.qualifiedTableName())
	options.Write(builder)

	return builder, nil
}
