package project

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	AggregateType    = "project"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}
