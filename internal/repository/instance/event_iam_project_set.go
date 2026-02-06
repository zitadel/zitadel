package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	ProjectSetEventType           eventstore.EventType = "instance.iam.project.set"
	ManagementConsoleSetEventType eventstore.EventType = "instance.iam.console.set"
)

type ProjectSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	ProjectID string `json:"iamProjectId"`
}

func (e *ProjectSetEvent) Payload() interface{} {
	return e
}

func (e *ProjectSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func ProjectSetMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &ProjectSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "INST-cdFZH", "unable to unmarshal global org set")
	}

	return e, nil
}

type ManagementConsoleSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	ClientID string `json:"clientId"`
	AppID    string `json:"appId"`
}

func (e *ManagementConsoleSetEvent) Payload() interface{} {
	return e
}

func (e *ManagementConsoleSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewIAMManagementConsoleSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	clientID,
	appID *string,
) *ManagementConsoleSetEvent {
	return &ManagementConsoleSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ManagementConsoleSetEventType,
		),
		ClientID: *clientID,
		AppID:    *appID,
	}
}

func ManagementConsoleSetMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &ManagementConsoleSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "INST-cdFZH", "unable to unmarshal management console set")
	}

	return e, nil
}
