package instance

import (
	"strings"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/policy"
)

var (
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	DomainPolicyAdded DomainPolicyAddedEvent
)

type DomainPolicyAddedEvent struct {
	*policy.DomainPolicyAddedEvent
}

func DomainPolicyAddedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*DomainPolicyAddedEvent, error) {
	event, err := policy.DomainPolicyAddedEventFromStorage(e)
	if err != nil {
		return nil, err
	}
	return &DomainPolicyAddedEvent{
		DomainPolicyAddedEvent: event,
	}, nil
}

func (e DomainPolicyAddedEvent) IsType(typ string) bool {
	return strings.HasPrefix(typ, "instance") && e.DomainPolicyAddedEvent.HasTypeSuffix(typ)
}

var (
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	DomainPolicyChanged DomainPolicyChangedEvent
)

type DomainPolicyChangedEvent struct {
	*policy.DomainPolicyChangedEvent
}

func DomainPolicyChangedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*DomainPolicyChangedEvent, error) {
	event, err := policy.DomainPolicyChangedEventFromStorage(e)
	if err != nil {
		return nil, err
	}
	return &DomainPolicyChangedEvent{
		DomainPolicyChangedEvent: event,
	}, nil
}

func (e DomainPolicyChangedEvent) IsType(typ string) bool {
	return strings.HasPrefix(typ, "instance") && e.DomainPolicyChangedEvent.HasTypeSuffix(typ)
}
