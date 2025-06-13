package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
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

const queryInstanceStmt = `SELECT id, name, default_org_id, iam_project_id, console_client_id, console_app_id, default_language, created_at, updated_at, deleted_at` +
	` FROM zitadel.instances`

// Get implements [domain.InstanceRepository].
func (i *instance) Get(ctx context.Context, opts ...database.Condition) (*domain.Instance, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(queryInstanceStmt)

	// return only non deleted isntances
	opts = append(opts, database.IsNull(i.DeletedAtColumn()))
	andCondition := database.And(opts...)
	i.writeCondition(&builder, andCondition)

	return scanInstance(ctx, i.client, &builder)
}

// List implements [domain.InstanceRepository].
func (i *instance) List(ctx context.Context, opts ...database.Condition) ([]*domain.Instance, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(queryInstanceStmt)

	// return only non deleted isntances
	opts = append(opts, database.IsNull(i.DeletedAtColumn()))
	andCondition := database.And(opts...)
	i.writeCondition(&builder, andCondition)

	return scanInstances(ctx, i.client, &builder)
}

const createInstanceStmt = `INSERT INTO zitadel.instances (id, name, default_org_id, iam_project_id, console_client_id, console_app_id, default_language)` +
	` VALUES ($1, $2, $3, $4, $5, $6, $7)` +
	` RETURNING created_at, updated_at`

// Create implements [domain.InstanceRepository].
func (i *instance) Create(ctx context.Context, instance *domain.Instance) error {
	builder := database.StatementBuilder{}
	builder.AppendArgs(instance.ID, instance.Name, instance.DefaultOrgID, instance.IAMProjectID, instance.ConsoleClientID, instance.ConsoleAppID, instance.DefaultLanguage)
	builder.WriteString(createInstanceStmt)

	err := i.client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&instance.CreatedAt, &instance.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// constraint violation
			if pgErr.Code == "23514" {
				if pgErr.ConstraintName == "instances_name_check" {
					return errors.New("instnace name not provided")
				}
				if pgErr.ConstraintName == "instances_id_check" {
					return errors.New("instnace id not provided")
				}
			}
			// duplicate
			if pgErr.Code == "23505" {
				if pgErr.ConstraintName == "instances_pkey" {
					return errors.New("instnace id already exists")
				}
			}
		}
	}
	return err
}

// Update implements [domain.InstanceRepository].
func (i instance) Update(ctx context.Context, condition database.Condition, changes ...database.Change) (int64, error) {
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.instances SET `)
	database.Changes(changes).Write(&builder)
	i.writeCondition(&builder, condition)

	stmt := builder.String()

	rowsAffected, err := i.client.Exec(ctx, stmt, builder.Args()...)
	return rowsAffected, err
}

// Delete implements [domain.InstanceRepository].
func (i instance) Delete(ctx context.Context, condition database.Condition) error {
	if condition == nil {
		return errors.New("Delete must contain a condition") // (otherwise ALL instances will be deleted)
	}
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.instances SET deleted_at = $1`)
	builder.AppendArgs(time.Now())

	i.writeCondition(&builder, condition)
	_, err := i.client.Exec(ctx, builder.String(), builder.Args()...)
	return err
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

// DeletedAtColumn implements [domain.instanceColumns].
func (instance) DeletedAtColumn() database.Column {
	return database.NewColumn("deleted_at")
}

func (i *instance) writeCondition(
	builder *database.StatementBuilder,
	condition database.Condition,
) {
	if condition == nil {
		return
	}
	builder.WriteString(" WHERE ")
	condition.Write(builder)
}

func scanInstance(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.Instance, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, nil
	}

	instance := new(domain.Instance)
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(instance); err != nil {
		// if no results returned, this is not a error
		// it just means the organization was not found
		// the caller should check if the returned organization is nil
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
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
		// if no results returned, this is not a error
		// it just means the organization was not found
		// the caller should check if the returned organization is nil
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return instances, nil
}
