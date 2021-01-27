package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

func (r *CommandSide) ChangeApplication(ctx context.Context, appChange *domain.Application, resourceOwner string) (*domain.Application, error) {
	if !appChange.IsValid(false) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4m9vS", "Errors.Project.App.Invalid")
	}

	existingApp, err := r.getApplicationWriteModel(ctx, appChange.AggregateID, appChange.AppID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingApp.State == domain.AppStateUnspecified || existingApp.State == domain.AppStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-28di9", "Errors.Project.App.NotExisting")
	}
	if existingApp.Name == appChange.Name {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-2m8vx", "Errors.NoChangesFound")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingApp.WriteModel)
	projectAgg.PushEvents(project.NewApplicationChangedEvent(ctx, appChange.AppID, existingApp.Name, appChange.Name, appChange.AggregateID))

	err = r.eventstore.PushAggregate(ctx, existingApp, projectAgg)
	if err != nil {
		return nil, err
	}

	return applicationWriteModelToApplication(existingApp), nil
}

func (r *CommandSide) DeactivateApplication(ctx context.Context, projectID, appID, resourceOwner string) error {
	if projectID == "" || appID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-88fi0", "Errors.IDMissing")
	}

	existingApp, err := r.getApplicationWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return err
	}
	if existingApp.State == domain.AppStateUnspecified || existingApp.State == domain.AppStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-ov9d3", "Errors.Project.App.NotExisting")
	}
	if existingApp.State != domain.AppStateActive {
		return caos_errs.ThrowNotFound(nil, "COMMAND-ov9d3", "Errors.Project.App.NotActive")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingApp.WriteModel)
	projectAgg.PushEvents(project.NewApplicationDeactivatedEvent(ctx, appID))

	return r.eventstore.PushAggregate(ctx, existingApp, projectAgg)
}

func (r *CommandSide) ReactivateApplication(ctx context.Context, projectID, appID, resourceOwner string) error {
	if projectID == "" || appID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-983dF", "Errors.IDMissing")
	}

	existingApp, err := r.getApplicationWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return err
	}
	if existingApp.State == domain.AppStateUnspecified || existingApp.State == domain.AppStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-ov9d3", "Errors.Project.App.NotExisting")
	}
	if existingApp.State != domain.AppStateInactive {
		return caos_errs.ThrowNotFound(nil, "COMMAND-1n8cM", "Errors.Project.App.NotInactive")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingApp.WriteModel)
	projectAgg.PushEvents(project.NewApplicationReactivatedEvent(ctx, appID))

	return r.eventstore.PushAggregate(ctx, existingApp, projectAgg)
}

func (r *CommandSide) RemoveApplication(ctx context.Context, projectID, appID, resourceOwner string) error {
	if projectID == "" || appID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-1b7Jf", "Errors.IDMissing")
	}

	existingApp, err := r.getApplicationWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return err
	}
	if existingApp.State == domain.AppStateUnspecified || existingApp.State == domain.AppStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-0po9s", "Errors.Project.App.NotExisting")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingApp.WriteModel)
	projectAgg.PushEvents(project.NewApplicationRemovedEvent(ctx, appID, existingApp.Name, projectID))

	return r.eventstore.PushAggregate(ctx, existingApp, projectAgg)
}

func (r *CommandSide) getApplicationWriteModel(ctx context.Context, projectID, appID, resourceOwner string) (*ApplicationWriteModel, error) {
	appWriteModel := NewApplicationWriteModelWithAppIDC(projectID, appID, resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, appWriteModel)
	if err != nil {
		return nil, err
	}
	return appWriteModel, nil
}
