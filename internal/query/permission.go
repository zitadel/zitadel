package query

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
)

const (
	// eventstore.permitted_orgs(instanceid text, userid text, perm text)
	wherePermittedOrgsClause = "%s = ANY(eventstore.permitted_orgs(?, ?, ?))"
)

func wherePermittedOrgs(ctx context.Context, query sq.SelectBuilder, orgIDColumn, permission string) sq.SelectBuilder {
	return query.Where(
		fmt.Sprintf(wherePermittedOrgsClause, orgIDColumn),
		authz.GetInstance(ctx).InstanceID(),
		authz.GetCtxData(ctx).UserID,
		permission,
	)
}
