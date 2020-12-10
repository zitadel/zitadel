package mfa

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	mfaEventPrefix          = eventstore.EventType("user.human.mfa.")
	HumanMFAInitSkippedType = mfaEventPrefix + "init.skiped"
)

type InitSkippedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *InitSkippedEvent) CheckPrevious() bool {
	return true
}

func (e *InitSkippedEvent) Data() interface{} {
	return e
}

func NewHumanMFAInitSkippedEvent(ctx context.Context) *InitSkippedEvent {
	return &InitSkippedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAInitSkippedType,
		),
	}
}

func HumanMFAInitSkippedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &InitSkippedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
