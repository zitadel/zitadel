package user

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type humanProfileChangedPayload struct {
	FirstName         string         `json:"firstName,omitempty"`
	LastName          string         `json:"lastName,omitempty"`
	NickName          *string        `json:"nickName,omitempty"`
	DisplayName       *string        `json:"displayName,omitempty"`
	PreferredLanguage *language.Tag  `json:"preferredLanguage,omitempty"`
	Gender            *domain.Gender `json:"gender,omitempty"`
}

type HumanProfileChangedEvent eventstore.Event[humanProfileChangedPayload]

const HumanProfileChangedType = humanPrefix + ".profile.changed"

var _ eventstore.TypeChecker = (*HumanProfileChangedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *HumanProfileChangedEvent) ActionType() string {
	return HumanProfileChangedType
}

func HumanProfileChangedEventFromStorage(event *eventstore.StorageEvent) (e *HumanProfileChangedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-Z1aFH", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[humanProfileChangedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &HumanProfileChangedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
