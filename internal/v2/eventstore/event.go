package eventstore

import "time"

type Event[P any] struct {
	Aggregate Aggregate
	CreatedAt time.Time
	Creator   string
	Position  GlobalPosition
	Revision  uint16
	Sequence  uint32
	Type      string
	Payload   P
}

type StoragePayload interface {
	Unmarshal(ptr any) error
}

func EventFromStorage[E Event[P], P any](event *Event[StoragePayload]) (*E, error) {
	var payload P

	if err := event.Payload.Unmarshal(&payload); err != nil {
		return nil, err
	}
	return &E{
		Aggregate: event.Aggregate,
		CreatedAt: event.CreatedAt,
		Creator:   event.Creator,
		Position:  event.Position,
		Revision:  event.Revision,
		Sequence:  event.Sequence,
		Type:      event.Type,
		Payload:   payload,
	}, nil
}
