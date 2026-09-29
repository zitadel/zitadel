package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/action"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/member"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/target"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
)

func TestUniqueConstraintOwners(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name string
		cmd  eventstore.Command
		want [][]string
	}{
		{
			name: "username",
			cmd: &user.HumanAddedEvent{
				BaseEvent: eventstore.BaseEvent{Agg: agg(user.AggregateType, "user-1", "org-1")},
				UserName:  "alice",
			},
			want: [][]string{{"org:org-1", "user:user-1"}},
		},
		{
			name: "external idp link",
			cmd:  user.NewUserIDPLinkAddedEvent(ctx, agg(user.AggregateType, "user-1", "org-1"), "idp-1", "display", "ext-1"),
			want: [][]string{{"org:org-1", "idp:idp-1", "user:user-1"}},
		},
		{
			name: "org name",
			cmd:  org.NewOrgAddedEvent(ctx, agg(org.AggregateType, "org-1", "org-1"), "acme"),
			want: [][]string{{"org:org-1"}},
		},
		{
			name: "org domain",
			cmd:  org.NewDomainVerifiedEvent(ctx, agg(org.AggregateType, "org-1", "org-1"), "acme.example"),
			want: [][]string{{"org:org-1"}},
		},
		{
			name: "app name",
			cmd:  project.NewApplicationAddedEvent(ctx, agg(project.AggregateType, "project-1", "org-1"), "app-1", "console"),
			want: [][]string{{"org:org-1", "project:project-1", "app:app-1"}},
		},
		{
			name: "project role",
			cmd:  project.NewRoleAddedEvent(ctx, agg(project.AggregateType, "project-1", "org-1"), "role.key", "Role", "group"),
			want: [][]string{{"org:org-1", "project:project-1"}},
		},
		{
			name: "saml entity id",
			cmd:  project.NewSAMLConfigAddedEvent(ctx, agg(project.AggregateType, "project-1", "org-1"), "app-1", "https://sp.example", nil, "", 0, ""),
			want: [][]string{{"org:org-1", "project:project-1", "app:app-1"}},
		},
		{
			name: "project name",
			cmd: &project.ProjectAddedEvent{
				BaseEvent: eventstore.BaseEvent{Agg: agg(project.AggregateType, "project-1", "org-1")},
				Name:      "docs",
			},
			want: [][]string{{"org:org-1", "project:project-1"}},
		},
		{
			name: "project grant two orgs",
			cmd:  project.NewGrantAddedEvent(ctx, agg(project.AggregateType, "project-1", "granting-org"), "grant-1", "granted-org", nil),
			want: [][]string{{"org:granting-org", "org:granted-org", "project:project-1", "grant:grant-1"}},
		},
		{
			name: "project grant member",
			cmd: project.NewProjectGrantMemberAddedEvent(ctx, agg(project.AggregateType, "project-1", "org-1"), "user-1", "grant-1", "ROLE").
				WithOwnerOrgs("user-org", "granted-org"),
			want: [][]string{{"org:org-1", "org:user-org", "org:granted-org", "user:user-1", "project:project-1", "grant:grant-1"}},
		},
		{
			name: "user grant",
			cmd: usergrant.NewUserGrantAddedEvent(ctx, agg(user.AggregateType, "grant-row", "org-1"), "user-1", "project-1", "grant-1", nil).
				WithOwnerOrgs("user-org", "project-org", "granted-org"),
			want: [][]string{{"org:org-1", "org:user-org", "org:project-org", "org:granted-org", "user:user-1", "project:project-1", "grant:grant-1"}},
		},
		{
			name: "org member",
			cmd: member.NewMemberAddedEvent(
				&eventstore.BaseEvent{Agg: agg(org.AggregateType, "org-1", "org-1")},
				"user-1",
			).WithUserResourceOwner("org-1"),
			want: [][]string{{"user:user-1", "org:org-1"}},
		},
		{
			name: "org member cross-org user",
			cmd: member.NewMemberAddedEvent(
				&eventstore.BaseEvent{Agg: agg(org.AggregateType, "org-1", "org-1")},
				"user-1",
			).WithUserResourceOwner("user-org"),
			want: [][]string{{"user:user-1", "org:user-org", "org:org-1"}},
		},
		{
			name: "project member",
			cmd: member.NewMemberAddedEvent(
				&eventstore.BaseEvent{Agg: agg(project.AggregateType, "project-1", "org-1")},
				"user-1",
			).WithUserResourceOwner("org-1"),
			want: [][]string{{"user:user-1", "org:org-1", "project:project-1"}},
		},
		{
			name: "project member cross-org user",
			cmd: member.NewMemberAddedEvent(
				&eventstore.BaseEvent{Agg: agg(project.AggregateType, "project-1", "org-1")},
				"user-1",
			).WithUserResourceOwner("user-org"),
			want: [][]string{{"user:user-1", "org:user-org", "org:org-1", "project:project-1"}},
		},
		{
			name: "instance member",
			cmd: member.NewMemberAddedEvent(
				&eventstore.BaseEvent{Agg: agg(instance.AggregateType, "instance-1", "instance-1")},
				"user-1",
			).WithUserResourceOwner("user-org"),
			want: [][]string{{"user:user-1", "org:user-org"}},
		},
		{
			name: "action name",
			cmd:  action.NewAddedEvent(ctx, agg(action.AggregateType, "action-1", "org-1"), "notify", "", 0, false),
			want: [][]string{{"org:org-1"}},
		},
		{
			name: "idp config name",
			cmd:  idpconfig.NewIDPConfigAddedEvent(&eventstore.BaseEvent{Agg: agg(org.AggregateType, "org-1", "org-1")}, "idp-1", "Google", 0, 0, false),
			want: [][]string{{"org:org-1", "idp:idp-1"}},
		},
		{
			name: "mail text org",
			cmd:  policy.NewMailTextAddedEvent(&eventstore.BaseEvent{Agg: agg(org.AggregateType, "org-1", "org-1")}, "InitCode", "en", "", "", "", "", "", ""),
			want: [][]string{{"org:org-1"}},
		},
		{
			name: "mail text instance empty owners",
			cmd:  policy.NewMailTextAddedEvent(&eventstore.BaseEvent{Agg: agg(instance.AggregateType, "instance-1", "instance-1")}, "InitCode", "en", "", "", "", "", "", ""),
			want: [][]string{nil},
		},
		{
			name: "instance domain global empty owners",
			cmd:  instance.NewDomainAddedEvent(ctx, agg(instance.AggregateType, "instance-1", "instance-1"), "acme.example", true),
			want: [][]string{nil},
		},
		{
			name: "target instance only empty owners",
			cmd:  target.NewAddedEvent(ctx, agg(target.AggregateType, "target-1", "instance-1"), "hook", 0, "", 0, false, nil, 0),
			want: [][]string{nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraints := tt.cmd.UniqueConstraints()
			require.Len(t, constraints, len(tt.want))
			for i, constraint := range constraints {
				assert.Equal(t, eventstore.UniqueConstraintAdd, constraint.Action)
				assert.ElementsMatch(t, tt.want[i], constraint.Owners)
			}
		})
	}
}

