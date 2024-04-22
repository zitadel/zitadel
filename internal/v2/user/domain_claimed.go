package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type domainClaimedPayload struct {
	Username          string `json:"userName"`
	TriggeredAtOrigin string `json:"triggerOrigin,omitempty"`
}

type DomainClaimedEvent domainClaimedEvent
type domainClaimedEvent = eventstore.Event[domainClaimedPayload]

func DomainClaimedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*DomainClaimedEvent, error) {
	event, err := eventstore.EventFromStorage[domainClaimedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*DomainClaimedEvent)(event), nil
}
