package query

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
)

func TestNewUserLoginNameExistsQuery_EqualsIgnoreCaseIsMarker(t *testing.T) {
	qry, err := NewUserLoginNameExistsQuery("User.Name@Org.Localhost", TextEqualsIgnoreCase)
	require.NoError(t, err)

	ln, ok := qry.(*loginNameEqualsFilter)
	require.True(t, ok)
	assert.Equal(t, "user.name", ln.username)
	assert.Equal(t, "org.localhost", ln.domain)
	assert.Equal(t, "user.name@org.localhost", ln.loginName)
	assert.True(t, ln.ignoreCase)
}

func TestNewUserLoginNameExistsQuery_EqualsIsMarker(t *testing.T) {
	qry, err := NewUserLoginNameExistsQuery("User.Name@Org.Localhost", TextEquals)
	require.NoError(t, err)

	ln, ok := qry.(*loginNameEqualsFilter)
	require.True(t, ok)
	assert.Equal(t, "User.Name", ln.username)
	assert.Equal(t, "Org.Localhost", ln.domain)
	assert.Equal(t, "User.Name@Org.Localhost", ln.loginName)
	assert.False(t, ln.ignoreCase)
}

func TestNewUserLoginNameExistsQuery_ContainsFallsBackToView(t *testing.T) {
	qry, err := NewUserLoginNameExistsQuery("user", TextContains)
	require.NoError(t, err)

	_, ok := qry.(*loginNameEqualsFilter)
	assert.False(t, ok)

	sql, _, err := qry.comp().ToSql()
	require.NoError(t, err)
	assert.Contains(t, sql, "projections.login_names3")
	assert.Contains(t, strings.ToLower(sql), "login_name")
}

func TestExtractLoginEqualitySeeks_LoginNameAndOrg(t *testing.T) {
	loginNameQuery, err := NewUserLoginNameExistsQuery("user@org.localhost", TextEqualsIgnoreCase)
	require.NoError(t, err)
	orgQuery, err := NewUserResourceOwnerSearchQuery("org1", TextEquals)
	require.NoError(t, err)

	seeks, remaining, ok := extractLoginEqualitySeeks("inst-1", []SearchQuery{loginNameQuery, orgQuery})
	require.True(t, ok)
	require.Len(t, seeks, 1)
	assert.Contains(t, seeks[0].sql, "login_name_matches")
	require.Len(t, remaining, 1)
	assert.Equal(t, orgQuery, remaining[0])
}

func TestNewUserEmailSearchQuery_EqualsIsLowerEqualsFilter(t *testing.T) {
	qry, err := NewUserEmailSearchQuery("Ada@Example.com", TextEquals)
	require.NoError(t, err)
	_, ok := qry.(*lowerEqualsFilter)
	require.True(t, ok)

	contains, err := NewUserEmailSearchQuery("Ada", TextContains)
	require.NoError(t, err)
	_, ok = contains.(*lowerEqualsFilter)
	assert.False(t, ok)
}

func TestExtractLoginEqualitySeeks_OrLoginNameAndEmail(t *testing.T) {
	loginNameQuery, err := NewUserLoginNameExistsQuery("user@org.localhost", TextEqualsIgnoreCase)
	require.NoError(t, err)
	emailQuery, err := NewUserEmailSearchQuery("user@example.com", TextEqualsIgnoreCase)
	require.NoError(t, err)
	orQuery, err := NewOrQuery(loginNameQuery, emailQuery)
	require.NoError(t, err)

	seeks, remaining, ok := extractLoginEqualitySeeks("inst-1", []SearchQuery{orQuery})
	require.True(t, ok)
	require.Len(t, seeks, 2)
	assert.Contains(t, seeks[0].sql, "login_name_matches")
	assert.Contains(t, seeks[1].sql, "LOWER("+HumanEmailCol.identifier()+")")
	assert.Empty(t, remaining)
}

