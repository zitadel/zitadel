package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	ProjectSetEventType eventstore.EventType = "iam.project.iam.set"
)

type ProjectSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	ProjectID string `json:"iamProjectId"`
}

func (e *ProjectSetEvent) CheckPrevious() bool {
	return e.Type() == SetupStartedEventType
}

func (e *ProjectSetEvent) Data() interface{} {
	return e
}

func NewProjectSetEvent(ctx context.Context, projectID string) *ProjectSetEvent {
	return &ProjectSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			SetupDoneEventType,
		),
		ProjectID: projectID,
	}
}
