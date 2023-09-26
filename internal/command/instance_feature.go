package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/repository/feature"
)

func (c *Commands) SetFeatureDefaultLoginInstance(ctx context.Context) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	writeModel := NewInstanceFeatureWriteModel[feature.Boolean](instanceID, domain.FeatureLoginDefaultOrg)
	cmds, err := preparation.PrepareCommands(ctx,
		c.eventstore.Filter,
		prepareSetFeature[feature.Boolean](writeModel, feature.Boolean{Boolean: true}, c.idGenerator, true))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) SetBooleanInstanceFeature(ctx context.Context, f domain.Feature, value, allowSetOnce bool) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	writeModel := NewInstanceFeatureWriteModel[feature.Boolean](instanceID, f)
	cmds, err := preparation.PrepareCommands(ctx,
		c.eventstore.Filter,
		prepareSetFeature[feature.Boolean](writeModel, feature.Boolean{Boolean: value}, c.idGenerator, allowSetOnce))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func prepareSetFeature[T feature.SetEventType](writeModel *InstanceFeatureWriteModel[T], value T, idGenerator id.Generator, allowSetOnce bool) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.feature == domain.FeatureUnspecified {
			return nil, errors.ThrowPreconditionFailed(nil, "FEAT-JK3td", "Errors.Feature.NotSpecified")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if allowSetOnce && len(events) > 0 {
				return nil, errors.ThrowPreconditionFailed(nil, "FEAT-3jklt", "Errors.Feature.AlreadySet") //TODO: i18n
			}
			if len(events) == 0 {
				writeModel.AggregateID, err = idGenerator.Next()
				if err != nil {
					return nil, err
				}
			}
			return []eventstore.Command{
				feature.NewSetEvent[T](
					ctx,
					&feature.NewAggregate(writeModel.AggregateID).Aggregate,
					writeModel.eventType(),
					value,
				),
			}, nil
		}, nil
	}
}
