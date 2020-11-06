package member

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

func NewMemberAddedEvent(ctx context.Context, userID string, roles ...string) *MemberAddedEvent {
	return &MemberAddedEvent{
		BaseEvent: eventstore.BaseEvent{},
		Roles:     roles,
		UserID:    userID,
	}
}

type MemberAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles  []string `json:"roles"`
	UserID string   `json:"userId"`
}

func (e *MemberAddedEvent) CheckPrevious() bool {
	return true
}

func (e *MemberAddedEvent) Data() interface{} {
	return e
}
