package repository

import (
	"context"
	"slices"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	internaldb "github.com/zitadel/zitadel/internal/database"
)

var _ domain.AdministratorRepository = (*administrator)(nil)

type administrator struct{}

func AdministratorRepository() domain.AdministratorRepository {
	return new(administrator)
}

func (a administrator) unqualifiedTableName() string {
	return "administrators"
}

func (a administrator) qualifiedTableName() string {
	return "zitadel." + a.unqualifiedTableName()
}

func (a administrator) unqualifiedRolesTableName() string {
	return "administrator_roles"
}

const queryAdministratorStmt = "SELECT administrators.instance_id " +
	", administrators.id " +
	", administrators.user_id " +
	", administrators.scope " +
	", administrators.organization_id " +
	", administrators.project_id " +
	", administrators.project_grant_id " +
	", administrators.created_at " +
	", administrators.updated_at " +
	", ARRAY_AGG(administrator_roles.role_name ORDER BY administrator_roles.role_name) FILTER (WHERE administrator_roles.administrator_id IS NOT NULL) AS roles " +
	"FROM zitadel.administrators " +
	"LEFT JOIN zitadel.administrator_roles" +
	" ON administrators.instance_id = administrator_roles.instance_id" +
	" AND administrators.id = administrator_roles.administrator_id"

func (a administrator) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Administrator, error) {
	builder, err := a.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return scanAdministrator(ctx, client, builder)
}

func (a administrator) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Administrator, error) {
	builder, err := a.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return scanAdministrators(ctx, client, builder)
}

func (a administrator) prepareQuery(opts []database.QueryOption) (*database.StatementBuilder, error) {
	opts = append(opts, database.WithGroupBy(a.InstanceIDColumn(), a.IDColumn()))
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if err := checkRestrictingColumns(options.Condition, a.InstanceIDColumn()); err != nil {
		return nil, err
	}
	builder := database.NewStatementBuilder(queryAdministratorStmt)
	options.Write(builder)
	return builder, nil
}

func (a administrator) Create(ctx context.Context, client database.QueryExecutor, administrator *domain.Administrator) error {
	var (
		createdAt any = database.DefaultInstruction
		updatedAt any = database.DefaultInstruction
		roles         = rolesOrEmpty(administrator.Roles)
	)
	if !administrator.CreatedAt.IsZero() {
		createdAt = administrator.CreatedAt
	}
	if !administrator.UpdatedAt.IsZero() {
		updatedAt = administrator.UpdatedAt
	}

	builder := database.NewStatementBuilder("WITH inserted_admin AS (INSERT INTO zitadel.administrators (instance_id, user_id, scope, organization_id, project_id, project_grant_id, created_at, updated_at) VALUES (")
	builder.WriteArgs(
		administrator.InstanceID,
		administrator.UserID,
		administrator.Scope,
		administrator.OrganizationID,
		administrator.ProjectID,
		administrator.ProjectGrantID,
		createdAt,
		updatedAt,
	)
	builder.WriteString(") RETURNING id, instance_id, created_at, updated_at), inserted_roles AS (INSERT INTO zitadel.administrator_roles (instance_id, administrator_id, role_name) SELECT inserted_admin.instance_id, inserted_admin.id, unnest(")
	builder.WriteArg(roles)
	builder.WriteString("::TEXT[]) FROM inserted_admin) SELECT id, created_at, updated_at FROM inserted_admin")

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(
		&administrator.ID,
		&administrator.CreatedAt,
		&administrator.UpdatedAt,
	)
}

func (a administrator) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}
	if err := checkRestrictingColumns(condition, a.InstanceIDColumn(), a.UserIDColumn(), a.ScopeColumn()); err != nil {
		return 0, err
	}

	if !database.Changes(changes).IsOnColumn(a.UpdatedAtColumn()) {
		changes = append(changes, a.clearUpdatedAt())
	}

	var builder database.StatementBuilder
	builder.WriteString("WITH existing_administrator AS (SELECT * FROM zitadel.administrators ")
	writeCondition(&builder, condition)
	builder.WriteString(") ")
	for i, change := range changes {
		sessionCTE(change, i, 0, &builder)
	}
	builder.WriteString("UPDATE zitadel.administrators SET ")
	if err := database.Changes(changes).Write(&builder); err != nil {
		return 0, err
	}
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (a administrator) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if err := checkRestrictingColumns(condition, a.InstanceIDColumn(), a.UserIDColumn(), a.ScopeColumn()); err != nil {
		return 0, err
	}

	builder := database.NewStatementBuilder("DELETE FROM ")
	builder.WriteString(a.qualifiedTableName())
	builder.WriteRune(' ')
	writeCondition(builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (a administrator) SetUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(a.UpdatedAtColumn(), updatedAt)
}

func (a administrator) AddRole(role string) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("INSERT INTO zitadel.administrator_roles (instance_id, administrator_id, role_name) SELECT instance_id, id, ")
			builder.WriteArg(role)
			builder.WriteString(" FROM existing_administrator ON CONFLICT (instance_id, administrator_id, role_name) DO NOTHING")
		}, nil,
	)
}

func (a administrator) RemoveRole(role string) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("DELETE FROM zitadel.administrator_roles ar USING existing_administrator ea WHERE ar.instance_id = ea.instance_id AND ar.administrator_id = ea.id AND ar.role_name = ")
			builder.WriteArg(role)
		}, nil,
	)
}

