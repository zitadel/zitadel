package setup

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/migration"
	"github.com/zitadel/zitadel/internal/query/projection"
)

var _ migration.RepeatableMigration = (*BackfillUniqueConstraintOwners)(nil)

func TestBackfillUniqueConstraintOwnersStmts(t *testing.T) {
	statements, err := backfillUniqueConstraintOwnersStatements()
	require.NoError(t, err)

	want := []struct {
		file       string
		uniqueType string
		field      string
		table      string
		extra      []string
	}{
		{"01_usernames_org_scoped.sql", "usernames", "u.username || u.resource_owner", projection.UserTable, []string{"'org:' || u.resource_owner"}},
		{"02_usernames.sql", "usernames", "lower(u.username)", projection.UserTable, []string{"'org:' || u.resource_owner", "user_login_must_be_domain"}},
		{"03_external_idps.sql", "external_idps", "l.idp_id || l.external_user_id", projection.IDPUserLinkTable, []string{"'idp:' || l.idp_id"}},
		{"04_org_name.sql", "org_name", "o.name", projection.OrgProjectionTable, nil},
		{"05_org_domain.sql", "org_domain", "d.domain", projection.OrgDomainTable, []string{"d.is_verified"}},
		{"06_project_names.sql", "project_names", "p.name || p.resource_owner", projection.ProjectProjectionTable, nil},
		{"07_appname.sql", "appname", "a.name || ':' || a.project_id", projection.AppProjectionTable, []string{"'app:' || a.id"}},
		{"08_project_role.sql", "project_role", "r.role_key || ':' || r.project_id", projection.ProjectRoleProjectionTable, nil},
		{"09_entity_ids.sql", "entity_ids", "s.entity_id", projection.AppProjectionTable, []string{projection.AppSAMLTable, "'app:' || a.id"}},
		{"10_project_grant.sql", "project_grant", "g.granted_org_id || ':' || g.project_id", projection.ProjectGrantProjectionTable, []string{"'org:' || g.granted_org_id", "'grant:' || g.grant_id"}},
		{"11_project_grant_member.sql", "project_grant_member", "m.project_id || ':' || m.user_id || ':' || m.grant_id", projection.ProjectGrantMemberProjectionTable, []string{"'grant:' || m.grant_id", "m.user_resource_owner", "m.granted_org"}},
		{"12_user_grant.sql", "user_grant", "g.resource_owner || ':' || g.user_id || ':' || g.project_id || ':' || g.grant_id", projection.UserGrantProjectionTable, []string{"'grant:' || g.grant_id", "g.resource_owner_user", "g.granted_org"}},
		{"13_member_org.sql", "member", "m.org_id || ':' || m.user_id", projection.OrgMemberProjectionTable, []string{"m.user_resource_owner"}},
		{"14_member_project.sql", "member", "m.project_id || ':' || m.user_id", projection.ProjectMemberProjectionTable, []string{"m.user_resource_owner"}},
		{"15_member_instance.sql", "member", "m.id || ':' || m.user_id", projection.InstanceMemberProjectionTable, []string{"m.user_resource_owner"}},
		{"17_action_names.sql", "action_names", "a.name || ':' || a.resource_owner", projection.ActionTable, nil},
		{"18_idp_config_names.sql", "idp_config_names", "i.name || i.resource_owner", projection.IDPTable, []string{"'idp:' || i.id"}},
		{"19_idp_config_names_templates.sql", "idp_config_names", "t.name || t.resource_owner", projection.IDPTemplateTable, []string{"'idp:' || t.id"}},
		{"20_mail_text.sql", "mail_text", "t.aggregate_id || ':' || t.type || ':' || t.language", projection.MessageTextTable, []string{"t.aggregate_id <> t.instance_id"}},
	}

	require.Len(t, statements, len(want))
	for i, stmt := range statements {
		assert.Equal(t, want[i].file, stmt.file)
		assert.Contains(t, stmt.query, "unique_type = '"+want[i].uniqueType+"'")
		assert.Contains(t, stmt.query, want[i].field)
		assert.Contains(t, stmt.query, want[i].table)
		assert.Contains(t, stmt.query, "uc.unique_field =")
		assert.Contains(t, stmt.query, "LIMIT $2")
		assert.Contains(t, stmt.query, "($1::text[])[1]")
		assert.NotContains(t, stmt.query, "uc.unique_field IN (")
		assert.NotContains(t, stmt.query, "lower(uc.unique_field)")
		assert.NotContains(t, stmt.query, "{{.")
		for _, extra := range want[i].extra {
			assert.Contains(t, stmt.query, extra)
		}
	}
}

