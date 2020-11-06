package member

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

type MemberChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles  []string `json:"roles"`
	UserID string   `json:"userId"`
}

func (e *MemberChangedEvent) CheckPrevious() bool {
	return true
}

func (e *MemberChangedEvent) Data() interface{} {
	return e
}

func NewMemberChangedEvent(
	ctx context.Context,
	eventType eventstore.EventType,
	userID string,
	roles ...string,
) *MemberChangedEvent {

	return &MemberChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			eventType,
		),
		Roles:  roles,
		UserID: userID,
	}
}
