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

func OrganizationRepository(client database.QueryExecutor) domain.OrganizationRepository {
	return &org{
		repository: repository{
			client: client,
		},
	}
}

const queryOrganizationStmt = `SELECT id, name, instance_id, state, created_at, updated_at, deleted_at` +
	` FROM zitadel.organizations`

// Get implements [domain.OrganizationRepository].
func (o *org) Get(ctx context.Context, id domain.OrgIdentifierCondition, instanceID string, conditions ...database.Condition) (*domain.Organization, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(queryOrganizationStmt)

	instanceIDCondition := o.InstanceIDCondition(instanceID)
	// don't update deleted organizations
	nonDeletedOrgs := database.IsNull(o.DeletedAtColumn())

	conditions = append(conditions, id, instanceIDCondition, nonDeletedOrgs)
	o.writeCondition(&builder, database.And(conditions...))

	return scanOrganization(ctx, o.client, &builder)
}

// List implements [domain.OrganizationRepository].
func (o *org) List(ctx context.Context, opts ...database.Condition) ([]*domain.Organization, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(queryOrganizationStmt)

	// return only non deleted organizations
	opts = append(opts, database.IsNull(o.DeletedAtColumn()))
	o.writeCondition(&builder, database.And(opts...))

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
					return errors.New("organization name not provided")
				}
				if pgErr.ConstraintName == "organizations_id_check" {
					return errors.New("organization id not provided")
				}
				if pgErr.ConstraintName == "organizations_instance_id_check" {
					return errors.New("instance id not provided")
				}
			}
			// duplicate
			if pgErr.Code == "23505" {
				if pgErr.ConstraintName == "organizations_pkey" {
					return errors.New("organization id already exists")
				}
				if pgErr.ConstraintName == "organizations_instance_id_name_key" {
					return errors.New("organization name already exists for instance")
				}
			}
			// invalid instance id
			if pgErr.Code == "23503" {
				if pgErr.ConstraintName == "organizations_instance_id_fkey" {
					return errors.New("invalid instance id")
				}
			}
		}
		return err
	}
	return nil
}

// Update implements [domain.OrganizationRepository].
func (o org) Update(ctx context.Context, id domain.OrgIdentifierCondition, instanceID string, changes ...database.Change) (int64, error) {
	if changes == nil {
		return 0, errors.New("Update must contain a condition") // (otherwise ALL organizations will be updated)
	}
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.organizations SET `)

	instanceIDCondition := o.InstanceIDCondition(instanceID)
	// don't update deleted organizations
	nonDeletedOrgs := database.IsNull(o.DeletedAtColumn())

	conditions := []database.Condition{id, instanceIDCondition, nonDeletedOrgs}
	database.Changes(changes).Write(&builder)
	o.writeCondition(&builder, database.And(conditions...))

	stmt := builder.String()

	rowsAffected, err := o.client.Exec(ctx, stmt, builder.Args()...)
	return rowsAffected, err
}

// Delete implements [domain.OrganizationRepository].
func (o org) Delete(ctx context.Context, id domain.OrgIdentifierCondition, instanceID string) (int64, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(`UPDATE zitadel.organizations SET deleted_at = $1`)
	builder.AppendArgs(time.Now())

	instanceIDCondition := o.InstanceIDCondition(instanceID)
	// don't update deleted organizations
	nonDeletedOrgs := database.IsNull(o.DeletedAtColumn())

	conditions := []database.Condition{id, instanceIDCondition, nonDeletedOrgs}
	o.writeCondition(&builder, database.And(conditions...))

	return o.client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetName implements [domain.organizationChanges].
func (o org) SetName(name string) database.Change {
	return database.NewChange(o.NameColumn(), name)
}

// SetState implements [domain.organizationChanges].
func (i org) SetState(state domain.OrgState) database.Change {
	return database.NewChange(i.StateColumn(), state)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// IDCondition implements [domain.organizationConditions].
func (o org) IDCondition(id string) domain.OrgIdentifierCondition {
	return database.NewTextCondition(o.IDColumn(), database.TextOperationEqual, id)
}

// NameCondition implements [domain.organizationConditions].
func (o org) NameCondition(name string) domain.OrgIdentifierCondition {
	// return database.NewTextCondition(o.NameColumn(), database.TextOperationEqualIgnoreCase, name)
	return database.NewTextCondition(o.NameColumn(), database.TextOperationEqual, name)
}

// InstanceIDCondition implements [domain.organizationConditions].
func (o org) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(o.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// StateCondition implements [domain.organizationConditions].
func (o org) StateCondition(state domain.OrgState) database.Condition {
	return database.NewTextCondition(o.StateColumn(), database.TextOperationEqual, state.String())
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
