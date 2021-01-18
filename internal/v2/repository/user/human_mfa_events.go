package user

import (
	"context"
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

func (e *HumanMFAInitSkippedEvent) Data() interface{} {
	return e
}

func (e *HumanMFAInitSkippedEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanMFAInitSkippedEvent(ctx context.Context) *HumanMFAInitSkippedEvent {
	return &HumanMFAInitSkippedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanMFAInitSkippedType,
		),
	}
}

func HumanMFAInitSkippedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanMFAInitSkippedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
