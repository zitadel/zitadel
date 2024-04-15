package org

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
	return strings.HasPrefix(typ, "org") && e.DomainPolicyAddedEvent.HasTypeSuffix(typ)
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
	return strings.HasPrefix(typ, "org") && e.DomainPolicyChangedEvent.HasTypeSuffix(typ)
}

var (
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	DomainPolicyRemoved DomainPolicyRemovedEvent
)

type DomainPolicyRemovedEvent struct {
	*policy.DomainPolicyRemovedEvent
}

func DomainPolicyRemovedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*DomainPolicyRemovedEvent, error) {
	event, err := policy.DomainPolicyRemovedEventFromStorage(e)
	if err != nil {
		return nil, err
	}
	return &DomainPolicyRemovedEvent{
		DomainPolicyRemovedEvent: event,
	}, nil
}

func (e DomainPolicyRemovedEvent) IsType(typ string) bool {
	return strings.HasPrefix(typ, "org") && e.DomainPolicyRemovedEvent.HasTypeSuffix(typ)
}