func TestBulkRemoveKeyedAndByOwnerAreExclusive(t *testing.T) {
	ctx := t.Context()
	orgAgg := agg(org.AggregateType, "org-1", "org-1")
	userAgg := agg(user.AggregateType, "user-1", "org-1")
	projectAgg := agg(project.AggregateType, "project-1", "org-1")

	assertExclusive := func(t *testing.T, constraints []*eventstore.UniqueConstraint) {
		t.Helper()
		var hasKeyed, hasByOwner bool
		for _, constraint := range constraints {
			switch constraint.Action {
			case eventstore.UniqueConstraintRemove:
				hasKeyed = true
			case eventstore.UniqueConstraintRemoveByOwner:
				hasByOwner = true
			case eventstore.UniqueConstraintAdd, eventstore.UniqueConstraintInstanceRemove:
			}
		}
		assert.False(t, hasKeyed && hasByOwner, "UniqueConstraints must not mix keyed remove and RemoveByOwner")
	}

	orgKeyed := org.NewOrgRemovedEvent(ctx, orgAgg, "acme", nil, false, nil, nil, nil).UniqueConstraints()
	require.Len(t, orgKeyed, 1)
	assert.Equal(t, eventstore.UniqueConstraintRemove, orgKeyed[0].Action)
	assertExclusive(t, orgKeyed)
	orgByOwner := org.NewOrgRemovedByOwnerEvent(ctx, orgAgg, "acme").UniqueConstraints()
	require.Len(t, orgByOwner, 1)
	assert.Equal(t, eventstore.UniqueConstraintRemoveByOwner, orgByOwner[0].Action)
	assert.Equal(t, []string{"org:org-1"}, orgByOwner[0].Owners)
	assertExclusive(t, orgByOwner)

	keyedUser := user.NewUserRemovedEvent(ctx, userAgg, "alice", nil, false).UniqueConstraints()
	require.Len(t, keyedUser, 1)
	assert.Equal(t, eventstore.UniqueConstraintRemove, keyedUser[0].Action)
	assertExclusive(t, keyedUser)
	userByOwner := user.NewUserRemovedByOwnerEvent(ctx, userAgg).UniqueConstraints()
	require.Len(t, userByOwner, 1)
	assert.Equal(t, eventstore.UniqueConstraintRemoveByOwner, userByOwner[0].Action)
	assert.Equal(t, []string{"user:user-1"}, userByOwner[0].Owners)
	assertExclusive(t, userByOwner)
	assert.Empty(t, user.NewUserRemovedEvent(ctx, userAgg, "", nil, false).UniqueConstraints())

	idpKeyed := idpconfig.NewIDPConfigRemovedEvent(&eventstore.BaseEvent{Agg: orgAgg}, "idp-1", "Google").UniqueConstraints()
	require.Len(t, idpKeyed, 1)
	assert.Equal(t, eventstore.UniqueConstraintRemove, idpKeyed[0].Action)
	assertExclusive(t, idpKeyed)
	idpByOwner := idpconfig.NewIDPConfigRemovedByOwnerEvent(&eventstore.BaseEvent{Agg: orgAgg}, "idp-1", "Google").UniqueConstraints()
	require.Len(t, idpByOwner, 1)
	assert.Equal(t, eventstore.UniqueConstraintRemoveByOwner, idpByOwner[0].Action)
	assert.Equal(t, []string{"idp:idp-1"}, idpByOwner[0].Owners)
	assertExclusive(t, idpByOwner)

	projectKeyed := project.NewProjectRemovedEvent(ctx, projectAgg, "docs", nil).UniqueConstraints()
	require.Len(t, projectKeyed, 1)
	assert.Equal(t, eventstore.UniqueConstraintRemove, projectKeyed[0].Action)
	assertExclusive(t, projectKeyed)
	projectByOwner := project.NewProjectRemovedByOwnerEvent(ctx, projectAgg, "docs").UniqueConstraints()
	require.Len(t, projectByOwner, 1)
	assert.Equal(t, eventstore.UniqueConstraintRemoveByOwner, projectByOwner[0].Action)
	assert.Equal(t, []string{"project:project-1"}, projectByOwner[0].Owners)
	assertExclusive(t, projectByOwner)

	grantKeyed := project.NewGrantRemovedEvent(ctx, projectAgg, "grant-1", "granted-org").UniqueConstraints()
	require.Len(t, grantKeyed, 1)
	assert.Equal(t, eventstore.UniqueConstraintRemove, grantKeyed[0].Action)
	assertExclusive(t, grantKeyed)
	grantByOwner := project.NewGrantRemovedByOwnerEvent(ctx, projectAgg, "grant-1", "granted-org").UniqueConstraints()
	require.Len(t, grantByOwner, 1)
	assert.Equal(t, eventstore.UniqueConstraintRemoveByOwner, grantByOwner[0].Action)
	assert.Equal(t, []string{"grant:grant-1"}, grantByOwner[0].Owners)
	assertExclusive(t, grantByOwner)

	appByOwner := project.NewApplicationRemovedByOwnerEvent(ctx, projectAgg, "app-1").UniqueConstraints()
	require.Len(t, appByOwner, 1)
	assert.Equal(t, eventstore.UniqueConstraintRemoveByOwner, appByOwner[0].Action)
	assertExclusive(t, appByOwner)
	appKeyed := project.NewApplicationRemovedEvent(ctx, projectAgg, "app-1", "console", "").UniqueConstraints()
	require.Len(t, appKeyed, 1)
	assert.Equal(t, eventstore.UniqueConstraintRemove, appKeyed[0].Action)
	assertExclusive(t, appKeyed)
}

