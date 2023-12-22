package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/repository/feature"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) SetBooleanInstanceFeature(ctx context.Context, f domain.Feature, value bool) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	writeModel, err := NewInstanceFeatureWriteModel[feature.Boolean](instanceID, f)
	if err != nil {
		return nil, err
	}
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter,
		prepareSetFeature(writeModel, feature.Boolean{Boolean: value}, c.idGenerator))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		return writeModelToObjectDetails(&writeModel.FeatureWriteModel.WriteModel), nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func prepareSetFeature[T feature.SetEventType](writeModel *InstanceFeatureWriteModel[T], value T, idGenerator id.Generator) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if !writeModel.feature.IsAFeature() || writeModel.feature == domain.FeatureUnspecified {
			return nil, zerrors.ThrowPreconditionFailed(nil, "FEAT-JK3td", "Errors.Feature.NotExisting")
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
			if len(events) == 0 {
				writeModel.AggregateID, err = idGenerator.Next()
				if err != nil {
					return nil, err
				}
			}
			setEvent, err := writeModel.Set(ctx, value)
			if err != nil || setEvent == nil {
				return nil, err
			}
			return []eventstore.Command{setEvent}, nil
		}, nil
	}
}
