package command

import (
	"context"
	"reflect"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddProjectGrantWithID(ctx context.Context, grant *domain.ProjectGrant, grantID string, resourceOwner string) (_ *domain.ProjectGrant, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return c.addProjectGrantWithID(ctx, grant, grantID, resourceOwner)
}

func (c *Commands) AddProjectGrant(ctx context.Context, grant *domain.ProjectGrant, resourceOwner string) (_ *domain.ProjectGrant, err error) {
	if !grant.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-3b8fs", "Errors.Project.Grant.Invalid")
	}
	err = c.checkProjectGrantPreCondition(ctx, grant)
	if err != nil {
		return nil, err
	}

	grantID, err := id_generator.Next()
	if err != nil {
		return nil, err
	}

	return c.addProjectGrantWithID(ctx, grant, grantID, resourceOwner)
}

func (c *Commands) addProjectGrantWithID(ctx context.Context, grant *domain.ProjectGrant, grantID string, resourceOwner string) (_ *domain.ProjectGrant, err error) {
	grant.GrantID = grantID

	addedGrant := NewProjectGrantWriteModel(grant.GrantID, grant.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedGrant.WriteModel)
	pushedEvents, err := c.eventstore.Push(
		ctx,
		project.NewGrantAddedEvent(ctx, projectAgg, grant.GrantID, grant.GrantedOrgID, grant.RoleKeys))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return projectGrantWriteModelToProjectGrant(addedGrant), nil
}

func (c *Commands) ChangeProjectGrant(ctx context.Context, grant *domain.ProjectGrant, resourceOwner string, cascadeUserGrantIDs ...string) (_ *domain.ProjectGrant, err error) {
	if grant.GrantID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-1j83s", "Errors.IDMissing")
	}
	existingGrant, err := c.projectGrantWriteModelByID(ctx, grant.GrantID, grant.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	grant.GrantedOrgID = existingGrant.GrantedOrgID
	err = c.checkProjectGrantPreCondition(ctx, grant)
	if err != nil {
		return nil, err
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingGrant.WriteModel)

	if reflect.DeepEqual(existingGrant.RoleKeys, grant.RoleKeys) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "PROJECT-0o0pL", "Errors.NoChangesFoundc")
	}

	events := []eventstore.Command{
		project.NewGrantChangedEvent(ctx, projectAgg, grant.GrantID, grant.RoleKeys),
	}

	removedRoles := domain.GetRemovedRoles(existingGrant.RoleKeys, grant.RoleKeys)
	if len(removedRoles) == 0 {
		pushedEvents, err := c.eventstore.Push(ctx, events...)
		if err != nil {
			return nil, err
		}
		err = AppendAndReduce(existingGrant, pushedEvents...)
		if err != nil {
			return nil, err
		}
		return projectGrantWriteModelToProjectGrant(existingGrant), nil
	}

	for _, userGrantID := range cascadeUserGrantIDs {
		event, err := c.removeRoleFromUserGrant(ctx, userGrantID, removedRoles, true)
		if err != nil {
			continue
		}
		events = append(events, event)
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return projectGrantWriteModelToProjectGrant(existingGrant), nil
}

func (c *Commands) removeRoleFromProjectGrant(ctx context.Context, projectAgg *eventstore.Aggregate, projectID, projectGrantID, roleKey string, cascade bool) (_ eventstore.Command, _ *ProjectGrantWriteModel, err error) {
	existingProjectGrant, err := c.projectGrantWriteModelByID(ctx, projectGrantID, projectID, "")
	if err != nil {
		return nil, nil, err
	}
	if existingProjectGrant.State == domain.ProjectGrantStateUnspecified || existingProjectGrant.State == domain.ProjectGrantStateRemoved {
		return nil, nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.Grant.NotFound")
	}
	keyExists := false
	for i, key := range existingProjectGrant.RoleKeys {
		if key == roleKey {
			keyExists = true
			copy(existingProjectGrant.RoleKeys[i:], existingProjectGrant.RoleKeys[i+1:])
			existingProjectGrant.RoleKeys[len(existingProjectGrant.RoleKeys)-1] = ""
			existingProjectGrant.RoleKeys = existingProjectGrant.RoleKeys[:len(existingProjectGrant.RoleKeys)-1]
			continue
		}
	}
	if !keyExists {
		return nil, nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-5m8g9", "Errors.Project.Grant.RoleKeyNotFound")
	}
	changedProjectGrant := NewProjectGrantWriteModel(projectGrantID, projectID, existingProjectGrant.ResourceOwner)

	if cascade {
		return project.NewGrantCascadeChangedEvent(ctx, projectAgg, projectGrantID, existingProjectGrant.RoleKeys), changedProjectGrant, nil
	}

	return project.NewGrantChangedEvent(ctx, projectAgg, projectGrantID, existingProjectGrant.RoleKeys), changedProjectGrant, nil
}

