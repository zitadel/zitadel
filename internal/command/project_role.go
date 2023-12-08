package command

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddProjectRole(ctx context.Context, projectRole *domain.ProjectRole, resourceOwner string) (_ *domain.ProjectRole, err error) {
	err = c.checkProjectExists(ctx, projectRole.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}

	roleWriteModel := NewProjectRoleWriteModelWithKey(projectRole.Key, projectRole.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&roleWriteModel.WriteModel)
	events, err := c.addProjectRoles(ctx, projectAgg, projectRole)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
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
	events, err := c.addProjectRoles(ctx, projectAgg, projectRoles...)
	if err != nil {
		return details, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(roleWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&roleWriteModel.WriteModel), nil
}

func (c *Commands) addProjectRoles(ctx context.Context, projectAgg *eventstore.Aggregate, projectRoles ...*domain.ProjectRole) ([]eventstore.Command, error) {
	var events []eventstore.Command
	for _, projectRole := range projectRoles {
		projectRole.AggregateID = projectAgg.ID
		if !projectRole.IsValid() {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4m9vS", "Errors.Project.Role.Invalid")
		}
		events = append(events, project.NewRoleAddedEvent(
			ctx,
			projectAgg,
			projectRole.Key,
			projectRole.DisplayName,
			projectRole.Group,
		))
	}

	return events, nil
}

func (c *Commands) ChangeProjectRole(ctx context.Context, projectRole *domain.ProjectRole, resourceOwner string) (_ *domain.ProjectRole, err error) {
	if !projectRole.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-2ilfW", "Errors.Project.Invalid")
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
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-vv8M9", "Errors.Project.Role.NotExisting")
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingRole.WriteModel)

	changeEvent, changed, err := existingRole.NewProjectRoleChangedEvent(ctx, projectAgg, projectRole.Key, projectRole.DisplayName, projectRole.Group)
	if err != nil {
		return nil, err
	}
	if !changed {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-5M0cs", "Errors.NoChangesFound")
	}

	pushedEvents, err := c.eventstore.Push(ctx, changeEvent)
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
		return details, zerrors.ThrowInvalidArgument(nil, "COMMAND-fl9eF", "Errors.Project.Role.Invalid")
	}
	existingRole, err := c.getProjectRoleWriteModelByID(ctx, key, projectID, resourceOwner)
	if err != nil {
		return details, err
	}
	if existingRole.State == domain.ProjectRoleStateUnspecified || existingRole.State == domain.ProjectRoleStateRemoved {
		return details, zerrors.ThrowNotFound(nil, "COMMAND-m9vMf", "Errors.Project.Role.NotExisting")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingRole.WriteModel)
	events := []eventstore.Command{
		project.NewRoleRemovedEvent(ctx, projectAgg, key),
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

	pushedEvents, err := c.eventstore.Push(ctx, events...)
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
