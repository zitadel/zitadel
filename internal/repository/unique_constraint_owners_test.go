package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/action"
	"github.com/zitadel/zitadel/internal/repository/group"
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
			want: [][]string{{"org:org-1", "project:project-1"}},
		},
		{
			name: "project role",
			cmd:  project.NewRoleAddedEvent(ctx, agg(project.AggregateType, "project-1", "org-1"), "role.key", "Role", "group"),
			want: [][]string{{"org:org-1", "project:project-1"}},
		},
		{
			name: "saml entity id",
			cmd:  project.NewSAMLConfigAddedEvent(ctx, agg(project.AggregateType, "project-1", "org-1"), "app-1", "https://sp.example", nil, "", 0, ""),
			want: [][]string{{"org:org-1", "project:project-1"}},
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
			want: [][]string{{"org:granting-org", "org:granted-org", "project:project-1"}},
		},
		{
			name: "project grant member",
			cmd:  project.NewProjectGrantMemberAddedEvent(ctx, agg(project.AggregateType, "project-1", "org-1"), "user-1", "grant-1", "ROLE"),
			want: [][]string{{"org:org-1", "user:user-1", "project:project-1", "grant:grant-1"}},
		},
		{
			name: "user grant",
			cmd:  usergrant.NewUserGrantAddedEvent(ctx, agg(user.AggregateType, "grant-row", "org-1"), "user-1", "project-1", "grant-1", nil),
			want: [][]string{{"org:org-1", "user:user-1", "project:project-1", "grant:grant-1"}},
		},
		{
			name: "org member",
			cmd: &member.MemberAddedEvent{
				BaseEvent: eventstore.BaseEvent{Agg: agg(org.AggregateType, "org-1", "org-1")},
				UserID:    "user-1",
			},
			want: [][]string{{"user:user-1", "org:org-1"}},
		},
		{
			name: "project member",
			cmd: &member.MemberAddedEvent{
				BaseEvent: eventstore.BaseEvent{Agg: agg(project.AggregateType, "project-1", "org-1")},
				UserID:    "user-1",
			},
			want: [][]string{{"user:user-1", "org:org-1", "project:project-1"}},
		},
		{
			name: "instance member",
			cmd: &member.MemberAddedEvent{
				BaseEvent: eventstore.BaseEvent{Agg: agg(instance.AggregateType, "instance-1", "instance-1")},
				UserID:    "user-1",
			},
			want: [][]string{{"user:user-1"}},
		},
		{
			name: "group name",
			cmd:  group.NewGroupAddedEvent(ctx, agg(group.AggregateType, "group-1", "org-1"), "engineering", ""),
			want: [][]string{{"org:org-1"}},
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

func TestBulkRemoveAlsoEmitsRemoveByOwner(t *testing.T) {
	ctx := t.Context()
	orgAgg := agg(org.AggregateType, "org-1", "org-1")
	orgConstraints := org.NewOrgRemovedEvent(ctx, orgAgg, "acme", nil, false, nil, nil, nil).UniqueConstraints()
	require.GreaterOrEqual(t, len(orgConstraints), 2)
	assert.Equal(t, eventstore.UniqueConstraintRemoveByOwner, orgConstraints[1].Action)
	assert.Equal(t, []string{"org:org-1"}, orgConstraints[1].Owners)

	userAgg := agg(user.AggregateType, "user-1", "org-1")
	userConstraints := user.NewUserRemovedEvent(ctx, userAgg, "alice", nil, false).UniqueConstraints()
	last := userConstraints[len(userConstraints)-1]
	assert.Equal(t, eventstore.UniqueConstraintRemoveByOwner, last.Action)
	assert.Equal(t, []string{"user:user-1"}, last.Owners)

	idpConstraints := idpconfig.NewIDPConfigRemovedEvent(&eventstore.BaseEvent{Agg: orgAgg}, "idp-1", "Google").UniqueConstraints()
	assert.Equal(t, eventstore.UniqueConstraintRemoveByOwner, idpConstraints[1].Action)
	assert.Equal(t, []string{"idp:idp-1"}, idpConstraints[1].Owners)

	projectAgg := agg(project.AggregateType, "project-1", "org-1")
	projectConstraints := project.NewProjectRemovedEvent(ctx, projectAgg, "docs", nil).UniqueConstraints()
	assert.Equal(t, eventstore.UniqueConstraintRemoveByOwner, projectConstraints[1].Action)
	assert.Equal(t, []string{"project:project-1"}, projectConstraints[1].Owners)
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