func (a administrator) SetRoles(roles []string) database.Change {
	normalizedRoles := rolesOrEmpty(roles)
	return database.NewChanges(
		database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.administrator_roles (instance_id, administrator_id, role_name) SELECT instance_id, id, unnest(")
				builder.WriteArg(normalizedRoles)
				builder.WriteString("::text[]) FROM existing_administrator ON CONFLICT (instance_id, administrator_id, role_name) DO NOTHING")
			}, nil,
		),
		database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("DELETE FROM zitadel.administrator_roles ar USING existing_administrator ea WHERE ar.instance_id = ea.instance_id AND ar.administrator_id = ea.id AND NOT ar.role_name = ANY (")
				builder.WriteArg(normalizedRoles)
				builder.WriteString("::text[])")
			}, nil,
		),
	)
}

func (a administrator) PrimaryKeyCondition(instanceID, administratorID string) database.Condition {
	return database.And(
		a.InstanceIDCondition(instanceID),
		a.IDCondition(administratorID),
	)
}

func (a administrator) IDCondition(administratorID string) database.Condition {
	return database.NewTextCondition(a.IDColumn(), database.TextOperationEqual, administratorID)
}

func (a administrator) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(a.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

func (a administrator) UserIDCondition(userID string) database.Condition {
	return database.NewTextCondition(a.UserIDColumn(), database.TextOperationEqual, userID)
}

func (a administrator) ScopeCondition(scope domain.AdministratorScope) database.Condition {
	return database.NewTextCondition(a.ScopeColumn(), database.TextOperationEqual, scope.String())
}

func (a administrator) InstanceAdministratorCondition(instanceID, userID string) database.Condition {
	return database.And(
		a.InstanceIDCondition(instanceID),
		a.ScopeCondition(domain.AdministratorScopeInstance),
		a.UserIDCondition(userID),
	)
}

func (a administrator) OrganizationAdministratorCondition(instanceID, organizationID, userID string) database.Condition {
	return database.And(
		a.InstanceIDCondition(instanceID),
		a.ScopeCondition(domain.AdministratorScopeOrganization),
		a.OrganizationIDCondition(organizationID),
		a.UserIDCondition(userID),
	)
}

func (a administrator) ProjectAdministratorCondition(instanceID, projectID, userID string) database.Condition {
	return database.And(
		a.InstanceIDCondition(instanceID),
		a.ScopeCondition(domain.AdministratorScopeProject),
		a.ProjectIDCondition(projectID),
		a.UserIDCondition(userID),
	)
}

func (a administrator) ProjectGrantAdministratorCondition(instanceID, projectGrantID, userID string) database.Condition {
	return database.And(
		a.InstanceIDCondition(instanceID),
		a.ScopeCondition(domain.AdministratorScopeProjectGrant),
		a.ProjectGrantIDCondition(projectGrantID),
		a.UserIDCondition(userID),
	)
}

func (a administrator) OrganizationIDCondition(organizationID string) database.Condition {
	return database.NewTextCondition(a.OrganizationIDColumn(), database.TextOperationEqual, organizationID)
}

func (a administrator) ProjectIDCondition(projectID string) database.Condition {
	return database.NewTextCondition(a.ProjectIDColumn(), database.TextOperationEqual, projectID)
}

func (a administrator) ProjectGrantIDCondition(projectGrantID string) database.Condition {
	return database.NewTextCondition(a.ProjectGrantIDColumn(), database.TextOperationEqual, projectGrantID)
}

func (a administrator) RoleCondition(op database.TextOperation, role string) database.Condition {
	return database.NewTextCondition(database.NewColumn(a.unqualifiedRolesTableName(), "role_name"), op, role)
}

func (a administrator) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		a.InstanceIDColumn(),
		a.IDColumn(),
	}
}

func (a administrator) IDColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "id")
}

func (a administrator) InstanceIDColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "instance_id")
}

func (a administrator) UserIDColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "user_id")
}

func (a administrator) ScopeColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "scope")
}

func (a administrator) OrganizationIDColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "organization_id")
}

func (a administrator) ProjectIDColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "project_id")
}

func (a administrator) ProjectGrantIDColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "project_grant_id")
}

func (a administrator) CreatedAtColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "created_at")
}

func (a administrator) UpdatedAtColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "updated_at")
}

func (a administrator) clearUpdatedAt() database.Change {
	return database.NewChange(a.UpdatedAtColumn(), database.NullInstruction)
}

// rawAdministrator is used for query collection because ARRAY_AGG roles cannot be
// scanned directly into the domain model's []string field when the client is backed
// by [database/sql]. The intermediate TextArray field works with both [database/sql] and pgx.
type rawAdministrator struct {
	domain.Administrator
	Roles internaldb.TextArray[string] `db:"roles"`
}

func (a *rawAdministrator) toDomain() *domain.Administrator {
	a.Administrator.Roles = []string(a.Roles)
	return &a.Administrator
}

func scanAdministrator(ctx context.Context, client database.QueryExecutor, builder *database.StatementBuilder) (*domain.Administrator, error) {
	rows, err := client.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var administrator rawAdministrator
	if err = rows.(database.CollectableRows).CollectExactlyOneRow(&administrator); err != nil {
		return nil, err
	}
	return administrator.toDomain(), nil
}

func scanAdministrators(ctx context.Context, client database.QueryExecutor, builder *database.StatementBuilder) ([]*domain.Administrator, error) {
	rows, err := client.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	var administrators []rawAdministrator
	if err = rows.(database.CollectableRows).Collect(&administrators); err != nil {
		return nil, err
	}

	result := make([]*domain.Administrator, len(administrators))
	for i := range administrators {
		result[i] = administrators[i].toDomain()
	}
	return result, nil
}

func uniqueStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if slices.Contains(result, value) {
			continue
		}
		result = append(result, value)
	}
	return result
}

func rolesOrEmpty(roles []string) []string {
	if roles == nil {
		return []string{}
	}
	return uniqueStrings(roles)
}
