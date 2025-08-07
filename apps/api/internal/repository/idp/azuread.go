package idp

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AzureADIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID              string              `json:"id"`
	Name            string              `json:"name,omitempty"`
	ClientID        string              `json:"client_id,omitempty"`
	ClientSecret    *crypto.CryptoValue `json:"client_secret,omitempty"`
	Scopes          []string            `json:"scopes,omitempty"`
	Tenant          string              `json:"tenant,omitempty"`
	IsEmailVerified bool                `json:"isEmailVerified,omitempty"`
	Options
}

func NewAzureADIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	tenant string,
	isEmailVerified bool,
	options Options,
) *AzureADIDPAddedEvent {
	return &AzureADIDPAddedEvent{
		BaseEvent:       *base,
		ID:              id,
		Name:            name,
		ClientID:        clientID,
		ClientSecret:    clientSecret,
		Scopes:          scopes,
		Tenant:          tenant,
		IsEmailVerified: isEmailVerified,
		Options:         options,
	}
}

func (e *AzureADIDPAddedEvent) Payload() interface{} {
	return e
}

func (e *AzureADIDPAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func AzureADIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &AzureADIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IDP-Grh2g", "unable to unmarshal event")
	}

	return e, nil
}

type AzureADIDPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID              string              `json:"id"`
	Name            *string             `json:"name,omitempty"`
	ClientID        *string             `json:"client_id,omitempty"`
	ClientSecret    *crypto.CryptoValue `json:"client_secret,omitempty"`
	Scopes          []string            `json:"scopes,omitempty"`
	Tenant          *string             `json:"tenant,omitempty"`
	IsEmailVerified *bool               `json:"isEmailVerified,omitempty"`
	OptionChanges
}

func NewAzureADIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []AzureADIDPChanges,
) (*AzureADIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IDP-BH3dl", "Errors.NoChangesFound")
	}
	changedEvent := &AzureADIDPChangedEvent{
		BaseEvent: *base,
		ID:        id,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent, nil
}

type AzureADIDPChanges func(*AzureADIDPChangedEvent)

func ChangeAzureADName(name string) func(*AzureADIDPChangedEvent) {
	return func(e *AzureADIDPChangedEvent) {
		e.Name = &name
	}
}

func ChangeAzureADClientID(clientID string) func(*AzureADIDPChangedEvent) {
	return func(e *AzureADIDPChangedEvent) {
		e.ClientID = &clientID
	}
}

func ChangeAzureADClientSecret(clientSecret *crypto.CryptoValue) func(*AzureADIDPChangedEvent) {
	return func(e *AzureADIDPChangedEvent) {
		e.ClientSecret = clientSecret
	}
}

func ChangeAzureADOptions(options OptionChanges) func(*AzureADIDPChangedEvent) {
	return func(e *AzureADIDPChangedEvent) {
		e.OptionChanges = options
	}
}

func ChangeAzureADScopes(scopes []string) func(*AzureADIDPChangedEvent) {
	return func(e *AzureADIDPChangedEvent) {
		e.Scopes = scopes
	}
}

func ChangeAzureADTenant(tenant string) func(*AzureADIDPChangedEvent) {
	return func(e *AzureADIDPChangedEvent) {
		e.Tenant = &tenant
	}
}

func ChangeAzureADIsEmailVerified(isEmailVerified bool) func(*AzureADIDPChangedEvent) {
	return func(e *AzureADIDPChangedEvent) {
		e.IsEmailVerified = &isEmailVerified
	}
}

func (e *AzureADIDPChangedEvent) Payload() interface{} {
	return e
}

func (e *AzureADIDPChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func AzureADIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &AzureADIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IDP-D3gjzh", "unable to unmarshal event")
	}

	return e, nil
}
