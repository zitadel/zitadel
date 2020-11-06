package member

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

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

func NewMemberAddedEvent(ctx context.Context, eventType eventstore.EventType, service, userID string, roles ...string) *MemberAddedEvent {
	return &MemberAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			service,
			eventType,
		),
		Roles:  roles,
		UserID: userID,
	}
}
