package intent

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	intentEventPrefix = "intent."
	IntentAddedType   = intentEventPrefix + "added"
)

type IntentAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPID      string `json:"idpID"`
	SuccessURL string `json:"successURL"`
	FailureURL string `json:"failureURL"`
}

func (e *IntentAddedEvent) Data() interface{} {
	return e
}

func (e *IntentAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewIntentAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,

	idpID string,
	successURL string,
	failureURL string,
) *IntentAddedEvent {
	return &IntentAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			IntentAddedType,
		),
		IDPID:      idpID,
		SuccessURL: successURL,
		FailureURL: failureURL,
	}
}

func IntentAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	humanAdded := &IntentAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, humanAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "INTENT-5Gm9s", "unable to unmarshal intent added")
	}

	return humanAdded, nil
}
