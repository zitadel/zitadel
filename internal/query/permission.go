package query

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	domain_pkg "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	// eventstore.permitted_orgs(req_instance_id text, auth_user_id text, system_user_perms JSONB, perm text, filter_org text)
	joinPermittedOrgsFunction = `INNER JOIN eventstore.permitted_orgs(?, ?, ?, ?, ?) permissions ON `

	// eventstore.permitted_projects(req_instance_id text, auth_user_id text, system_user_perms JSONB, perm text, filter_org text)
	joinPermittedProjectsFunction = `INNER JOIN eventstore.permitted_projects(?, ?, ?, ?, ?) permissions ON `
)

// permissionClauseBuilder is used to build the SQL clause for permission checks.
// Don't use it directly, use the [PermissionClause] function with proper options instead.
type permissionClauseBuilder struct {
	orgIDColumn       Column
	instanceID        string
	userID            string
	systemPermissions []authz.SystemUserPermissions
	permission        string

	// optional fields
	orgID           *string
	projectIDColumn *Column
	connections     []sq.Eq
}

func (b *permissionClauseBuilder) appendConnection(column string, value any) {
	b.connections = append(b.connections, sq.Eq{column: value})
}

// joinFunction picks the correct SQL function and return the required arguments for that function.
func (b *permissionClauseBuilder) joinFunction() (sql string, args []any) {
	sql = joinPermittedOrgsFunction
	if b.projectIDColumn != nil {
		sql = joinPermittedProjectsFunction
	}
	return sql, []any{
		b.instanceID,
		b.userID,
		database.NewJSONArray(b.systemPermissions),
		b.permission,
		b.orgID,
	}
}

// joinConditions returns the conditions for the join,
// which are dynamic based on the provided options.
func (b *permissionClauseBuilder) joinConditions() sq.Or {
	conditions := make(sq.Or, 2, len(b.connections)+3)
	conditions[0] = sq.Expr("permissions.instance_permitted")
	conditions[1] = sq.Expr(b.orgIDColumn.identifier() + " = ANY(permissions.org_ids)")
	if b.projectIDColumn != nil {
		conditions = append(conditions,
			sq.Expr(b.projectIDColumn.identifier()+" = ANY(permissions.project_ids)"),
		)
	}
	for _, c := range b.connections {
		conditions = append(conditions, c)
	}
	return conditions
}

type PermissionOption func(b *permissionClauseBuilder)

// OwnedRowsPermissionOption allows rows to be returned of which the current user is the owner.
// Even if the user does not have an explicit permission for the organization.
// For example an authenticated user can always see his own user account.
// This option may be provided multiple times to allow matching with multiple columns.
// See [ConnectionPermissionOption] for more details.
func OwnedRowsPermissionOption(userIDColumn Column) PermissionOption {
	return func(b *permissionClauseBuilder) {
		b.appendConnection(userIDColumn.identifier(), b.userID)
	}
}

// ConnectionPermissionOption allows returning of rows where the value is matched.
// Even if the user does not have an explicit permission for the resource.
// Multiple connections may be provided.
// Each connection is applied in a OR condition, so if previous permissions are not met,
// matching rows are still returned for a later match.
func ConnectionPermissionOption(column Column, value any) PermissionOption {
	return func(b *permissionClauseBuilder) {
		b.appendConnection(column.identifier(), value)
	}
}

// SingleOrgPermissionOption may be used to optimize the permitted orgs function by limiting the
// returned organizations, to the one used in the requested filters.
func SingleOrgPermissionOption(queries []SearchQuery) PermissionOption {
	return func(b *permissionClauseBuilder) {
		orgID, ok := findTextEqualsQuery(b.orgIDColumn, queries)
		if ok {
			b.orgID = &orgID
		}
	}
}

// WithProjectsPermissionOption sets an additional filter against the project ID column,
// allowing for project specific permissions.
func WithProjectsPermissionOption(projectIDColumn Column) PermissionOption {
	return func(b *permissionClauseBuilder) {
		b.projectIDColumn = &projectIDColumn
	}
}

// PermissionClause builds a `INNER JOIN` clause which can be applied to a query builder.
// It filters returned rows the current authenticated user has the requested permission to.
// See permission_example_test.go for examples.
//
// Experimental: Work in progress. Currently only organization and project permissions are supported
// TODO: Add support for project grants.
func PermissionClause(ctx context.Context, orgIDCol Column, permission string, options ...PermissionOption) (string, []any) {
	ctxData := authz.GetCtxData(ctx)
	b := &permissionClauseBuilder{
		orgIDColumn:       orgIDCol,
		instanceID:        authz.GetInstance(ctx).InstanceID(),
		userID:            ctxData.UserID,
		systemPermissions: ctxData.SystemUserPermissions,
		permission:        permission,
	}
	for _, opt := range options {
		opt(b)
	}
	logging.WithFields(
		"org_id_column", b.orgIDColumn,
		"instance_id", b.instanceID,
		"user_id", b.userID,
		"system_user_permissions", b.systemPermissions,
		"permission", b.permission,
		"org_id", b.orgID,
		"project_id_column", b.projectIDColumn,
		"connections", b.connections,
	).Debug("permitted orgs check used")

	sql, args := b.joinFunction()
	conditions, conditionArgs, err := b.joinConditions().ToSql()
	if err != nil {
		// all cases are tested, no need to return an error.
		// If an error does happen, it's a bug and not a user error.
		panic(zerrors.ThrowInternal(err, "PERMISSION-OoS5o", "Errors.Internal"))
	}
	return sql + conditions, append(args, conditionArgs...)
}

// PermissionV2 checks are enabled when the feature flag is set and the permission check function is not nil.
// When the permission check function is nil, it indicates a v1 API and no resource based permission check is needed.
func PermissionV2(ctx context.Context, cf domain_pkg.PermissionCheck) bool {
	return authz.GetFeatures(ctx).PermissionCheckV2 && cf != nil
}
