package command

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) checkProjectExistsOld(ctx context.Context, projectID, resourceOwner string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	projectWriteModel, err := c.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return err
	}
	if projectWriteModel.State == domain.ProjectStateUnspecified || projectWriteModel.State == domain.ProjectStateRemoved {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-EbFMN", "Errors.Project.NotFound")
	}
	return nil
}

func (c *Commands) changeProjectOld(ctx context.Context, projectChange *domain.Project, resourceOwner string) (*domain.Project, error) {
	if !projectChange.IsValid() || projectChange.AggregateID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4m9vS", "Errors.Project.Invalid")
	}

	existingProject, err := c.getProjectWriteModelByID(ctx, projectChange.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingProject.State == domain.ProjectStateUnspecified || existingProject.State == domain.ProjectStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.NotFound")
	}

	//nolint: contextcheck
	projectAgg := ProjectAggregateFromWriteModel(&existingProject.WriteModel)
	changedEvent, hasChanged, err := existingProject.NewChangedEvent(
		ctx,
		projectAgg,
		projectChange.Name,
		projectChange.ProjectRoleAssertion,
		projectChange.ProjectRoleCheck,
		projectChange.HasProjectCheck,
		projectChange.PrivateLabelingSetting)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-2M0fs", "Errors.NoChangesFound")
	}
	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingProject, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return projectWriteModelToProject(existingProject), nil
}

func (c *Commands) deactivateProjectOld(ctx context.Context, projectID string, resourceOwner string) (*domain.ObjectDetails, error) {
	existingProject, err := c.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingProject.State == domain.ProjectStateUnspecified || existingProject.State == domain.ProjectStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-112M9", "Errors.Project.NotFound")
	}
	if existingProject.State != domain.ProjectStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-mki55", "Errors.Project.NotActive")
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
	if existingProject.State == domain.ProjectStateUnspecified || existingProject.State == domain.ProjectStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.NotFound")
	}
	if existingProject.State != domain.ProjectStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-5M9bs", "Errors.Project.NotInactive")
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

func (c *Commands) removeProjectOld(ctx context.Context, projectID, resourceOwner string, cascadingUserGrantIDs ...string) (*domain.ObjectDetails, error) {
	existingProject, err := c.getProjectWriteModelByID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingProject.State == domain.ProjectStateUnspecified || existingProject.State == domain.ProjectStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-3M9sd", "Errors.Project.NotFound")
	}

	samlEntityIDsAgg, err := c.getSAMLEntityIdsWriteModelByProjectID(ctx, projectID, resourceOwner)
	if err != nil {
		return nil, err
	}

	uniqueConstraints := make([]*eventstore.UniqueConstraint, len(samlEntityIDsAgg.EntityIDs))
	for i, entityID := range samlEntityIDsAgg.EntityIDs {
		uniqueConstraints[i] = project.NewRemoveSAMLConfigEntityIDUniqueConstraint(entityID.EntityID)
	}

	//nolint: contextcheck
	projectAgg := ProjectAggregateFromWriteModel(&existingProject.WriteModel)
	events := []eventstore.Command{
		project.NewProjectRemovedEvent(ctx, projectAgg, existingProject.Name, uniqueConstraints),
	}

	for _, grantID := range cascadingUserGrantIDs {
		event, _, err := c.removeUserGrant(ctx, grantID, "", true)
		if err != nil {
			logging.WithFields("usergrantid", grantID).WithError(err).Warn("could not cascade remove user grant")
			continue
		}
		events = append(events, event)
	}

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingProject, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingProject.WriteModel), nil
}

func (c *Commands) checkProjectGrantPreConditionOld(ctx context.Context, projectGrant *domain.ProjectGrant, resourceOwner string) error {
	preConditions := NewProjectGrantPreConditionReadModel(projectGrant.AggregateID, projectGrant.GrantedOrgID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, preConditions)
	if err != nil {
		return err
	}
	if !preConditions.ProjectExists {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-m9gsd", "Errors.Project.NotFound")
	}
	if !preConditions.GrantedOrgExists {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-3m9gg", "Errors.Org.NotFound")
	}
	if projectGrant.HasInvalidRoles(preConditions.ExistingRoleKeys) {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-6m9gd", "Errors.Project.Role.NotFound")
	}
	return nil
}
