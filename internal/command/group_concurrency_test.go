package command

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	repo "github.com/zitadel/zitadel/internal/repository/group"
)

// The existing fixtures replay events at sequence zero. Explicitly expect that
// revision on group writes, while leaving other aggregates in the batch alone.
func groupCommandsAtFixtureSequence(commands []eventstore.Command) {
	for _, command := range commands {
		if aggregate := command.Aggregate(); aggregate.Type == repo.AggregateType && aggregate.ExpectedSequence == nil {
			aggregate.ExpectedSequence = gu.Ptr(uint64(0))
		}
	}
}

func expectGroupPush(commands ...eventstore.Command) expect {
	groupCommandsAtFixtureSequence(commands)
	return expectPush(commands...)
}

func expectGroupPushFailed(err error, commands ...eventstore.Command) expect {
	groupCommandsAtFixtureSequence(commands)
	return expectPushFailed(err, commands...)
}

func TestCommands_MutateGroup_recomputesAfterConflict(t *testing.T) {
	ctx := context.Background()
	added := eventFromEventPusher(addNewGroupEvent("group1", "org1"))
	alice := eventFromEventPusher(newGroupUserAddedEvent("group1", "org1", "alice"))
	carol := eventFromEventPusher(newGroupUserAddedEvent("group1", "org1", "carol"))
	bob := eventFromEventPusher(addNewUserEvent("bob", "org1"))
	c := &Commands{eventstore: expectEventstore(
		expectFilterScoped(t, "org1", added, alice),
		expectFilterScoped(t, "org1", bob),
		expectGroupPushFailed(eventstore.ErrConcurrentModification,
			newGroupUserAddedEvent("group1", "org1", "bob"), newGroupUserRemovedEvent("group1", "org1", "alice")),
		// Another writer added Carol before our first push. The replacement
		// must now remove Carol too, rather than replay the original events.
		expectFilterScoped(t, "org1", added, alice, carol),
		expectFilterScoped(t, "org1", bob),
		expectGroupPush(newGroupUserAddedEvent("group1", "org1", "bob"),
			newGroupUserRemovedEvent("group1", "org1", "alice"), newGroupUserRemovedEvent("group1", "org1", "carol")),
	)(t), checkPermission: newMockPermissionCheckAllowed()}
	_, err := c.MutateGroup(ctx, "group1", "org1", []GroupMutation{{Members: gu.Ptr([]string{"bob"})}})
	require.NoError(t, err)
}

func TestRetryGroupChange_bounded(t *testing.T) {
	calls := 0
	c := &Commands{eventstore: expectEventstore()(t)}
	_, err := c.retryGroupChange(context.Background(), func(context.Context) (*domain.ObjectDetails, error) {
		calls++
		return nil, eventstore.ErrConcurrentModification
	})
	require.ErrorIs(t, err, eventstore.ErrConcurrentModification)
	require.Equal(t, 3, calls)
}

func TestCommands_MutateGroup_reloadsReservationAfterConflict(t *testing.T) {
	ctx := context.Background()
	added := eventFromEventPusher(addNewGroupEvent("group1", "org1"))
	binding := func(old, value string) *repo.GroupExternalIDSetEvent {
		return repo.NewGroupExternalIDSetEvent(ctx, &repo.NewAggregate("group1", "org1").Aggregate, "entra", old, value)
	}
	c := &Commands{eventstore: expectEventstore(
		expectFilterScoped(t, "org1", added, eventFromEventPusher(binding("", "a"))),
		expectGroupPushFailed(eventstore.ErrConcurrentModification, binding("a", "c")),
		expectFilterScoped(t, "org1", added, eventFromEventPusher(binding("", "b"))),
		expectGroupPush(binding("b", "c")),
	)(t), checkPermission: newMockPermissionCheckAllowed()}
	_, err := c.MutateGroup(ctx, "group1", "org1", []GroupMutation{{ExternalID: &GroupExternalID{Scope: "entra", Value: "c"}}})
	require.NoError(t, err)
}

func TestGroupWriteModel_revisionIncludesManagerRoles(t *testing.T) {
	wm := NewGroupWriteModel("group1", "org1", nil)
	added := eventFromEventPusher(addNewGroupEvent("group1", "org1"))
	added.Seq = 1
	roles := eventFromEventPusher(repo.NewGroupManagerRolesSetEvent(context.Background(), &repo.NewAggregate("group1", "org1").Aggregate, nil))
	roles.Seq = 4
	wm.AppendEvents(added, roles)
	require.NoError(t, wm.Reduce())
	require.Equal(t, uint64(4), *GroupAggregateFromWriteModel(context.Background(), &wm.WriteModel).ExpectedSequence)
}
