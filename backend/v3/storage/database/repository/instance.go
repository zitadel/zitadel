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

const queryInstanceStmt = `SELECT id, name, created_at, updated_at, deleted_at` +
	` FROM zitadel.instances`

// Get implements [domain.InstanceRepository].
func (i *instance) Get(ctx context.Context, opts ...database.QueryOption) (*domain.Instance, error) {
	i.builder = database.StatementBuilder{}
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	i.builder.WriteString(queryInstanceStmt)
	options.WriteCondition(&i.builder)
	options.WriteOrderBy(&i.builder)
	options.WriteLimit(&i.builder)
	options.WriteOffset(&i.builder)

	return scanInstance(i.client.QueryRow(ctx, i.builder.String(), i.builder.Args()...))
}

const createInstanceStmt = `INSERT INTO zitadel.instances (id, name)` +
	` VALUES ($1, $2)` +
	` RETURNING created_at, updated_at`

// Create implements [domain.InstanceRepository].
func (i *instance) Create(ctx context.Context, instance *domain.Instance) error {
	i.builder = database.StatementBuilder{}
	i.builder.AppendArgs(instance.ID, instance.Name)
	i.builder.WriteString(createInstanceStmt)

	return i.client.QueryRow(ctx, i.builder.String(), i.builder.Args()...).Scan(&instance.CreatedAt, &instance.UpdatedAt)
}

// Update implements [domain.InstanceRepository].
func (i instance) Update(ctx context.Context, condition database.Condition, changes ...database.Change) error {
	i.builder = database.StatementBuilder{}
	i.builder.WriteString(`UPDATE human_users SET `)
	database.Changes(changes).Write(&i.builder)
	i.writeCondition(condition)

	stmt := i.builder.String()

	return i.client.Exec(ctx, stmt, i.builder.Args()...)
}

// Delete implements [domain.InstanceRepository].
func (i instance) Delete(ctx context.Context, condition database.Condition) error {
	i.builder.WriteString("DELETE FROM instance")

	if condition == nil {
		return errors.New("Delete must contain a condition") // (otherwise ALL instances will be deleted)
	}
	i.writeCondition(condition)
	return i.client.Exec(ctx, i.builder.String(), i.builder.Args()...)
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

// UpdatedAtColumn implements [domain.instanceColumns].
func (instance) UpdatedAtColumn() database.Column {
	return database.NewColumn("updated_at")
}

// DeletedAtColumn implements [domain.instanceColumns].
func (instance) DeletedAtColumn() database.Column {
	return database.NewColumn("deleted_at")
}

func (i *instance) writeCondition(condition database.Condition) {
	if condition == nil {
		return
	}
	i.builder.WriteString(" WHERE ")
	condition.Write(&i.builder)
}

func scanInstance(scanner database.Scanner) (*domain.Instance, error) {
	var instance domain.Instance
	err := scanner.Scan(
		&instance.ID,
		&instance.Name,
		&instance.CreatedAt,
		&instance.UpdatedAt,
		&instance.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}
