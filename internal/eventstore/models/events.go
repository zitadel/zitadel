package models

import (
	"github.com/caos/eventstore-lib/pkg/models"
)

var _ models.Events = (*Events)(nil)

type Events []*Event

func InitEvents() *Events {
	events := make(Events, 0)
	return &events
}

func (e *Events) Len() int {
	return len(*e)
}

func (e *Events) Get(index int) models.Event {
	if e.Len() < index {
		return nil
	}
	return (*e)[index]
}

func (e *Events) GetAll() []models.Event {
	events := make([]models.Event, e.Len())
	for idx := 0; idx < e.Len(); idx++ {
		events[idx] = e.Get(idx)
	}
	return events
}

func (e *Events) Append(event models.Event) {
	model, ok := event.(*Event)
	if !ok {
		return
	}
	*e = append(*e, model)
}

func (e *Events) Insert(position int, event models.Event) {
	if position > e.Len() {
		e.Append(event)
		return
	}
	model, ok := event.(*Event)
	if !ok {
		return
	}

	events := (*e)[:position]
	events = append(events, model)
	events = append(events, (*e)[position+1:]...)
	*e = events
}
