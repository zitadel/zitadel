package eventstore

import (
	"time"
)

type Unmarshal func(ptr any) error

type Payload interface {
	Unmarshal | any
}

type Action[P Payload] struct {
	Creator  string
	Type     string
	Revision uint16
	Payload  P
}

type Command struct {
	Action[any]
	UniqueConstraints []*UniqueConstraint
}

type StorageEvent struct {
	Action[Unmarshal]

	Aggregate Aggregate
	CreatedAt time.Time
	Position  GlobalPosition
	Sequence  uint32
}

type Event[P any] struct {
	*StorageEvent
	Payload P
}

func UnmarshalPayload[P any](unmarshal Unmarshal) (P, error) {
	var payload P
	err := unmarshal(&payload)
	return payload, err
}

type EmptyPayload struct{}

type TypeChecker interface {
	ActionType() string
}

func Type[T TypeChecker]() string {
	var t T
	return t.ActionType()
}

func IsType[T TypeChecker](types ...string) bool {
	gotten := Type[T]()

	for _, typ := range types {
		if gotten == typ {
			return true
		}
	}

	return false
}
