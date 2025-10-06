package command

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	repo "github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type CreateGroup struct {
	models.ObjectRoot

	Name        string
	Description string
}

func (g *CreateGroup) IsValid() error {
	if strings.TrimSpace(g.Name) == "" {
		return zerrors.ThrowInvalidArgument(nil, "GROUP-m177lN", "Errors.Group.InvalidName")
	}
	return nil
}

// CreateGroup creates a new user group in an organization
func (c *Commands) CreateGroup(ctx context.Context, group *CreateGroup) (details *domain.ObjectDetails, err error) {
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

type UpdateGroup struct {
	models.ObjectRoot

	Name        *string
	Description *string
}

func (g *UpdateGroup) IsValid() error {
	if g.Name != nil && strings.TrimSpace(*g.Name) == "" {
		return zerrors.ThrowInvalidArgument(nil, "GROUP-dUNd3r", "Errors.Group.InvalidName")
	}
	return nil
}

// UpdateGroup updates a user group in an organization
func (c *Commands) UpdateGroup(ctx context.Context, groupUpdate *UpdateGroup) (details *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err = groupUpdate.IsValid(); err != nil {
		return nil, err
	}

	// todo: check permissions
	existingGroup, err := c.getGroupWriteModelByID(ctx, groupUpdate.AggregateID, groupUpdate.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingGroup.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "CMDGRP-b33zly", "Errors.Group.NotFound")
	}

	changedEvent := existingGroup.NewChangedEvent(
		ctx,
		GroupAggregateFromWriteModel(ctx, &existingGroup.WriteModel),
		groupUpdate.Name,
		groupUpdate.Description,
	)
	if changedEvent == nil {
		return writeModelToObjectDetails(&existingGroup.WriteModel), nil
	}

	err = c.pushAppendAndReduce(ctx,
		existingGroup,
		changedEvent)
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
