package idp

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type DingTalkIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string              `json:"id"`
	Name         string              `json:"name,omitempty"`
	ClientID     string              `json:"clientId"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret"`
	Scopes       []string            `json:"scopes,omitempty"`
	Options
}

func NewDingTalkIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	options Options,
) *DingTalkIDPAddedEvent {
	return &DingTalkIDPAddedEvent{
		BaseEvent:    *base,
		ID:           id,
		Name:         name,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Options:      options,
	}
}

func (e *DingTalkIDPAddedEvent) Payload() interface{} {
	return e
}

func (e *DingTalkIDPAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func DingTalkIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &DingTalkIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IDP-SAff1", "unable to unmarshal event")
	}

	return e, nil
}

type DingTalkIDPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string              `json:"id"`
	Name         *string             `json:"name,omitempty"`
	ClientID     *string             `json:"clientId,omitempty"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Scopes       []string            `json:"scopes,omitempty"`
	OptionChanges
}

func NewDingTalkIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []DingTalkIDPChanges,
) (*DingTalkIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IDP-Dg3qs", "Errors.NoChangesFound")
	}
	changedEvent := &DingTalkIDPChangedEvent{
		BaseEvent: *base,
		ID:        id,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent, nil
}

type DingTalkIDPChanges func(*DingTalkIDPChangedEvent)

func ChangeDingTalkName(name string) func(*DingTalkIDPChangedEvent) {
	return func(e *DingTalkIDPChangedEvent) {
		e.Name = &name
	}
}
func ChangeDingTalkClientID(clientID string) func(*DingTalkIDPChangedEvent) {
	return func(e *DingTalkIDPChangedEvent) {
		e.ClientID = &clientID
	}
}

func ChangeDingTalkClientSecret(clientSecret *crypto.CryptoValue) func(*DingTalkIDPChangedEvent) {
	return func(e *DingTalkIDPChangedEvent) {
		e.ClientSecret = clientSecret
	}
}

func ChangeDingTalkScopes(scopes []string) func(*DingTalkIDPChangedEvent) {
	return func(e *DingTalkIDPChangedEvent) {
		e.Scopes = scopes
	}
}

func ChangeDingTalkOptions(options OptionChanges) func(*DingTalkIDPChangedEvent) {
	return func(e *DingTalkIDPChangedEvent) {
		e.OptionChanges = options
	}
}

func (e *DingTalkIDPChangedEvent) Payload() interface{} {
	return e
}

func (e *DingTalkIDPChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func DingTalkIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &DingTalkIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IDP-SF3t2", "unable to unmarshal event")
	}

	return e, nil
}
