package instance

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	ProjectSetEventType eventstore.EventType = "instance.iam.project.set"
	ConsoleSetEventType eventstore.EventType = "instance.iam.console.set"
)

type ProjectSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	ProjectID string `json:"iamProjectId"`
}

func (e *ProjectSetEvent) Data() interface{} {
	return e
}

func (e *ProjectSetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewIAMProjectSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	projectID string,
) *ProjectSetEvent {
	return &ProjectSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ProjectSetEventType,
		),
		ProjectID: projectID,
	}
}

func ProjectSetMapper(event *repository.Event) (eventstore.Event, error) {
	e := &ProjectSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-cdFZH", "unable to unmarshal global org set")
	}

	return e, nil
}

type ConsoleSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	ClientID string `json:"clientId"`
	AppID    string `json:"appId"`
}

func (e *ConsoleSetEvent) Data() interface{} {
	return e
}

func (e *ConsoleSetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewIAMConsoleSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	clientID,
	appID *string,
) *ConsoleSetEvent {
	return &ConsoleSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ConsoleSetEventType,
		),
		ClientID: *clientID,
		AppID:    *appID,
	}
}

func ConsoleSetMapper(event *repository.Event) (eventstore.Event, error) {
	e := &ConsoleSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-cdFZH", "unable to unmarshal console set")
	}

	return e, nil
}
