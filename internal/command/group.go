package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	repo "github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// CreateGroup creates a new user group in an organization
func (c *Commands) CreateGroup(ctx context.Context, group *domain.Group) (details *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// todo: check permissions

	if err = group.IsValid(); err != nil {
		return nil, err
	}
	if group.ResourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "CMDGRP-msc0Tt", "Errors.Group.MissingOrganizationID")
	}
	// check whether the organization where the group should be created exists
	err = c.checkOrgExists(ctx, group.ResourceOwner)
	if err != nil {
		return nil, zerrors.ThrowPreconditionFailed(nil, "CMDGRP-j1mH8l", "Errors.Org.NotFound")
	}

	// create a unique group ID if not provided
	if group.AggregateID == "" {
		group.AggregateID, err = c.idGenerator.Next()
		if err != nil {
			return nil, err
		}
	}

	// check if a group with the same ID already exists
	groupWriteModel, err := c.getGroupWriteModelByID(ctx, group.AggregateID, group.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if groupWriteModel.State.Exists() {
		return nil, zerrors.ThrowAlreadyExists(nil, "CMDGRP-shRut3", "Errors.Group.AlreadyExists")
	}

	err = c.pushAppendAndReduce(ctx,
		groupWriteModel,
		repo.NewGroupAddedEvent(ctx,
			GroupAggregateFromWriteModel(ctx, &groupWriteModel.WriteModel),
			group.Name,
			group.Description,
		))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&groupWriteModel.WriteModel), nil
}

// UpdateGroup updates a user group in an organization
func (c *Commands) UpdateGroup(ctx context.Context, groupUpdate *domain.Group) (details *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// todo: check permissions
	existingGroup, err := c.getGroupWriteModelByID(ctx, groupUpdate.AggregateID, groupUpdate.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingGroup.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "CMDGRP-b33zly", "Errors.Group.NotFound")
	}

	// if there are no updates being made, return a successful response as the desired state is achieved.
	if groupUpdate.Name == existingGroup.Name && groupUpdate.Description == existingGroup.Description {
		return writeModelToObjectDetails(&existingGroup.WriteModel), nil
	}

	// validate the group name if it is being updated
	if groupUpdate.Name != "" && existingGroup.Name != groupUpdate.Name {
		if err = groupUpdate.IsValid(); err != nil {
			return nil, err
		}
	}

	err = c.pushAppendAndReduce(ctx,
		existingGroup,
		repo.NewGroupChangedEvent(ctx,
			GroupAggregateFromWriteModel(ctx, &existingGroup.WriteModel),
			existingGroup.Name,
			groupUpdate.Name,
			groupUpdate.Description,
		))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGroup.WriteModel), nil
}

// DeleteGroup deletes a user group from an organization
func (c *Commands) DeleteGroup(ctx context.Context, groupID string) (details *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	existingGroup, err := c.getGroupWriteModelByID(ctx, groupID, "")
	if err != nil {
		return nil, err
	}
	if !existingGroup.State.Exists() {
		return writeModelToObjectDetails(&existingGroup.WriteModel), nil
	}

	err = c.pushAppendAndReduce(ctx,
		existingGroup,
		repo.NewGroupRemovedEvent(ctx,
			GroupAggregateFromWriteModel(ctx, &existingGroup.WriteModel),
			existingGroup.Name,
		))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGroup.WriteModel), nil
}

func (c *Commands) getGroupWriteModelByID(ctx context.Context, id, orgID string) (*GroupWriteModel, error) {
	groupWriteModel := NewGroupWriteModel(id, orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, groupWriteModel)
	if err != nil {
		return nil, err
	}
	return groupWriteModel, nil
}
