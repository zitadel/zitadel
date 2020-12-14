package eventstore

import (
	"time"
)

type EventPusher interface {
	// EditorService is the service who wants to push the event
	EditorService() string
	//EditorUser is the user who wants to push the event
	EditorUser() string
	//Type must return an event type which should be unique in the aggregate
	Type() EventType
	//Data returns the payload of the event. It represent the changed fields by the event
	// valid types are:
	// * nil (no payload),
	// * json byte array
	// * struct which can be marshalled to json
	// * pointer to struct which can be marshalled to json
	Data() interface{}
}

type EventReader interface {
	// EditorService is the service who pushed the event
	EditorService() string
	//EditorUser is the user who pushed the event
	EditorUser() string
	//Type is the type of the event
	Type() EventType

	AggregateID() string
	AggregateType() AggregateType
	ResourceOwner() string
	AggregateVersion() Version
	Sequence() uint64
	PreviousSequence() uint64
	CreationDate() time.Time
}
