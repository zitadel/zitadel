package user

import (
	"github.com/caos/zitadel/internal/eventstore"
)

const (
	AggregateType    = "user"
	AggregateVersion = "v2"
)

type Aggregate struct {
	eventstore.Aggregate
}
