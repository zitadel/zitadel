package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/v2"

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

	projectRole := domain.RoleProjectOwner
	iam, err := r.GetIAM(ctx)
	if err != nil {
		return nil, nil, err
	}
	if iam.GlobalOrgID == resourceOwner {
		projectRole = domain.RoleProjectOwnerGlobal
	}
	projectAgg.PushEvents(
		project.NewProjectAddedEvent(ctx, projectAdd.Name, resourceOwner),
		project.NewProjectMemberAddedEvent(ctx, ownerUserID, projectRole),
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

func (r *CommandSide) checkProjectExists(ctx context.Context, projectID, resourceOwner string) error {
	projectWriteModel, err := r.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return err
	}
	if projectWriteModel.State == domain.ProjectStateUnspecified || projectWriteModel.State == domain.ProjectStateRemoved {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M0fs", "Errors.Project.NotFound")
	}
	return nil
}

func (r *CommandSide) ChangeProject(ctx context.Context, projectChange *domain.Project, resourceOwner string) (*domain.Project, error) {
	if !projectChange.IsValid() && projectChange.AggregateID != "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4m9vS", "Errors.Project.Invalid")
	}

	existingProject, err := r.getProjectWriteModelByID(ctx, projectChange.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingProject.State == domain.ProjectStateUnspecified || existingProject.State == domain.ProjectStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.NotFound")
	}

	changedEvent, hasChanged, err := existingProject.NewChangedEvent(ctx, projectChange.Name, projectChange.ProjectRoleAssertion, projectChange.ProjectRoleCheck)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M0fs", "Errors.NoChangesFound")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingProject.WriteModel)
	projectAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingProject, projectAgg)
	if err != nil {
		return nil, err
	}

	return projectWriteModelToProject(existingProject), nil
}

func (r *CommandSide) DeactivateProject(ctx context.Context, projectID string, resourceOwner string) error {
	if projectID == "" || resourceOwner == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-88iF0", "Errors.Project.ProjectIDMissing")
	}

	existingProject, err := r.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return err
	}
	if existingProject.State == domain.ProjectStateUnspecified || existingProject.State == domain.ProjectStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-112M9", "Errors.Project.NotFound")
	}
	if existingProject.State != domain.ProjectStateActive {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-mki55", "Errors.Project.NotActive")
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingProject.WriteModel)
	projectAgg.PushEvents(project.NewProjectDeactivatedEvent(ctx))

	return r.eventstore.PushAggregate(ctx, existingProject, projectAgg)
}

func (r *CommandSide) ReactivateProject(ctx context.Context, projectID string, resourceOwner string) error {
	if projectID == "" || resourceOwner == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4m9vS", "Errors.Project.ProjectIDMissing")
	}

	existingProject, err := r.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return err
	}
	if existingProject.State == domain.ProjectStateUnspecified || existingProject.State == domain.ProjectStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.NotFound")
	}
	if existingProject.State != domain.ProjectStateInactive {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5M9bs", "Errors.Project.NotInctive")
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingProject.WriteModel)
	projectAgg.PushEvents(project.NewProjectDeactivatedEvent(ctx))

	return r.eventstore.PushAggregate(ctx, existingProject, projectAgg)
}

func (r *CommandSide) RemoveProject(ctx context.Context, projectID, resourceOwner string, cascadingUserGrantIDs ...string) error {
	if projectID == "" || resourceOwner == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-66hM9", "Errors.Project.ProjectIDMissing")
	}

	existingProject, err := r.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return err
	}
	if existingProject.State == domain.ProjectStateUnspecified || existingProject.State == domain.ProjectStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.NotFound")
	}

	aggregates := make([]eventstore.Aggregater, 0)
	projectAgg := ProjectAggregateFromWriteModel(&existingProject.WriteModel)
	projectAgg.PushEvents(project.NewProjectRemovedEvent(ctx, existingProject.Name, existingProject.ResourceOwner))
	aggregates = append(aggregates, projectAgg)

	for _, grantID := range cascadingUserGrantIDs {
		grantAgg, _, err := r.removeUserGrant(ctx, grantID, "", true)
		if err != nil {
			logging.LogWithFields("COMMAND-b8Djf", "usergrantid", grantID).WithError(err).Warn("could not cascade remove user grant")
			continue
		}
		aggregates = append(aggregates, grantAgg)
	}

	_, err = r.eventstore.PushAggregates(ctx, aggregates...)
	return err
}

func (r *CommandSide) getProjectWriteModelByID(ctx context.Context, projectID, resourceOwner string) (*ProjectWriteModel, error) {
	projectWriteModel := NewProjectWriteModel(projectID, resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, projectWriteModel)
	if err != nil {
		return nil, err
	}
	return projectWriteModel, nil
}
