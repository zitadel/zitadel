package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"reflect"
)

func (c *Commands) AddProjectGrant(ctx context.Context, grant *domain.ProjectGrant, resourceOwner string) (_ *domain.ProjectGrant, err error) {
	if !grant.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "PROJECT-3b8fs", "Errors.Project.Grant.Invalid")
	}
	err = c.checkProjectGrantPreCondition(ctx, grant)
	if err != nil {
		return nil, err
	}
	grant.GrantID, err = c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	addedGrant := NewProjectGrantWriteModel(grant.GrantID, grant.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedGrant.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(
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
		return nil, caos_errs.ThrowInvalidArgument(nil, "PROJECT-1j83s", "Errors.IDMissing")
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
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-0o0pL", "Errors.NoChangesFoundc")
	}

	events := []eventstore.EventPusher{
		project.NewGrantChangedEvent(ctx, projectAgg, grant.GrantID, grant.RoleKeys),
	}

	removedRoles := domain.GetRemovedRoles(existingGrant.RoleKeys, grant.RoleKeys)
	if len(removedRoles) == 0 {
		pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
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
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingGrant, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return projectGrantWriteModelToProjectGrant(existingGrant), nil
}

func (c *Commands) removeRoleFromProjectGrant(ctx context.Context, projectAgg *eventstore.Aggregate, projectID, projectGrantID, roleKey string, cascade bool) (_ eventstore.EventPusher, _ *ProjectGrantWriteModel, err error) {
	existingProjectGrant, err := c.projectGrantWriteModelByID(ctx, projectGrantID, projectID, "")
	if err != nil {
		return nil, nil, err
	}
	if existingProjectGrant.State == domain.ProjectGrantStateUnspecified || existingProjectGrant.State == domain.ProjectGrantStateRemoved {
		return nil, nil, caos_errs.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.Grant.NotFound")
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
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5m8g9", "Errors.Project.Grant.RoleKeyNotFound")
	}
	changedProjectGrant := NewProjectGrantWriteModel(projectGrantID, projectID, existingProjectGrant.ResourceOwner)

	if cascade {
		return project.NewGrantCascadeChangedEvent(ctx, projectAgg, projectGrantID, existingProjectGrant.RoleKeys), changedProjectGrant, nil
	}

	return project.NewGrantChangedEvent(ctx, projectAgg, projectGrantID, existingProjectGrant.RoleKeys), changedProjectGrant, nil
}

func (c *Commands) DeactivateProjectGrant(ctx context.Context, projectID, grantID, resourceOwner string) (details *domain.ObjectDetails, err error) {
	if grantID == "" || projectID == "" {
		return details, caos_errs.ThrowInvalidArgument(nil, "PROJECT-p0s4V", "Errors.IDMissing")
	}
	err = c.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return details, err
	}
	existingGrant, err := c.projectGrantWriteModelByID(ctx, grantID, projectID, resourceOwner)
	if err != nil {
		return details, err
	}
	if existingGrant.State != domain.ProjectGrantStateActive {
		return details, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-47fu8", "Errors.Project.Grant.NotActive")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingGrant.WriteModel)

	pushedEvents, err := c.eventstore.PushEvents(ctx, project.NewGrantDeactivateEvent(ctx, projectAgg, grantID))
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
		return details, caos_errs.ThrowInvalidArgument(nil, "PROJECT-p0s4V", "Errors.IDMissing")
	}
	err = c.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return details, err
	}
	existingGrant, err := c.projectGrantWriteModelByID(ctx, grantID, projectID, resourceOwner)
	if err != nil {
		return details, err
	}
	if existingGrant.State != domain.ProjectGrantStateInactive {
		return details, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-47fu8", "Errors.Project.Grant.NotInactive")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingGrant.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, project.NewGrantReactivatedEvent(ctx, projectAgg, grantID))
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
		return details, caos_errs.ThrowInvalidArgument(nil, "PROJECT-1m9fJ", "Errors.IDMissing")
	}
	err = c.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return details, caos_errs.ThrowPreconditionFailed(err, "PROJECT-6mf9s", "Errors.Project.NotFound")
	}
	existingGrant, err := c.projectGrantWriteModelByID(ctx, grantID, projectID, resourceOwner)
	if err != nil {
		return details, err
	}
	events := make([]eventstore.EventPusher, 0)
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
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
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
		return nil, caos_errs.ThrowNotFound(nil, "PROJECT-D8JxR", "Errors.Project.Grant.NotFound")
	}

	return writeModel, nil
}

func (c *Commands) checkProjectGrantPreCondition(ctx context.Context, projectGrant *domain.ProjectGrant) error {
	preConditions := NewProjectGrantPreConditionReadModel(projectGrant.AggregateID, projectGrant.GrantedOrgID)
	err := c.eventstore.FilterToQueryReducer(ctx, preConditions)
	if err != nil {
		return err
	}
	if !preConditions.ProjectExists {
		return caos_errs.ThrowPreconditionFailed(err, "COMMAND-m9gsd", "Errors.Project.NotFound")
	}
	if !preConditions.GrantedOrgExists {
		return caos_errs.ThrowPreconditionFailed(err, "COMMAND-3m9gg", "Errors.Org.NotFound")
	}
	if projectGrant.HasInvalidRoles(preConditions.ExistingRoleKeys) {
		return caos_errs.ThrowPreconditionFailed(err, "COMMAND-6m9gd", "Errors.Project.Role.NotFound")
	}
	return nil
}
