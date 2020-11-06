package member

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

type MemberRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID string `json:"userId"`
}

func (e *MemberRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *MemberRemovedEvent) Data() interface{} {
	return e
}

func NewMemberRemovedEvent(ctx context.Context, eventType eventstore.EventType, service, userID string) *MemberRemovedEvent {
	return &MemberRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			service,
			eventType,
		),
		UserID: userID,
	}
}
