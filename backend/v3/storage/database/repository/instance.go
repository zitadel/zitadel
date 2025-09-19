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

func (instance) qualifiedTableName() string {
	return "zitadel.instances"
}

func (instance) unqualifiedTableName() string {
	return "instances"
}

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

const (
	queryInstanceStmt = `SELECT instances.id, instances.name, instances.default_org_id, instances.iam_project_id, instances.console_client_id, instances.console_app_id, instances.default_language, instances.created_at, instances.updated_at` +
		` , jsonb_agg(json_build_object('domain', instance_domains.domain, 'isPrimary', instance_domains.is_primary, 'isGenerated', instance_domains.is_generated, 'createdAt', instance_domains.created_at, 'updatedAt', instance_domains.updated_at)) FILTER (WHERE instance_domains.instance_id IS NOT NULL) AS domains` +
		` FROM zitadel.instances`
)

// Get implements [domain.InstanceRepository].
func (i instance) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Instance, error) {
	opts = append(opts,
		i.joinDomains(),
		database.WithGroupBy(i.IDColumn()),
	)

	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	var builder database.StatementBuilder
	builder.WriteString(queryInstanceStmt)
	options.Write(&builder)

	return scanInstance(ctx, client, &builder)
}

// List implements [domain.InstanceRepository].
func (i instance) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Instance, error) {
	opts = append(opts,
		i.joinDomains(),
		database.WithGroupBy(i.IDColumn()),
	)

	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	var builder database.StatementBuilder
	builder.WriteString(queryInstanceStmt)
	options.Write(&builder)

	return scanInstances(ctx, client, &builder)
}

// Create implements [domain.InstanceRepository].
func (i instance) Create(ctx context.Context, client database.QueryExecutor, instance *domain.Instance) error {
	var (
		builder              database.StatementBuilder
		createdAt, updatedAt any = database.DefaultInstruction, database.DefaultInstruction
	)
	if !instance.CreatedAt.IsZero() {
		createdAt = instance.CreatedAt
	}
	if !instance.UpdatedAt.IsZero() {
		updatedAt = instance.UpdatedAt
	}

	builder.WriteString(`INSERT INTO `)
	builder.WriteString(i.qualifiedTableName())
	builder.WriteString(` (id, name, default_org_id, iam_project_id, console_client_id, console_app_id, default_language, created_at, updated_at) VALUES (`)
	builder.WriteArgs(instance.ID, instance.Name, instance.DefaultOrgID, instance.IAMProjectID, instance.ConsoleClientID, instance.ConsoleAppID, instance.DefaultLanguage, createdAt, updatedAt)
	builder.WriteString(`) RETURNING created_at, updated_at`)

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&instance.CreatedAt, &instance.UpdatedAt)
}

// Update implements [domain.InstanceRepository].
func (i instance) Update(ctx context.Context, client database.QueryExecutor, id string, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}
	var builder database.StatementBuilder

	builder.WriteString(`UPDATE zitadel.instances SET `)

	database.Changes(changes).Write(&builder)

	idCondition := i.IDCondition(id)
	writeCondition(&builder, idCondition)

	stmt := builder.String()

	return client.Exec(ctx, stmt, builder.Args()...)
}

// Delete implements [domain.InstanceRepository].
func (i instance) Delete(ctx context.Context, client database.QueryExecutor, id string) (int64, error) {
	var builder database.StatementBuilder

	builder.WriteString(`DELETE FROM `)
	builder.WriteString(i.qualifiedTableName())

	idCondition := i.IDCondition(id)
	writeCondition(&builder, idCondition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
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

// IDCondition implements [domain.instanceConditions].
func (i instance) IDCondition(id string) database.Condition {
	return database.NewTextCondition(i.IDColumn(), database.TextOperationEqual, id)
}

// NameCondition implements [domain.instanceConditions].
func (i instance) NameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(i.NameColumn(), op, name)
}

func (i instance) ExistsDomain(cond database.Condition) database.Condition {
	// Build a correlated subquery: EXISTS (SELECT 1 FROM zitadel.org_domains WHERE
	//   organizations.instance_id = org_domains.instance_id AND organizations.id = org_domains.org_id AND <cond>)
	correlated := database.And(
		database.NewColumnCondition(i.IDColumn(), i.domainRepo.InstanceIDColumn()),
		cond,
	)
	return existsCondition{
		table:     i.qualifiedTableName(),
		condition: correlated,
	}
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

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
	return database.NewColumn(i.unqualifiedTableName(), "default_org_id")
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
	return database.NewColumn(i.unqualifiedTableName(), "console_app_id")
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
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var instance rawInstance
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(&instance); err != nil {
		return nil, err
	}
	instance.Instance.Domains = instance.Domains

	return instance.Instance, nil
}

func scanInstances(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.Instance, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var instances []*rawInstance
	if err := rows.(database.CollectableRows).Collect(&instances); err != nil {
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
