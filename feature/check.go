package feature

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/feature"
)

type Checker interface {
	CheckInstanceBooleanFeature(ctx context.Context, f domain.Feature) (feature.Boolean, error)
}

func NewCheck(eventstore *eventstore.Eventstore) *Check {
	return &Check{eventstore: eventstore}
}

type Check struct {
	eventstore *eventstore.Eventstore
}

func (c *Check) CheckInstanceBooleanFeature(ctx context.Context, f domain.Feature) (feature.Boolean, error) {
	return getInstanceFeature[feature.Boolean](ctx, f, c.eventstore.Filter)
}

func getInstanceFeature[T feature.SetEventType](ctx context.Context, f domain.Feature, filter preparation.FilterToQueryReducer) (T, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	writeModel, err := command.NewInstanceFeatureWriteModel[T](instanceID, f)
	if err != nil {
		return writeModel.Value, err
	}
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return writeModel.Value, err
	}
	writeModel.AppendEvents(events...)
	if err = writeModel.Reduce(); err != nil {
		return writeModel.Value, err
	}
	return writeModel.Value, nil
}
