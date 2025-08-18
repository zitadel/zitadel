package idp

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type GoogleIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string              `json:"id"`
	Name         string              `json:"name,omitempty"`
	ClientID     string              `json:"clientId"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret"`
	Scopes       []string            `json:"scopes,omitempty"`
	Options
}

func NewGoogleIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	options Options,
) *GoogleIDPAddedEvent {
	return &GoogleIDPAddedEvent{
		BaseEvent:    *base,
		ID:           id,
		Name:         name,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Options:      options,
	}
}

func (e *GoogleIDPAddedEvent) Payload() any {
	return e
}

func (e *GoogleIDPAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func GoogleIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GoogleIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IDP-SAff1", "unable to unmarshal event")
	}

	return e, nil
}

type GoogleIDPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string              `json:"id"`
	Name         *string             `json:"name,omitempty"`
	ClientID     *string             `json:"clientId,omitempty"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Scopes       []string            `json:"scopes,omitempty"`
	OptionChanges
}

func NewGoogleIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []GoogleIDPChanges,
) (*GoogleIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IDP-Dg3qs", "Errors.NoChangesFound")
	}
	changedEvent := &GoogleIDPChangedEvent{
		BaseEvent: *base,
		ID:        id,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent, nil
}

type GoogleIDPChanges func(*GoogleIDPChangedEvent)

func ChangeGoogleName(name string) func(*GoogleIDPChangedEvent) {
	return func(e *GoogleIDPChangedEvent) {
		e.Name = &name
	}
}
func ChangeGoogleClientID(clientID string) func(*GoogleIDPChangedEvent) {
	return func(e *GoogleIDPChangedEvent) {
		e.ClientID = &clientID
	}
}

func ChangeGoogleClientSecret(clientSecret *crypto.CryptoValue) func(*GoogleIDPChangedEvent) {
	return func(e *GoogleIDPChangedEvent) {
		e.ClientSecret = clientSecret
	}
}

func ChangeGoogleScopes(scopes []string) func(*GoogleIDPChangedEvent) {
	return func(e *GoogleIDPChangedEvent) {
		e.Scopes = scopes
	}
}

func ChangeGoogleOptions(options OptionChanges) func(*GoogleIDPChangedEvent) {
	return func(e *GoogleIDPChangedEvent) {
		e.OptionChanges = options
	}
}

func (e *GoogleIDPChangedEvent) Payload() any {
	return e
}

func (e *GoogleIDPChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func GoogleIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GoogleIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IDP-SF3t2", "unable to unmarshal event")
	}

	return e, nil
}
