package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	IAMProjectSetEventType eventstore.EventType = "iam.project.iam.set"
)

type IAMProjectSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	ProjectID string `json:"iamProjectId"`
}

func (e *IAMProjectSetEvent) CheckPrevious() bool {
	return e.Type() == SetupStartedEventType
}

func (e *IAMProjectSetEvent) Data() interface{} {
	return e
}

func NewIAMProjectSetEvent(ctx context.Context, service, projectID string) *IAMProjectSetEvent {
	return &IAMProjectSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			service,
			SetupDoneEventType,
		),
		ProjectID: projectID,
	}
}
