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
		{name: "rename and membership events share one push", mutations: []GroupMutation{{Name: gu.Ptr("Engineering"), AddMembers: []string{"bob"}, RemoveMembers: []string{"alice"}}}, expects: []expect{initial(), expectFilterScoped(t, "org1", eventFromEventPusher(addNewUserEvent("bob", "org1"))), expectPush(renamed, newGroupUserAddedEvent("group1", "org1", "bob"), newGroupUserRemovedEvent("group1", "org1", "alice"))}},
		{name: "failed member validation never pushes rename", mutations: []GroupMutation{{Name: gu.Ptr("Engineering")}, {AddMembers: []string{"missing"}}}, expects: []expect{initial(), expectFilterScoped(t, "org1")}, wantErr: true},
		{name: "duplicate add is idempotent", mutations: []GroupMutation{{AddMembers: []string{"alice", "alice"}}}, expects: []expect{initial(), expectFilterScoped(t, "org1", eventFromEventPusher(addNewUserEvent("alice", "org1")))}},
		{name: "missing removal is idempotent", mutations: []GroupMutation{{RemoveMembers: []string{"missing"}}}, expects: []expect{initial()}},
		{name: "add then remove has no net event", mutations: []GroupMutation{{AddMembers: []string{"bob"}}, {RemoveMembers: []string{"bob"}}}, expects: []expect{initial(), expectFilterScoped(t, "org1", eventFromEventPusher(addNewUserEvent("bob", "org1")))}},
		{name: "replace removes omitted members", mutations: []GroupMutation{{Members: gu.Ptr([]string{"bob"})}}, expects: []expect{initial(), expectFilterScoped(t, "org1", eventFromEventPusher(addNewUserEvent("bob", "org1"))), expectPush(newGroupUserAddedEvent("group1", "org1", "bob"), newGroupUserRemovedEvent("group1", "org1", "alice"))}},
		{name: "empty replacement clears members", mutations: []GroupMutation{{Members: gu.Ptr([]string{})}}, expects: []expect{initial(), expectPush(newGroupUserRemovedEvent("group1", "org1", "alice"))}},
		{name: "empty removal preserves members", mutations: []GroupMutation{{RemoveMembers: []string{}}}, expects: []expect{initial()}},
		{name: "organization boundary", mutations: []GroupMutation{{Name: gu.Ptr("Engineering")}}, expects: []expect{expectFilterScoped(t, "org1")}, wantErr: true},
		{name: "permission denied", mutations: []GroupMutation{{RemoveMembers: []string{"alice"}}}, expects: []expect{initial()}, denied: true, wantErr: true},
		{name: "push failure returned", mutations: []GroupMutation{{RemoveMembers: []string{"alice"}}}, expects: []expect{initial(), expectPushFailed(errors.New("push failed"), newGroupUserRemovedEvent("group1", "org1", "alice"))}, wantErr: true},
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
