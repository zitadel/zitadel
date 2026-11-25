package query

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
)

const (
	// loginNamesCorrelation is how the view query joins the outer user.
	loginNamesCorrelation = "projections.login_names3.user_id = projections.users14.id"

	loginNamePatternMatches        = "SELECT projections.login_names3.user_id FROM projections.login_names3 WHERE projections.login_names3.instance_id = $%d AND "
	loginNamePatternMatchesInOrg   = loginNamePatternMatches + "projections.login_names3.resource_owner = $%d AND "
	loginNamePatternLikeIgnoreCase = "LOWER(projections.login_names3.login_name) LIKE $%d"
)

func TestNewUserLoginNameExistsQuery_PatternIsFilter(t *testing.T) {
	for _, comparison := range []TextComparison{
		TextStartsWith,
		TextStartsWithIgnoreCase,
		TextEndsWith,
		TextEndsWithIgnoreCase,
		TextContains,
		TextContainsIgnoreCase,
	} {
		t.Run(fmt.Sprintf("comparison_%d", comparison), func(t *testing.T) {
			qry, err := NewUserLoginNameExistsQuery("User", comparison)
			require.NoError(t, err)
			_, ok := qry.(*loginNamePatternFilter)
			assert.True(t, ok)
		})
	}
}

func TestLoginNamePatternFilter_UnboundIsCorrelated(t *testing.T) {
	qry, err := NewUserLoginNameExistsQuery("user", TextContainsIgnoreCase)
	require.NoError(t, err)

	sql, args, err := qry.comp().ToSql()
	require.NoError(t, err)
	assert.Contains(t, sql, loginNamesCorrelation)
	assert.Equal(t, []any{"%user%"}, args)
}

func TestPrepareUsersQuery_LoginNamePatternIsUncorrelated(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	tests := []struct {
		comparison TextComparison
		value      string
		like       string
		wantLike   string
	}{
		{TextContains, "User_1", "projections.login_names3.login_name LIKE", `%User\_1%`},
		{TextContainsIgnoreCase, "User_1", "LOWER(projections.login_names3.login_name) LIKE", `%user\_1%`},
		{TextStartsWith, "User", "projections.login_names3.login_name LIKE", "User%"},
		{TextStartsWithIgnoreCase, "User", "LOWER(projections.login_names3.login_name) LIKE", "user%"},
		{TextEndsWith, "@Org.Localhost", "projections.login_names3.login_name LIKE", "%@Org.Localhost"},
		{TextEndsWithIgnoreCase, "@Org.Localhost", "LOWER(projections.login_names3.login_name) LIKE", "%@org.localhost"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("comparison_%d", tt.comparison), func(t *testing.T) {
			loginNameQuery, err := NewUserLoginNameExistsQuery(tt.value, tt.comparison)
			require.NoError(t, err)
			q := &UserSearchQueries{
				SearchRequest: SearchRequest{Limit: 20},
				Queries:       []SearchQuery{loginNameQuery},
			}

			builder, _ := q.prepareUsersQuery(ctx, false)
			stmt, args, err := builder.ToSql()
			require.NoError(t, err)

			inner, outer := splitUsersPageSQL(t, stmt)
			// one array in the total, one in the page
			assert.Equal(t, 2, strings.Count(stmt, UserIDCol.identifier()+" = ANY(ARRAY(SELECT projections.login_names3.user_id FROM projections.login_names3 WHERE projections.login_names3.instance_id = $"))
			assert.Contains(t, inner, UserIDCol.identifier()+" = ANY(ARRAY(")
			assert.Equal(t, 2, strings.Count(stmt, tt.like+" $"))
			assert.NotContains(t, stmt, loginNamesCorrelation)
			assert.NotContains(t, stmt, "login_names3.resource_owner = projections.users14.resource_owner")
			assert.NotContains(t, outer, "ANY(ARRAY(")
			assert.Equal(t, 2, countArg(args, tt.wantLike))
			assert.Equal(t, 5, countArg(args, "inst-1"), "outer and pattern instance in the total and in the page, plus the page join")
		})
	}
}

func TestPrepareUsersQuery_LoginNamePatternWithOrg(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	loginNameQuery, err := NewUserLoginNameExistsQuery("User", TextContainsIgnoreCase)
	require.NoError(t, err)
	orgQuery, err := NewUserResourceOwnerSearchQuery("org1", TextEquals)
	require.NoError(t, err)

	q := &UserSearchQueries{
		SearchRequest: SearchRequest{Limit: 20},
		Queries:       []SearchQuery{loginNameQuery, orgQuery},
	}
	builder, _ := q.prepareUsersCountQuery(ctx, false)
	stmt, args, err := builder.ToSql()
	require.NoError(t, err)

	assert.Equal(t,
		"SELECT COUNT(*) FROM projections.users14 WHERE projections.users14.instance_id = $1"+
			" AND projections.users14.id = ANY(ARRAY("+fmt.Sprintf(loginNamePatternMatchesInOrg+loginNamePatternLikeIgnoreCase, 2, 3, 4)+"))"+
			" AND projections.users14.resource_owner = $5",
		stmt,
	)
	assert.Equal(t, []any{"inst-1", "inst-1", "org1", "%user%", "org1"}, args)
}

