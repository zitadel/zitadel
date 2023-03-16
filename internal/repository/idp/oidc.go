package idp

import (
	"encoding/json"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

type OIDCIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID               string              `json:"id"`
	Name             string              `json:"name"`
	Issuer           string              `json:"issuer"`
	ClientID         string              `json:"clientId"`
	ClientSecret     *crypto.CryptoValue `json:"clientSecret"`
	Scopes           []string            `json:"scopes,omitempty"`
	IsIDTokenMapping bool                `json:"idTokenMapping,omitempty"`
	Options
}

func NewOIDCIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	name,
	issuer,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	isIDTokenMapping bool,
	options Options,
) *OIDCIDPAddedEvent {
	return &OIDCIDPAddedEvent{
		BaseEvent:        *base,
		ID:               id,
		Name:             name,
		Issuer:           issuer,
		ClientID:         clientID,
		ClientSecret:     clientSecret,
		Scopes:           scopes,
		IsIDTokenMapping: isIDTokenMapping,
		Options:          options,
	}
}

func (e *OIDCIDPAddedEvent) Data() interface{} {
	return e
}

func (e *OIDCIDPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func OIDCIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &OIDCIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-Et1dq", "unable to unmarshal event")
	}

	return e, nil
}

type OIDCIDPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID               string              `json:"id"`
	Name             *string             `json:"name,omitempty"`
	Issuer           *string             `json:"issuer,omitempty"`
	ClientID         *string             `json:"clientId,omitempty"`
	ClientSecret     *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Scopes           []string            `json:"scopes,omitempty"`
	IsIDTokenMapping *bool               `json:"idTokenMapping,omitempty"`
	OptionChanges
}

func NewOIDCIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []OIDCIDPChanges,
) (*OIDCIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-BH3dl", "Errors.NoChangesFound")
	}
	changedEvent := &OIDCIDPChangedEvent{
		BaseEvent: *base,
		ID:        id,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent, nil
}

type OIDCIDPChanges func(*OIDCIDPChangedEvent)

func ChangeOIDCName(name string) func(*OIDCIDPChangedEvent) {
	return func(e *OIDCIDPChangedEvent) {
		e.Name = &name
	}
}

func ChangeOIDCIssuer(issuer string) func(*OIDCIDPChangedEvent) {
	return func(e *OIDCIDPChangedEvent) {
		e.Issuer = &issuer
	}
}

func ChangeOIDCClientID(clientID string) func(*OIDCIDPChangedEvent) {
	return func(e *OIDCIDPChangedEvent) {
		e.ClientID = &clientID
	}
}

func ChangeOIDCClientSecret(clientSecret *crypto.CryptoValue) func(*OIDCIDPChangedEvent) {
	return func(e *OIDCIDPChangedEvent) {
		e.ClientSecret = clientSecret
	}
}

func ChangeOIDCOptions(options OptionChanges) func(*OIDCIDPChangedEvent) {
	return func(e *OIDCIDPChangedEvent) {
		e.OptionChanges = options
	}
}

func ChangeOIDCScopes(scopes []string) func(*OIDCIDPChangedEvent) {
	return func(e *OIDCIDPChangedEvent) {
		e.Scopes = scopes
	}
}

func ChangeOIDCIsIDTokenMapping(idTokenMapping bool) func(*OIDCIDPChangedEvent) {
	return func(e *OIDCIDPChangedEvent) {
		e.IsIDTokenMapping = &idTokenMapping
	}
}

func (e *OIDCIDPChangedEvent) Data() interface{} {
	return e
}

func (e *OIDCIDPChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func OIDCIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &OIDCIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-D3gjzh", "unable to unmarshal event")
	}

	return e, nil
}
