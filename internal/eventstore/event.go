package eventstore

import (
	"encoding/json"
	"log/slog"
	"maps"
	"reflect"
	"slices"
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

// Command is the intent to store an event into the eventstore
type Command interface {
	action
	// Payload returns the payload of the event. It represents the changed fields by the event
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

type logValue struct {
	event Event
}

func eventToLogValue(event Event) slog.LogValuer {
	return &logValue{
		event: event,
	}
}

func (lv *logValue) LogValue() slog.Value {
	attributes := make([]slog.Attr, 0, 12)
	aggregate := lv.event.Aggregate()
	attributes = append(attributes,
		slog.String("aggregate_id", aggregate.ID),
		slog.String("aggregate_type", string(aggregate.Type)),
		slog.String("resource_owner", aggregate.ResourceOwner),
		slog.String("instance_id", aggregate.InstanceID),
		slog.String("version", string(aggregate.Version)),
		slog.String("creator", lv.event.Creator()),
		slog.String("event_type", string(lv.event.Type())),
		slog.Uint64("revision", uint64(lv.event.Revision())),
		slog.Uint64("sequence", lv.event.Sequence()),
		slog.Time("created_at", lv.event.CreatedAt()),
		slog.String("position", lv.event.Position().String()),
	)

	var m map[string]any
	err := lv.event.Unmarshal(&m)
	if err != nil {
		attributes = append(attributes,
			slog.String("msg", "failed to unmarshal event for logging"),
			slog.String("err", err.Error()),
		)
		return slog.GroupValue(attributes...)
	}
	attributes = append(attributes,
		slog.Any("data", mapToLogValue(m)),
	)
	return slog.GroupValue(attributes...)
}

// mapToLogValue converts a map[string]any to a [slog.Value], handling nested maps recursively.
func mapToLogValue(m map[string]any) slog.Value {
	attributes := make([]slog.Attr, 0, len(m))
	for _, key := range slices.Sorted(maps.Keys(m)) {
		value := m[key]
		if nestedMap, ok := value.(map[string]any); ok {
			value = mapToLogValue(nestedMap)
		}
		attributes = append(attributes, slog.Any(key, value))
	}
	return slog.GroupValue(attributes...)
}
