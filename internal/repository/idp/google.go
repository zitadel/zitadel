package idp

import (
	"encoding/json"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

type GoogleIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string              `json:"id"`
	ClientID     string              `json:"client_id"`
	ClientSecret *crypto.CryptoValue `json:"client_secret"`
	Options
}

func NewGoogleIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	clientID string,
	clientSecret *crypto.CryptoValue,
	options Options,
) *GoogleIDPAddedEvent {
	return &GoogleIDPAddedEvent{
		BaseEvent:    *base,
		ID:           id,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Options:      options,
	}
}

func (e *GoogleIDPAddedEvent) Data() interface{} {
	return e
}

func (e *GoogleIDPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
	//return []*eventstore.EventUniqueConstraint{NewAddIDPConfigNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner)} //TODO: ?
}

func GoogleIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &GoogleIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}

type GoogleIDPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string              `json:"id"`
	ClientID     *string             `json:"client_id,omitempty"`
	ClientSecret *crypto.CryptoValue `json:"client_secret,omitempty"`
	OptionChanges
}

func NewGoogleIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []GoogleIDPChanges,
) (*GoogleIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-S2fa1", "Errors.NoChangesFound")
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

func ChangeGoogleOptions(options OptionChanges) func(*GoogleIDPChangedEvent) {
	return func(e *GoogleIDPChangedEvent) {
		e.OptionChanges = options
	}
}

func (e *GoogleIDPChangedEvent) Data() interface{} {
	return e
}

func (e *GoogleIDPChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	//return []*eventstore.EventUniqueConstraint{NewAddGoogleIDPConfigNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner)} //TODO: ?
	return nil
}

func GoogleIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &GoogleIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-Df3f2", "unable to unmarshal event")
	}

	return e, nil
}
