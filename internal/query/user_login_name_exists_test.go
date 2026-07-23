package query

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
)

func TestNewUserLoginNameExistsQuery_EqualsIgnoreCase(t *testing.T) {
	qry, err := NewUserLoginNameExistsQuery("User.Name@Org.Localhost", TextEqualsIgnoreCase)
	require.NoError(t, err)

	ln, ok := qry.(*loginNameExistsQuery)
	require.True(t, ok)
	assert.Equal(t, "user.name", ln.username)
	assert.Equal(t, "org.localhost", ln.domain)
	assert.Equal(t, "user.name@org.localhost", ln.loginName)
	assert.True(t, ln.ignoreCase)

	ln.instanceID = "inst-1"
	sql, args, err := ln.comp().ToSql()
	require.NoError(t, err)

	assert.Contains(t, sql, "projections.login_names3_users")
	assert.Contains(t, sql, "user_name_lower")
	assert.NotContains(t, sql, "login_name_lower")
	// outer IN uses users14.id, but the subquery must not correlate on it
	assert.True(t, strings.Count(sql, "projections.users14.id") == 1)
	assert.NotContains(t, sql, "u.user_id = projections.users14.id")
	assert.Equal(t, []interface{}{
		"inst-1", "inst-1", "org.localhost", "inst-1",
		"user.name", "user.name@org.localhost",
		"user.name", "org.localhost", "user.name@org.localhost",
	}, args)
}

func TestNewUserLoginNameExistsQuery_Equals(t *testing.T) {
	qry, err := NewUserLoginNameExistsQuery("User.Name@Org.Localhost", TextEquals)
	require.NoError(t, err)

	ln, ok := qry.(*loginNameExistsQuery)
	require.True(t, ok)
	assert.Equal(t, "User.Name", ln.username)
	assert.Equal(t, "Org.Localhost", ln.domain)
	assert.Equal(t, "User.Name@Org.Localhost", ln.loginName)
	assert.False(t, ln.ignoreCase)

	ln.instanceID = "inst-1"
	sql, args, err := ln.comp().ToSql()
	require.NoError(t, err)

	assert.Contains(t, sql, "u.user_name IN")
	assert.NotContains(t, sql, "user_name_lower")
	assert.Equal(t, []interface{}{
		"inst-1", "inst-1", "Org.Localhost", "inst-1",
		"User.Name", "User.Name@Org.Localhost",
		"User.Name", "Org.Localhost", "User.Name@Org.Localhost",
	}, args)
}

func TestNewUserLoginNameExistsQuery_ContainsFallsBackToView(t *testing.T) {
	qry, err := NewUserLoginNameExistsQuery("user", TextContains)
	require.NoError(t, err)

	_, ok := qry.(*loginNameExistsQuery)
	assert.False(t, ok)

	sql, _, err := qry.comp().ToSql()
	require.NoError(t, err)
	assert.Contains(t, sql, "projections.login_names3")
	assert.Contains(t, strings.ToLower(sql), "login_name")
}

func TestPrepareUsersQuery_LoginNameEqualsUsesIndexFriendlyFilter(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "inst-1")
	loginNameQuery, err := NewUserLoginNameExistsQuery("user165000@org.localhost", TextEqualsIgnoreCase)
	require.NoError(t, err)

	q := &UserSearchQueries{
		Queries: []SearchQuery{loginNameQuery},
	}
	builder, _ := q.prepareUsersQuery(ctx, false)
	sql, args, err := builder.ToSql()
	require.NoError(t, err)

	assert.Contains(t, sql, "login_names3_users")
	assert.Contains(t, sql, "user_name_lower")
	assert.NotContains(t, sql, "login_name_lower")
	assert.NotContains(t, sql, "user_metadata5")
	assert.NotContains(t, sql, "SELECT DISTINCT")
	assert.Contains(t, args, "inst-1")
	assert.Contains(t, args, "user165000")
	assert.Contains(t, args, "org.localhost")
}
