package query

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	domain_pkg "github.com/zitadel/zitadel/internal/domain"
)

const (
	// eventstore.permitted_orgs(instanceid text, userid text, system_user_perms JSONB, perm text, filter_org text)
	wherePermittedOrgsExpr = "%s = ANY(eventstore.permitted_orgs(?, ?, ?, ?, ?))"
)

type permissionClauseBuilder struct {
	orgIDColumn       Column
	instanceID        string
	userID            string
	systemPermissions []authz.SystemUserPermissions
	permission        string
	orgID             string
	overrides         []sq.Eq
}

func (b *permissionClauseBuilder) appendOverride(column string, value any) {
	b.overrides = append(b.overrides, sq.Eq{column: value})
}

func (b *permissionClauseBuilder) clauses() sq.Or {
	clauses := make(sq.Or, 1, len(b.overrides)+1)
	clauses[0] = sq.Expr(
		fmt.Sprintf(wherePermittedOrgsExpr, b.orgIDColumn.identifier()),
		b.instanceID,
		b.userID,
		database.NewJSONArray(b.systemPermissions),
		b.permission,
		b.orgID,
	)
	for _, include := range b.overrides {
		clauses = append(clauses, include)
	}
	return clauses
}

type PermissionOption func(b *permissionClauseBuilder)

// OwnedRowsPermissionOption allows rows to be returned of which the current user is the owner.
// Even if the user does not have an explicit permission for the organization.
// For example an authenticated user can always see his own user account.
func OwnedRowsPermissionOption(userIDColumn Column) PermissionOption {
	return func(b *permissionClauseBuilder) {
		b.appendOverride(userIDColumn.identifier(), b.userID)
	}
}

// OverridePermissionOption allows returning of rows where the value is matched.
// Even if the user does not have an explicit permission for the organization.
func OverridePermissionOption(column Column, value any) PermissionOption {
	return func(b *permissionClauseBuilder) {
		b.appendOverride(column.identifier(), value)
	}
}

// SingleOrgPermissionOption may be used to optimize the permitted orgs function by limiting the
// returned organizations, to the one used in the requested filters.
func SingleOrgPermissionOption(queries []SearchQuery) PermissionOption {
	return func(b *permissionClauseBuilder) {
		b.orgID = findTextEqualsQuery(b.orgIDColumn, queries)
	}
}

// PermissionClause sets a `WHERE` clause to query,
// which filters returned rows the current authenticated user has the requested permission to.
//
// Experimental: Work in progress. Currently only organization permissions are supported
func PermissionClause(ctx context.Context, query sq.SelectBuilder, enabled bool, orgIDCol Column, permission string, options ...PermissionOption) sq.SelectBuilder {
	if !enabled {
		return query
	}

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
		"overrides", b.overrides,
	).Debug("permitted orgs check used")

	return query.Where(b.clauses())
}

// PermissionV2 checks are enabled when the feature flag is set and the permission check function is not nil.
// When the permission check function is nil, it indicates a v1 API and no resource based permission check is needed.
func PermissionV2(ctx context.Context, cf domain_pkg.PermissionCheck) bool {
	return authz.GetFeatures(ctx).PermissionCheckV2 && cf != nil
}
