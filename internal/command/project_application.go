package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddApp struct {
	Aggregate project.Aggregate
	ID        string
	Name      string
}

func (c *Commands) UpdateApplicationName(ctx context.Context, projectID string, appChange domain.Application, resourceOwner string) (*domain.ObjectDetails, error) {
	if projectID == "" || appChange.GetAppID() == "" || appChange.GetApplicationName() == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4m9vS", "Errors.Project.App.Invalid")
	}

	existingApp, err := c.getApplicationWriteModel(ctx, projectID, appChange.GetAppID(), resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingApp.State == domain.AppStateUnspecified || existingApp.State == domain.AppStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-28di9", "Errors.Project.App.NotExisting")
	}
	if existingApp.Name == appChange.GetApplicationName() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-2m8vx", "Errors.NoChangesFound")
	}
	if err := c.eventstore.FilterToQueryReducer(ctx, existingApp); err != nil {
		return nil, err
	}
	if err := c.checkPermissionUpdateApplication(ctx, existingApp.ResourceOwner, existingApp.AggregateID); err != nil {
		return nil, err
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingApp.WriteModel)
	pushedEvents, err := c.eventstore.Push(
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
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-88fi0", "Errors.IDMissing")
	}

	existingApp, err := c.getApplicationWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingApp.State == domain.AppStateUnspecified || existingApp.State == domain.AppStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-ov9d3", "Errors.Project.App.NotExisting")
	}
	if existingApp.State != domain.AppStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-dsh35", "Errors.Project.App.NotActive")
	}
	if err := c.eventstore.FilterToQueryReducer(ctx, existingApp); err != nil {
		return nil, err
	}
	if err := c.checkPermissionUpdateApplication(ctx, existingApp.ResourceOwner, existingApp.AggregateID); err != nil {
		return nil, err
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingApp.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, project.NewApplicationDeactivatedEvent(ctx, projectAgg, appID))
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
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-983dF", "Errors.IDMissing")
	}

	existingApp, err := c.getApplicationWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingApp.State == domain.AppStateUnspecified || existingApp.State == domain.AppStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-ov9d3", "Errors.Project.App.NotExisting")
	}
	if existingApp.State != domain.AppStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-1n8cM", "Errors.Project.App.NotInactive")
	}

	if err := c.checkPermissionUpdateApplication(ctx, existingApp.ResourceOwner, existingApp.AggregateID); err != nil {
		return nil, err
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingApp.WriteModel)

	pushedEvents, err := c.eventstore.Push(ctx, project.NewApplicationReactivatedEvent(ctx, projectAgg, appID))
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
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-1b7Jf", "Errors.IDMissing")
	}

	existingApp, err := c.getApplicationWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingApp.State == domain.AppStateUnspecified || existingApp.State == domain.AppStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-0po9s", "Errors.Project.App.NotExisting")
	}
	if err := c.eventstore.FilterToQueryReducer(ctx, existingApp); err != nil {
		return nil, err
	}
	if err := c.checkPermissionDeleteApp(ctx, existingApp.ResourceOwner, existingApp.AggregateID); err != nil {
		return nil, err
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingApp.WriteModel)

	entityID := ""
	samlWriteModel, err := c.getSAMLAppWriteModel(ctx, projectID, appID, resourceOwner)
	if err == nil && samlWriteModel.State != domain.AppStateUnspecified && samlWriteModel.State != domain.AppStateRemoved && samlWriteModel.saml {
		entityID = samlWriteModel.EntityID
	}

	pushedEvents, err := c.eventstore.Push(ctx, project.NewApplicationRemovedEvent(ctx, projectAgg, appID, existingApp.Name, entityID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingApp, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingApp.WriteModel), nil
}

func (c *Commands) getApplicationWriteModel(ctx context.Context, projectID, appID, resourceOwner string) (_ *ApplicationWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	appWriteModel := NewApplicationWriteModelWithAppIDC(projectID, appID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, appWriteModel)
	if err != nil {
		return nil, err
	}
	return appWriteModel, nil
}
