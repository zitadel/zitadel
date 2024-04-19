package eventstore

import (
	"context"
)

func NewEventstore(querier Querier, pusher Pusher) *EventStore {
	return &EventStore{
		Pusher:  pusher,
		Querier: querier,
	}
}

func NewEventstoreFromOne(o one) *EventStore {
	return NewEventstore(o, o)
}

type EventStore struct {
	Pusher
	Querier
}

type one interface {
	Pusher
	Querier
}

type healther interface {
	Health(ctx context.Context) error
}

type Pusher interface {
	healther
	// Push writes the intents to the storage
	// if an intent implements [PushReducerIntent] [PushReducerIntent.Reduce] is called after
	// the intent was stored
	Push(ctx context.Context, intents ...PushIntent) error
}

type Querier interface {
	healther
	Query(ctx context.Context, query *Query) (eventCount int, err error)
}

type PushIntent interface {
	// Aggregate describes the object the commands will live in
	Aggregate() *Aggregate
	// Commands is an ordered list of changes on the aggregate
	Commands() []Command
	// CurrentSequence checks the current state of the aggregate.
	// The following types match the current sequence of the aggregate as described:
	// * nil or [SequenceIgnore]: Not relevant to add the commands
	// * [SequenceMatches]: Must exactly match
	// * [SequenceAtLeast]: Must be >= the given sequence
	CurrentSequence() CurrentSequence
}

// PushIntentReducer calls the [Reducer.Reduce] method after the events got created
type PushIntentReducer interface {
	PushIntent
	Reducer
}

type Command interface {
	action

	// Payload returns the payload of the event. It represent the changed fields by the event
	// valid types are:
	// * nil: no payload
	// * struct: which can be marshalled to json
	// * pointer to struct: which can be marshalled to json
	// * []byte: json marshalled data
	Payload() any
	// UniqueConstraints should be added for unique attributes of an event, if nil constraints will not be checked
	UniqueConstraints() []*UniqueConstraint
}

type GlobalPosition struct {
	Position        float64
	InPositionOrder uint32
}

type action interface {
	// Creator is the id of the user which created the action
	Creator() string
	// Type describes the action it's in the past (e.g. user.created)
	Type() string
	// Revision of the action
	Revision() uint16
}

type Reducer interface {
	Reduce(events ...*Event[StoragePayload]) error
}

type Reduce func(events ...*Event[StoragePayload]) error
