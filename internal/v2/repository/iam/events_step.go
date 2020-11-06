package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	SetupDoneEventType    eventstore.EventType = "iam.setup.done"
	SetupStartedEventType eventstore.EventType = "iam.setup.started"
)

type Step int8

type SetupStepEvent struct {
	eventstore.BaseEvent `json:"-"`

	Step Step `json:"Step"`
}

func (e *SetupStepEvent) CheckPrevious() bool {
	return e.Type() == SetupStartedEventType
}

func (e *SetupStepEvent) Data() interface{} {
	return e
}

func NewSetupStepDoneEvent(ctx context.Context, service string) *SetupStepEvent {
	return &SetupStepEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			service,
			SetupDoneEventType,
		),
	}
}

func NewSetupStepStartedEvent(ctx context.Context, service string) *SetupStepEvent {
	return &SetupStepEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			service,
			SetupStartedEventType,
		),
	}
}
