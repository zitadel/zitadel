package member

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	ChangedEventType = "member.changed"
)

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles  []string `json:"roles"`
	UserID string   `json:"userId"`
}

func (e *ChangedEvent) CheckPrevious() bool {
	return true
}

func (e *ChangedEvent) Data() interface{} {
	return e
}

func NewChangedEvent(
	base *eventstore.BaseEvent,
	userID string,
	roles ...string,
) *ChangedEvent {

	return &ChangedEvent{
		BaseEvent: *base,
		Roles:     roles,
		UserID:    userID,
	}
}

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-puqv4", "unable to unmarshal label policy")
	}

	return e, nil
}
