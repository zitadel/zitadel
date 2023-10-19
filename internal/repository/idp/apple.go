package idp

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type AppleIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID         string              `json:"id"`
	Name       string              `json:"name,omitempty"`
	ClientID   string              `json:"clientId"`
	TeamID     string              `json:"teamId"`
	KeyID      string              `json:"keyId"`
	PrivateKey *crypto.CryptoValue `json:"privateKey"`
	Scopes     []string            `json:"scopes,omitempty"`
	Options
}

func NewAppleIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	name,
	clientID,
	teamID,
	keyID string,
	privateKey *crypto.CryptoValue,
	scopes []string,
	options Options,
) *AppleIDPAddedEvent {
	return &AppleIDPAddedEvent{
		BaseEvent:  *base,
		ID:         id,
		Name:       name,
		ClientID:   clientID,
		TeamID:     teamID,
		KeyID:      keyID,
		PrivateKey: privateKey,
		Scopes:     scopes,
		Options:    options,
	}
}

func (e *AppleIDPAddedEvent) Payload() interface{} {
	return e
}

func (e *AppleIDPAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func AppleIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &AppleIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-Beqss", "unable to unmarshal event")
	}

	return e, nil
}

type AppleIDPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID         string              `json:"id"`
	Name       *string             `json:"name,omitempty"`
	ClientID   *string             `json:"clientId,omitempty"`
	TeamID     *string             `json:"teamId,omitempty"`
	KeyID      *string             `json:"keyId,omitempty"`
	PrivateKey *crypto.CryptoValue `json:"privateKey,omitempty"`
	Scopes     []string            `json:"scopes,omitempty"`
	OptionChanges
}

func NewAppleIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []AppleIDPChanges,
) (*AppleIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-SF3h2", "Errors.NoChangesFound")
	}
	changedEvent := &AppleIDPChangedEvent{
		BaseEvent: *base,
		ID:        id,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent, nil
}

type AppleIDPChanges func(*AppleIDPChangedEvent)

func ChangeAppleName(name string) func(*AppleIDPChangedEvent) {
	return func(e *AppleIDPChangedEvent) {
		e.Name = &name
	}
}

func ChangeAppleClientID(clientID string) func(*AppleIDPChangedEvent) {
	return func(e *AppleIDPChangedEvent) {
		e.ClientID = &clientID
	}
}

func ChangeAppleTeamID(teamID string) func(*AppleIDPChangedEvent) {
	return func(e *AppleIDPChangedEvent) {
		e.TeamID = &teamID
	}
}

func ChangeAppleKeyID(keyID string) func(*AppleIDPChangedEvent) {
	return func(e *AppleIDPChangedEvent) {
		e.KeyID = &keyID
	}
}

func ChangeApplePrivateKey(privateKey *crypto.CryptoValue) func(*AppleIDPChangedEvent) {
	return func(e *AppleIDPChangedEvent) {
		e.PrivateKey = privateKey
	}
}

func ChangeAppleScopes(scopes []string) func(*AppleIDPChangedEvent) {
	return func(e *AppleIDPChangedEvent) {
		e.Scopes = scopes
	}
}

func ChangeAppleOptions(options OptionChanges) func(*AppleIDPChangedEvent) {
	return func(e *AppleIDPChangedEvent) {
		e.OptionChanges = options
	}
}

func (e *AppleIDPChangedEvent) Payload() interface{} {
	return e
}

func (e *AppleIDPChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func AppleIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &AppleIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-NBe1s", "unable to unmarshal event")
	}

	return e, nil
}
