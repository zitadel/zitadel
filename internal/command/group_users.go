package command

import (
	"context"
	"slices"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/tracing"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	repo "github.com/zitadel/zitadel/internal/repository/group"
)

func (c *Commands) AddUsersToGroup(ctx context.Context, groupID string, userIDs []string) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// precondition: check whether the group exists
	group, err := c.checkGroupExists(ctx, groupID, userIDs)
	if err != nil {
		return nil, err
	}

	// check whether the requester has permissions to add users to the group
	err = c.checkPermissionAddUserToGroup(ctx, group.ResourceOwner, group.AggregateID)
	if err != nil {
		return nil, err
	}

	// add the users to the group
	return c.addUsersToGroup(ctx, group)
}

func (c *Commands) RemoveUsersFromGroup(ctx context.Context, groupID string, userIDs []string) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// precondition: check whether the group exists
	group, err := c.checkGroupExists(ctx, groupID, userIDs)
	if err != nil {
		return nil, err
	}

	// check whether the requester has permissions to remove users from the group
	err = c.checkPermissionRemoveUserFromGroup(ctx, group.ResourceOwner, group.AggregateID)
	if err != nil {
		return nil, err
	}

	userIDsToRemove := group.getUserIDsToRemove()
	if len(userIDsToRemove) == 0 {
		// the userIDs are not present in the group; desired state achieved
		return writeModelToObjectDetails(&group.WriteModel), nil
	}

	// remove users from the group
	return c.pushAppendAndReduceDetails(ctx,
		group,
		repo.NewGroupUsersRemovedEvent(
			ctx,
			GroupAggregateFromWriteModel(ctx, &group.WriteModel),
			userIDsToRemove,
		))
}

func (c *Commands) addUsersToGroup(ctx context.Context, group *GroupWriteModel) (*domain.ObjectDetails, error) {
	userIDsToAdd := group.getUserIDsToAdd()
	if len(userIDsToAdd) == 0 {
		// no new users to add
		return writeModelToObjectDetails(&group.WriteModel), nil
	}

	// precondition: check whether the users exist
	for _, userID := range userIDsToAdd {
		// check whether the user exists in the same organization as the group
		_, err := c.checkUserExists(ctx, userID, group.ResourceOwner)
		if err != nil {
			return nil, err
		}
	}

	// add users to the group
	return c.pushAppendAndReduceDetails(ctx,
		group,
		repo.NewGroupUsersAddedEvent(
			ctx,
			GroupAggregateFromWriteModel(ctx, &group.WriteModel),
			userIDsToAdd,
		),
	)
}

// getUserIDsToAdd returns the userIDs that are not already in the group
func (g *GroupWriteModel) getUserIDsToAdd() []string {
	userIDsToAdd := make([]string, 0)
	for _, userID := range g.UserIDs {
		if _, ok := g.existingUserIDs[userID]; !ok && !slices.Contains(userIDsToAdd, userID) {
			userIDsToAdd = append(userIDsToAdd, userID)
		}
	}
	return userIDsToAdd
}

// getUserIDsToRemove returns the userIDs that are in the group and should be removed
// if a userID is not in the group, the desired state has already been achieved
func (g *GroupWriteModel) getUserIDsToRemove() []string {
	userIDsToRemove := make([]string, 0)
	for _, userID := range g.UserIDs {
		if _, ok := g.existingUserIDs[userID]; ok && !slices.Contains(userIDsToRemove, userID) {
			userIDsToRemove = append(userIDsToRemove, userID)
		}
	}
	return userIDsToRemove
}

// removeUserFromGroups returns the events to remove a user from multiple groups.
// This is needed when a user is deleted and subsequently needs to be removed from all groups.
// Note: Ensure that the groupIDs are retrieved via SearchGroupUsers before calling this method
func (c *Commands) removeUserFromGroups(ctx context.Context, userID string, groupIDs []string, resourceOwner string) ([]eventstore.Command, error) {
	events := make([]eventstore.Command, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		events = append(
			events,
			repo.NewGroupUsersRemovedEvent(
				ctx,
				&repo.NewAggregate(groupID, resourceOwner).Aggregate,
				[]string{userID},
			),
		)
	}
	return events, nil
}
