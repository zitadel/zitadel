package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/project"
)

func (c *Commands) ChangeApplication(ctx context.Context, projectID string, appChange domain.Application, resourceOwner string) (*domain.ObjectDetails, error) {
	if projectID == "" || appChange.GetAppID() == "" || appChange.GetApplicationName() == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-4m9vS", "Errors.Project.App.Invalid")
	}

	existingApp, err := c.getApplicationWriteModel(ctx, projectID, appChange.GetAppID(), resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingApp.State == domain.AppStateUnspecified || existingApp.State == domain.AppStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-28di9", "Errors.Project.App.NotExisting")
	}
	if existingApp.Name == appChange.GetApplicationName() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2m8vx", "Errors.NoChangesFound")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingApp.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(
		ctx,
		project.NewApplicationChangedEvent(ctx, projectAgg, appChange.GetAppID(), existingApp.Name, appChange.GetApplicationName()))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingApp, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingApp.WriteModel), nil
}

func (c *Commands) DeactivateApplication(ctx context.Context, projectID, appID, resourceOwner string) (*domain.ObjectDetails, error) {
	if projectID == "" || appID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-88fi0", "Errors.IDMissing")
	}

	existingApp, err := c.getApplicationWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingApp.State == domain.AppStateUnspecified || existingApp.State == domain.AppStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-ov9d3", "Errors.Project.App.NotExisting")
	}
	if existingApp.State != domain.AppStateActive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-dsh35", "Errors.Project.App.NotActive")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingApp.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, project.NewApplicationDeactivatedEvent(ctx, projectAgg, appID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingApp, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingApp.WriteModel), nil
}

func (c *Commands) ReactivateApplication(ctx context.Context, projectID, appID, resourceOwner string) (*domain.ObjectDetails, error) {
	if projectID == "" || appID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-983dF", "Errors.IDMissing")
	}

	existingApp, err := c.getApplicationWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingApp.State == domain.AppStateUnspecified || existingApp.State == domain.AppStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-ov9d3", "Errors.Project.App.NotExisting")
	}
	if existingApp.State != domain.AppStateInactive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-1n8cM", "Errors.Project.App.NotInactive")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingApp.WriteModel)

	pushedEvents, err := c.eventstore.PushEvents(ctx, project.NewApplicationReactivatedEvent(ctx, projectAgg, appID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingApp, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingApp.WriteModel), nil
}

func (c *Commands) RemoveApplication(ctx context.Context, projectID, appID, resourceOwner string) (*domain.ObjectDetails, error) {
	if projectID == "" || appID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-1b7Jf", "Errors.IDMissing")
	}

	existingApp, err := c.getApplicationWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingApp.State == domain.AppStateUnspecified || existingApp.State == domain.AppStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-0po9s", "Errors.Project.App.NotExisting")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingApp.WriteModel)

	pushedEvents, err := c.eventstore.PushEvents(ctx, project.NewApplicationRemovedEvent(ctx, projectAgg, appID, existingApp.Name))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingApp, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingApp.WriteModel), nil
}

func (c *Commands) getApplicationWriteModel(ctx context.Context, projectID, appID, resourceOwner string) (*ApplicationWriteModel, error) {
	appWriteModel := NewApplicationWriteModelWithAppIDC(projectID, appID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, appWriteModel)
	if err != nil {
		return nil, err
	}
	return appWriteModel, nil
}
