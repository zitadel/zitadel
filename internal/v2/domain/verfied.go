package domain

import (
	"strings"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type verifiedPayload struct {
	Name string `json:"domain"`
}

type VerifiedEvent verifiedEvent
type verifiedEvent = eventstore.Event[verifiedPayload]

func VerifiedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*VerifiedEvent, error) {
	event, err := eventstore.EventFromStorage[verifiedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*VerifiedEvent)(event), nil
}

func (e *VerifiedEvent) HasTypeSuffix(typ string) bool {
	return strings.HasSuffix(typ, "domain.verified")
}
