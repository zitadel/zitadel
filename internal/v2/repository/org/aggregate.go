package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	iamEventTypePrefix = eventstore.EventType("org.")
)

const (
	AggregateType    = "org"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}
