package command

import (
	"context"
	"errors"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	repo "github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_MutateGroup(t *testing.T) {
	ctx := context.Background()
	initial := func() expect {
		return expectFilterScoped(t, "org1", eventFromEventPusher(addNewGroupEvent("group1", "org1")), eventFromEventPusher(newGroupUserAddedEvent("group1", "org1", "alice")))
	}
	renamed := repo.NewGroupChangedEvent(ctx, &repo.NewAggregate("group1", "org1").Aggregate, "group1", []repo.GroupChanges{repo.ChangeName(gu.Ptr("Engineering"))})
	tests := []struct {
		name      string
		mutations []GroupMutation
		expects   []expect
		denied    bool
		wantErr   bool
	}{
		{name: "rename and membership events share one push", mutations: []GroupMutation{{Name: gu.Ptr("Engineering"), AddMembers: []string{"bob"}, RemoveMembers: []string{"alice"}}}, expects: []expect{initial(), expectFilterScoped(t, "org1", eventFromEventPusher(addNewUserEvent("bob", "org1"))), expectGroupPush(renamed, newGroupUserAddedEvent("group1", "org1", "bob"), newGroupUserRemovedEvent("group1", "org1", "alice"))}},
		{name: "failed member validation never pushes rename", mutations: []GroupMutation{{Name: gu.Ptr("Engineering")}, {AddMembers: []string{"missing"}}}, expects: []expect{initial(), expectFilterScoped(t, "org1")}, wantErr: true},
		{name: "duplicate add is idempotent", mutations: []GroupMutation{{AddMembers: []string{"alice", "alice"}}}, expects: []expect{initial(), expectFilterScoped(t, "org1", eventFromEventPusher(addNewUserEvent("alice", "org1")))}},
		{name: "missing removal is idempotent", mutations: []GroupMutation{{RemoveMembers: []string{"missing"}}}, expects: []expect{initial()}},
		{name: "add then remove has no net event", mutations: []GroupMutation{{AddMembers: []string{"bob"}}, {RemoveMembers: []string{"bob"}}}, expects: []expect{initial(), expectFilterScoped(t, "org1", eventFromEventPusher(addNewUserEvent("bob", "org1")))}},
		{name: "replace removes omitted members", mutations: []GroupMutation{{Members: gu.Ptr([]string{"bob"})}}, expects: []expect{initial(), expectFilterScoped(t, "org1", eventFromEventPusher(addNewUserEvent("bob", "org1"))), expectGroupPush(newGroupUserAddedEvent("group1", "org1", "bob"), newGroupUserRemovedEvent("group1", "org1", "alice"))}},
		{name: "empty replacement clears members", mutations: []GroupMutation{{Members: gu.Ptr([]string{})}}, expects: []expect{initial(), expectGroupPush(newGroupUserRemovedEvent("group1", "org1", "alice"))}},
		{name: "empty removal preserves members", mutations: []GroupMutation{{RemoveMembers: []string{}}}, expects: []expect{initial()}},
		{name: "organization boundary", mutations: []GroupMutation{{Name: gu.Ptr("Engineering")}}, expects: []expect{expectFilterScoped(t, "org1")}, wantErr: true},
		{name: "permission denied", mutations: []GroupMutation{{RemoveMembers: []string{"alice"}}}, expects: []expect{initial()}, denied: true, wantErr: true},
		{name: "push failure returned", mutations: []GroupMutation{{RemoveMembers: []string{"alice"}}}, expects: []expect{initial(), expectGroupPushFailed(errors.New("push failed"), newGroupUserRemovedEvent("group1", "org1", "alice"))}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{eventstore: expectEventstore(tt.expects...)(t), checkPermission: newMockPermissionCheckAllowed()}
			if tt.denied {
				c.checkPermission = newMockPermissionCheckNotAllowed()
			}
			_, err := c.MutateGroup(ctx, "group1", "org1", tt.mutations)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCommands_CreateGroup_invalidMemberDoesNotCreateGroup(t *testing.T) {
	c := &Commands{eventstore: expectEventstore(
		expectFilter(eventFromEventPusher(org.NewOrgAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate, "org1"))),
		expectFilter(),
		expectFilterScoped(t, "org1"),
	)(t), checkPermission: newMockPermissionCheckAllowed()}
	_, err := c.CreateGroup(context.Background(), &CreateGroup{ObjectRoot: models.ObjectRoot{AggregateID: "group1", ResourceOwner: "org1"}, Name: "Sales", UserIDs: []string{"missing"}})
	require.True(t, zerrors.IsPreconditionFailed(err))
}

func TestCommands_MutateGroup_externalID(t *testing.T) {
	ctx := context.Background()
	binding := func(scope, old, value string) *repo.GroupExternalIDSetEvent {
		return repo.NewGroupExternalIDSetEvent(ctx, &repo.NewAggregate("group1", "org1").Aggregate, scope, old, value)
	}
	for _, tt := range []struct {
		name, value string
		push        bool
	}{
		{"same identity is a no-op", "old", false}, {"replace identity", "new", true}, {"remove identity", "", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			expectations := []expect{expectFilterScoped(t, "org1", eventFromEventPusher(repo.NewGroupAddedEvent(ctx, &repo.NewAggregate("group1", "org1").Aggregate, "group1", "group description")), eventFromEventPusher(binding("entra", "", "old")), eventFromEventPusher(binding("okta", "", "other")))}
			if tt.push {
				expectations = append(expectations, expectGroupPush(binding("entra", "old", tt.value)))
			}
			c := &Commands{eventstore: expectEventstore(expectations...)(t), checkPermission: newMockPermissionCheckAllowed()}
			_, err := c.MutateGroup(ctx, "group1", "org1", []GroupMutation{{ExternalID: &GroupExternalID{Scope: "entra", Value: tt.value}}})
			require.NoError(t, err)
		})
	}
	t.Run("failed membership cannot commit identity", func(t *testing.T) {
		c := &Commands{eventstore: expectEventstore(expectFilterScoped(t, "org1", eventFromEventPusher(repo.NewGroupAddedEvent(ctx, &repo.NewAggregate("group1", "org1").Aggregate, "group1", "group description"))), expectFilterScoped(t, "org1"))(t), checkPermission: newMockPermissionCheckAllowed()}
		_, err := c.MutateGroup(ctx, "group1", "org1", []GroupMutation{{ExternalID: &GroupExternalID{Scope: "entra", Value: "external"}}, {AddMembers: []string{"missing"}}})
		require.Error(t, err)
	})
}

func TestCommands_MutateGroup_validatesOncePerCall(t *testing.T) {
	ctx := context.Background()
	mutations := make([]GroupMutation, maxGroupMutations)
	for i := range mutations {
		mutations[i] = GroupMutation{Name: gu.Ptr("Engineering"), AddMembers: []string{"bob"}, RemoveMembers: []string{"alice"}}
	}
	checks := make(map[string]int)
	c := &Commands{eventstore: expectEventstore(
		expectFilterScoped(t, "org1", eventFromEventPusher(addNewGroupEvent("group1", "org1")), eventFromEventPusher(newGroupUserAddedEvent("group1", "org1", "alice"))),
		// One existence query for all mutations; a second query fails as an unexpected call.
		expectFilterScoped(t, "org1", eventFromEventPusher(addNewUserEvent("bob", "org1"))),
		expectGroupPush(
			repo.NewGroupChangedEvent(ctx, &repo.NewAggregate("group1", "org1").Aggregate, "group1", []repo.GroupChanges{repo.ChangeName(gu.Ptr("Engineering"))}),
			newGroupUserAddedEvent("group1", "org1", "bob"), newGroupUserRemovedEvent("group1", "org1", "alice")),
	)(t), checkPermission: func(_ context.Context, permission, _, _ string) error {
		checks[permission]++
		return nil
	}}
	_, err := c.MutateGroup(ctx, "group1", "org1", mutations)
	require.NoError(t, err)
	require.Len(t, checks, 3)
	for permission, count := range checks {
		require.Equal(t, 1, count, permission)
	}

	c = &Commands{eventstore: expectEventstore()(t), checkPermission: newMockPermissionCheckAllowed()}
	_, err = c.MutateGroup(ctx, "group1", "org1", make([]GroupMutation, maxGroupMutations+1))
	require.True(t, zerrors.IsErrorInvalidArgument(err), err)
}
