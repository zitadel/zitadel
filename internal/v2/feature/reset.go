package feature

import (
	"strings"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const resetEventSuffix = "reset"

type ResetEvent struct {
	eventstore.Event[eventstore.EmptyPayload]

	Level Level
	Key   Key
}

var _ eventstore.TypeChecker = (*ResetEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *ResetEvent) ActionType() string {
	return strings.Join([]string{AggregateType, c.Level.String(), c.Key.String(), resetEventSuffix}, ".")
}

func ResetEventFromStorage(event *eventstore.StorageEvent) (e *ResetEvent, _ error) {
	eventTypeParts := strings.Split(event.Type, ".")
	if len(eventTypeParts) != 4 || eventTypeParts[0] != AggregateType || eventTypeParts[3] != resetEventSuffix {
		return nil, zerrors.ThrowInvalidArgument(nil, "FEATU-hUu3q", "Errors.Invalid.Event.Type")
	}

	level := levelFromString(eventTypeParts[1])
	key := keyFromString(eventTypeParts[2])

	if !level.IsALevel() || !key.IsAKey() {
		return nil, zerrors.ThrowInvalidArgument(nil, "FEATU-Uy0TJ", "Errors.Invalid.Event.Type")
	}

	return &ResetEvent{
		Event: eventstore.Event[eventstore.EmptyPayload]{
			StorageEvent: event,
		},
		Level: level,
		Key:   key,
	}, nil
}

func ResetEventType(level Level) string {
	return strings.Join([]string{AggregateType, level.String(), setEventSuffix}, ".")
}
