package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) checkProjectExistsOld(ctx context.Context, projectID, resourceOwner string) (_ string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	projectWriteModel, err := c.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return "", err
	}
	if !isProjectStateExists(projectWriteModel.State) {
		return "", zerrors.ThrowPreconditionFailed(nil, "COMMAND-EbFMN", "Errors.Project.NotFound")
	}
	return projectWriteModel.ResourceOwner, nil
}

func (c *Commands) deactivateProjectOld(ctx context.Context, projectID string, resourceOwner string) (*domain.ObjectDetails, error) {
	existingProject, err := c.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isProjectStateExists(existingProject.State) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-112M9", "Errors.Project.NotFound")
	}
	if existingProject.State != domain.ProjectStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-mki55", "Errors.Project.NotActive")
	}
	if err := c.checkPermissionUpdateProject(ctx, existingProject.ResourceOwner, existingProject.AggregateID); err != nil {
		return nil, err
	}

	//nolint: contextcheck
	projectAgg := ProjectAggregateFromWriteModel(&existingProject.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, project.NewProjectDeactivatedEvent(ctx, projectAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingProject, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingProject.WriteModel), nil
}

func (c *Commands) reactivateProjectOld(ctx context.Context, projectID string, resourceOwner string) (*domain.ObjectDetails, error) {
	existingProject, err := c.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isProjectStateExists(existingProject.State) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.NotFound")
	}
	if existingProject.State != domain.ProjectStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-5M9bs", "Errors.Project.NotInactive")
	}
	if err := c.checkPermissionUpdateProject(ctx, existingProject.ResourceOwner, existingProject.AggregateID); err != nil {
		return nil, err
	}

	//nolint: contextcheck
	projectAgg := ProjectAggregateFromWriteModel(&existingProject.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, project.NewProjectReactivatedEvent(ctx, projectAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingProject, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingProject.WriteModel), nil
}

func (c *Commands) checkProjectGrantPreConditionOld(ctx context.Context, projectID, grantedOrgID, resourceOwner string, roles []string) (string, error) {
	preConditions := NewProjectGrantPreConditionReadModel(projectID, grantedOrgID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, preConditions)
	if err != nil {
		return "", err
	}
	if !preConditions.ProjectExists {
		return "", zerrors.ThrowPreconditionFailed(err, "COMMAND-m9gsd", "Errors.Project.NotFound")
	}
	if !preConditions.GrantedOrgExists {
		return "", zerrors.ThrowPreconditionFailed(err, "COMMAND-3m9gg", "Errors.Org.NotFound")
	}
	if domain.HasInvalidRoles(preConditions.ExistingRoleKeys, roles) {
		return "", zerrors.ThrowPreconditionFailed(err, "COMMAND-6m9gd", "Errors.Project.Role.NotFound")
	}
	return preConditions.ProjectResourceOwner, nil
}
