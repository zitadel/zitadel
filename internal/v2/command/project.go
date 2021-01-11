package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

func (r *CommandSide) AddProject(ctx context.Context, project *domain.Project, ownerUserID string) (_ *domain.Project, err error) {
	projectAgg, addedProject, err := r.addProject(ctx, project, ownerUserID)
	if err != nil {
		return nil, err
	}
	err = r.eventstore.PushAggregate(ctx, addedProject, projectAgg)
	if err != nil {
		return nil, err
	}

	return projectWriteModelToProject(addedProject), nil
}

func (r *CommandSide) addProject(ctx context.Context, projectAdd *domain.Project, ownerUserID string) (_ *project.Aggregate, _ *ProjectWriteModel, err error) {
	if !projectAdd.IsValid() {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-IOVCC", "Errors.Project.Invalid")
	}
	projectAdd.AggregateID, err = r.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}
	//project.State = proj_model.ProjectStateActive

	addedProject := NewProjectWriteModel(projectAdd.AggregateID)
	projectAgg := ProjectAggregateFromWriteModel(&addedProject.WriteModel)

	projectRole := domain.RoleOrgOwner
	//if global { //TODO: !
	//	projectRole = domain.RoleProjectOwnerGlobal
	//}
	projectAgg.PushEvents(
		project.NewProjectAddedEvent(ctx, projectAdd.Name),
		project.NewMemberAddedEvent(ctx, ownerUserID, projectRole),
	)
	return projectAgg, addedProject, nil
}

func (r *CommandSide) getProjectByID(ctx context.Context, projectID string) (*domain.Project, error) {
	projectWriteModel, err := r.getProjectWriteModelByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return projectWriteModelToProject(projectWriteModel), nil
}

func (r *CommandSide) getProjectWriteModelByID(ctx context.Context, projectID string) (*ProjectWriteModel, error) {
	projectWriteModel := NewProjectWriteModel(projectID)
	err := r.eventstore.FilterToQueryReducer(ctx, projectWriteModel)
	if err != nil {
		return nil, err
	}
	return projectWriteModel, nil
}
