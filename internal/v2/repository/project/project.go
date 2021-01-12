package project

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	projectEventTypePrefix = eventstore.EventType("project.")
	ProjectAdded           = projectEventTypePrefix + "added"
	ProjectChanged         = projectEventTypePrefix + "changed"
	ProjectDeactivated     = projectEventTypePrefix + "deactivated"
	ProjectReactivated     = projectEventTypePrefix + "reactivated"
	ProjectRemoved         = projectEventTypePrefix + "removed"
)

type ProjectAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name                 string `json:"name,omitempty"`
	ProjectRoleAssertion bool   `json:"projectRoleAssertion,omitempty"`
	ProjectRoleCheck     bool   `json:"projectRoleCheck,omitempty"`
}

func (e *ProjectAddedEvent) Data() interface{} {
	return e
}

func NewProjectAddedEvent(ctx context.Context, name string) *ProjectAddedEvent {
	return &ProjectAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			ProjectAdded,
		),
		Name: name,
	}
}

func ProjectAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ProjectAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-Bfg2f", "unable to unmarshal project")
	}

	return e, nil
}
