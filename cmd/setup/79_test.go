package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/migration"
)

var _ migration.RepeatableMigration = (*BackfillUniqueConstraintOwners)(nil)

func TestBackfillUniqueConstraintOwnersStmts(t *testing.T) {
	statements, err := readStatements(backfillUniqueConstraintOwnersFS, "79")
	require.NoError(t, err)

	want := []struct {
		file       string
		uniqueType string
		field      string
		extra      []string
	}{
		{"01_usernames_org_scoped.sql", "usernames", "u.username || u.resource_owner", nil},
		{"02_usernames.sql", "usernames", "u.username)", nil},
		{"03_external_idps.sql", "external_idps", "l.idp_id || l.external_user_id", []string{"'idp:' || l.idp_id"}},
		{"04_org_name.sql", "org_name", "o.name", nil},
		{"05_org_domain.sql", "org_domain", "d.domain", nil},
		{"06_project_names.sql", "project_names", "p.name || p.resource_owner", nil},
		{"07_appname.sql", "appname", "a.name || ':' || a.project_id", nil},
		{"08_project_role.sql", "project_role", "r.role_key || ':' || r.project_id", nil},
		{"09_entity_ids.sql", "entity_ids", "s.entity_id", nil},
		{"10_project_grant.sql", "project_grant", "g.granted_org_id || ':' || g.project_id", []string{"'org:' || g.granted_org_id"}},
		{"11_project_grant_member.sql", "project_grant_member", "m.project_id || ':' || m.user_id || ':' || m.grant_id", []string{"'grant:' || m.grant_id"}},
		{"12_user_grant.sql", "user_grant", "g.resource_owner || ':' || g.user_id || ':' || g.project_id || ':' || g.grant_id", []string{"'grant:' || g.grant_id"}},
		{"13_member_org.sql", "member", "m.org_id || ':' || m.user_id", nil},
		{"14_member_project.sql", "member", "m.project_id || ':' || m.user_id", nil},
		{"15_member_instance.sql", "member", "m.id || ':' || m.user_id", nil},
		{"16_group_name.sql", "group_name", "g.name || ':' || g.resource_owner", nil},
		{"17_action_names.sql", "action_names", "a.name || ':' || a.resource_owner", nil},
		{"18_idp_config_names.sql", "idp_config_names", "i.name || i.resource_owner", []string{"'idp:' || i.id", "projections.idps3"}},
		{"19_idp_config_names_templates.sql", "idp_config_names", "t.name || t.resource_owner", []string{"'idp:' || t.id", "projections.idp_templates6"}},
		{"20_mail_text.sql", "mail_text", "t.aggregate_id || ':' || t.type || ':' || t.language", []string{"t.aggregate_id <> t.instance_id"}},
	}

	require.Len(t, statements, len(want))
	for i, stmt := range statements {
		assert.Equal(t, want[i].file, stmt.file)
		assert.Contains(t, stmt.query, "unique_type = '"+want[i].uniqueType+"'")
		assert.Contains(t, stmt.query, want[i].field)
		for _, extra := range want[i].extra {
			assert.Contains(t, stmt.query, extra)
		}
	}
}

func TestBackfillUniqueConstraintOwners_Check(t *testing.T) {
	mig := &BackfillUniqueConstraintOwners{Version: "v2.0.0"}

	assert.True(t, mig.Check(nil), "missing lastRun should run")
	assert.True(t, mig.Check(map[string]interface{}{}), "empty lastRun should run")
	assert.True(t, mig.Check(map[string]interface{}{"version": "v1.0.0"}), "different version should run")
	assert.False(t, mig.Check(map[string]interface{}{"version": "v2.0.0"}), "same version should skip")
}

func TestUniqueTypesWithOwnersOmitsMailText(t *testing.T) {
	assert.NotContains(t, eventstore.UniqueTypesWithOwners, "mail_text")
	assert.Contains(t, eventstore.UniqueTypesWithOwners, "usernames")
	assert.Contains(t, eventstore.UniqueTypesWithOwners, "idp_config_names")
}
