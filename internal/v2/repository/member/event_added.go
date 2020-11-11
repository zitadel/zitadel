package member

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	AddedEventType = "member.added"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles  []string `json:"roles"`
	UserID string   `json:"userId"`
}

func (e *AddedEvent) CheckPrevious() bool {
	return true
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func NewMemberAddedEvent(
	base *eventstore.BaseEvent,
	userID string,
	roles ...string,
) *AddedEvent {

	return &AddedEvent{
		BaseEvent: *base,
		Roles:     roles,
		UserID:    userID,
	}
}
