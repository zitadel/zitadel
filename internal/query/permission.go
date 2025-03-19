package query

import (
	"context"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
)

const (
	// eventstore.permitted_orgs(instanceid text, userid text, system_user_perms JSONB, perm text filter_orgs text)
	wherePermittedOrgsClause              = "%s = ANY(eventstore.permitted_orgs(?, ?, ?, ?, ?))"
	wherePermittedOrgsOrCurrentUserClause = "(" + wherePermittedOrgsClause + " OR %s = ?" + ")"
)

// wherePermittedOrgs sets a `WHERE` clause to the query that filters the orgs
// for which the authenticated user has the requested permission for.
// The user ID is taken from the context.
// The `orgIDColumn` specifies the table column to which this filter must be applied,
// and is typically the `resource_owner` column in ZITADEL.
// We use full identifiers in the query builder so this function should be
// called with something like `UserResourceOwnerCol.identifier()` for example.
func wherePermittedOrgs(ctx context.Context, query sq.SelectBuilder, systemUserPermission []authz.SystemUserPermissionsDBQuery, filterOrgIds, orgIDColumn, permission string) (sq.SelectBuilder, error) {
	userID := authz.GetCtxData(ctx).UserID
	logging.WithFields("permission_check_v2_flag", authz.GetFeatures(ctx).PermissionCheckV2, "org_id_column", orgIDColumn, "permission", permission, "user_id", userID).Debug("permitted orgs check used")

	systemUserPermissionsJson := "[]"
	if systemUserPermission != nil {
		systemUserPermissionsBytes, err := json.Marshal(systemUserPermission)
		if err != nil {
			return query, err
		}
		systemUserPermissionsJson = string(systemUserPermissionsBytes)
	}

	return query.Where(
		fmt.Sprintf(wherePermittedOrgsClause, orgIDColumn),
		authz.GetInstance(ctx).InstanceID(),
		userID,
		systemUserPermissionsJson,
		systemUserPermission,
		permission,
		filterOrgIds,
	), nil
}

func wherePermittedOrgsOrCurrentUser(ctx context.Context, query sq.SelectBuilder, systemUserPermission []authz.SystemUserPermissionsDBQuery, filterOrgIds, orgIDColumn, userIdColum, permission string) (sq.SelectBuilder, error) {
	userID := authz.GetCtxData(ctx).UserID
	logging.WithFields("permission_check_v2_flag", authz.GetFeatures(ctx).PermissionCheckV2, "org_id_column", orgIDColumn, "user_id_colum", userIdColum, "permission", permission, "user_id", userID).Debug("permitted orgs check used")

	systemUserPermissionsJson := "[]"
	if systemUserPermission != nil {
		systemUserPermissionsBytes, err := json.Marshal(systemUserPermission)
		if err != nil {
			return query, err
		}
		systemUserPermissionsJson = string(systemUserPermissionsBytes)
	}
	fmt.Printf("@@ >>>>>>>>>>>>>>>>>>>>>>>>>>>> authz.GetInstance(ctx).InstanceID() = %+v\n", authz.GetInstance(ctx).InstanceID())
	fmt.Printf("@@ >>>>>>>>>>>>>>>>>>>>>>>>>>>> systemUserPermissionsJson = %+v\n", systemUserPermissionsJson)
	fmt.Printf("@@ >>>>>>>>>>>>>>>>>>>>>>>>>>>> permission = %+v\n", permission)

	return query.Where(
		fmt.Sprintf(wherePermittedOrgsOrCurrentUserClause, orgIDColumn, userIdColum),
		authz.GetInstance(ctx).InstanceID(),
		userID,
		systemUserPermissionsJson,
		permission,
		filterOrgIds,
		userID,
	), nil
}
