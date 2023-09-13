package handler

import "github.com/zitadel/zitadel/internal/eventstore"

// EventReducer represents the required data
// to work with events
type EventReducer struct {
	Event  eventstore.EventType
	Reduce Reduce
}

// Reduce reduces the given event to a statement
// which is used to update the projection
type Reduce func(eventstore.Event) (*Statement, error)

// EventReducer represents the required data
// to work with aggregates
type AggregateReducer struct {
	Aggregate     eventstore.AggregateType
	EventReducers []EventReducer
}
