package query

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
)

func TestGroupExternalIDSearchQuery(t *testing.T) {
	query, err := NewGroupExternalIDSearchQuery("source:entra", "CaseSensitive' OR true --")
	require.NoError(t, err)
	sql, args, err := query.toQuery(sq.Select("id").From("projections.groups1")).ToSql()
	require.NoError(t, err)
	require.Contains(t, sql, "projections.group_external_ids1.instance_id = projections.groups1.instance_id")
	require.Contains(t, sql, "projections.group_external_ids1.resource_owner = projections.groups1.resource_owner")
	require.Contains(t, sql, "projections.group_external_ids1.scope = ?")
	require.Contains(t, sql, "projections.group_external_ids1.external_id = ?")
	require.NotContains(t, sql, "CaseSensitive")
	require.Equal(t, []interface{}{"source:entra", "CaseSensitive' OR true --"}, args)
}
