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

type healthier interface {
	Health(ctx context.Context) error
}

type GlobalPosition struct {
	Position        float64
	InPositionOrder uint32
}

func (gp GlobalPosition) IsLess(other GlobalPosition) bool {
	return gp.Position < other.Position || (gp.Position == other.Position && gp.InPositionOrder < other.InPositionOrder)
}

type Reducer interface {
	Reduce(events ...*StorageEvent) error
}

type Reduce func(events ...*StorageEvent) error
