package user

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type humanProfileChangedPayload struct {
	FirstName         string         `json:"firstName,omitempty"`
	LastName          string         `json:"lastName,omitempty"`
	NickName          *string        `json:"nickName,omitempty"`
	DisplayName       *string        `json:"displayName,omitempty"`
	PreferredLanguage *language.Tag  `json:"preferredLanguage,omitempty"`
	Gender            *domain.Gender `json:"gender,omitempty"`
}

type HumanProfileChangedEvent humanProfileChangedEvent
type humanProfileChangedEvent = eventstore.Event[humanProfileChangedPayload]

func HumanProfileChangedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*HumanProfileChangedEvent, error) {
	event, err := eventstore.EventFromStorage[humanProfileChangedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*HumanProfileChangedEvent)(event), nil
}
