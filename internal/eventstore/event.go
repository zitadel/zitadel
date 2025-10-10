package eventstore

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type action interface {
	Aggregate() *Aggregate

	// Creator is the userid of the user which created the action
	Creator() string
	// Type describes the action
	Type() EventType
	// Revision of the action
	Revision() uint16
}

// Command is the intend to store an event into the eventstore
type Command interface {
	action
	// Payload returns the payload of the event. It represent the changed fields by the event
	// valid types are:
	// * nil: no payload
	// * struct: which can be marshalled to json
	// * pointer: to struct which can be marshalled to json
	// * []byte: json marshalled data
	Payload() any
	// UniqueConstraints should be added for unique attributes of an event, if nil constraints will not be checked
	UniqueConstraints() []*UniqueConstraint
	// Fields should be added for fields which should be indexed for lookup, if nil fields will not be indexed
	Fields() []*FieldOperation
}

// Event is a stored activity
type Event interface {
	action

	// Sequence of the event in the aggregate
	Sequence() uint64
	// CreatedAt is the time the event was created at
	CreatedAt() time.Time
	// Position is the global position of the event
	Position() decimal.Decimal

	// Unmarshal parses the payload and stores the result
	// in the value pointed to by ptr. If ptr is nil or not a pointer,
	// Unmarshal returns an error
	Unmarshal(ptr any) error

	// Deprecated: only use for migration
	DataAsBytes() []byte
}

type EventType string

func EventData(event Command) ([]byte, error) {
	switch data := event.Payload().(type) {
	case nil:
		return nil, nil
	case []byte:
		if json.Valid(data) {
			return data, nil
		}
		return nil, zerrors.ThrowInvalidArgument(nil, "V2-6SbbS", "data bytes are not json")
	}
	dataType := reflect.TypeOf(event.Payload())
	if dataType.Kind() == reflect.Ptr {
		dataType = dataType.Elem()
	}
	if dataType.Kind() == reflect.Struct {
		dataBytes, err := json.Marshal(event.Payload())
		if err != nil {
			return nil, zerrors.ThrowInvalidArgument(err, "V2-xG87M", "could  not marshal data")
		}
		return dataBytes, nil
	}
	return nil, zerrors.ThrowInvalidArgument(nil, "V2-91NRm", "wrong type of event data")
}

type BaseEventSetter[T any] interface {
	Event
	SetBaseEvent(*BaseEvent)
	*T
}

func GenericEventMapper[T any, PT BaseEventSetter[T]](event Event) (Event, error) {
	e := PT(new(T))
	e.SetBaseEvent(BaseEventFromRepo(event))

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ES-Thai6", "unable to unmarshal event")
	}

	return e, nil
}

func isEventTypes(command Command, types ...EventType) bool {
	for _, typ := range types {
		if command.Type() == typ {
			return true
		}
	}
	return false
}
