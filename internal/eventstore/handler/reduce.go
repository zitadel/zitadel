package handler

import "github.com/caos/zitadel/internal/eventstore"

//EventReducer represents the required data
//to work with events
type EventReducer struct {
	Aggregate eventstore.AggregateType
	Event     eventstore.EventType
	Reduce    Reduce
}