func TestPrepareUsersQuery_LoginNamePatternOrgOnlyFromTopLevelEquals(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	loginNameQuery, err := NewUserLoginNameExistsQuery("user", TextContainsIgnoreCase)
	require.NoError(t, err)
	org1, err := NewUserResourceOwnerSearchQuery("org1", TextEquals)
	require.NoError(t, err)
	org2, err := NewUserResourceOwnerSearchQuery("org2", TextEquals)
	require.NoError(t, err)
	orgs, err := NewOrQuery(org1, org2)
	require.NoError(t, err)
	orgIgnoreCase, err := NewUserResourceOwnerSearchQuery("ORG1", TextEqualsIgnoreCase)
	require.NoError(t, err)

	for name, orgFilter := range map[string]SearchQuery{
		"or of orgs":             orgs,
		"org equals ignore case": orgIgnoreCase,
	} {
		t.Run(name, func(t *testing.T) {
			q := &UserSearchQueries{Queries: []SearchQuery{loginNameQuery, orgFilter}}
			builder, _ := q.prepareUsersCountQuery(ctx, false)
			stmt, _, err := builder.ToSql()
			require.NoError(t, err)

			assert.Contains(t, stmt, "ANY(ARRAY("+fmt.Sprintf(loginNamePatternMatches+loginNamePatternLikeIgnoreCase, 2, 3)+"))")
			assert.NotContains(t, stmt, "projections.login_names3.resource_owner")
		})
	}
}

func TestPrepareUsersQuery_LoginNamePatternInOrIsHashableIn(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	loginNameQuery, err := NewUserLoginNameExistsQuery("user", TextContainsIgnoreCase)
	require.NoError(t, err)
	emailQuery, err := NewUserEmailSearchQuery("user", TextContainsIgnoreCase)
	require.NoError(t, err)
	orQuery, err := NewOrQuery(loginNameQuery, emailQuery)
	require.NoError(t, err)
	orgQuery, err := NewUserResourceOwnerSearchQuery("org1", TextEquals)
	require.NoError(t, err)

	q := &UserSearchQueries{Queries: []SearchQuery{orQuery, orgQuery}}
	builder, _ := q.prepareUsersCountQuery(ctx, false)
	stmt, args, err := builder.ToSql()
	require.NoError(t, err)

	assert.Equal(t,
		"SELECT COUNT(*) FROM projections.users14"+
			" LEFT JOIN projections.users14_humans ON projections.users14.id = projections.users14_humans.user_id AND projections.users14.instance_id = projections.users14_humans.instance_id"+
			" WHERE projections.users14.instance_id = $1"+
			" AND (projections.users14.id IN ("+fmt.Sprintf(loginNamePatternMatchesInOrg+loginNamePatternLikeIgnoreCase, 2, 3, 4)+")"+
			" OR LOWER(projections.users14_humans.email) LIKE $5)"+
			" AND projections.users14.resource_owner = $6",
		stmt,
	)
	assert.Equal(t, []any{"inst-1", "inst-1", "org1", "%user%", "%user%", "org1"}, args)
}

func TestPrepareUsersQuery_LoginNamePatternInNotIsHashableIn(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	loginNameQuery, err := NewUserLoginNameExistsQuery("user", TextContainsIgnoreCase)
	require.NoError(t, err)
	notQuery, err := NewNotQuery(loginNameQuery)
	require.NoError(t, err)

	q := &UserSearchQueries{Queries: []SearchQuery{notQuery}}
	builder, _ := q.prepareUsersCountQuery(ctx, false)
	stmt, args, err := builder.ToSql()
	require.NoError(t, err)

	assert.Equal(t,
		"SELECT COUNT(*) FROM projections.users14 WHERE projections.users14.instance_id = $1"+
			" AND NOT (projections.users14.id IN ("+fmt.Sprintf(loginNamePatternMatches+loginNamePatternLikeIgnoreCase, 2, 3)+"))",
		stmt,
	)
	assert.Equal(t, []any{"inst-1", "inst-1", "%user%"}, args)
}

func TestPrepareUsersQuery_LoginNamePatternInAndIsConjunct(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	loginNameQuery, err := NewUserLoginNameExistsQuery("user", TextContainsIgnoreCase)
	require.NoError(t, err)
	typeQuery, err := NewUserTypeSearchQuery(1)
	require.NoError(t, err)
	andQuery, err := NewAndQuery(typeQuery, loginNameQuery)
	require.NoError(t, err)

	q := &UserSearchQueries{Queries: []SearchQuery{andQuery}}
	builder, _ := q.prepareUsersCountQuery(ctx, false)
	stmt, _, err := builder.ToSql()
	require.NoError(t, err)

	assert.Contains(t, stmt, UserIDCol.identifier()+" = ANY(ARRAY(")
	assert.NotContains(t, stmt, " IN (")
}

func TestPrepareUsersQuery_LoginNamePatternKeepsCallerQueries(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	loginNameQuery, err := NewUserLoginNameExistsQuery("user", TextContainsIgnoreCase)
	require.NoError(t, err)
	emailQuery, err := NewUserEmailSearchQuery("user", TextContainsIgnoreCase)
	require.NoError(t, err)
	orQuery, err := NewOrQuery(loginNameQuery, emailQuery)
	require.NoError(t, err)

	q := &UserSearchQueries{Queries: []SearchQuery{loginNameQuery, orQuery}}
	_, _ = q.prepareUsersQuery(ctx, false)

	assert.Same(t, loginNameQuery, q.Queries[0])
	assert.Same(t, orQuery, q.Queries[1])
	assert.Same(t, loginNameQuery, orQuery.queries[0])
	assert.Nil(t, loginNameQuery.(*loginNamePatternFilter).binding)
}

func countArg(args []any, want any) int {
	n := 0
	for _, arg := range args {
		if arg == want {
			n++
		}
	}
	return n
}
