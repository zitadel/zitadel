package iam

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/domain"
)

const (
	SetupDoneEventType    eventstore.EventType = "iam.setup.done"
	SetupStartedEventType eventstore.EventType = "iam.setup.started"
)

type SetupStepEvent struct {
	eventstore.BaseEvent `json:"-"`

	Step domain.Step `json:"Step"`
	Done bool        `json:"-"`
}

func (e *SetupStepEvent) Data() interface{} {
	return e
}

func (e *SetupStepEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
	return nil
}

func SetupStepMapper(event *repository.Event) (eventstore.EventReader, error) {
	step := &SetupStepEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
		Done:      eventstore.EventType(event.Type) == SetupDoneEventType,
	}
	err := json.Unmarshal(event.Data, step)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-O6rVg", "unable to unmarshal step")
	}

	return step, nil
}

func NewSetupStepDoneEvent(
	ctx context.Context,
	step domain.Step,
) *SetupStepEvent {

	return &SetupStepEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			SetupDoneEventType,
		),
		Step: step,
	}
}

func NewSetupStepStartedEvent(
	ctx context.Context,
	step domain.Step,
) *SetupStepEvent {

	return &SetupStepEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			SetupStartedEventType,
		),
		Step: step,
	}
}
