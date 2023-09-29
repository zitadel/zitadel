package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/feature"
)

type FeatureWriteModel[T feature.SetEventType] struct {
	eventstore.WriteModel

	feature domain.Feature

	Value T
}

func NewFeatureWriteModel[T feature.SetEventType](instanceID, resourceOwner string, feature domain.Feature) (*FeatureWriteModel[T], error) {
	wm := &FeatureWriteModel[T]{
		WriteModel: eventstore.WriteModel{
			InstanceID:    instanceID,
			ResourceOwner: resourceOwner,
		},
		feature: feature,
	}
	if wm.Value.FeatureType() != feature.Type() {
		return nil, errors.ThrowPreconditionFailed(nil, "FEAT-AS4k1", "Errors.Feature.InvalidValue")
	}
	return wm, nil
}
func (wm *FeatureWriteModel[T]) Set(ctx context.Context, value T) (event *feature.SetEvent[T], err error) {
	if wm.Value == value {
		return nil, nil
	}
	return feature.NewSetEvent[T](
		ctx,
		&feature.NewAggregate(wm.AggregateID, wm.ResourceOwner).Aggregate,
		wm.eventType(),
		value,
	), nil
}

func (wm *FeatureWriteModel[T]) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(feature.AggregateType).
		EventTypes(wm.eventType()).
		Builder()
}

func (wm *FeatureWriteModel[T]) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *feature.SetEvent[T]:
			wm.Value = e.Value
		default:
			return errors.ThrowPreconditionFailed(nil, "FEAT-SDfjk", "Errors.Feature.TypeNotSupported")
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *FeatureWriteModel[T]) eventType() eventstore.EventType {
	return feature.EventTypeFromFeature(wm.feature)
}

type InstanceFeatureWriteModel[T feature.SetEventType] struct {
	FeatureWriteModel[T]
}

func NewInstanceFeatureWriteModel[T feature.SetEventType](instanceID string, feature domain.Feature) (*InstanceFeatureWriteModel[T], error) {
	wm, err := NewFeatureWriteModel[T](instanceID, instanceID, feature)
	if err != nil {
		return nil, err
	}
	return &InstanceFeatureWriteModel[T]{
		FeatureWriteModel: *wm,
	}, nil
}
