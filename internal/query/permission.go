package query

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
)

const (
	// eventstore.permitted_orgs(instanceid text, userid text, perm text, filter_orgs text)
	wherePermittedOrgsClause              = "%s = ANY(eventstore.permitted_orgs(?, ?, ?, ?))"
	wherePermittedOrgsOrCurrentUserClause = "(" + wherePermittedOrgsClause + " OR %s = '%s'" + ")"
)

// wherePermittedOrgs sets a `WHERE` clause to the query that filters the orgs
// for which the authenticated user has the requested permission for.
// The user ID is taken from the context.
//
// The `orgIDColumn` specifies the table column to which this filter must be applied,
// and is typically the `resource_owner` column in ZITADEL.
// We use full identifiers in the query builder so this function should be
// called with something like `UserResourceOwnerCol.identifier()` for example.
func wherePermittedOrgs(ctx context.Context, query sq.SelectBuilder, filterOrgIds, orgIDColumn, permission string) sq.SelectBuilder {
	userID := authz.GetCtxData(ctx).UserID
	logging.WithFields("permission_check_v2_flag", authz.GetFeatures(ctx).PermissionCheckV2, "org_id_column", orgIDColumn, "permission", permission, "user_id", userID).Debug("permitted orgs check used")

	return query.Where(
		fmt.Sprintf(wherePermittedOrgsClause, orgIDColumn),
		authz.GetInstance(ctx).InstanceID(),
		userID,
		permission,
		filterOrgIds,
	)
}

func wherePermittedOrgsOrCurrentUser(ctx context.Context, query sq.SelectBuilder, filterOrgIds, orgIDColumn, userIdColum, permission string) sq.SelectBuilder {
	userID := authz.GetCtxData(ctx).UserID
	fmt.Printf("userID = %+v\n", userID)
	logging.WithFields("permission_check_v2_flag", authz.GetFeatures(ctx).PermissionCheckV2, "org_id_column", orgIDColumn, "user_id_colum", userIdColum, "permission", permission, "user_id", userID).Debug("permitted orgs check used")

	return query.Where(
		fmt.Sprintf(wherePermittedOrgsOrCurrentUserClause, orgIDColumn, userIdColum, userID),
		authz.GetInstance(ctx).InstanceID(),
		userID,
		permission,
		filterOrgIds,
	)
}
