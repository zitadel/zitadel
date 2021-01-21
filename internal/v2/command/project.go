package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

func (r *CommandSide) AddProject(ctx context.Context, project *domain.Project, resourceOwner, ownerUserID string) (_ *domain.Project, err error) {
	projectAgg, addedProject, err := r.addProject(ctx, project, resourceOwner, ownerUserID)
	if err != nil {
		return nil, err
	}
	err = r.eventstore.PushAggregate(ctx, addedProject, projectAgg)
	if err != nil {
		return nil, err
	}

	return projectWriteModelToProject(addedProject), nil
}

func (r *CommandSide) addProject(ctx context.Context, projectAdd *domain.Project, resourceOwner, ownerUserID string) (_ *project.Aggregate, _ *ProjectWriteModel, err error) {
	if !projectAdd.IsValid() {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-IOVCC", "Errors.Project.Invalid")
	}
	projectAdd.AggregateID, err = r.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}
	addedProject := NewProjectWriteModel(projectAdd.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedProject.WriteModel)

	projectRole := domain.RoleOrgOwner
	//if global { //TODO: !
	//	projectRole = domain.RoleProjectOwnerGlobal
	//}
	projectAgg.PushEvents(
		project.NewProjectAddedEvent(ctx, projectAdd.Name, resourceOwner),
		project.NewMemberAddedEvent(ctx, ownerUserID, projectRole),
	)
	return projectAgg, addedProject, nil
}

func (r *CommandSide) getProjectByID(ctx context.Context, projectID, resourceOwner string) (*domain.Project, error) {
	projectWriteModel, err := r.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if projectWriteModel.State == domain.ProjectStateUnspecified || projectWriteModel.State == domain.ProjectStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "PROJECT-Gd2hh", "Errors.Project.NotFound")
	}
	return projectWriteModelToProject(projectWriteModel), nil
}

func (r *CommandSide) getProjectWriteModelByID(ctx context.Context, projectID, resourceOwner string) (*ProjectWriteModel, error) {
	projectWriteModel := NewProjectWriteModel(projectID, resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, projectWriteModel)
	if err != nil {
		return nil, err
	}
	return projectWriteModel, nil
}
