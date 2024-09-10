package projection

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/org"
)

type OrgState struct {
	projection

	id string

	org.State
}

func NewStateProjection(id string) *OrgState {
	// TODO: check buffer for id and return from buffer if exists
	return &OrgState{
		id: id,
	}
}

func (s *OrgState) Reducers() map[string]map[string]eventstore.ReduceEvent {
	if s.reducers != nil {
		return s.reducers
	}

	s.reducers = map[string]map[string]eventstore.ReduceEvent{
		org.AggregateType: {
			org.AddedType:       s.reduceAdded,
			org.DeactivatedType: s.reduceDeactivated,
			org.ReactivatedType: s.reduceReactivated,
			org.RemovedType:     s.reduceRemoved,
		},
	}

	return s.reducers
}

func (s *OrgState) reduceAdded(event *eventstore.StorageEvent) error {
	if !s.ShouldReduce(event) {
		return nil
	}

	s.State = org.ActiveState
	s.set(event)
	return nil
}

func (s *OrgState) reduceDeactivated(event *eventstore.StorageEvent) error {
	if !s.ShouldReduce(event) {
		return nil
	}

	s.State = org.InactiveState
	s.set(event)
	return nil
}

func (s *OrgState) reduceReactivated(event *eventstore.StorageEvent) error {
	if !s.ShouldReduce(event) {
		return nil
	}

	s.State = org.ActiveState
	s.set(event)
	return nil
}

func (s *OrgState) reduceRemoved(event *eventstore.StorageEvent) error {
	if !s.ShouldReduce(event) {
		return nil
	}

	s.State = org.RemovedState
	s.set(event)
	return nil
}

func (s *OrgState) ShouldReduce(event *eventstore.StorageEvent) bool {
	return event.Aggregate.ID == s.id && s.projection.ShouldReduce(event)
}
