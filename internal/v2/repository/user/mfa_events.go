package user

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	mfaEventPrefix          = humanEventPrefix + "mfa."
	HumanMFAInitSkippedType = mfaEventPrefix + "init.skiped"
)

type HumanMFAInitSkippedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanMFAInitSkippedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanMFAInitSkippedEvent) Data() interface{} {
	return e
}

func NewHumanMFAInitSkippedEvent(base *eventstore.BaseEvent) *HumanMFAInitSkippedEvent {
	return &HumanMFAInitSkippedEvent{
		BaseEvent: *base,
	}
}

func HumanMFAInitSkippedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanMFAInitSkippedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
