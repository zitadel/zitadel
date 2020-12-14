package repository

import (
	"time"
)

//Event represents all information about a manipulation of an aggregate
type Event struct {
	//ID is a generated uuid for this event
	ID string

	//Sequence is the sequence of the event
	Sequence uint64

	//PreviousSequence is the sequence of the previous sequence
	// if it's 0 then it's the first event of this aggregate
	PreviousSequence uint64

	//CheckPreviousSequence indicates if the given PreviousSequence should be checked
	CheckPreviousSequence bool

	//CreationDate is the time the event is created
	// it's used for human readability.
	// Don't use it for event ordering,
	// time drifts in different services could cause integrity problems
	CreationDate time.Time

	//Type describes the cause of the event (e.g. user.added)
	// it should always be in past-form
	Type EventType

	//Data describe the changed fields (e.g. userName = "hodor")
	// data must always a pointer to a struct, a struct or a byte array containing json bytes
	Data []byte

	//EditorService should be a unique identifier for the service which created the event
	// it's meant for maintainability
	EditorService string
	//EditorUser should be a unique identifier for the user which created the event
	// it's meant for maintainability.
	// It's recommend to use the aggregate id of the user
	EditorUser string

	//Version describes the definition of the aggregate at a certain point in time
	// it's used in read models to reduce the events in the correct definition
	Version Version
	//AggregateID id is the unique identifier of the aggregate
	// the client must generate it by it's own
	AggregateID string
	//AggregateType describes the meaning of the aggregate for this event
	// it could an object like user
	AggregateType AggregateType
	//ResourceOwner is the organisation which owns this aggregate
	// an aggregate can only be managed by one organisation
	// use the ID of the org
	ResourceOwner string
}

//EventType is the description of the change
type EventType string

//AggregateType is the object name
type AggregateType string
