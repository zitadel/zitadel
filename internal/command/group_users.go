package command

import (
	"context"
	"slices"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	repo "github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddUsersToGroupResponse struct {
	*domain.ObjectDetails
	FailedUserIDs []string
}

func (c *Commands) AddUsersToGroup(ctx context.Context, groupID string, userIDs []string) (_ *AddUsersToGroupResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// precondition: check whether the group exists
	group, err := c.getGroupWriteModelByID(ctx, groupID, "")
	if err != nil {
		return nil, err
	}
	if !group.State.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "CMDGRP-eQfeur", "Errors.Group.NotFound")
	}

	// check whether the requester has permissions to add users to the group
	err = c.checkPermissionAddUserToGroup(ctx, group.ResourceOwner, group.AggregateID)
	if err != nil {
		return nil, err
	}

	// add the users to the group
	details, failedUserIDs, err := c.addUsersToGroup(ctx, group.ResourceOwner, group.AggregateID, userIDs)
	if err != nil {
		return nil, err
	}

	return &AddUsersToGroupResponse{
		FailedUserIDs: failedUserIDs,
		ObjectDetails: details,
	}, nil
}

func (c *Commands) RemoveUsersFromGroup(ctx context.Context, groupID string, userIDs []string) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// precondition: check whether the group exists
	group, err := c.getGroupWriteModelByID(ctx, groupID, "")
	if err != nil {
		return nil, err
	}
	if !group.State.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "CMDGRP-JRBnLw", "Errors.Group.NotFound")
	}

	// check whether the requester has permissions to remove users from the group
	err = c.checkPermissionRemoveUserFromGroup(ctx, group.ResourceOwner, group.AggregateID)
	if err != nil {
		return nil, err
	}

	// remove duplicate userIDs
	uniqueUserIDs := removeDuplicateUserIDs(userIDs)

	groupUsersWriteModel, err := c.getGroupUsersWriteModel(ctx, group.ResourceOwner, groupID, uniqueUserIDs)
	if err != nil {
		return nil, err
	}
	userIDsToRemove := groupUsersWriteModel.userIDsToRemove()
	if len(userIDsToRemove) == 0 {
		// the userIDs are already removed from the group; desired state achieved
		return writeModelToObjectDetails(&groupUsersWriteModel.WriteModel), nil
	}

	// remove users from the group
	err = c.pushAppendAndReduce(ctx,
		groupUsersWriteModel,
		repo.NewGroupUsersRemovedEvent(
			ctx,
			GroupAggregateFromWriteModel(ctx, &groupUsersWriteModel.WriteModel),
			userIDsToRemove,
		),
	)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&groupUsersWriteModel.WriteModel), nil
}

func (c *Commands) addUsersToGroup(ctx context.Context, resourceOwner, groupID string, userIDs []string) (*domain.ObjectDetails, []string, error) {
	var failedUserIDs, userIDsToAdd []string

	// remove duplicate userIDs
	uniqueUserIDs := removeDuplicateUserIDs(userIDs)

	for _, userID := range uniqueUserIDs {
		// check whether the user exists in the same organization as the group
		userResourceOwner, err := c.checkUserExists(ctx, userID, "")
		if err != nil || userResourceOwner != resourceOwner {
			logging.WithFields(
				"user_id", userID,
				"group_id", groupID,
				"user_resource_owner", userResourceOwner,
				"group_resource_owner", resourceOwner,
			).WithError(err).Error("user does not exist or is not in the same organization as the group")
			failedUserIDs = append(failedUserIDs, userID)
			continue
		}
		userIDsToAdd = append(userIDsToAdd, userID)
	}

	if len(userIDsToAdd) == 0 {
		// todo: or send an error?
		return &domain.ObjectDetails{EventDate: time.Now().UTC()}, failedUserIDs, nil
	}

	groupUsersWriteModel, err := c.getGroupUsersWriteModel(ctx, resourceOwner, groupID, userIDsToAdd)
	if err != nil {
		return nil, failedUserIDs, err
	}
	// filter out users who already exist in the group
	userIDsToAdd = groupUsersWriteModel.userIDsToAdd()
	if len(userIDsToAdd) == 0 {
		// all users already exist in the group; desired state achieved
		return writeModelToObjectDetails(&groupUsersWriteModel.WriteModel), failedUserIDs, nil
	}

	// add users to the group
	err = c.pushAppendAndReduce(ctx,
		groupUsersWriteModel,
		repo.NewGroupUsersAddedEvent(
			ctx,
			GroupAggregateFromWriteModel(ctx, &groupUsersWriteModel.WriteModel),
			userIDsToAdd,
		),
	)
	if err != nil {
		return nil, failedUserIDs, err
	}

	return writeModelToObjectDetails(&groupUsersWriteModel.WriteModel), failedUserIDs, nil
}

func (c *Commands) getGroupUsersWriteModel(ctx context.Context, resourceOwner, groupID string, userIDs []string) (*GroupUsersWriteModel, error) {
	groupUserWriteModel := NewGroupUsersWriteModel(resourceOwner, groupID, userIDs)
	err := c.eventstore.FilterToQueryReducer(ctx, groupUserWriteModel)
	if err != nil {
		return nil, err
	}
	return groupUserWriteModel, nil
}

func removeDuplicateUserIDs(userIDs []string) []string {
	uniqueUserIDs := append([]string(nil), userIDs...)
	slices.Sort(uniqueUserIDs)
	return slices.Compact(uniqueUserIDs)
}
