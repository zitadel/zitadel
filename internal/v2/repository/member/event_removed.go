package member

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	RemovedEventType = "member.removed"
)

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID string `json:"userId"`
}

func (e *RemovedEvent) Data() interface{} {
	return e
}

func NewRemovedEvent(
	base *eventstore.BaseEvent,
	userID string,
) *RemovedEvent {

	return &RemovedEvent{
		BaseEvent: *base,
		UserID:    userID,
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-Ep4ip", "unable to unmarshal label policy")
	}

	return e, nil
}
