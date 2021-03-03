package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/project"
)

func (c *Commands) AddProjectRole(ctx context.Context, projectRole *domain.ProjectRole, resourceOwner string) (_ *domain.ProjectRole, err error) {
	err = c.checkProjectExists(ctx, projectRole.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}

	roleWriteModel := NewProjectRoleWriteModelWithKey(projectRole.Key, projectRole.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&roleWriteModel.WriteModel)
	events, err := c.addProjectRoles(ctx, projectAgg, projectRole.AggregateID, projectRole)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(roleWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return roleWriteModelToRole(roleWriteModel), nil
}

func (c *Commands) BulkAddProjectRole(ctx context.Context, projectID, resourceOwner string, projectRoles []*domain.ProjectRole) (details *domain.ObjectDetails, err error) {
	err = c.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return details, err
	}

	roleWriteModel := NewProjectRoleWriteModel(projectID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&roleWriteModel.WriteModel)
	events, err := c.addProjectRoles(ctx, projectAgg, projectID, projectRoles...)
	if err != nil {
		return details, err
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(roleWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&roleWriteModel.WriteModel), nil
}

func (c *Commands) addProjectRoles(ctx context.Context, projectAgg *eventstore.Aggregate, projectID string, projectRoles ...*domain.ProjectRole) ([]eventstore.EventPusher, error) {
	var events []eventstore.EventPusher
	for _, projectRole := range projectRoles {
		if !projectRole.IsValid() {
			return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4m9vS", "Errors.Project.Invalid")
		}
		events = append(events, project.NewRoleAddedEvent(
			ctx,
			projectAgg,
			projectRole.Key,
			projectRole.DisplayName,
			projectRole.Group,
			projectID,
		))
	}

	return events, nil
}

func (c *Commands) ChangeProjectRole(ctx context.Context, projectRole *domain.ProjectRole, resourceOwner string) (_ *domain.ProjectRole, err error) {
	if !projectRole.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4m9vS", "Errors.Project.Invalid")
	}
	err = c.checkProjectExists(ctx, projectRole.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}

	existingRole, err := c.getProjectRoleWriteModelByID(ctx, projectRole.Key, projectRole.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingRole.State == domain.ProjectRoleStateUnspecified || existingRole.State == domain.ProjectRoleStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-vv8M9", "Errors.Project.Role.NotExisting")
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingRole.WriteModel)

	changeEvent, changed, err := existingRole.NewProjectRoleChangedEvent(ctx, projectAgg, projectRole.Key, projectRole.DisplayName, projectRole.Group)
	if err != nil {
		return nil, err
	}
	if !changed {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5M0cs", "Errors.NoChangesFound")
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, changeEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingRole, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return roleWriteModelToRole(existingRole), nil
}

func (c *Commands) RemoveProjectRole(ctx context.Context, projectID, key, resourceOwner string, cascadingProjectGrantIds []string, cascadeUserGrantIDs ...string) (details *domain.ObjectDetails, err error) {
	if projectID == "" || key == "" {
		return details, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4m9vS", "Errors.Project.Role.Invalid")
	}
	existingRole, err := c.getProjectRoleWriteModelByID(ctx, key, projectID, resourceOwner)
	if err != nil {
		return details, err
	}
	if existingRole.State == domain.ProjectRoleStateUnspecified || existingRole.State == domain.ProjectRoleStateRemoved {
		return details, caos_errs.ThrowNotFound(nil, "COMMAND-m9vMf", "Errors.Project.Role.NotExisting")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingRole.WriteModel)
	events := []eventstore.EventPusher{
		project.NewRoleRemovedEvent(ctx, projectAgg, key, projectID),
	}

	for _, projectGrantID := range cascadingProjectGrantIds {
		event, _, err := c.removeRoleFromProjectGrant(ctx, projectAgg, projectID, projectGrantID, key, true)
		if err != nil {
			logging.LogWithFields("COMMAND-6n77g", "projectgrantid", projectGrantID).WithError(err).Warn("could not cascade remove role from project grant")
			continue
		}
		events = append(events, event)
	}

	for _, grantID := range cascadeUserGrantIDs {
		event, err := c.removeRoleFromUserGrant(ctx, grantID, []string{key}, true)
		if err != nil {
			logging.LogWithFields("COMMAND-mK0of", "usergrantid", grantID).WithError(err).Warn("could not cascade remove role on user grant")
			continue
		}
		events = append(events, event)
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingRole, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingRole.WriteModel), nil
}

func (c *Commands) getProjectRoleWriteModelByID(ctx context.Context, key, projectID, resourceOwner string) (*ProjectRoleWriteModel, error) {
	projectRoleWriteModel := NewProjectRoleWriteModelWithKey(key, projectID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, projectRoleWriteModel)
	if err != nil {
		return nil, err
	}
	return projectRoleWriteModel, nil
}