func TestSkipCommandUniqueConstraints(t *testing.T) {
	ctx := t.Context()
	orgAgg := agg(org.AggregateType, "org-1", "org-1")
	userAgg := agg(user.AggregateType, "user-1", "org-1")
	projectAgg := agg(project.AggregateType, "project-1", "org-1")
	grantAgg := agg(usergrant.AggregateType, "grant-row", "org-1")

	cmds := []eventstore.Command{
		usergrant.NewUserGrantCascadeRemovedEvent(ctx, grantAgg, "user-1", "project-1", "grant-1"),
		org.NewMemberCascadeRemovedEvent(ctx, orgAgg, "user-1"),
		project.NewProjectMemberCascadeRemovedEvent(ctx, projectAgg, "user-1"),
		instance.NewMemberCascadeRemovedEvent(ctx, agg(instance.AggregateType, "instance-1", "instance-1"), "user-1"),
		project.NewProjectGrantMemberCascadeRemovedEvent(ctx, projectAgg, "user-1", "grant-1"),
		user.NewUserIDPLinkCascadeRemovedEvent(ctx, userAgg, "idp-1", "ext-1"),
	}
	for _, cmd := range cmds {
		require.NotEmpty(t, cmd.UniqueConstraints())
		eventstore.SkipCommandUniqueConstraints(cmd)
		assert.Empty(t, cmd.UniqueConstraints())
	}
}

func TestUsernameScopeRewriteKeepsUserTag(t *testing.T) {
	constraints := user.NewUsernameChangedEvent(
		t.Context(),
		agg(user.AggregateType, "user-1", "org-1"),
		"alice",
		"alice",
		true,
		user.UsernameChangedEventWithPolicyChange(),
	).UniqueConstraints()
	require.Len(t, constraints, 2)
	assert.Equal(t, eventstore.UniqueConstraintRemove, constraints[0].Action)
	assert.Equal(t, eventstore.UniqueConstraintAdd, constraints[1].Action)
	assert.ElementsMatch(t, []string{"org:org-1", "user:user-1"}, constraints[1].Owners)
}

func agg(typ eventstore.AggregateType, id, resourceOwner string) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            id,
		Type:          typ,
		ResourceOwner: resourceOwner,
		InstanceID:    "instance-1",
		Version:       "v1",
	}
}
