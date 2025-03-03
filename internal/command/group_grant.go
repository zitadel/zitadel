package command

import (
	"context"
	"reflect"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/groupgrant"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddGroupGrant(ctx context.Context, groupgrant *domain.GroupGrant, resourceOwner string) (_ *domain.GroupGrant, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	event, addedGroupGrant, err := c.addGroupGrant(ctx, groupgrant, resourceOwner)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(addedGroupGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return groupGrantWriteModelToGroupGrant(addedGroupGrant), nil
}

func (c *Commands) addGroupGrant(ctx context.Context, groupGrant *domain.GroupGrant, resourceOwner string) (command eventstore.Command, _ *GroupGrantWriteModel, err error) {
	if !groupGrant.IsValid() {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-LqfMa", "Errors.GroupGrant.Invalid")
	}
	err = c.checkGroupGrantPreCondition(ctx, groupGrant, resourceOwner)
	if err != nil {
		return nil, nil, err
	}
	groupGrant.AggregateID, err = c.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}

	addedGroupGrant := NewGroupGrantWriteModel(groupGrant.AggregateID, resourceOwner)
	groupGrantAgg := GroupGrantAggregateFromWriteModel(&addedGroupGrant.WriteModel)
	command = groupgrant.NewGroupGrantAddedEvent(
		ctx,
		groupGrantAgg,
		groupGrant.GroupID,
		groupGrant.ProjectID,
		groupGrant.ProjectGrantID,
		groupGrant.RoleKeys,
	)
	return command, addedGroupGrant, nil
}

func (c *Commands) ChangeGroupGrant(ctx context.Context, groupGrant *domain.GroupGrant, resourceOwner string) (_ *domain.GroupGrant, err error) {
	event, changedGroupGrant, err := c.changeGroupGrant(ctx, groupGrant, resourceOwner, false)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(changedGroupGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return groupGrantWriteModelToGroupGrant(changedGroupGrant), nil
}

func (c *Commands) changeGroupGrant(ctx context.Context, groupGrant *domain.GroupGrant, resourceOwner string, cascade bool) (_ eventstore.Command, _ *GroupGrantWriteModel, err error) {
	if groupGrant.AggregateID == "" {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4Nsd", "Errors.GroupGrant.Invalid")
	}
	existingGroupGrant, err := c.groupGrantWriteModelByID(ctx, groupGrant.AggregateID, groupGrant.ResourceOwner)
	if err != nil {
		return nil, nil, err
	}
	err = checkExplicitProjectPermission(ctx, existingGroupGrant.ProjectGrantID, existingGroupGrant.ProjectID)
	if err != nil {
		return nil, nil, err
	}
	if existingGroupGrant.State == domain.GroupGrantStateUnspecified || existingGroupGrant.State == domain.GroupGrantStateRemoved {
		return nil, nil, zerrors.ThrowNotFound(nil, "COMMAND-4AM1l", "Errors.GroupGrant.NotFound")
	}
	if reflect.DeepEqual(existingGroupGrant.RoleKeys, groupGrant.RoleKeys) {
		return nil, nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Ta9fy", "Errors.GroupGrant.NotChanged")
	}
	groupGrant.ProjectID = existingGroupGrant.ProjectID
	groupGrant.ProjectGrantID = existingGroupGrant.ProjectGrantID
	err = c.checkGroupGrantPreCondition(ctx, groupGrant, resourceOwner)
	if err != nil {
		return nil, nil, err
	}

	changedGroupGrant := NewGroupGrantWriteModel(groupGrant.AggregateID, resourceOwner)
	groupGrantAgg := GroupGrantAggregateFromWriteModel(&changedGroupGrant.WriteModel)

	if cascade {
		return groupgrant.NewGroupGrantCascadeChangedEvent(ctx, groupGrantAgg, groupGrant.RoleKeys), existingGroupGrant, nil
	}
	return groupgrant.NewGroupGrantChangedEvent(ctx, groupGrantAgg, groupGrant.RoleKeys), existingGroupGrant, nil
}

func (c *Commands) removeRoleFromGroupGrant(ctx context.Context, groupGrantID string, roleKeys []string, cascade bool) (_ eventstore.Command, err error) {
	existingGroupGrant, err := c.groupGrantWriteModelByID(ctx, groupGrantID, "")
	if err != nil {
		return nil, err
	}
	if existingGroupGrant.State == domain.GroupGrantStateUnspecified || existingGroupGrant.State == domain.GroupGrantStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-4AM1l", "Errors.GroupGrant.NotFound")
	}
	keyExists := false
	for i, key := range existingGroupGrant.RoleKeys {
		for _, roleKey := range roleKeys {
			if key == roleKey {
				keyExists = true
				copy(existingGroupGrant.RoleKeys[i:], existingGroupGrant.RoleKeys[i+1:])
				existingGroupGrant.RoleKeys[len(existingGroupGrant.RoleKeys)-1] = ""
				existingGroupGrant.RoleKeys = existingGroupGrant.RoleKeys[:len(existingGroupGrant.RoleKeys)-1]
				continue
			}
		}
	}
	if !keyExists {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-6m9a1", "Errors.GroupGrant.RoleKeyNotFound")
	}
	changedGroupGrant := NewGroupGrantWriteModel(groupGrantID, existingGroupGrant.ResourceOwner)
	groupGrantAgg := GroupGrantAggregateFromWriteModel(&changedGroupGrant.WriteModel)

	if cascade {
		return groupgrant.NewGroupGrantCascadeChangedEvent(ctx, groupGrantAgg, existingGroupGrant.RoleKeys), nil
	}

	return groupgrant.NewGroupGrantChangedEvent(ctx, groupGrantAgg, existingGroupGrant.RoleKeys), nil
}

