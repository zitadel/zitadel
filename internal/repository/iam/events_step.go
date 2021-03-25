package iam

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
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

func (e *SetupStepEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SetupStepMapper(event *repository.Event) (eventstore.EventReader, error) {
	step := &SetupStepEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
		Done:      eventstore.EventType(event.Type) == SetupDoneEventType,
		Step:      domain.Step1,
	}
	if len(event.Data) == 0 {
		return step, nil
	}
	err := json.Unmarshal(event.Data, step)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-O6rVg", "unable to unmarshal step")
	}

	return step, nil
}

func NewSetupStepDoneEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	step domain.Step,
) *SetupStepEvent {

	return &SetupStepEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SetupDoneEventType,
		),
		Step: step,
	}
}

func NewSetupStepStartedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	step domain.Step,
) *SetupStepEvent {

	return &SetupStepEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SetupStartedEventType,
		),
		Step: step,
	}
}
