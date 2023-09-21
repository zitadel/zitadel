package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/feature"
)

func (c *Commands) SetFeatureDefaultLoginInstance(ctx context.Context) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	id, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	featureAgg := feature.NewAggregate(id)
	writeModel := NewInstanceFeatureWriteModel[feature.Boolean](instanceID, feature.DefaultLoginInstanceEventType)
	cmds, err := preparation.PrepareCommands(ctx,
		c.eventstore.Filter,
		prepareSingleSetFeature[feature.Boolean](featureAgg, writeModel, feature.Boolean{B: true}))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func prepareSingleSetFeature[T feature.SetEventType](a *feature.Aggregate, writeModel *InstanceFeatureWriteModel[T], value T) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if writeModel.Type != nil {
				return nil, errors.ThrowPreconditionFailed(nil, "FEAT-3jklt", "Errors.Feature.AlreadySet") //TODO: i18n
			}
			return []eventstore.Command{
				feature.NewSetEvent[T](
					ctx,
					&a.Aggregate,
					writeModel.eventType,
					value,
				),
			}, nil
		}, nil
	}
}
