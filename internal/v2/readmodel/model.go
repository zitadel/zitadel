package readmodel

import "github.com/zitadel/zitadel/internal/v2/eventstore"

type readModel interface{}

type Model[T readModel] struct {
	latestPosition eventstore.GlobalPosition
	ReadModel      T
}
