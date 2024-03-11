package eventstore

import (
	"context"
	"time"
)

type EventStore struct {
	Pusher
	Querier
}

func NewEventstore(push Pusher, query Querier) *EventStore {
	return &EventStore{
		Pusher:  push,
		Querier: query,
	}
}

type one interface {
	Pusher
	Querier
}

func NewEventstoreFromOne(o one) *EventStore {
	return NewEventstore(o, o)
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
	Query(ctx context.Context, instance string, reducer Reducer, filters ...*Filter) error
}

type Aggregate struct {
	ID       string
	Type     string
	Instance string
	Owner    string
}

func (agg *Aggregate) Equals(aggregate *Aggregate) bool {
	if aggregate.ID != "" && aggregate.ID != agg.ID {
		return false
	}
	if aggregate.Type != "" && aggregate.Type != agg.Type {
		return false
	}
	if aggregate.Instance != "" && aggregate.Instance != agg.Instance {
		return false
	}
	if aggregate.Owner != "" && aggregate.Owner != agg.Owner {
		return false
	}
	return true
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

func Unmarshal[T any](event Event) (*T, error) {
	value := new(T)

	if err := event.Unmarshal(value); err != nil {
		return nil, err
	}

	return value, nil
}

type GlobalPosition struct {
	Position        float64
	InPositionOrder uint32
}

type Event interface {
	action

	// Aggregate represents the object the event lives in
	Aggregate() *Aggregate
	// Sequence of the event in the aggregate
	Sequence() uint32
	// CreatedAt is the time the event was created at
	CreatedAt() time.Time
	// Position is the global position of the event
	Position() GlobalPosition

	// Unmarshal parses the payload and stores the result
	// in the value pointed to by ptr. If ptr is nil or not a pointer,
	// Unmarshal returns an error
	Unmarshal(ptr any) error
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
	Reduce(events ...Event) error
}
