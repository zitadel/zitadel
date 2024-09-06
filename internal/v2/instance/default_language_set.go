package instance

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const DefaultLanguageSetType = eventTypePrefix + "default.language.set"

type defaultLanguageSetPayload struct {
	Language language.Tag `json:"language"`
}

type DefaultLanguageSetEvent eventstore.Event[defaultLanguageSetPayload]

var _ eventstore.TypeChecker = (*DefaultLanguageSetEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *DefaultLanguageSetEvent) ActionType() string {
	return DefaultLanguageSetType
}

func DefaultLanguageSetEventFromStorage(event *eventstore.StorageEvent) (e *DefaultLanguageSetEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTA-kXlDN", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[defaultLanguageSetPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &DefaultLanguageSetEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
