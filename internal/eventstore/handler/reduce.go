package handler

import "github.com/caos/zitadel/internal/eventstore"

//EventReducer represents the required data
//to work with events
type EventReducer struct {
	Event  eventstore.EventType
	Reduce Reduce
}

//EventReducer represents the required data
//to work with aggregates
type AggregateReducer struct {
	Aggregate     eventstore.AggregateType
	EventRedusers []EventReducer
}