func TestExtractLoginEqualitySeeks_EmptyOrArmDoesNotExtract(t *testing.T) {
	emailQuery, err := NewUserEmailSearchQuery("", TextEqualsIgnoreCase)
	require.NoError(t, err)
	phoneQuery, err := NewUserPhoneSearchQuery("+41000000000", TextEquals)
	require.NoError(t, err)
	orQuery, err := NewOrQuery(emailQuery, phoneQuery)
	require.NoError(t, err)

	seeks, remaining, ok := extractLoginEqualitySeeks("inst-1", []SearchQuery{orQuery})
	assert.False(t, ok)
	assert.Nil(t, seeks)
	require.Len(t, remaining, 1)
	assert.Equal(t, orQuery, remaining[0])
}

func TestExtractLoginEqualitySeeks_NestedAndExtractsLoginName(t *testing.T) {
	loginNameQuery, err := NewUserLoginNameExistsQuery("user@org.localhost", TextEqualsIgnoreCase)
	require.NoError(t, err)
	orgQuery, err := NewUserResourceOwnerSearchQuery("org1", TextEquals)
	require.NoError(t, err)
	andQuery, err := NewAndQuery(loginNameQuery, orgQuery)
	require.NoError(t, err)

	seeks, remaining, ok := extractLoginEqualitySeeks("inst-1", []SearchQuery{andQuery})
	require.True(t, ok)
	require.Len(t, seeks, 1)
	assert.Contains(t, seeks[0].sql, "login_name_matches")
	require.Len(t, remaining, 1)
	assert.Equal(t, orgQuery, remaining[0])
}