func TestBackfillUniqueConstraintOwners_CheckAndFinalize(t *testing.T) {
	tests := []struct {
		name          string
		forceFinalize bool
		lastRun       map[string]interface{}
		wantRun       bool
		wantFinalized bool
	}{
		{
			name:          "missing lastRun runs without finalizing",
			lastRun:       nil,
			wantRun:       true,
			wantFinalized: false,
		},
		{
			name:          "empty lastRun runs without finalizing",
			lastRun:       map[string]interface{}{},
			wantRun:       true,
			wantFinalized: false,
		},
		{
			name:          "missing lastRun with ForceFinalize finalizes",
			forceFinalize: true,
			lastRun:       nil,
			wantRun:       true,
			wantFinalized: true,
		},
		{
			name:          "same version not finalized skips",
			lastRun:       map[string]interface{}{"version": "v2.0.0", "finalized": false},
			wantRun:       false,
			wantFinalized: false,
		},
		{
			name:          "same version missing finalized skips",
			lastRun:       map[string]interface{}{"version": "v2.0.0"},
			wantRun:       false,
			wantFinalized: false,
		},
		{
			name:          "same version not finalized with ForceFinalize runs and finalizes",
			forceFinalize: true,
			lastRun:       map[string]interface{}{"version": "v2.0.0", "finalized": false},
			wantRun:       true,
			wantFinalized: true,
		},
		{
			name:          "same version already finalized skips",
			lastRun:       map[string]interface{}{"version": "v2.0.0", "finalized": true},
			wantRun:       false,
			wantFinalized: true,
		},
		{
			name:          "same version already finalized with ForceFinalize skips",
			forceFinalize: true,
			lastRun:       map[string]interface{}{"version": "v2.0.0", "finalized": true},
			wantRun:       false,
			wantFinalized: true,
		},
		{
			name:          "version changed runs and finalizes",
			lastRun:       map[string]interface{}{"version": "v1.0.0", "finalized": false},
			wantRun:       true,
			wantFinalized: true,
		},
		{
			name:          "version changed keeps finalized true",
			lastRun:       map[string]interface{}{"version": "v1.0.0", "finalized": true},
			wantRun:       true,
			wantFinalized: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mig := &BackfillUniqueConstraintOwners{Version: "v2.0.0", ForceFinalize: tt.forceFinalize}
			assert.Equal(t, tt.wantRun, mig.Check(tt.lastRun))
			assert.False(t, mig.Finalized, "Check must not set finalized before Execute")
			mig.finalizeAfterSuccess()
			assert.Equal(t, tt.wantFinalized, mig.Finalized)
		})
	}
}

func TestBackfillUniqueConstraintOwners_SequentialVersionChanges(t *testing.T) {
	mig := &BackfillUniqueConstraintOwners{Version: "v1.0.0"}
	require.True(t, mig.Check(nil))
	mig.finalizeAfterSuccess()
	assert.False(t, mig.Finalized)

	lastRun := map[string]interface{}{"version": mig.Version, "finalized": mig.Finalized}
	mig = &BackfillUniqueConstraintOwners{Version: "v2.0.0"}
	require.True(t, mig.Check(lastRun))
	mig.finalizeAfterSuccess()
	assert.True(t, mig.Finalized)

	lastRun = map[string]interface{}{"version": mig.Version, "finalized": mig.Finalized}
	mig = &BackfillUniqueConstraintOwners{Version: "v3.0.0"}
	require.True(t, mig.Check(lastRun))
	mig.finalizeAfterSuccess()
	assert.True(t, mig.Finalized)
}

func TestBackfillUniqueConstraintOwners_ForceFinalizeRestart(t *testing.T) {
	mig := &BackfillUniqueConstraintOwners{Version: "v1.0.0"}
	require.True(t, mig.Check(nil))
	mig.finalizeAfterSuccess()
	require.False(t, mig.Finalized)

	lastRun := map[string]interface{}{"version": "v1.0.0", "finalized": false}
	mig = &BackfillUniqueConstraintOwners{Version: "v1.0.0", ForceFinalize: true}
	require.True(t, mig.Check(lastRun), "restart with ForceFinalize should re-run")
	assert.False(t, mig.Finalized, "Check must not set finalized before Execute")
	mig.finalizeAfterSuccess()
	assert.True(t, mig.Finalized)

	lastRun = map[string]interface{}{"version": "v1.0.0", "finalized": true}
	mig = &BackfillUniqueConstraintOwners{Version: "v1.0.0", ForceFinalize: true}
	assert.False(t, mig.Check(lastRun), "later restart with ForceFinalize should skip")
}

func TestUniqueTypesWithOwnersOmitsMailText(t *testing.T) {
	assert.NotContains(t, eventstore.UniqueTypesWithOwners, "mail_text")
	assert.Contains(t, eventstore.UniqueTypesWithOwners, "usernames")
	assert.Contains(t, eventstore.UniqueTypesWithOwners, "idp_config_names")
}

func TestBackfillUniqueConstraintOwnersJSONOmitsForceFinalize(t *testing.T) {
	mig := &BackfillUniqueConstraintOwners{
		Version:       "v2.0.0",
		Finalized:     true,
		ForceFinalize: true,
		BatchSize:     10,
		lastVersion:   "v1.0.0",
		lastFinalized: true,
	}
	data, err := json.Marshal(mig)
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(data, &payload))
	assert.Equal(t, "v2.0.0", payload["version"])
	assert.Equal(t, true, payload["finalized"])
	assert.NotContains(t, payload, "ForceFinalize")
	assert.NotContains(t, payload, "forceFinalize")
	assert.NotContains(t, payload, "BatchSize")
	assert.NotContains(t, payload, "batchSize")
	assert.NotContains(t, payload, "lastVersion")
	assert.NotContains(t, payload, "lastFinalized")
	assert.False(t, strings.Contains(string(data), "v1.0.0"))
}
