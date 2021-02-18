package iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	iamEventTypePrefix = eventstore.EventType("iam.")
)

const (
	AggregateType    = "iam"
	AggregateVersion = "v1"
)

type Aggregate struct {
	eventstore.Aggregate
}
