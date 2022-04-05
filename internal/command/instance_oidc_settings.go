package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/instance"
)

func (c *Commands) AddOIDCSettings(ctx context.Context, settings *domain.OIDCSettings) (*domain.ObjectDetails, error) {
	oidcSettingWriteModel, err := c.getOIDCSettings(ctx)
	if err != nil {
		return nil, err
	}
	if oidcSettingWriteModel.State.Exists() {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-d9nlw", "Errors.OIDCSettings.AlreadyExists")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&oidcSettingWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewOIDCSettingsAddedEvent(
		ctx,
		instanceAgg,
		settings.AccessTokenLifetime,
		settings.IdTokenLifetime,
		settings.RefreshTokenIdleExpiration,
		settings.RefreshTokenExpiration))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(oidcSettingWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&oidcSettingWriteModel.WriteModel), nil
}

func (c *Commands) ChangeOIDCSettings(ctx context.Context, settings *domain.OIDCSettings) (*domain.ObjectDetails, error) {
	oidcSettingWriteModel, err := c.getOIDCSettings(ctx)
	if err != nil {
		return nil, err
	}
	if !oidcSettingWriteModel.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-8snEr", "Errors.OIDCSettings.NotFound")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&oidcSettingWriteModel.WriteModel)

	changedEvent, hasChanged, err := oidcSettingWriteModel.NewChangedEvent(
		ctx,
		instanceAgg,
		settings.AccessTokenLifetime,
		settings.IdTokenLifetime,
		settings.RefreshTokenIdleExpiration,
		settings.RefreshTokenExpiration)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-398uF", "Errors.NoChangesFound")
	}
	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(oidcSettingWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&oidcSettingWriteModel.WriteModel), nil
}

func (c *Commands) getOIDCSettings(ctx context.Context) (_ *InstanceOIDCSettingsWriteModel, err error) {
	writeModel := NewInstanceOIDCSettingsWriteModel(ctx)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}
