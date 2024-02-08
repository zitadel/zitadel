package org

import "github.com/zitadel/zitadel/internal/eventstore"

var (
	_           eventstore.Event = (*reactivated)(nil)
	Reactivated *reactivated
)

type reactivated struct {
	eventstore.BaseEvent
}

func NewReactivated() *reactivated {
	return new(reactivated)
}

func (*reactivated) Type() eventstore.EventType {
	return "org.reactivated"
}