func (c *Commands) DeactivateGroupGrant(ctx context.Context, grantID, resourceOwner string) (objectDetails *domain.ObjectDetails, err error) {
	if grantID == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-N1QhH", "Errors.GroupGrant.IDMissing")
	}

	existingGroupGrant, err := c.groupGrantWriteModelByID(ctx, grantID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingGroupGrant.State == domain.GroupGrantStateUnspecified || existingGroupGrant.State == domain.GroupGrantStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-4AM1l", "Errors.GroupGrant.NotFound")
	}
	if existingGroupGrant.State != domain.GroupGrantStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-2qzgx", "Errors.GroupGrant.NotActive")
	}
	err = checkExplicitProjectPermission(ctx, existingGroupGrant.ProjectGrantID, existingGroupGrant.ProjectID)
	if err != nil {
		return nil, err
	}

	deactivateGroupGrant := NewGroupGrantWriteModel(grantID, resourceOwner)
	groupGrantAgg := GroupGrantAggregateFromWriteModel(&deactivateGroupGrant.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, groupgrant.NewGroupGrantDeactivatedEvent(ctx, groupGrantAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingGroupGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGroupGrant.WriteModel), nil
}

func (c *Commands) ReactivateGroupGrant(ctx context.Context, grantID, resourceOwner string) (objectDetails *domain.ObjectDetails, err error) {
	if grantID == "" || resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-1zz8v", "Errors.GroupGrant.IDMissing")
	}

	existingGroupGrant, err := c.groupGrantWriteModelByID(ctx, grantID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingGroupGrant.State == domain.GroupGrantStateUnspecified || existingGroupGrant.State == domain.GroupGrantStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-NQ1qs", "Errors.GroupGrant.NotFound")
	}
	if existingGroupGrant.State != domain.GroupGrantStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-1nQ1v", "Errors.GroupGrant.NotInactive")
	}
	err = checkExplicitProjectPermission(ctx, existingGroupGrant.ProjectGrantID, existingGroupGrant.ProjectID)
	if err != nil {
		return nil, err
	}
	deactivateGroupGrant := NewGroupGrantWriteModel(grantID, resourceOwner)
	groupGrantAgg := GroupGrantAggregateFromWriteModel(&deactivateGroupGrant.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, groupgrant.NewGroupGrantReactivatedEvent(ctx, groupGrantAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingGroupGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGroupGrant.WriteModel), nil
}

func (c *Commands) RemoveGroupGrant(ctx context.Context, grantID, resourceOwner string) (objectDetails *domain.ObjectDetails, err error) {
	event, existingGroupGrant, err := c.removeGroupGrant(ctx, grantID, resourceOwner, false)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingGroupGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGroupGrant.WriteModel), nil
}

