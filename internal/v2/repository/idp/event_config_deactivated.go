package idp

import "github.com/caos/zitadel/internal/eventstore/v2"

type DeactivatedEvent struct {
	eventstore.BaseEvent

	ID string `idpConfigId`
}

func NewDeactivatedEvent(
	base *eventstore.BaseEvent,
	configID string,
) *DeactivatedEvent {

	return &DeactivatedEvent{
		BaseEvent: *base,
		ID:        configID,
	}
}

func (e *DeactivatedEvent) CheckPrevious() bool {
	return true
}

func (e *DeactivatedEvent) Data() interface{} {
	return e
}
