package setup

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackfillUniqueConstraintOwnersStmts(t *testing.T) {
	var (
		projectNames  string
		idpConfigIDP  string
		idpConfigTmpl string
	)
	for _, stmt := range backfillUniqueConstraintOwnersStmts {
		switch {
		case strings.Contains(stmt, "unique_type = 'project_names'"):
			projectNames = stmt
		case strings.Contains(stmt, "projections.idps3"):
			idpConfigIDP = stmt
		case strings.Contains(stmt, "projections.idp_templates6"):
			idpConfigTmpl = stmt
		}
	}
	require.NotEmpty(t, projectNames)
	assert.Contains(t, projectNames, "p.name || p.resource_owner")
	require.NotEmpty(t, idpConfigIDP)
	assert.Contains(t, idpConfigIDP, "'idp:' || i.id")
	require.NotEmpty(t, idpConfigTmpl)
	assert.Contains(t, idpConfigTmpl, "'idp:' || t.id")
}
