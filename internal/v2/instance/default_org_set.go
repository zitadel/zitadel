package instance

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const DefaultOrgSetType = eventTypePrefix + "default.org.set"

type defaultOrgSetPayload struct {
	OrgID string `json:"orgId"`
}

type DefaultOrgSetEvent eventstore.Event[defaultOrgSetPayload]

var _ eventstore.TypeChecker = (*DefaultOrgSetEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *DefaultOrgSetEvent) ActionType() string {
	return DefaultOrgSetType
}

func DefaultOrgSetEventFromStorage(event *eventstore.StorageEvent) (e *DefaultOrgSetEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTA-81mI1", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[defaultOrgSetPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &DefaultOrgSetEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
