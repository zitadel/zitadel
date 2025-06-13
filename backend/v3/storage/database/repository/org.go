package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

var _ domain.OrganizationRepository = (*org)(nil)

type org struct {
	repository
}

func OrgRepository(client database.QueryExecutor) domain.OrganizationRepository {
	return &org{
		repository: repository{
			client: client,
		},
	}
}

const queryOrganizationStmt = `SELECT id, name, instance_id, state, created_at, updated_at, deleted_at` +
	` FROM zitadel.organizations`

// Get implements [domain.OrganizationRepository].
func (o *org) Get(ctx context.Context, opts ...database.Condition) (*domain.Organization, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(queryOrganizationStmt)

	// return only non deleted isntances
	opts = append(opts, database.IsNull(o.DeletedAtColumn()))
	andCondition := database.And(opts...)
	o.writeCondition(&builder, andCondition)

	// rows, err := o.client.Query(ctx, builder.String(), builder.Args()...)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	return scanOrganization(ctx, o.client, &builder)
}

// List implements [domain.OrganizationRepository].
func (o *org) List(ctx context.Context, opts ...database.Condition) ([]*domain.Organization, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(queryOrganizationStmt)

	// return only non deleted isntances
	opts = append(opts, database.IsNull(o.DeletedAtColumn()))
	andCondition := database.And(opts...)
	o.writeCondition(&builder, andCondition)

	// rows, err := o.client.Query(ctx, builder.String(), builder.Args()...)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	return scanOrganizations(ctx, o.client, &builder)
}

const createOrganizationStmt = `INSERT INTO zitadel.organizations (id, name, instance_id, state)` +
	` VALUES ($1, $2, $3, $4)` +
	` RETURNING created_at, updated_at`

// Create implements [domain.OrganizationRepository].
func (o *org) Create(ctx context.Context, organization *domain.Organization) error {
	builder := database.StatementBuilder{}
	builder.AppendArgs(organization.ID, organization.Name, organization.InstanceID, organization.State)
	builder.WriteString(createOrganizationStmt)

	err := o.client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&organization.CreatedAt, &organization.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// constraint violation
			if pgErr.Code == "23514" {
				if pgErr.ConstraintName == "organizations_name_check" {
					return errors.New("instnace name not provided")
				}
				if pgErr.ConstraintName == "organizations_id_check" {
					return errors.New("instnace id not provided")
				}
			}
			// duplicate
			if pgErr.Code == "23505" {
				if pgErr.ConstraintName == "organizations_pkey" {
					return errors.New("instnace id already exists")
				}
			}
		}
	}
	return err
}

// Update implements [domain.OrganizationRepository].
func (o org) Update(ctx context.Context, condition database.Condition, changes ...database.Change) (int64, error) {
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.organizations SET `)
	database.Changes(changes).Write(&builder)
	o.writeCondition(&builder, condition)

	stmt := builder.String()

	rowsAffected, err := o.client.Exec(ctx, stmt, builder.Args()...)
	return rowsAffected, err
}

// Delete implements [domain.OrganizationRepository].
func (o org) Delete(ctx context.Context, condition database.Condition) error {
	if condition == nil {
		return errors.New("Delete must contain a condition") // (otherwise ALL organizations will be deleted)
	}
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.organizations SET deleted_at = $1`)
	builder.AppendArgs(time.Now())

	o.writeCondition(&builder, condition)
	_, err := o.client.Exec(ctx, builder.String(), builder.Args()...)
	return err
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetName implements [domain.organizationChanges].
func (i org) SetName(name string) database.Change {
	return database.NewChange(i.NameColumn(), name)
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

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// IDColumn implements [domain.organizationColumns].
func (org) IDColumn() database.Column {
	return database.NewColumn("id")
}

// NameColumn implements [domain.organizationColumns].
func (org) NameColumn() database.Column {
	return database.NewColumn("name")
}

// InstanceIDColumn implements [domain.organizationColumns].
func (org) InstanceIDColumn() database.Column {
	return database.NewColumn("instance_id")
}

// StateColumn implements [domain.organizationColumns].
func (org) StateColumn() database.Column {
	return database.NewColumn("state")
}

// CreatedAtColumn implements [domain.organizationColumns].
func (org) CreatedAtColumn() database.Column {
	return database.NewColumn("created_at")
}

// UpdatedAtColumn implements [domain.organizationColumns].
func (org) UpdatedAtColumn() database.Column {
	return database.NewColumn("updated_at")
}

// DeletedAtColumn implements [domain.organizationColumns].
func (org) DeletedAtColumn() database.Column {
	return database.NewColumn("deleted_at")
}

func (o *org) writeCondition(
	builder *database.StatementBuilder,
	condition database.Condition,
) {
	if condition == nil {
		return
	}
	builder.WriteString(" WHERE ")
	condition.Write(builder)
}

func scanOrganization(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.Organization, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	organization := &domain.Organization{}
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(organization); err != nil {
		// if no results returned, this is not a error
		// it just means the organization was not found
		// the caller should check if the returned organization is nil
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return organization, nil
}

func scanOrganizations(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.Organization, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	organizations := []*domain.Organization{}
	if err := rows.(database.CollectableRows).Collect(&organizations); err != nil {
		// if no results returned, this is not a error
		// it just means the organization was not found
		// the caller should check if the returned organization is nil
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return organizations, nil
}
