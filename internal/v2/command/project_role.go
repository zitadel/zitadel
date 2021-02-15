package command

import (
	"context"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
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
	r.addProjectRoles(ctx, projectAgg, projectRole.AggregateID, resourceOwner, projectRole)

	err = r.eventstore.PushAggregate(ctx, roleWriteModel, projectAgg)
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
	r.addProjectRoles(ctx, projectAgg, projectID, resourceOwner, projectRoles...)

	return r.eventstore.PushAggregate(ctx, roleWriteModel, projectAgg)
}

func (r *CommandSide) addProjectRoles(ctx context.Context, projectAgg *project.Aggregate, projectID, resourceOwner string, projectRoles ...*domain.ProjectRole) error {
	for _, projectRole := range projectRoles {
		if !projectRole.IsValid() {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4m9vS", "Errors.Project.Invalid")
		}
		projectAgg.PushEvents(
			project.NewRoleAddedEvent(
				ctx,
				projectRole.Key,
				projectRole.DisplayName,
				projectRole.Group,
				projectID,
			),
		)
	}

	return nil
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

	changeEvent, changed, err := existingRole.NewProjectRoleChangedEvent(ctx, projectRole.Key, projectRole.DisplayName, projectRole.Group)
	if err != nil {
		return nil, err
	}
	if !changed {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5M0cs", "Errors.NoChangesFound")
	}
	projectAgg.PushEvents(changeEvent)

	err = r.eventstore.PushAggregate(ctx, existingRole, projectAgg)
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
	aggregates := make([]eventstore.Aggregater, 0)
	projectAgg := ProjectAggregateFromWriteModel(&existingRole.WriteModel)
	projectAgg.PushEvents(project.NewRoleRemovedEvent(ctx, key, projectID))
	for _, projectGrantID := range cascadingProjectGrantIds {
		_, err = r.removeRoleFromProjectGrant(ctx, projectAgg, projectID, projectGrantID, key, true)
		if err != nil {
			logging.LogWithFields("COMMAND-6n77g", "projectgrantid", projectGrantID).WithError(err).Warn("could not cascade remove role from project grant")
			continue
		}
	}
	aggregates = append(aggregates, projectAgg)

	for _, grantID := range cascadeUserGrantIDs {
		grantAgg, _, err := r.removeRoleFromUserGrant(ctx, grantID, []string{key}, true)
		if err != nil {
			logging.LogWithFields("COMMAND-mK0of", "usergrantid", grantID).WithError(err).Warn("could not cascade remove role on user grant")
			continue
		}
		aggregates = append(aggregates, grantAgg)
	}

	_, err = r.eventstore.PushAggregates(ctx, aggregates...)
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
