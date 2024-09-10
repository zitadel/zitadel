package feature

import (
	"encoding/json"
	"strings"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const setEventSuffix = "set"

type setPayload struct {
	Value json.RawMessage
}

type SetEvent struct {
	eventstore.Event[setPayload]

	Level Level
	Key   Key
}

var _ eventstore.TypeChecker = (*SetEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *SetEvent) ActionType() string {
	return strings.Join([]string{AggregateType, c.Level.String(), c.Key.String(), setEventSuffix}, ".")
}

func SetEventFromStorage(event *eventstore.StorageEvent) (e *SetEvent, _ error) {
	eventTypeParts := strings.Split(event.Type, ".")
	if len(eventTypeParts) != 4 || eventTypeParts[0] != AggregateType || eventTypeParts[3] != setEventSuffix {
		return nil, zerrors.ThrowInvalidArgument(nil, "FEATU-0wJ8n", "Errors.Invalid.Event.Type")
	}

	level := levelFromString(eventTypeParts[1])
	key := keyFromString(eventTypeParts[2])

	if !level.IsALevel() || !key.IsAKey() {
		return nil, zerrors.ThrowInvalidArgument(nil, "FEATU-H8iYd", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[setPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &SetEvent{
		Event: eventstore.Event[setPayload]{
			StorageEvent: event,
			Payload:      payload,
		},
		Level: level,
		Key:   key,
	}, nil
}

func SetEventType(level Level, key Key) string {
	return strings.Join([]string{AggregateType, level.String(), key.String(), setEventSuffix}, ".")
}
