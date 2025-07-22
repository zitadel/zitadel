package repository

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.InstanceRepository = (*instance)(nil)

type instance struct {
	repository
	shouldLoadDomains bool
	domainRepo        *instanceDomain
}

func InstanceRepository(client database.QueryExecutor) domain.InstanceRepository {
	return &instance{
		repository: repository{
			client: client,
		},
	}
}


// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

const (
	queryInstanceStmt = `SELECT instances.id, instances.name, instances.default_org_id, instances.iam_project_id, instances.console_client_id, instances.console_app_id, instances.default_language, instances.created_at, instances.updated_at` +
	` , CASE WHEN count(instance_domains.domain) > 0 THEN jsonb_agg(json_build_object('domain', instance_domains.domain, 'isVerified', instance_domains.is_verified, 'isPrimary', instance_domains.is_primary, 'isGenerated', instance_domains.is_generated, 'validationType', instance_domains.validation_type, 'createdAt', instance_domains.created_at, 'updatedAt', instance_domains.updated_at)) ELSE NULL::JSONB END domains` +
	` FROM zitadel.instances` 
)

// Get implements [domain.InstanceRepository].
func (i *instance) Get(ctx context.Context, opts ...database.QueryOption) (*domain.Instance, error) {
	opts = append(opts, 
		i.joinDomains(), 
		database.WithGroupBy(i.IDColumn(true)),
	)
	
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	var builder database.StatementBuilder
	builder.WriteString(queryInstanceStmt)
	options.Write(&builder)

	return scanInstance(ctx, i.client, &builder)
}

// List implements [domain.InstanceRepository].
func (i *instance) List(ctx context.Context, opts ...database.QueryOption) ([]*domain.Instance, error) {
	opts = append(opts, 
		i.joinDomains(), 
		database.WithGroupBy(i.IDColumn(true)),
	)
	
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	var builder database.StatementBuilder
	builder.WriteString(queryInstanceStmt)
	options.Write(&builder)

	return scanInstances(ctx, i.client, &builder)
}

func (i *instance) joinDomains() database.QueryOption {
	columns := make([]database.Condition, 0, 2)
	columns = append(columns, database.NewColumnCondition(i.IDColumn(true), i.Domains(false).InstanceIDColumn(true)))

	// If domains should not be joined, we make sure to return null for the domain columns
	// the query optimizer of the dialect should optimize this away if no domains are requested
	if !i.shouldLoadDomains {
		columns = append(columns, database.IsNull(i.Domains(false).InstanceIDColumn(true)))
	}

	return database.WithLeftJoin(
		"zitadel.instance_domains",
		database.And(columns...),
	)
}

const createInstanceStmt = `INSERT INTO zitadel.instances (id, name, default_org_id, iam_project_id, console_client_id, console_app_id, default_language)` +
	` VALUES ($1, $2, $3, $4, $5, $6, $7)` +
	` RETURNING created_at, updated_at`

// Create implements [domain.InstanceRepository].
func (i *instance) Create(ctx context.Context, instance *domain.Instance) error {
	var builder database.StatementBuilder

	builder.AppendArgs(instance.ID, instance.Name, instance.DefaultOrgID, instance.IAMProjectID, instance.ConsoleClientID, instance.ConsoleAppID, instance.DefaultLanguage)
	builder.WriteString(createInstanceStmt)

	return i.client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&instance.CreatedAt, &instance.UpdatedAt)
}

// Update implements [domain.InstanceRepository].
func (i instance) Update(ctx context.Context, id string, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.NoChangesError
	}
	var builder database.StatementBuilder

	builder.WriteString(`UPDATE zitadel.instances SET `)

	database.Changes(changes).Write(&builder)

	idCondition := i.IDCondition(id)
	writeCondition(&builder, idCondition)

	stmt := builder.String()

	return i.client.Exec(ctx, stmt, builder.Args()...)
}

// Delete implements [domain.InstanceRepository].
func (i instance) Delete(ctx context.Context, id string) (int64, error) {
	var builder database.StatementBuilder

	builder.WriteString(`DELETE FROM zitadel.instances`)

	idCondition := i.IDCondition(id)
	writeCondition(&builder, idCondition)

	return i.client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetName implements [domain.instanceChanges].
func (i instance) SetName(name string) database.Change {
	return database.NewChange(i.NameColumn(false), name)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// IDCondition implements [domain.instanceConditions].
func (i instance) IDCondition(id string) database.Condition {
	return database.NewTextCondition(i.IDColumn(true), database.TextOperationEqual, id)
}

// NameCondition implements [domain.instanceConditions].
func (i instance) NameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(i.NameColumn(true), op, name)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// IDColumn implements [domain.instanceColumns].
func (instance) IDColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instances.id")
	}
	return database.NewColumn("id")
}

// NameColumn implements [domain.instanceColumns].
func (instance) NameColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instances.name")
	}
	return database.NewColumn("name")
}

// CreatedAtColumn implements [domain.instanceColumns].
func (instance) CreatedAtColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instances.created_at")
	}
	return database.NewColumn("created_at")
}

// DefaultOrgIdColumn implements [domain.instanceColumns].
func (instance) DefaultOrgIDColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instances.default_org_id")
	}
	return database.NewColumn("default_org_id")
}

// IAMProjectIDColumn implements [domain.instanceColumns].
func (instance) IAMProjectIDColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instances.iam_project_id")
	}
	return database.NewColumn("iam_project_id")
}

// ConsoleClientIDColumn implements [domain.instanceColumns].
func (instance) ConsoleClientIDColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instances.console_client_id")
	}
	return database.NewColumn("console_client_id")
}

// ConsoleAppIDColumn implements [domain.instanceColumns].
func (instance) ConsoleAppIDColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instances.console_app_id")
	}
	return database.NewColumn("console_app_id")
}

// DefaultLanguageColumn implements [domain.instanceColumns].
func (instance) DefaultLanguageColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instances.default_language")
	}
	return database.NewColumn("default_language")
}

// UpdatedAtColumn implements [domain.instanceColumns].
func (instance) UpdatedAtColumn(qualified bool) database.Column {
	if qualified {
		return database.NewColumn("instances.updated_at")
	}
	return database.NewColumn("updated_at")
}

type rawInstance struct {
	*domain.Instance
	RawDomains json.RawMessage `json:"domains,omitempty" db:"domains"`
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

	if len(instance.RawDomains) > 0 {
		if err := json.Unmarshal(instance.RawDomains, &instance.Domains); err != nil {
			return nil, err
		}
	}

	return instance.Instance, nil
}

func scanInstances(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.Instance, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var rawInstances []*rawInstance
	if err := rows.(database.CollectableRows).Collect(&rawInstances); err != nil {
		return nil, err
	}

	instances := make([]*domain.Instance, len(rawInstances))
	for i, instance := range rawInstances {
		if len(instance.RawDomains) > 0 {
			if err := json.Unmarshal(instance.RawDomains, &instance.Domains); err != nil {
				return nil, err
			}
		}
		instances[i] = instance.Instance
	}
	return instances, nil
}

// -------------------------------------------------------------
// sub repositories
// -------------------------------------------------------------

// Domains implements [domain.InstanceRepository].
func (i *instance) Domains(shouldLoad bool) domain.InstanceDomainRepository {
	if !i.shouldLoadDomains {
		i.shouldLoadDomains = shouldLoad
	}

	if i.domainRepo != nil {
		return i.domainRepo
	}

	i.domainRepo = &instanceDomain{
		repository: i.repository,
		instance:   i,
	}
	return i.domainRepo
}
