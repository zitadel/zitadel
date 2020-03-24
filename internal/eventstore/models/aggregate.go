package models

import "github.com/caos/eventstore-lib/pkg/models"

var _ models.Aggregate = (*Aggregate)(nil)

type Aggregate struct {
	id             string
	typ            string
	events         []*Event
	latestSequence uint64
}

func NewAggregate(id, typ string, latestSequence uint64, events ...*Event) *Aggregate {
	return &Aggregate{id: id, typ: typ, events: events, latestSequence: latestSequence}
}

func (a *Aggregate) Type() string {
	return a.typ
}

func (a *Aggregate) ID() string {
	return a.id
}

func (a *Aggregate) Events() models.Events {
	events := make(Events, len(a.events))
	for idx, event := range a.events {
		events[idx] = event
	}

	return &events
}

func (a *Aggregate) LatestSequence() uint64 {
	return a.latestSequence
}
