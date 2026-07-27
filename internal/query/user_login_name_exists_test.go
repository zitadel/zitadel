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

func TestExtractLoginNameEqualsFilter_TopLevel(t *testing.T) {
	loginNameQuery, err := NewUserLoginNameExistsQuery("user@org.localhost", TextEqualsIgnoreCase)
	require.NoError(t, err)
	orgQuery, err := NewUserResourceOwnerSearchQuery("org1", TextEquals)
	require.NoError(t, err)

	filter, remaining, ok := extractLoginNameEqualsFilter([]SearchQuery{loginNameQuery, orgQuery})
	require.True(t, ok)
	require.NotNil(t, filter)
	assert.Equal(t, "user", filter.username)
	require.Len(t, remaining, 1)
	assert.Equal(t, orgQuery, remaining[0])
}

func TestExtractLoginNameEqualsFilter_SkipsOrQuery(t *testing.T) {
	loginNameQuery, err := NewUserLoginNameExistsQuery("user@org.localhost", TextEqualsIgnoreCase)
	require.NoError(t, err)
	emailQuery, err := NewUserEmailSearchQuery("user@example.com", TextEqualsIgnoreCase)
	require.NoError(t, err)
	orQuery, err := NewOrQuery(loginNameQuery, emailQuery)
	require.NoError(t, err)

	_, remaining, ok := extractLoginNameEqualsFilter([]SearchQuery{orQuery})
	assert.False(t, ok)
	require.Len(t, remaining, 1)
	assert.Equal(t, orQuery, remaining[0])
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

	assert.Contains(t, sql, "INNER JOIN")
	assert.Contains(t, sql, "login_name_matches")
	assert.Contains(t, sql, "login_names3_users")
	assert.Contains(t, sql, "user_name_lower")
	assert.NotContains(t, sql, "login_name_lower")
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

	assert.Contains(t, sql, "login_name_matches")
	assert.Contains(t, sql, "u.user_name IN")
	assert.NotContains(t, sql, "user_name_lower")
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

	assert.Contains(t, sql, "login_name_matches")
	assert.Contains(t, sql, "resource_owner")
	assert.Contains(t, args, "org1")
	assert.Contains(t, args, "user")
}

func TestPrepareUsersQuery_LoginNameOrEmailDoesNotRewrite(t *testing.T) {
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
	sql, _, err := builder.ToSql()
	require.NoError(t, err)

	assert.NotContains(t, sql, "login_name_matches")
	assert.Contains(t, sql, "projections.login_names3")
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

	assert.Contains(t, sql, "SELECT DISTINCT")
	assert.Contains(t, sql, "user_metadata5")
}