func (c *Commands) DeactivateProjectGrant(ctx context.Context, projectID, grantID, resourceOwner string) (details *domain.ObjectDetails, err error) {
	if grantID == "" || projectID == "" {
		return details, zerrors.ThrowInvalidArgument(nil, "PROJECT-p0s4V", "Errors.IDMissing")
	}

	err = c.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}

	existingGrant, err := c.projectGrantWriteModelByID(ctx, grantID, projectID, resourceOwner)
	if err != nil {
		return details, err
	}
	if existingGrant.State != domain.ProjectGrantStateActive {
		return details, zerrors.ThrowPreconditionFailed(nil, "PROJECT-47fu8", "Errors.Project.Grant.NotActive")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingGrant.WriteModel)

	pushedEvents, err := c.eventstore.Push(ctx, project.NewGrantDeactivateEvent(ctx, projectAgg, grantID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGrant.WriteModel), nil
}

func (c *Commands) ReactivateProjectGrant(ctx context.Context, projectID, grantID, resourceOwner string) (details *domain.ObjectDetails, err error) {
	if grantID == "" || projectID == "" {
		return details, zerrors.ThrowInvalidArgument(nil, "PROJECT-p0s4V", "Errors.IDMissing")
	}

	err = c.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}

	existingGrant, err := c.projectGrantWriteModelByID(ctx, grantID, projectID, resourceOwner)
	if err != nil {
		return details, err
	}
	if existingGrant.State != domain.ProjectGrantStateInactive {
		return details, zerrors.ThrowPreconditionFailed(nil, "PROJECT-47fu8", "Errors.Project.Grant.NotInactive")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingGrant.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, project.NewGrantReactivatedEvent(ctx, projectAgg, grantID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGrant.WriteModel), nil
}

func (c *Commands) RemoveProjectGrant(ctx context.Context, projectID, grantID, resourceOwner string, cascadeUserGrantIDs ...string) (details *domain.ObjectDetails, err error) {
	if grantID == "" || projectID == "" {
		return details, zerrors.ThrowInvalidArgument(nil, "PROJECT-1m9fJ", "Errors.IDMissing")
	}

	err = c.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}

	existingGrant, err := c.projectGrantWriteModelByID(ctx, grantID, projectID, resourceOwner)
	if err != nil {
		return details, err
	}
	events := make([]eventstore.Command, 0)
	projectAgg := ProjectAggregateFromWriteModel(&existingGrant.WriteModel)
	events = append(events, project.NewGrantRemovedEvent(ctx, projectAgg, grantID, existingGrant.GrantedOrgID))

	for _, userGrantID := range cascadeUserGrantIDs {
		event, _, err := c.removeUserGrant(ctx, userGrantID, "", true)
		if err != nil {
			logging.LogWithFields("COMMAND-3m8sG", "usergrantid", grantID).WithError(err).Warn("could not cascade remove user grant")
			continue
		}
		events = append(events, event)
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingGrant.WriteModel), nil
}

func (c *Commands) projectGrantWriteModelByID(ctx context.Context, grantID, projectID, resourceOwner string) (member *ProjectGrantWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewProjectGrantWriteModel(grantID, projectID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.ProjectGrantStateUnspecified || writeModel.State == domain.ProjectGrantStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "PROJECT-D8JxR", "Errors.Project.Grant.NotFound")
	}

	return writeModel, nil
}

func (c *Commands) checkProjectGrantPreCondition(ctx context.Context, projectGrant *domain.ProjectGrant) error {
	if !authz.GetFeatures(ctx).ShouldUseImprovedPerformance(feature.ImprovedPerformanceTypeProjectGrant) {
		return c.checkProjectGrantPreConditionOld(ctx, projectGrant)
	}
	existingRoleKeys, err := c.searchProjectGrantState(ctx, projectGrant.AggregateID, projectGrant.GrantedOrgID)
	if err != nil {
		return err
	}

	if projectGrant.HasInvalidRoles(existingRoleKeys) {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-6m9gd", "Errors.Project.Role.NotFound")
	}
	return nil
}

func (c *Commands) searchProjectGrantState(ctx context.Context, projectID, grantedOrgID string) (existingRoleKeys []string, err error) {
	results, err := c.eventstore.Search(
		ctx,
		// project state query
		map[eventstore.FieldType]any{
			eventstore.FieldTypeAggregateType: project.AggregateType,
			eventstore.FieldTypeAggregateID:   projectID,
			eventstore.FieldTypeFieldName:     project.ProjectStateSearchField,
			eventstore.FieldTypeObjectType:    project.ProjectSearchType,
		},
		// granted org query
		map[eventstore.FieldType]any{
			eventstore.FieldTypeAggregateType: org.AggregateType,
			eventstore.FieldTypeAggregateID:   grantedOrgID,
			eventstore.FieldTypeFieldName:     org.OrgStateSearchField,
			eventstore.FieldTypeObjectType:    org.OrgSearchType,
		},
		// role query
		map[eventstore.FieldType]any{
			eventstore.FieldTypeAggregateType: project.AggregateType,
			eventstore.FieldTypeAggregateID:   projectID,
			eventstore.FieldTypeFieldName:     project.ProjectRoleKeySearchField,
			eventstore.FieldTypeObjectType:    project.ProjectRoleSearchType,
		},
	)
	if err != nil {
		return nil, err
	}

	var (
		existsProject    bool
		existsGrantedOrg bool
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
		}
	}

	if !existsProject {
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-m9gsd", "Errors.Project.NotFound")
	}
	if !existsGrantedOrg {
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-3m9gg", "Errors.Org.NotFound")
	}
	return existingRoleKeys, nil
}
