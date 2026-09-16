package setup

import (
	"encoding/json"
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
		{"10_project_grant.sql", "project_grant", "g.granted_org_id || ':' || g.project_id", []string{"'org:' || g.granted_org_id", "'grant:' || g.grant_id"}},
		{"11_project_grant_member.sql", "project_grant_member", "m.project_id || ':' || m.user_id || ':' || m.grant_id", []string{"'grant:' || m.grant_id", "m.user_resource_owner", "m.granted_org"}},
		{"12_user_grant.sql", "user_grant", "g.resource_owner || ':' || g.user_id || ':' || g.project_id || ':' || g.grant_id", []string{"'grant:' || g.grant_id", "g.resource_owner_user", "g.granted_org"}},
		{"13_member_org.sql", "member", "m.org_id || ':' || m.user_id", []string{"m.user_resource_owner"}},
		{"14_member_project.sql", "member", "m.project_id || ':' || m.user_id", []string{"m.user_resource_owner"}},
		{"15_member_instance.sql", "member", "m.id || ':' || m.user_id", []string{"m.user_resource_owner"}},
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
			assert.Equal(t, tt.wantFinalized, mig.Finalized)
		})
	}
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
	}
	data, err := json.Marshal(mig)
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(data, &payload))
	assert.Equal(t, "v2.0.0", payload["version"])
	assert.Equal(t, true, payload["finalized"])
	assert.NotContains(t, payload, "ForceFinalize")
	assert.NotContains(t, payload, "forceFinalize")
}
