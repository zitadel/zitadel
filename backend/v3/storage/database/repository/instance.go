package repository

import (
	"context"
	"errors"
	"time"

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
// func (i *instance) Get(ctx context.Context, opts ...database.QueryOption) (*domain.Instance, error) {
func (i *instance) Get(ctx context.Context, opts ...database.Condition) (*domain.Instance, error) {
	i.builder = database.StatementBuilder{}

	i.builder.WriteString(queryInstanceStmt)

	isNotDeletedCondition := database.IsNull(i.DeletedAtColumn())
	opts = append(opts, isNotDeletedCondition)
	andCondition := database.And(opts...)
	andCondition.Write(&i.builder)

	return scanInstance(i.client.QueryRow(ctx, i.builder.String(), i.builder.Args()...))
}

const createInstanceStmt = `INSERT INTO zitadel.instances (id, name, default_org_id, iam_project_id, console_client_id, console_app_id, default_language)` +
	` VALUES ($1, $2, $3, $4, $5, $6, $7)` +
	` RETURNING created_at, updated_at`

// Create implements [domain.InstanceRepository].
func (i *instance) Create(ctx context.Context, instance *domain.Instance) error {
	i.builder = database.StatementBuilder{}
	i.builder.AppendArgs(instance.ID, instance.Name, instance.DefaultOrgID, instance.IAMProjectID, instance.ConsoleClientId, instance.ConsoleAppID, instance.DefaultLanguage)
	i.builder.WriteString(createInstanceStmt)

	return i.client.QueryRow(ctx, i.builder.String(), i.builder.Args()...).Scan(&instance.CreatedAt, &instance.UpdatedAt)
}

// Update implements [domain.InstanceRepository].
func (i instance) Update(ctx context.Context, condition database.Condition, changes ...database.Change) error {
	i.builder = database.StatementBuilder{}
	i.builder.WriteString(`UPDATE zitadel.instances SET `)
	database.Changes(changes).Write(&i.builder)
	i.writeCondition(condition)

	stmt := i.builder.String()

	return i.client.Exec(ctx, stmt, i.builder.Args()...)
}

// Delete implements [domain.InstanceRepository].
func (i instance) Delete(ctx context.Context, condition database.Condition) error {
	if condition == nil {
		return errors.New("Delete must contain a condition") // (otherwise ALL instances will be deleted)
	}
	i.builder = database.StatementBuilder{}
	i.builder.WriteString(`UPDATE zitadel.instances SET deleted_at = $1`)
	i.builder.AppendArgs(time.Now())

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

// DefaultOrgIdColumn implements [domain.instanceColumns].
func (instance) DefaultOrgIdColumn() database.Column {
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
		&instance.DefaultOrgID,
		&instance.IAMProjectID,
		&instance.ConsoleClientId,
		&instance.ConsoleAppID,
		&instance.DefaultLanguage,
		&instance.CreatedAt,
		&instance.UpdatedAt,
		&instance.DeletedAt,
	)
	if err != nil {
		// if no results returned, this is not a error
		// it just means the instance was not found
		// the caller should check if the returned instance is nil
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return &instance, nil
}