func (c *Commands) BulkRemoveGroupGrant(ctx context.Context, grantIDs []string, resourceOwner string) (err error) {
	if len(grantIDs) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-8M0sd", "Errors.GroupGrant.IDMissing")
	}
	events := make([]eventstore.Command, len(grantIDs))
	for i, grantID := range grantIDs {
		event, _, err := c.removeGroupGrant(ctx, grantID, resourceOwner, false)
		if err != nil {
			return err
		}
		events[i] = event
	}
	_, err = c.eventstore.Push(ctx, events...)
	return err
}

func (c *Commands) removeGroupGrant(ctx context.Context, grantID, resourceOwner string, cascade bool) (_ eventstore.Command, writeModel *GroupGrantWriteModel, err error) {
	if grantID == "" {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-I1qd5", "Errors.GroupGrant.IDMissing")
	}

	existingGroupGrant, err := c.groupGrantWriteModelByID(ctx, grantID, resourceOwner)
	if err != nil {
		return nil, nil, err
	}
	if existingGroupGrant.State == domain.GroupGrantStateUnspecified || existingGroupGrant.State == domain.GroupGrantStateRemoved {
		return nil, nil, zerrors.ThrowNotFound(nil, "COMMAND-2Nz0t", "Errors.GroupGrant.NotFound")
	}
	if !cascade {
		err = checkExplicitProjectPermission(ctx, existingGroupGrant.ProjectGrantID, existingGroupGrant.ProjectID)
		if err != nil {
			return nil, nil, err
		}
	}

	removeGroupGrant := NewGroupGrantWriteModel(grantID, existingGroupGrant.ResourceOwner)
	groupGrantAgg := GroupGrantAggregateFromWriteModel(&removeGroupGrant.WriteModel)
	if cascade {
		return groupgrant.NewGroupGrantCascadeRemovedEvent(
			ctx,
			groupGrantAgg,
			existingGroupGrant.GroupID,
			existingGroupGrant.ProjectID,
			existingGroupGrant.ProjectGrantID), existingGroupGrant, nil
	}
	return groupgrant.NewGroupGrantRemovedEvent(
		ctx,
		groupGrantAgg,
		existingGroupGrant.GroupID,
		existingGroupGrant.ProjectID,
		existingGroupGrant.ProjectGrantID), existingGroupGrant, nil
}

