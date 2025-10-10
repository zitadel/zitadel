package command

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
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

	// check whether the group exists
	group, err := c.getGroupWriteModelByID(ctx, groupID, "")
	if err != nil {
		return nil, err
	}
	if !group.State.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "CMDGRP-eQfeur", "Errors.Group.NotFound")
	}

	// check whether the requester has permissions to add users to the group
	err = c.checkPermissionCreateGroup(ctx, group.ResourceOwner, group.AggregateID)
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

func (c *Commands) addUsersToGroup(ctx context.Context, resourceOwner, groupID string, userIDs []string) (*domain.ObjectDetails, []string, error) {
	var failedUserIDs []string
	var groupUserWriteModel *GroupUserWriteModel
	var events []eventstore.Command

	for _, userID := range userIDs {
		// check whether the user exists in the same organization as the group
		userResourceOwner, err := c.checkUserExists(ctx, userID, "")
		if err != nil || userResourceOwner != resourceOwner {
			logging.WithFields(
				"user_id", userID,
				"group_id", groupID,
				"user_resource_owner", userResourceOwner,
				"group_resource_owner", resourceOwner,
			).WithError(err).Error("failed to add user to group")
			failedUserIDs = append(failedUserIDs, userID)
			continue
		}

		// check whether the user is already a member of the group
		wm, err := c.getGroupUserWriteModel(ctx, resourceOwner, groupID, userID)
		if err != nil {
			// failed to get the writemodel
			logging.WithFields("user_id", userID, "group_id", groupID).WithError(err).Error("failed to add user to group")
			failedUserIDs = append(failedUserIDs, userID)
			continue
		}
		if wm.State.Exists() {
			// the user is already a member of the group
			continue
		}

		// add the user to the group
		events = append(
			events,
			repo.NewGroupUserAddedEvent(
				ctx,
				GroupAggregateFromWriteModel(ctx, &wm.WriteModel),
				userID,
			),
		)
		groupUserWriteModel = wm
	}

	// no users were added
	if len(events) == 0 {
		return nil, failedUserIDs, nil
	}

	details, err := c.pushAppendAndReduceDetails(ctx,
		groupUserWriteModel,
		events...)
	if err != nil {
		return nil, failedUserIDs, err
	}
	return details, failedUserIDs, nil
}

func (c *Commands) getGroupUserWriteModel(ctx context.Context, resourceOwner, groupID, userID string) (*GroupUserWriteModel, error) {
	groupUserWriteModel := NewGroupUserWriteModel(resourceOwner, groupID, userID)
	err := c.eventstore.FilterToQueryReducer(ctx, groupUserWriteModel)
	if err != nil {
		return nil, err
	}
	return groupUserWriteModel, nil
}
