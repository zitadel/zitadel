package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	repo "github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddUsersToGroup(ctx context.Context, groupID string, users []repo.GroupUser) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	groupUsers, err := c.getGroupUsersWriteModel(ctx, groupID, "")
	if err != nil {
		return nil, err
	}
	if !groupUsers.State.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "CMDGRP-eQfeur", "Errors.Group.NotFound")
	}

	if err = c.checkPermissionAddUserToGroup(ctx, groupUsers.ResourceOwner, groupUsers.AggregateID); err != nil {
		return nil, err
	}

	toAdd := groupUsers.UsersToAdd(users)
	if len(toAdd) == 0 {
		// all requested users are already in the group; desired state achieved
		return writeModelToObjectDetails(&groupUsers.WriteModel), nil
	}

	// precondition: all users must exist in the same organization as the group
	for _, u := range toAdd {
		if _, err := c.checkUserExists(ctx, u.UserID, groupUsers.ResourceOwner); err != nil {
			return nil, err
		}
	}

	return c.pushAppendAndReduceDetails(ctx,
		groupUsers,
		repo.NewGroupUsersAddedEvent(
			ctx,
			GroupAggregateFromWriteModel(ctx, &groupUsers.WriteModel),
			toAdd,
		),
	)
}

func (c *Commands) RemoveUsersFromGroup(ctx context.Context, groupID string, userIDs []string) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	groupUsers, err := c.getGroupUsersWriteModel(ctx, groupID, "")
	if err != nil {
		return nil, err
	}
	if !groupUsers.State.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "CMDGRP-eQfeur", "Errors.Group.NotFound")
	}

	if err = c.checkPermissionRemoveUserFromGroup(ctx, groupUsers.ResourceOwner, groupUsers.AggregateID); err != nil {
		return nil, err
	}

	toRemove := groupUsers.UserIDsToRemove(userIDs)
	if len(toRemove) == 0 {
		// none of the requested userIDs are in the group; desired state achieved
		return writeModelToObjectDetails(&groupUsers.WriteModel), nil
	}

	return c.pushAppendAndReduceDetails(ctx,
		groupUsers,
		repo.NewGroupUsersRemovedEvent(
			ctx,
			GroupAggregateFromWriteModel(ctx, &groupUsers.WriteModel),
			toRemove,
		),
	)
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
