package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	GlobalOrgSetEventType eventstore.EventType = "iam.global.org.set"
)

type GlobalOrgSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	OrgID string `json:"globalOrgId"`
}

func (e *GlobalOrgSetEvent) CheckPrevious() bool {
	return true
}

func (e *GlobalOrgSetEvent) Data() interface{} {
	return e
}

func NewGlobalOrgSetEventEvent(ctx context.Context, orgID string) *GlobalOrgSetEvent {
	return &GlobalOrgSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			GlobalOrgSetEventType,
		),
		OrgID: orgID,
	}
}
