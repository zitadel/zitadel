package command

import (
	"context"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

func (r *CommandSide) AddProjectRole(ctx context.Context, projectRole *domain.ProjectRole, resourceOwner string) (_ *domain.ProjectRole, err error) {
	err = r.checkProjectExists(ctx, projectRole.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}

	roleWriteModel := NewProjectRoleWriteModelWithKey(projectRole.Key, projectRole.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&roleWriteModel.WriteModel)
	events, err := r.addProjectRoles(ctx, projectAgg, projectRole.AggregateID, projectRole)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := r.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(roleWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return roleWriteModelToRole(roleWriteModel), nil
}

func (r *CommandSide) BulkAddProjectRole(ctx context.Context, projectID, resourceOwner string, projectRoles []*domain.ProjectRole) (err error) {
	err = r.checkProjectExists(ctx, projectID, resourceOwner)
	if err != nil {
		return err
	}

	roleWriteModel := NewProjectRoleWriteModel(projectID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&roleWriteModel.WriteModel)
	events, err := r.addProjectRoles(ctx, projectAgg, projectID, projectRoles...)
	if err != nil {
		return err
	}

	_, err = r.eventstore.PushEvents(ctx, events...)
	return err
}

func (r *CommandSide) addProjectRoles(ctx context.Context, projectAgg *eventstore.Aggregate, projectID string, projectRoles ...*domain.ProjectRole) ([]eventstore.EventPusher, error) {
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

func (r *CommandSide) ChangeProjectRole(ctx context.Context, projectRole *domain.ProjectRole, resourceOwner string) (_ *domain.ProjectRole, err error) {
	if !projectRole.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4m9vS", "Errors.Project.Invalid")
	}
	err = r.checkProjectExists(ctx, projectRole.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}

	existingRole, err := r.getProjectRoleWriteModelByID(ctx, projectRole.Key, projectRole.AggregateID, resourceOwner)
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

	pushedEvents, err := r.eventstore.PushEvents(ctx, changeEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingRole, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return roleWriteModelToRole(existingRole), nil
}

func (r *CommandSide) RemoveProjectRole(ctx context.Context, projectID, key, resourceOwner string, cascadingProjectGrantIds []string, cascadeUserGrantIDs ...string) (err error) {
	if projectID == "" || key == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4m9vS", "Errors.Project.Role.Invalid")
	}
	existingRole, err := r.getProjectRoleWriteModelByID(ctx, key, projectID, resourceOwner)
	if err != nil {
		return err
	}
	if existingRole.State == domain.ProjectRoleStateUnspecified || existingRole.State == domain.ProjectRoleStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-m9vMf", "Errors.Project.Role.NotExisting")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingRole.WriteModel)
	events := []eventstore.EventPusher{
		project.NewRoleRemovedEvent(ctx, projectAgg, key, projectID),
	}

	for _, projectGrantID := range cascadingProjectGrantIds {
		event, _, err := r.removeRoleFromProjectGrant(ctx, projectAgg, projectID, projectGrantID, key, true)
		if err != nil {
			logging.LogWithFields("COMMAND-6n77g", "projectgrantid", projectGrantID).WithError(err).Warn("could not cascade remove role from project grant")
			continue
		}
		events = append(events, event)
	}

	for _, grantID := range cascadeUserGrantIDs {
		event, err := r.removeRoleFromUserGrant(ctx, grantID, []string{key}, true)
		if err != nil {
			logging.LogWithFields("COMMAND-mK0of", "usergrantid", grantID).WithError(err).Warn("could not cascade remove role on user grant")
			continue
		}
		events = append(events, event)
	}

	_, err = r.eventstore.PushEvents(ctx, events...)
	return err
}

func (r *CommandSide) getProjectRoleWriteModelByID(ctx context.Context, key, projectID, resourceOwner string) (*ProjectRoleWriteModel, error) {
	projectRoleWriteModel := NewProjectRoleWriteModelWithKey(key, projectID, resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, projectRoleWriteModel)
	if err != nil {
		return nil, err
	}
	return projectRoleWriteModel, nil
}
