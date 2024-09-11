package projection

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/org"
)

type OrgState struct {
	Projection

	id string

	org.State
}

func NewOrgState(id string) *OrgState {
	// TODO: check buffer for id and return from buffer if exists
	return &OrgState{
		id: id,
	}
}

func (s *OrgState) Reducers() map[string]map[string]eventstore.ReduceEvent {
	if s.Projection.Reducers != nil {
		return s.Projection.Reducers
	}

	s.Projection.Reducers = map[string]map[string]eventstore.ReduceEvent{
		org.AggregateType: {
			org.AddedType:       s.reduceAdded,
			org.DeactivatedType: s.reduceDeactivated,
			org.ReactivatedType: s.reduceReactivated,
			org.RemovedType:     s.reduceRemoved,
		},
	}

	return s.Projection.Reducers
}

func (s *OrgState) reduceAdded(event *eventstore.StorageEvent) error {
	if !s.ShouldReduce(event) {
		return nil
	}

	s.State = org.ActiveState
	s.Set(event)
	return nil
}

func (s *OrgState) reduceDeactivated(event *eventstore.StorageEvent) error {
	if !s.ShouldReduce(event) {
		return nil
	}

	s.State = org.InactiveState
	s.Set(event)
	return nil
}

func (s *OrgState) reduceReactivated(event *eventstore.StorageEvent) error {
	if !s.ShouldReduce(event) {
		return nil
	}

	s.State = org.ActiveState
	s.Set(event)
	return nil
}

func (s *OrgState) reduceRemoved(event *eventstore.StorageEvent) error {
	if !s.ShouldReduce(event) {
		return nil
	}

	s.State = org.RemovedState
	s.Set(event)
	return nil
}

func (s *OrgState) ShouldReduce(event *eventstore.StorageEvent) bool {
	return event.Aggregate.ID == s.id && s.Projection.ShouldReduce(event)
}
