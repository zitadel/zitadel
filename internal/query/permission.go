package query

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
)

const (
	// eventstore.permitted_orgs(instanceid text, userid text, perm text)
	inSelectPermittedOrgs = "%s IN (SELECT eventstore.permitted_orgs(?, ?, ?))"
)

func whereInPermittedOrgs(ctx context.Context, query sq.SelectBuilder, orgIDColumn, permission string) sq.SelectBuilder {
	return query.Where(
		fmt.Sprintf(inSelectPermittedOrgs, orgIDColumn),
		authz.GetInstance(ctx).InstanceID(),
		authz.GetCtxData(ctx).UserID,
		permission,
	)
}