func TestPrepareUsersQuery_LoginNameEqualsUsesIndexedJoin(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	loginNameQuery, err := NewUserLoginNameExistsQuery("user165000@org.localhost", TextEqualsIgnoreCase)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{loginNameQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, args, err := builder.ToSql()
	require.NoError(t, err)

	inner, outer := splitUsersPageSQL(t, sql)
	assert.Contains(t, inner, "INNER JOIN")
	assert.Contains(t, inner, "login_name_matches")
	assert.GreaterOrEqual(t, strings.Count(sql, "login_name_matches"), 2)
	assert.Contains(t, inner, "login_names3_users")
	assert.Contains(t, inner, "user_name_lower")
	assert.NotContains(t, inner, "login_name_lower")
	assert.NotContains(t, outer, "login_name_matches")
	assert.Contains(t, sql, "(SELECT COUNT(*)")
	assert.NotContains(t, sql, "COUNT(*) OVER ()")
	assert.NotContains(t, sql, "user_metadata5")
	assert.NotContains(t, sql, "SELECT DISTINCT")
	assert.Contains(t, args, "inst-1")
	assert.Contains(t, args, "user165000")
	assert.Contains(t, args, "org.localhost")
}

func TestPrepareUsersQuery_LoginNameEqualsCaseSensitive(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	loginNameQuery, err := NewUserLoginNameExistsQuery("User165000@Org.Localhost", TextEquals)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{loginNameQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, args, err := builder.ToSql()
	require.NoError(t, err)

	inner := usersPageSubquery(t, sql)
	assert.Contains(t, inner, "login_name_matches")
	assert.Contains(t, inner, "u.user_name IN")
	assert.NotContains(t, inner, "user_name_lower")
	assert.Contains(t, args, "User165000")
	assert.Contains(t, args, "Org.Localhost")
}

func TestPrepareUsersQuery_LoginNameEqualsWithOrgFilter(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	loginNameQuery, err := NewUserLoginNameExistsQuery("user@org.localhost", TextEqualsIgnoreCase)
	require.NoError(t, err)
	orgQuery, err := NewUserResourceOwnerSearchQuery("org1", TextEquals)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{loginNameQuery, orgQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, args, err := builder.ToSql()
	require.NoError(t, err)

	inner := usersPageSubquery(t, sql)
	assert.Contains(t, inner, "login_name_matches")
	assert.Contains(t, inner, "resource_owner")
	assert.Contains(t, args, "org1")
	assert.Contains(t, args, "user")
}

func TestPrepareUsersQuery_LoginNameOrEmailUsesUnion(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	loginNameQuery, err := NewUserLoginNameExistsQuery("user@org.localhost", TextEqualsIgnoreCase)
	require.NoError(t, err)
	emailQuery, err := NewUserEmailSearchQuery("user@example.com", TextEqualsIgnoreCase)
	require.NoError(t, err)
	orQuery, err := NewOrQuery(loginNameQuery, emailQuery)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{orQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, args, err := builder.ToSql()
	require.NoError(t, err)

	inner, outer := splitUsersPageSQL(t, sql)
	assert.Contains(t, inner, "UNION")
	assert.Contains(t, inner, "AS matches")
	assert.Contains(t, inner, "login_name_matches")
	assert.Contains(t, inner, "login_names3_users")
	assert.Contains(t, inner, "LOWER("+HumanEmailCol.identifier()+")")
	assert.NotContains(t, inner, "login_name_lower")
	assert.NotContains(t, outer, "AS matches")
	assert.Contains(t, args, "user@example.com")
}

func TestPrepareUsersQuery_UsernameOrEmailUsesUnion(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	usernameQuery, err := NewUserUsernameSearchQuery("user165000", TextEqualsIgnoreCase)
	require.NoError(t, err)
	emailQuery, err := NewUserEmailSearchQuery("user@example.com", TextEqualsIgnoreCase)
	require.NoError(t, err)
	orQuery, err := NewOrQuery(usernameQuery, emailQuery)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{orQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, args, err := builder.ToSql()
	require.NoError(t, err)

	inner, outer := splitUsersPageSQL(t, sql)
	assert.Contains(t, inner, "UNION")
	assert.Contains(t, inner, "AS matches")
	assert.Contains(t, inner, "LOWER("+UserUsernameCol.identifier()+")")
	assert.Contains(t, inner, "LOWER("+HumanEmailCol.identifier()+")")
	assert.NotContains(t, inner, "login_name_matches")
	assert.NotContains(t, outer, "AS matches")
	assert.Contains(t, args, "user165000")
	assert.Contains(t, args, "user@example.com")
}

func TestPrepareUsersQuery_EmailOrPhoneUsesUnion(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	emailQuery, err := NewUserEmailSearchQuery("user@example.com", TextEqualsIgnoreCase)
	require.NoError(t, err)
	phoneQuery, err := NewUserPhoneSearchQuery("+41000000000", TextEquals)
	require.NoError(t, err)
	orQuery, err := NewOrQuery(emailQuery, phoneQuery)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{orQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, args, err := builder.ToSql()
	require.NoError(t, err)

	inner, outer := splitUsersPageSQL(t, sql)
	assert.Contains(t, inner, "UNION")
	assert.Contains(t, inner, "AS matches")
	assert.Contains(t, inner, "LOWER("+HumanEmailCol.identifier()+")")
	assert.Contains(t, inner, "LOWER("+HumanPhoneCol.identifier()+")")
	assert.NotContains(t, outer, "AS matches")
	assert.Contains(t, args, "user@example.com")
	assert.Contains(t, args, "+41000000000")
}

func TestPrepareUsersQuery_EmailEqualsRechecksExactColumn(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	emailQuery, err := NewUserEmailSearchQuery("Ada@Example.com", TextEquals)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{emailQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, args, err := builder.ToSql()
	require.NoError(t, err)

	inner := usersPageSubquery(t, sql)
	assert.Contains(t, inner, "AS matches")
	assert.Contains(t, inner, "LOWER("+HumanEmailCol.identifier()+")")
	assert.Contains(t, inner, "LOWER($")
	assert.Contains(t, inner, HumanEmailCol.identifier()+" =")
	assert.Contains(t, args, "Ada@Example.com")
	assert.NotContains(t, args, "ada@example.com")
}

func TestPrepareUsersQuery_EmailIgnoreCaseDoesNotRecheckExactColumn(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	emailQuery, err := NewUserEmailSearchQuery("Ada@Example.com", TextEqualsIgnoreCase)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{emailQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, args, err := builder.ToSql()
	require.NoError(t, err)

	assert.Contains(t, sql, "AS matches")
	assert.Contains(t, sql, "LOWER("+HumanEmailCol.identifier()+")")
	assert.NotContains(t, sql, "LOWER($")
	assert.NotContains(t, sql, HumanEmailCol.identifier()+" =")
	assert.Contains(t, args, "ada@example.com")
	assert.NotContains(t, args, "Ada@Example.com")
}

func TestPrepareUsersQuery_PhoneIgnoreCaseUsesUnion(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	phoneQuery, err := NewUserPhoneSearchQuery("+41000000000", TextEqualsIgnoreCase)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{phoneQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, args, err := builder.ToSql()
	require.NoError(t, err)

	inner := usersPageSubquery(t, sql)
	assert.Contains(t, inner, "AS matches")
	assert.Contains(t, inner, "LOWER("+HumanPhoneCol.identifier()+")")
	assert.Contains(t, args, "+41000000000")
}

func TestPrepareUsersQuery_ContainsDoesNotUseUnion(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	usernameQuery, err := NewUserUsernameSearchQuery("user", TextContainsIgnoreCase)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{usernameQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, _, err := builder.ToSql()
	require.NoError(t, err)

	assert.NotContains(t, sql, "AS matches")
	assert.NotContains(t, sql, " UNION ")
}

func TestPrepareUsersQuery_DisplayNameDoesNotUseUnion(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	displayNameQuery, err := NewUserDisplayNameSearchQuery("Ada Lovelace", TextEqualsIgnoreCase)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{displayNameQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, _, err := builder.ToSql()
	require.NoError(t, err)

	assert.NotContains(t, sql, "AS matches")
	assert.NotContains(t, sql, " UNION ")
}

func TestPrepareUsersQuery_UsernameOrContainsDoesNotUseUnion(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	usernameQuery, err := NewUserUsernameSearchQuery("user165000", TextEqualsIgnoreCase)
	require.NoError(t, err)
	displayNameQuery, err := NewUserDisplayNameSearchQuery("user", TextContainsIgnoreCase)
	require.NoError(t, err)
	orQuery, err := NewOrQuery(usernameQuery, displayNameQuery)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{orQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, _, err := builder.ToSql()
	require.NoError(t, err)

	assert.NotContains(t, sql, "AS matches")
	assert.NotContains(t, sql, " UNION ")
}

func TestPrepareUsersQuery_UsernameOrUsernameUsesUnion(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	alice, err := NewUserUsernameSearchQuery("alice", TextEqualsIgnoreCase)
	require.NoError(t, err)
	bob, err := NewUserUsernameSearchQuery("bob", TextEqualsIgnoreCase)
	require.NoError(t, err)
	orQuery, err := NewOrQuery(alice, bob)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{orQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, args, err := builder.ToSql()
	require.NoError(t, err)

	inner := usersPageSubquery(t, sql)
	assert.Contains(t, inner, "UNION")
	assert.Contains(t, inner, "AS matches")
	assert.Equal(t, 2, strings.Count(inner, "LOWER("+UserUsernameCol.identifier()+")"))
	assert.Contains(t, args, "alice")
	assert.Contains(t, args, "bob")
}

func TestPrepareUsersQuery_SiblingUsernameAndKeepsRemainingWhere(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	alice, err := NewUserUsernameSearchQuery("alice", TextEqualsIgnoreCase)
	require.NoError(t, err)
	bob, err := NewUserUsernameSearchQuery("bob", TextEqualsIgnoreCase)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{alice, bob},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, args, err := builder.ToSql()
	require.NoError(t, err)

	assert.Contains(t, sql, "AS matches")
	assert.NotContains(t, sql, " UNION ")
	assert.Contains(t, args, "alice")
	assert.Contains(t, args, "bob")
}

func TestPrepareUsersQuery_NestedAndOfOrUsesUnion(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	alice, err := NewUserUsernameSearchQuery("alice", TextEqualsIgnoreCase)
	require.NoError(t, err)
	bob, err := NewUserUsernameSearchQuery("bob", TextEqualsIgnoreCase)
	require.NoError(t, err)
	orQuery, err := NewOrQuery(alice, bob)
	require.NoError(t, err)
	orgQuery, err := NewUserResourceOwnerSearchQuery("org1", TextEquals)
	require.NoError(t, err)
	andQuery, err := NewAndQuery(orQuery, orgQuery)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{andQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, args, err := builder.ToSql()
	require.NoError(t, err)

	inner := usersPageSubquery(t, sql)
	assert.Contains(t, inner, "AS matches")
	assert.Contains(t, inner, " UNION ")
	assert.Contains(t, inner, "resource_owner")
	assert.Contains(t, args, "org1")
	assert.Contains(t, args, "alice")
	assert.Contains(t, args, "bob")
}

func TestPrepareUsersQuery_NestedOrOfAndDoesNotUseUnion(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	alice, err := NewUserUsernameSearchQuery("alice", TextEqualsIgnoreCase)
	require.NoError(t, err)
	orgQuery, err := NewUserResourceOwnerSearchQuery("org1", TextEquals)
	require.NoError(t, err)
	andQuery, err := NewAndQuery(alice, orgQuery)
	require.NoError(t, err)
	bob, err := NewUserUsernameSearchQuery("bob", TextEqualsIgnoreCase)
	require.NoError(t, err)
	orQuery, err := NewOrQuery(andQuery, bob)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{orQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, _, err := builder.ToSql()
	require.NoError(t, err)

	assert.NotContains(t, sql, "AS matches")
	assert.NotContains(t, sql, " UNION ")
}

func TestPrepareUsersQuery_MetadataFilterKeepsDistinctJoin(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	metadataQuery, err := NewUserMetadataKeySearchQuery("key", TextContains)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{metadataQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, _, err := builder.ToSql()
	require.NoError(t, err)

	inner := usersPageSubquery(t, sql)
	assert.Contains(t, inner, "SELECT DISTINCT")
	assert.Contains(t, inner, "user_metadata5")
	assert.Contains(t, sql, "COUNT(DISTINCT")
	assert.NotContains(t, inner, "COUNT(")
	assert.NotContains(t, sql, "COUNT(*) OVER ()")
	distinctIdx := strings.Index(inner, "SELECT DISTINCT")
	idsIdx := strings.Index(inner, ") AS ids")
	if distinctIdx < 0 || idsIdx < 0 || idsIdx <= distinctIdx {
		t.Fatalf("distinct ids subquery not found: %s", inner)
	}
	assert.NotContains(t, inner[distinctIdx:idsIdx], "COUNT(")
}

func TestPrepareUsersQuery_PaginatesBeforeLoginNamesJoin(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	q := &UserSearchQueries{
		SearchRequest: SearchRequest{Limit: 20, Offset: 17500},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, _, err := builder.ToSql()
	require.NoError(t, err)

	inner, outer := splitUsersPageSQL(t, sql)
	assert.Contains(t, inner, "LIMIT 20")
	assert.Contains(t, inner, "OFFSET 17500")
	assert.Contains(t, sql, "(SELECT COUNT(*)")
	assert.NotContains(t, inner, "SELECT COUNT(*)")
	assert.NotContains(t, sql, "COUNT(*) OVER ()")
	assert.NotContains(t, inner, "login_names3")
	assert.NotContains(t, inner, "users14_humans")
	assert.Contains(t, outer, "LEFT JOIN LATERAL")
	assert.Contains(t, outer, "login_names3")
	assert.NotContains(t, outer, "SELECT COUNT(*)")
}

func TestPrepareUsersQuery_DisplayNameFilterJoinsHumansInPage(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	displayNameQuery, err := NewUserDisplayNameSearchQuery("Ada Lovelace", TextEqualsIgnoreCase)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{displayNameQuery},
		SearchRequest: SearchRequest{
			SortingColumn: HumanDisplayNameCol,
			Limit:         20,
		},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, _, err := builder.ToSql()
	require.NoError(t, err)

	inner := usersPageSubquery(t, sql)
	assert.Contains(t, inner, "users14_humans")
	assert.Contains(t, inner, "LOWER(projections.users14_humans.display_name)")
	assert.Contains(t, inner, "ORDER BY sort_col DESC, projections.users14.id DESC")
	assert.Contains(t, sql, "ORDER BY page.sort_col DESC, page.id DESC")
}

func TestPrepareUsersCountQuery_SkipsDisplayJoins(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	orgQuery, err := NewUserResourceOwnerSearchQuery("org1", TextEquals)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{orgQuery},
	}
	builder, _ := q.prepareUsersCountQuery(ctx, false)
	sql, args, err := builder.ToSql()
	require.NoError(t, err)

	assert.Contains(t, sql, "SELECT COUNT(*)")
	assert.Contains(t, sql, "resource_owner")
	assert.NotContains(t, sql, "login_names3")
	assert.NotContains(t, sql, "users14_humans")
	assert.NotContains(t, sql, "users14_machines")
	assert.NotContains(t, sql, "COUNT(*) OVER ()")
	assert.Contains(t, args, "org1")
	assert.Contains(t, args, "inst-1")
}

func TestPrepareUsersCountQuery_MetadataUsesDistinct(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	metadataQuery, err := NewUserMetadataKeySearchQuery("key", TextContains)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{metadataQuery},
	}
	builder, _ := q.prepareUsersCountQuery(ctx, false)
	sql, _, err := builder.ToSql()
	require.NoError(t, err)

	assert.Contains(t, sql, "COUNT(DISTINCT")
	assert.Contains(t, sql, "user_metadata5")
	assert.NotContains(t, sql, "login_names3")
}

func TestPrepareUsersCountQuery_LoginNameEqualsUsesIndexedJoin(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	loginNameQuery, err := NewUserLoginNameExistsQuery("user@org.localhost", TextEqualsIgnoreCase)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{loginNameQuery},
	}
	builder, _ := q.prepareUsersCountQuery(ctx, false)
	sql, args, err := builder.ToSql()
	require.NoError(t, err)

	assert.Contains(t, sql, "SELECT COUNT(*)")
	assert.Contains(t, sql, "AS matches")
	assert.Contains(t, sql, "login_name_matches")
	assert.Contains(t, sql, "login_names3_users")
	assert.NotContains(t, sql, "ARRAY_AGG")
	assert.NotContains(t, sql, "AS login_names ON TRUE")
	assert.NotContains(t, sql, "COUNT(*) OVER ()")
	assert.Contains(t, args, "user")
	assert.Contains(t, args, "org.localhost")
}

func TestPrepareUsersCountQuery_LoginNameOrEmailUsesUnion(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	loginNameQuery, err := NewUserLoginNameExistsQuery("user@org.localhost", TextEqualsIgnoreCase)
	require.NoError(t, err)
	emailQuery, err := NewUserEmailSearchQuery("user@example.com", TextEqualsIgnoreCase)
	require.NoError(t, err)
	orQuery, err := NewOrQuery(loginNameQuery, emailQuery)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{orQuery},
	}
	builder, _ := q.prepareUsersCountQuery(ctx, false)
	sql, args, err := builder.ToSql()
	require.NoError(t, err)

	assert.Contains(t, sql, "SELECT COUNT(*)")
	assert.Contains(t, sql, "UNION")
	assert.Contains(t, sql, "AS matches")
	assert.Contains(t, sql, "login_name_matches")
	assert.Contains(t, sql, "LOWER("+HumanEmailCol.identifier()+")")
	assert.NotContains(t, sql, "ARRAY_AGG")
	assert.NotContains(t, sql, "AS login_names ON TRUE")
	assert.Contains(t, args, "user@example.com")
}

func usersPageSubquery(t *testing.T, stmt string) string {
	t.Helper()
	inner, _ := splitUsersPageSQL(t, stmt)
	return inner
}

func splitUsersPageSQL(t *testing.T, stmt string) (inner, outer string) {
	t.Helper()
	fromIdx := strings.Index(stmt, "FROM (SELECT")
	pageIdx := strings.Index(stmt, ") AS page")
	if fromIdx < 0 || pageIdx < 0 || pageIdx <= fromIdx {
		t.Fatalf("page subquery not found: %s", stmt)
	}
	return stmt[fromIdx:pageIdx], stmt[pageIdx:]
}
