package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	repo "github.com/zitadel/zitadel/internal/repository/group"
)

func expectGroupMembershipsFixture(userID, orgID string, ids ...string) expect {
	var events []eventstore.Event
	for _, id := range ids {
		aggregate := &repo.NewAggregate(id, orgID).Aggregate
		events = append(events, eventFromEventPusher(repo.NewGroupAddedEvent(context.Background(), aggregate, id, "")),
			eventFromEventPusher(repo.NewGroupUserAddedEvent(context.Background(), aggregate, userID)))
	}
	return expectFilter(events...)
}

func TestRemoveUserGroupMemberships_authoritativeState(t *testing.T) {
	ctx := context.Background()
	aggregate := &repo.NewAggregate("g", "org").Aggregate
	c := &Commands{eventstore: expectEventstore(expectFilter(
		eventFromEventPusher(repo.NewGroupAddedEvent(ctx, aggregate, "g", "")),
		eventFromEventPusher(repo.NewGroupUserAddedEvent(ctx, aggregate, "removed")),
		eventFromEventPusher(repo.NewGroupUserRemovedEvent(ctx, aggregate, "removed")),
		eventFromEventPusher(repo.NewGroupUserAddedEvent(ctx, aggregate, "current")),
	))(t)}
	events, err := c.removeUserGroupMemberships(ctx, "current", "org")
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "current", events[0].(*repo.GroupUserRemovedEvent).UserID)
	require.Equal(t, "org", events[0].Aggregate().ResourceOwner)
	// No projected list is supplied; the committed addition alone determines the cascade.
}
