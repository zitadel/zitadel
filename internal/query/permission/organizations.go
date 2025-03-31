package permission

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
)

const (
	// eventstore.permitted_orgs(instanceid text, userid text, system_user_perms JSONB, perm text, filter_org text)
	wherePermittedOrgsExpr = "%s = ANY(eventstore.permitted_orgs(?, ?, ?, ?, ?))"
)

type permittedOrgsBuilder struct {
	orgIDColumn       string
	instanceID        string
	userID            string
	systemPermissions []authz.SystemUserPermissions
	permission        string
	orgID             string

	// options
	userIDCol string
}

func (b *permittedOrgsBuilder) clauses() sq.Or {
	clauses := sq.Or{
		sq.Expr(
			fmt.Sprintf(wherePermittedOrgsExpr, b.orgIDColumn),
			b.instanceID,
			b.userID,
			database.NewJSONArray(b.systemPermissions),
			b.permission,
			b.orgID,
		),
	}
	if b.userIDCol != "" {
		clauses = append(clauses, sq.Eq{
			b.userIDCol: b.userID,
		})
	}
	return clauses
}

type OrgsOption func(b *permittedOrgsBuilder)

// OwnedRowsOption allows rows to be returned of which the current user is the owner.
// Even if the user does not have an explicit permission for the organization.
// For example an authenticated user can always see his own user account.
func OwnedRowsOption(userIDColumn string) OrgsOption {
	return func(b *permittedOrgsBuilder) {
		b.userIDCol = userIDColumn
	}
}

// OrgsFilter sets a `WHERE` clause to query, which filters returned rows against organizations the
// current authenticated user has the requested permission to.
// filterOrgID may be used to optimize the permitted orgs function by limiting the returned organizations,
func OrgsFilter(ctx context.Context, query sq.SelectBuilder, orgIDCol, filterOrgID, permission string, options ...OrgsOption) sq.SelectBuilder {
	ctxData := authz.GetCtxData(ctx)

	b := &permittedOrgsBuilder{
		orgIDColumn:       orgIDCol,
		instanceID:        authz.GetInstance(ctx).InstanceID(),
		userID:            ctxData.UserID,
		systemPermissions: ctxData.SystemUserPermissions,
		permission:        permission,
		orgID:             filterOrgID,
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
		"user_id_colum", b.userIDCol,
	).Debug("permitted orgs check used")

	return query.Where(b.clauses())
}
