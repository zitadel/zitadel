package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func (c *Commands) prepareAddOIDCSettings(a *instance.Aggregate, accessTokenLifetime, idTokenLifetime, refreshTokenIdleExpiration, refreshTokenExpiration time.Duration) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if accessTokenLifetime == time.Duration(0) ||
			idTokenLifetime == time.Duration(0) ||
			refreshTokenIdleExpiration == time.Duration(0) ||
			refreshTokenExpiration == time.Duration(0) {
			return nil, errors.ThrowInvalidArgument(nil, "INST-10s82j", "Errors.Invalid.Argument")
		}

		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := c.getOIDCSettingsWriteModel(ctx, filter)
			if err != nil {
				return nil, err
			}
			if writeModel.State == domain.OIDCSettingsStateActive {
				return nil, errors.ThrowAlreadyExists(nil, "INST-0aaj1o", "Errors.OIDCSettings.AlreadyExists")
			}
			return []eventstore.Command{
				instance.NewOIDCSettingsAddedEvent(
					ctx,
					&a.Aggregate,
					accessTokenLifetime,
					idTokenLifetime,
					refreshTokenIdleExpiration,
					refreshTokenExpiration,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOIDCSettings(a *instance.Aggregate, accessTokenLifetime, idTokenLifetime, refreshTokenIdleExpiration, refreshTokenExpiration time.Duration) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if accessTokenLifetime == time.Duration(0) ||
			idTokenLifetime == time.Duration(0) ||
			refreshTokenIdleExpiration == time.Duration(0) ||
			refreshTokenExpiration == time.Duration(0) {
			return nil, errors.ThrowInvalidArgument(nil, "INST-10sxks", "Errors.Invalid.Argument")
		}

		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := c.getOIDCSettingsWriteModel(ctx, filter)
			if err != nil {
				return nil, err
			}
			if writeModel.State != domain.OIDCSettingsStateActive {
				return nil, errors.ThrowNotFound(nil, "INST-90s32oj", "Errors.OIDCSettings.NotFound")
			}
			changedEvent, hasChanged, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				accessTokenLifetime,
				idTokenLifetime,
				refreshTokenIdleExpiration,
				refreshTokenExpiration,
			)
			if err != nil {
				return nil, err
			}
			if !hasChanged {
				return nil, errors.ThrowPreconditionFailed(nil, "COMMAND-0pk2nu", "Errors.NoChangesFound")
			}
			return []eventstore.Command{
				changedEvent,
			}, nil
		}, nil
	}
}

func (c *Commands) AddOIDCSettings(ctx context.Context, settings *domain.OIDCSettings) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	validation := c.prepareAddOIDCSettings(instanceAgg, settings.AccessTokenLifetime, settings.IdTokenLifetime, settings.RefreshTokenIdleExpiration, settings.RefreshTokenExpiration)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) ChangeOIDCSettings(ctx context.Context, settings *domain.OIDCSettings) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	validation := c.prepareUpdateOIDCSettings(instanceAgg, settings.AccessTokenLifetime, settings.IdTokenLifetime, settings.RefreshTokenIdleExpiration, settings.RefreshTokenExpiration)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) getOIDCSettingsWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer) (_ *InstanceOIDCSettingsWriteModel, err error) {
	writeModel := NewInstanceOIDCSettingsWriteModel(ctx)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return writeModel, nil
	}
	writeModel.AppendEvents(events...)
	err = writeModel.Reduce()
	return writeModel, err
}
