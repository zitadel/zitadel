package member

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	RemovedEventType = "member.removed"
)

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID string `json:"userId"`
}

func (e *RemovedEvent) CheckPrevious() bool {
	return true
}

func (e *RemovedEvent) Data() interface{} {
	return e
}

func NewMemberRemovedEvent(
	base *eventstore.BaseEvent,
	userID string,
) *RemovedEvent {

	return &RemovedEvent{
		BaseEvent: *base,
		UserID:    userID,
	}
}
