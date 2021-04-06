package iam

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
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

func NewAggregate() *Aggregate {
	return &Aggregate{
		Aggregate: eventstore.Aggregate{
			Typ:           AggregateType,
			Version:       AggregateVersion,
			ID:            domain.IAMID,
			ResourceOwner: domain.IAMID,
		},
	}
}
