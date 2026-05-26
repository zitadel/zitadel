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

	// precondition: all users must exist within the same instance
	for _, u := range toAdd {
		if _, err := c.checkUserExists(ctx, u.UserID, ""); err != nil {
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

// GroupUserRef identifies a group membership for cascading operations (e.g. user deletion).
type GroupUserRef struct {
	GroupID       string
	ResourceOwner string
}

// removeUserFromGroups returns the events to remove a user from multiple groups.
// This is needed when a user is deleted and subsequently needs to be removed from all groups.
// Note: Ensure that the groupRefs are retrieved via SearchGroupUsers before calling this method
func (c *Commands) removeUserFromGroups(ctx context.Context, userID string, groupRefs []GroupUserRef) ([]eventstore.Command, error) {
	events := make([]eventstore.Command, 0, len(groupRefs))
	for _, ref := range groupRefs {
		events = append(
			events,
			repo.NewGroupUsersRemovedEvent(
				ctx,
				&repo.NewAggregate(ref.GroupID, ref.ResourceOwner).Aggregate,
				[]string{userID},
			),
		)
	}
	return events, nil
}
