package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	repo "github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// CreateGroup creates a new user group in an organization
func (c *Commands) CreateGroup(ctx context.Context, group *domain.Group) (details *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// todo: check permissions

	if !group.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "CMDGRP-dUnd3r", "Errors.Group.InvalidName")
	}
	if group.OrganizationID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "CMDGRP-msc0Tt", "Errors.Group.MissingOrganizationID")
	}
	// check whether the organization where the group should be created exists
	err = c.checkOrgExists(ctx, group.OrganizationID)
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
	groupWriteModel, err := c.getGroupWriteModelByID(ctx, group)
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
			group.OrganizationID,
		))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&groupWriteModel.WriteModel), nil
}

// UpdateGroup updates a user group in an organization
func (c *Commands) UpdateGroup(ctx context.Context, group *domain.Group) (details *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// todo: check permissions
	if !group.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "CMDGRP-m177lN", "Errors.Group.InvalidName")
	}
	existingGroup, err := c.getGroupWriteModelByID(ctx, group)
	if err != nil {
		return nil, err
	}
	if !existingGroup.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "CMDGRP-b33zly", "Errors.Group.NotFound")
	}

	err = c.pushAppendAndReduce(ctx,
		existingGroup,
		repo.NewGroupChangedEvent(ctx,
			GroupAggregateFromWriteModel(ctx, &existingGroup.WriteModel),
			group.Name,
			group.Description,
			existingGroup.ResourceOwner,
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

	if groupID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "CMDGRP-aNg318", "Errors.Group.MissingID")
	}
	existingGroup, err := c.getGroupWriteModelByID(ctx, &domain.Group{ObjectRoot: models.ObjectRoot{AggregateID: groupID, ResourceOwner: authz.GetCtxData(ctx).ResourceOwner}})
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
		))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGroup.WriteModel), nil
}

func (c *Commands) getGroupWriteModelByID(ctx context.Context, group *domain.Group) (*GroupWriteModel, error) {
	groupWriteModel := NewGroupWriteModel(group)
	err := c.eventstore.FilterToQueryReducer(ctx, groupWriteModel)
	if err != nil {
		return nil, err
	}
	return groupWriteModel, nil
}
