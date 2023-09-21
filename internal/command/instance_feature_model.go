package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/feature"
)

type FeatureWriteModel[T feature.SetEventType] struct {
	eventstore.WriteModel

	eventType eventstore.EventType

	Type T
}

func (wm *FeatureWriteModel[T]) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(feature.AggregateType).
		EventTypes(wm.eventType).
		Builder()
}

func (wm *FeatureWriteModel[T]) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *feature.SetEvent[T]:
			wm.Type = e.Type
		}
	}
	return wm.WriteModel.Reduce()
}

type InstanceFeatureWriteModel[T feature.SetEventType] struct {
	FeatureWriteModel[T]
}

func NewInstanceFeatureWriteModel[T feature.SetEventType](instanceID string, eventType eventstore.EventType) *InstanceFeatureWriteModel[T] {
	return &InstanceFeatureWriteModel[T]{
		FeatureWriteModel[T]{
			WriteModel: eventstore.WriteModel{
				InstanceID:    instanceID,
				ResourceOwner: instanceID,
			},
			eventType: eventType,
		},
	}
}
