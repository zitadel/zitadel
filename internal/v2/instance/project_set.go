package instance

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const ProjectSetType = eventTypePrefix + "iam.project.set"

type projectSetPayload struct {
	ProjectID string `json:"projectId"`
}

type ProjectSetEvent eventstore.Event[projectSetPayload]

var _ eventstore.TypeChecker = (*ProjectSetEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *ProjectSetEvent) ActionType() string {
	return ProjectSetType
}

func ProjectSetEventFromStorage(event *eventstore.StorageEvent) (e *ProjectSetEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTA-kXlDN", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[projectSetPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &ProjectSetEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
