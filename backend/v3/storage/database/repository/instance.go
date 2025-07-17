package repository

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.InstanceRepository = (*instance)(nil)

type instance struct {
	repository
	shouldJoinDomains bool
	domainRepo        domain.InstanceDomainRepository
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

const queryInstanceStmt = `SELECT id, name, default_org_id, iam_project_id, console_client_id, console_app_id, default_language, created_at, updated_at` +
	` FROM zitadel.instances`

// Get implements [domain.InstanceRepository].
func (i *instance) Get(ctx context.Context, id string) (*domain.Instance, error) {
	var builder database.StatementBuilder

	builder.WriteString(queryInstanceStmt)

	idCondition := i.IDCondition(id)
	writeCondition(&builder, idCondition)

	return scanInstance(ctx, i.client, &builder)
}

// List implements [domain.InstanceRepository].
func (i *instance) List(ctx context.Context, conditions ...database.Condition) ([]*domain.Instance, error) {
	var builder database.StatementBuilder

	builder.WriteString(queryInstanceStmt)

	if conditions != nil {
		writeCondition(&builder, database.And(conditions...))
	}

	return scanInstances(ctx, i.client, &builder)
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
	if changes == nil {
		return 0, errors.New("Update must contain a change")
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
	return database.NewChange(i.NameColumn(), name)
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

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// IDColumn implements [domain.instanceColumns].
func (instance) IDColumn() database.Column {
	return database.NewColumn("id")
}

// NameColumn implements [domain.instanceColumns].
func (instance) NameColumn() database.Column {
	return database.NewColumn("name")
}

// CreatedAtColumn implements [domain.instanceColumns].
func (instance) CreatedAtColumn() database.Column {
	return database.NewColumn("created_at")
}

// DefaultOrgIdColumn implements [domain.instanceColumns].
func (instance) DefaultOrgIDColumn() database.Column {
	return database.NewColumn("default_org_id")
}

// IAMProjectIDColumn implements [domain.instanceColumns].
func (instance) IAMProjectIDColumn() database.Column {
	return database.NewColumn("iam_project_id")
}

// ConsoleClientIDColumn implements [domain.instanceColumns].
func (instance) ConsoleClientIDColumn() database.Column {
	return database.NewColumn("console_client_id")
}

// ConsoleAppIDColumn implements [domain.instanceColumns].
func (instance) ConsoleAppIDColumn() database.Column {
	return database.NewColumn("console_app_id")
}

// DefaultLanguageColumn implements [domain.instanceColumns].
func (instance) DefaultLanguageColumn() database.Column {
	return database.NewColumn("default_language")
}

// UpdatedAtColumn implements [domain.instanceColumns].
func (instance) UpdatedAtColumn() database.Column {
	return database.NewColumn("updated_at")
}

func scanInstance(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.Instance, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	instance := new(domain.Instance)
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(instance); err != nil {
		return nil, err
	}

	return instance, nil
}

func scanInstances(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (instances []*domain.Instance, err error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	if err := rows.(database.CollectableRows).Collect(&instances); err != nil {
		return nil, err
	}

	return instances, nil
}

// -------------------------------------------------------------
// sub repositories
// -------------------------------------------------------------

// Domains implements [domain.InstanceRepository].
func (i *instance) Domains() domain.InstanceDomainRepository {
	i.shouldJoinDomains = true

	if i.domainRepo != nil {
		return i.domainRepo
	}

	i.domainRepo = &instanceDomain{
		repository: i.repository,
		instance:   i,
	}
	return i.domainRepo
}
