package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

func (r *CommandSide) AddProjectRole(ctx context.Context, projectRole *domain.ProjectRole, resourceOwner string) (_ *domain.ProjectRole, err error) {
	err = r.checkProjectExists(ctx, projectRole.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}

	roleWriteModel := NewProjectRoleWriteModel(projectRole.AggregateID, resourceOwner)
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
				resourceOwner,
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

	existingRole, err := r.getProjectRoleWriteModelByID(ctx, projectRole.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingRole.State == domain.ProjectRoleStateUnspecified || existingRole.State == domain.ProjectRoleStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-vv8M9", "Errors.Project.NotFound")
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

func (r *CommandSide) getProjectRoleWriteModelByID(ctx context.Context, projectID, resourceOwner string) (*ProjectRoleWriteModel, error) {
	projectRoleWriteModel := NewProjectRoleWriteModel(projectID, resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, projectRoleWriteModel)
	if err != nil {
		return nil, err
	}
	return projectRoleWriteModel, nil
}
