package command

import (
	"context"
	"slices"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/groupgrant"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddGroupGrant struct {
	GroupID        string
	ProjectID      string
	ProjectGrantID string
	RoleKeys       []string
}

// AddGroupGrant authorizes all members of a group for a project with the given role keys.
// The project must be owned by or granted to the organization of the group.
func (c *Commands) AddGroupGrant(ctx context.Context, grant *AddGroupGrant) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if grant.GroupID == "" || grant.ProjectID == "" || len(grant.RoleKeys) == 0 {
		return nil, zerrors.ThrowInvalidArgument(nil, "CMDGG-9aFlqx", "Errors.GroupGrant.Invalid")
	}

	group, err := c.checkGroupExists(ctx, grant.GroupID, nil)
	if err != nil {
		return nil, err
	}
	if err = c.checkPermissionWriteGroupGrant(ctx, group.ResourceOwner, group.AggregateID); err != nil {
		return nil, err
	}
	preConditions, err := c.checkGroupGrantPreConditions(ctx, grant.GroupID, grant.ProjectID, grant.ProjectGrantID, group.ResourceOwner, grant.RoleKeys)
	if err != nil {
		return nil, err
	}

	grantID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}

	addedGrant := NewGroupGrantWriteModel(grantID, group.ResourceOwner)
	err = c.pushAppendAndReduce(ctx,
		addedGrant,
		groupgrant.NewGroupGrantAddedEvent(ctx,
			GroupGrantAggregateFromWriteModel(&addedGrant.WriteModel),
			grant.GroupID,
			grant.ProjectID,
			preConditions.FoundGrantID,
			grant.RoleKeys,
		))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&addedGrant.WriteModel), nil
}

// ChangeGroupGrant updates the role keys of a group grant.
// If the roles are unchanged, the request succeeds as the desired state is already achieved.
func (c *Commands) ChangeGroupGrant(ctx context.Context, grantID string, roleKeys []string) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if grantID == "" || len(roleKeys) == 0 {
		return nil, zerrors.ThrowInvalidArgument(nil, "CMDGG-1bXkVe", "Errors.GroupGrant.Invalid")
	}

	existingGrant, err := c.groupGrantWriteModelByID(ctx, grantID)
	if err != nil {
		return nil, err
	}
	if !existingGrant.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "CMDGG-pWqJ8d", "Errors.GroupGrant.NotFound")
	}
	if err = c.checkPermissionWriteGroupGrant(ctx, existingGrant.ResourceOwner, existingGrant.GroupID); err != nil {
		return nil, err
	}

	slices.Sort(roleKeys)
	existingRoleKeys := slices.Clone(existingGrant.RoleKeys)
	slices.Sort(existingRoleKeys)
	if slices.Equal(existingRoleKeys, roleKeys) {
		return writeModelToObjectDetails(&existingGrant.WriteModel), nil
	}

	if _, err = c.checkGroupGrantPreConditions(ctx, existingGrant.GroupID, existingGrant.ProjectID, existingGrant.ProjectGrantID, existingGrant.ResourceOwner, roleKeys); err != nil {
		return nil, err
	}

	err = c.pushAppendAndReduce(ctx,
		existingGrant,
		groupgrant.NewGroupGrantChangedEvent(ctx,
			GroupGrantAggregateFromWriteModel(&existingGrant.WriteModel),
			roleKeys,
		))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGrant.WriteModel), nil
}

// RemoveGroupGrant removes a group grant.
// If the grant is not found, the request succeeds as the desired state is already achieved.
func (c *Commands) RemoveGroupGrant(ctx context.Context, grantID string) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if grantID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "CMDGG-4hYwTb", "Errors.GroupGrant.Invalid")
	}

	existingGrant, err := c.groupGrantWriteModelByID(ctx, grantID)
	if err != nil {
		return nil, err
	}
	if !existingGrant.State.Exists() {
		return writeModelToObjectDetails(&existingGrant.WriteModel), nil
	}
	if err = c.checkPermissionDeleteGroupGrant(ctx, existingGrant.ResourceOwner, existingGrant.GroupID); err != nil {
		return nil, err
	}

	err = c.pushAppendAndReduce(ctx,
		existingGrant,
		groupgrant.NewGroupGrantRemovedEvent(ctx,
			GroupGrantAggregateFromWriteModel(&existingGrant.WriteModel),
			existingGrant.GroupID,
			existingGrant.ProjectID,
			existingGrant.ProjectGrantID,
		))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGrant.WriteModel), nil
}

// removeGroupGrantsFromGroup returns the events to cascade remove all grants of a group,
// releasing their unique constraints. It is used when a group is deleted.
func (c *Commands) removeGroupGrantsFromGroup(ctx context.Context, grantIDs []string) ([]eventstore.Command, error) {
	events := make([]eventstore.Command, 0, len(grantIDs))
	for _, grantID := range grantIDs {
		existingGrant, err := c.groupGrantWriteModelByID(ctx, grantID)
		if err != nil {
			return nil, err
		}
		if !existingGrant.State.Exists() {
			continue
		}
		events = append(events, groupgrant.NewGroupGrantCascadeRemovedEvent(ctx,
			GroupGrantAggregateFromWriteModel(&existingGrant.WriteModel),
			existingGrant.GroupID,
			existingGrant.ProjectID,
			existingGrant.ProjectGrantID,
		))
	}
	return events, nil
}

func (c *Commands) checkGroupGrantPreConditions(ctx context.Context, groupID, projectID, projectGrantID, resourceOwner string, roleKeys []string) (*GroupGrantPreConditionReadModel, error) {
	preConditions := NewGroupGrantPreConditionReadModel(groupID, projectID, projectGrantID, resourceOwner)
	if err := c.eventstore.FilterToQueryReducer(ctx, preConditions); err != nil {
		return nil, err
	}
	if !preConditions.GroupExists {
		return nil, zerrors.ThrowPreconditionFailed(nil, "CMDGG-jJt6Yk", "Errors.Group.NotFound")
	}
	projectIsOwned := resourceOwner == preConditions.ProjectResourceOwner
	if projectIsOwned && !preConditions.ProjectExists {
		return nil, zerrors.ThrowPreconditionFailed(nil, "CMDGG-xF4Vqz", "Errors.Project.NotFound")
	}
	if !projectIsOwned && preConditions.FoundGrantID == "" {
		return nil, zerrors.ThrowPreconditionFailed(nil, "CMDGG-Tr5GHl", "Errors.Project.Grant.NotFound")
	}
	existingRoles := preConditions.existingRoles()
	for _, roleKey := range roleKeys {
		if !slices.Contains(existingRoles, roleKey) {
			return nil, zerrors.ThrowPreconditionFailed(nil, "CMDGG-8nDgWq", "Errors.Project.Role.NotFound")
		}
	}
	return preConditions, nil
}

func (c *Commands) groupGrantWriteModelByID(ctx context.Context, grantID string) (writeModel *GroupGrantWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewGroupGrantWriteModel(grantID, "")
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
