package iam

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	ProjectSetEventType eventstore.EventType = "iam.project.iam.set"
)

type ProjectSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	ProjectID string `json:"iamProjectId"`
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

func ProjectSetMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ProjectSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-cdFZH", "unable to unmarshal global org set")
	}

	return e, nil
}
