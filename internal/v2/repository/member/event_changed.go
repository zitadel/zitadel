package member

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
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

func NewMemberChangedEvent(
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
