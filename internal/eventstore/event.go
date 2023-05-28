package eventstore

import (
	"encoding/json"
	"reflect"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/v3"
)

// Command is the intend to store an event into the eventstore
type Command = eventstore.Command

// Event is a stored activity
type Event = eventstore.Event

func EventData(event Command) ([]byte, error) {
	switch data := event.Payload().(type) {
	case nil:
		return nil, nil
	case []byte:
		if json.Valid(data) {
			return data, nil
		}
		return nil, errors.ThrowInvalidArgument(nil, "V2-6SbbS", "data bytes are not json")
	}
	dataType := reflect.TypeOf(event.Payload())
	if dataType.Kind() == reflect.Ptr {
		dataType = dataType.Elem()
	}
	if dataType.Kind() == reflect.Struct {
		dataBytes, err := json.Marshal(event.Payload())
		if err != nil {
			return nil, errors.ThrowInvalidArgument(err, "V2-xG87M", "could  not marshal data")
		}
		return dataBytes, nil
	}
	return nil, errors.ThrowInvalidArgument(nil, "V2-91NRm", "wrong type of event data")
}

type BaseEventSetter[T any] interface {
	Event
	SetBaseEvent(*BaseEvent)
	*T
}

func GenericEventMapper[T any, PT BaseEventSetter[T]](event *repository.Event) (Event, error) {
	e := PT(new(T))
	e.SetBaseEvent(BaseEventFromRepo(event))

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "V2-Thai6", "unable to unmarshal event")
	}

	return e, nil
}

func isEventTypes(event Event, types ...EventType) bool {
	for _, typ := range types {
		if event.Type() == typ {
			return true
		}
	}
	return false
}