func (c *Commands) groupGrantWriteModelByID(ctx context.Context, groupGrantID, resourceOwner string) (writeModel *GroupGrantWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewGroupGrantWriteModel(groupGrantID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) checkGroupGrantPreCondition(ctx context.Context, groupgrant *domain.GroupGrant, resourceOwner string) (err error) {
	/* Why?
	if !authz.GetFeatures(ctx).ShouldUseImprovedPerformance(feature.ImprovedPerformanceTypeGroupGrant) {
		return c.checkGroupGrantPreConditionOld(ctx, groupgrant, resourceOwner)
	}
	*/
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err := c.checkGroupExists(ctx, groupgrant.GroupID, resourceOwner); err != nil {
		return err
	}
	existingRoleKeys, err := c.searchGroupGrantPreConditionState(ctx, groupgrant, resourceOwner)
	if err != nil {
		return err
	}
	if groupgrant.HasInvalidRoles(existingRoleKeys) {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-nn1F4", "Errors.Project.Role.NotFound")
	}
	return nil
}

// this code needs to be rewritten anyways as soon as we improved the fields handling
//
//nolint:gocognit
func (c *Commands) searchGroupGrantPreConditionState(ctx context.Context, groupGrant *domain.GroupGrant, resourceOwner string) (existingRoleKeys []string, err error) {
	criteria := []map[eventstore.FieldType]any{
		// project state query
		{
			eventstore.FieldTypeAggregateType: project.AggregateType,
			eventstore.FieldTypeAggregateID:   groupGrant.ProjectID,
			eventstore.FieldTypeFieldName:     project.ProjectStateSearchField,
			eventstore.FieldTypeObjectType:    project.ProjectSearchType,
		},
		// granted org query
		{
			eventstore.FieldTypeAggregateType: org.AggregateType,
			eventstore.FieldTypeAggregateID:   resourceOwner,
			eventstore.FieldTypeFieldName:     org.OrgStateSearchField,
			eventstore.FieldTypeObjectType:    org.OrgSearchType,
		},
	}
	if groupGrant.ProjectGrantID != "" {
		criteria = append(criteria, map[eventstore.FieldType]any{
			eventstore.FieldTypeAggregateType: project.AggregateType,
			eventstore.FieldTypeAggregateID:   groupGrant.ProjectID,
			eventstore.FieldTypeObjectType:    project.ProjectGrantSearchType,
			eventstore.FieldTypeObjectID:      groupGrant.ProjectGrantID,
		})
	} else {
		criteria = append(criteria, map[eventstore.FieldType]any{
			eventstore.FieldTypeAggregateType: project.AggregateType,
			eventstore.FieldTypeAggregateID:   groupGrant.ProjectID,
			eventstore.FieldTypeObjectType:    project.ProjectRoleSearchType,
			eventstore.FieldTypeFieldName:     project.ProjectRoleKeySearchField,
		})
	}
	results, err := c.eventstore.Search(ctx, criteria...)
	if err != nil {
		return nil, err
	}

	var (
		existsProject    bool
		existsGrantedOrg bool
		existsGrant      bool
	)

	for _, result := range results {
		switch result.Object.Type {
		case project.ProjectRoleSearchType:
			var role string
			err := result.Value.Unmarshal(&role)
			if err != nil {
				return nil, err
			}
			existingRoleKeys = append(existingRoleKeys, role)
		case org.OrgSearchType:
			var state domain.OrgState
			err := result.Value.Unmarshal(&state)
			if err != nil {
				return nil, err
			}
			existsGrantedOrg = state.Valid() && state != domain.OrgStateRemoved
		case project.ProjectSearchType:
			var state domain.ProjectState
			err := result.Value.Unmarshal(&state)
			if err != nil {
				return nil, err
			}
			existsProject = state.Valid() && state != domain.ProjectStateRemoved
		case project.ProjectGrantSearchType:
			switch result.FieldName {
			case project.ProjectGrantGrantedOrgIDSearchField:
				var orgID string
				err := result.Value.Unmarshal(&orgID)
				if err != nil || orgID != resourceOwner {
					return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-PPqla", "Errors.Org.NotFound")
				}
			case project.ProjectGrantStateSearchField:
				var state domain.ProjectGrantState
				err := result.Value.Unmarshal(&state)
				if err != nil {
					return nil, err
				}
				existsGrant = state.Valid() && state != domain.ProjectGrantStateRemoved
			case project.ProjectGrantRoleKeySearchField:
				var role string
				err := result.Value.Unmarshal(&role)
				if err != nil {
					return nil, err
				}
				existingRoleKeys = append(existingRoleKeys, role)
			case project.ProjectGrantGrantIDSearchField:
				var grantID string
				err := result.Value.Unmarshal(&grantID)
				if err != nil || grantID != groupGrant.ProjectGrantID {
					return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-gVulj", "Errors.Project.Grant.NotFound")
				}
			}
		}
	}

	if !existsProject {
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-m9gsd", "Errors.Project.NotFound")
	}
	if !existsGrantedOrg {
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-PPqla", "Errors.Org.NotFound")
	}
	if groupGrant.ProjectGrantID != "" && !existsGrant {
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-gVulj", "Errors.Project.Grant.NotFound")
	}
	return existingRoleKeys, nil
}

func (c *Commands) checkGroupGrantPreConditionOld(ctx context.Context, groupgrant *domain.GroupGrant, resourceOwner string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	preConditions := NewGroupGrantPreConditionReadModel(groupgrant.GroupID, groupgrant.ProjectID, groupgrant.ProjectGrantID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, preConditions)
	if err != nil {
		return err
	}
	if !preConditions.GroupExists {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-31psg", "Errors.Group.NotFound")
	}
	if groupgrant.ProjectGrantID == "" && !preConditions.ProjectExists {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-1q78S", "Errors.Project.NotFound")
	}
	if groupgrant.ProjectGrantID != "" && !preConditions.ProjectGrantExists {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-4l9ff", "Errors.Project.Grant.NotFound")
	}
	if groupgrant.HasInvalidRoles(preConditions.ExistingRoleKeys) {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-nn8F4", "Errors.Project.Role.NotFound")
	}
	return nil
}
