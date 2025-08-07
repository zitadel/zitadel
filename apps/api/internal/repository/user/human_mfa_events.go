package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	mfaEventPrefix          = humanEventPrefix + "mfa."
	HumanMFAInitSkippedType = mfaEventPrefix + "init.skipped"
)

type HumanMFAInitSkippedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanMFAInitSkippedEvent) Payload() interface{} {
	return e
}

func (e *HumanMFAInitSkippedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanMFAInitSkippedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanMFAInitSkippedEvent {
	return &HumanMFAInitSkippedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanMFAInitSkippedType,
		),
	}
}

func HumanMFAInitSkippedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &HumanMFAInitSkippedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
