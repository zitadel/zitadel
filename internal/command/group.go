package command

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	repo "github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type CreateGroup struct {
	models.ObjectRoot

	Name        string
	Description string
	UserIDs     []string
	ExternalID  *GroupExternalID
}

type GroupExternalID struct {
	Scope string
	Value string
}

func (g *GroupExternalID) IsValid() error {
	if g != nil && g.Scope == "" {
		return zerrors.ThrowInvalidArgument(nil, "GROUP-ext1", "Provisioning source is required")
	}
	return nil
}

func (g *CreateGroup) IsValid() error {
	if err := g.ExternalID.IsValid(); err != nil {
		return err
	}
	if strings.TrimSpace(g.Name) == "" {
		return zerrors.ThrowInvalidArgument(nil, "GROUP-m177lN", "Errors.Group.InvalidName")
	}
	return nil
}

// CreateGroup creates a new user group in an organization
func (c *Commands) CreateGroup(ctx context.Context, group *CreateGroup) (*domain.ObjectDetails, error) {
	return c.retryGroupChange(ctx, func(ctx context.Context) (*domain.ObjectDetails, error) {
		return c.createGroupOnce(ctx, group)
	})
}

func (c *Commands) createGroupOnce(ctx context.Context, group *CreateGroup) (details *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// create a unique group ID if not provided
	if group.AggregateID == "" {
		group.AggregateID, err = c.idGenerator.Next()
		if err != nil {
			return nil, err
		}
	}

	if err = group.IsValid(); err != nil {
		return nil, err
	}
	if group.ResourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "CMDGRP-msc0Tt", "Errors.Group.MissingOrganizationID")
	}

	if err = c.checkPermissionCreateGroup(ctx, group.ResourceOwner, group.AggregateID); err != nil {
		return nil, err
	}

	// check whether the organization where the group should be created exists
	err = c.checkOrgExists(ctx, group.ResourceOwner)
	if err != nil {
		return nil, zerrors.ThrowPreconditionFailed(nil, "CMDGRP-j1mH8l", "Errors.Org.NotFound")
	}

	// check if a group with the same ID already exists
	groupWriteModel, err := c.getGroupWriteModelByID(ctx, group.AggregateID, group.ResourceOwner, nil)
	if err != nil {
		return nil, err
	}
	if groupWriteModel.State.Exists() {
		return nil, zerrors.ThrowAlreadyExists(nil, "CMDGRP-shRut3", "Errors.Group.AlreadyExists")
	}

	if len(group.UserIDs) > 0 {
		if err := c.checkPermissionAddUserToGroup(ctx, group.ResourceOwner, group.AggregateID); err != nil {
			return nil, err
		}
		if err := c.checkUsersExist(ctx, group.UserIDs, group.ResourceOwner); err != nil {
			return nil, err
		}
	}
	aggregate := GroupAggregateFromWriteModel(ctx, &groupWriteModel.WriteModel)
	events := []eventstore.Command{repo.NewGroupAddedEvent(ctx, aggregate, group.Name, group.Description)}
	if group.ExternalID != nil && group.ExternalID.Value != "" {
		events = append(events, repo.NewGroupExternalIDSetEvent(ctx, aggregate, group.ExternalID.Scope, "", group.ExternalID.Value))
	}
	seen := make(map[string]bool, len(group.UserIDs))
	for _, id := range group.UserIDs {
		if !seen[id] {
			events = append(events, repo.NewGroupUserAddedEvent(ctx, aggregate, id))
			seen[id] = true
		}
	}
	err = c.pushAppendAndReduce(ctx, groupWriteModel, events...)
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
func (c *Commands) UpdateGroup(ctx context.Context, groupUpdate *UpdateGroup) (*domain.ObjectDetails, error) {
	return c.retryGroupChange(ctx, func(ctx context.Context) (*domain.ObjectDetails, error) {
		return c.updateGroupOnce(ctx, groupUpdate)
	})
}

func (c *Commands) updateGroupOnce(ctx context.Context, groupUpdate *UpdateGroup) (details *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err = groupUpdate.IsValid(); err != nil {
		return nil, err
	}

	existingGroup, err := c.getGroupWriteModelByID(ctx, groupUpdate.AggregateID, groupUpdate.ResourceOwner, nil)
	if err != nil {
		return nil, err
	}
	if !existingGroup.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "CMDGRP-b33zly", "Errors.Group.NotFound")
	}

	if err = c.checkPermissionUpdateGroup(ctx, existingGroup.ResourceOwner, existingGroup.AggregateID); err != nil {
		return nil, err
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

// DeleteGroup deletes a user group from an organization.
// Grants of the group, provided as cascadingGrantIDs, are cascade removed
// so their unique constraints are released.
func (c *Commands) DeleteGroup(ctx context.Context, groupID string, cascadingGrantIDs ...string) (*domain.ObjectDetails, error) {
	return c.retryGroupChange(ctx, func(ctx context.Context) (*domain.ObjectDetails, error) {
		return c.deleteGroupOnce(ctx, groupID, cascadingGrantIDs...)
	})
}

func (c *Commands) deleteGroupOnce(ctx context.Context, groupID string, cascadingGrantIDs ...string) (details *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	existingGroup, err := c.getGroupWriteModelByID(ctx, groupID, "", nil)
	if err != nil {
		return nil, err
	}
	if !existingGroup.State.Exists() {
		return writeModelToObjectDetails(&existingGroup.WriteModel), nil
	}

	if err = c.checkPermissionDeleteGroup(ctx, existingGroup.ResourceOwner, existingGroup.AggregateID); err != nil {
		return nil, err
	}

	events := []eventstore.Command{groupRemoval(ctx, existingGroup)}
	if len(cascadingGrantIDs) > 0 {
		cascadeEvents, err := c.removeGroupGrantsFromGroup(ctx, cascadingGrantIDs)
		if err != nil {
			return nil, err
		}
		events = append(events, cascadeEvents...)
	}

	err = c.pushAppendAndReduce(ctx, existingGroup, events...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGroup.WriteModel), nil
}

func (c *Commands) getGroupWriteModelByID(ctx context.Context, groupID, orgID string, userIDs []string) (*GroupWriteModel, error) {
	groupWriteModel := NewGroupWriteModel(groupID, orgID, userIDs)
	err := c.eventstore.FilterToQueryReducer(ctx, groupWriteModel)
	if err != nil {
		return nil, err
	}
	return groupWriteModel, nil
}

func (c *Commands) checkGroupExists(ctx context.Context, groupID string, userIDs []string) (*GroupWriteModel, error) {
	group, err := c.getGroupWriteModelByID(ctx, groupID, "", userIDs)
	if err != nil {
		return nil, err
	}
	if !group.State.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "CMDGRP-eQfeur", "Errors.Group.NotFound")
	}
	return group, nil
}
