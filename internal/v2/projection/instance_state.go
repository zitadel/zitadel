package projection

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/instance"
)

type InstanceState struct {
	projection
	id string

	instance.State
}

func NewInstanceStateProjection(id string) *InstanceState {
	return &InstanceState{
		id: id,
	}
}

func (s *InstanceState) Reducers() map[string]map[string]func(*eventstore.StorageEvent) error {
	return map[string]map[string]func(*eventstore.StorageEvent) error{
		instance.AggregateType: {
			instance.AddedType:   s.reduceAdded,
			instance.RemovedType: s.reduceRemoved,
		},
	}
}

func (s *InstanceState) reduceAdded(event *eventstore.StorageEvent) error {
	if !s.ShouldReduce(event) {
		return nil
	}
	s.State = instance.ActiveState
	s.projection.set(event)
	return nil
}

func (s *InstanceState) reduceRemoved(event *eventstore.StorageEvent) error {
	if !s.ShouldReduce(event) {
		return nil
	}
	s.State = instance.RemovedState
	s.projection.set(event)
	return nil
}

func (s *InstanceState) ShouldReduce(event *eventstore.StorageEvent) bool {
	return event.Aggregate.ID == s.id && s.projection.ShouldReduce(event)
}
