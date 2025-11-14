package query

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
)

// ExamplePermissionClause_org shows how to use the PermissionClause function to filter
// permitted records based on the resource owner and the user's instance or organization membership.
func ExamplePermissionClause_org() {
	// These variables are typically set in the middleware of Zitadel.
	// They do not influence the generation of the clause, just what
	// the function does in Postgres.
	ctx := authz.WithInstanceID(context.Background(), "instanceID")
	ctx = authz.SetCtxData(ctx, authz.CtxData{
		UserID: "userID",
	})

	join, args := PermissionClause(
		ctx,
		UserResourceOwnerCol, // match the resource owner column
		domain.PermissionUserRead,
		SingleOrgPermissionOption([]SearchQuery{
			mustSearchQuery(NewUserDisplayNameSearchQuery("zitadel", TextContains)),
			mustSearchQuery(NewUserResourceOwnerSearchQuery("orgID", TextEquals)),
		}), // If the request had an orgID filter, it can be used to optimize the SQL function.
		OwnedRowsPermissionOption(UserIDCol), // allow user to find themselves.
	)

	sql, _, _ := sq.Select("*").
		From(userTable.identifier()).
		JoinClause(join, args...).
		Where(sq.Eq{
			UserInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
		}).ToSql()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM projections.users14 INNER JOIN eventstore.permitted_orgs(?, ?, ?, ?, ?) permissions ON (permissions.instance_permitted OR projections.users14.resource_owner = ANY(permissions.org_ids) OR projections.users14.id = ?) WHERE projections.users14.instance_id = ?
}

// ExamplePermissionClause_project shows how to use the PermissionClause function to filter
// permitted records based on the resource owner and the user's instance or organization membership.
// Additionally, it allows returning records based on the project ID and project membership.
func ExamplePermissionClause_project() {
	// These variables are typically set in the middleware of Zitadel.
	// They do not influence the generation of the clause, just what
	// the function does in Postgres.
	ctx := authz.WithInstanceID(context.Background(), "instanceID")
	ctx = authz.SetCtxData(ctx, authz.CtxData{
		UserID: "userID",
	})

	join, args := PermissionClause(
		ctx,
		ProjectColumnResourceOwner, // match the resource owner column
		"project.read",
		WithProjectsPermissionOption(ProjectColumnID),
		SingleOrgPermissionOption([]SearchQuery{
			mustSearchQuery(NewUserDisplayNameSearchQuery("zitadel", TextContains)),
			mustSearchQuery(NewUserResourceOwnerSearchQuery("orgID", TextEquals)),
		}), // If the request had an orgID filter, it can be used to optimize the SQL function.
	)

	sql, _, _ := sq.Select("*").
		From(projectsTable.identifier()).
		JoinClause(join, args...).
		Where(sq.Eq{
			ProjectColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		}).ToSql()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM projections.projects4 INNER JOIN eventstore.permitted_projects(?, ?, ?, ?, ?) permissions ON (permissions.instance_permitted OR projections.projects4.resource_owner = ANY(permissions.org_ids) OR projections.projects4.id = ANY(permissions.project_ids)) WHERE projections.projects4.instance_id = ?
}
