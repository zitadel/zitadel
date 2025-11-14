package idp

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	UsePKCE          bool                `json:"usePKCE,omitempty"`
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
	isIDTokenMapping, usePKCE bool,
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
		UsePKCE:          usePKCE,
		Options:          options,
	}
}

func (e *OIDCIDPAddedEvent) Payload() interface{} {
	return e
}

func (e *OIDCIDPAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func OIDCIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &OIDCIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IDP-Et1dq", "unable to unmarshal event")
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
	UsePKCE          *bool               `json:"usePKCE,omitempty"`
	OptionChanges
}

func NewOIDCIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []OIDCIDPChanges,
) (*OIDCIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IDP-BH3dl", "Errors.NoChangesFound")
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

func ChangeOIDCUsePKCE(usePKCE bool) func(*OIDCIDPChangedEvent) {
	return func(e *OIDCIDPChangedEvent) {
		e.UsePKCE = &usePKCE
	}
}

func (e *OIDCIDPChangedEvent) Payload() interface{} {
	return e
}

func (e *OIDCIDPChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func OIDCIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &OIDCIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IDP-D3gjzh", "unable to unmarshal event")
	}

	return e, nil
}

type OIDCIDPMigratedAzureADEvent struct {
	AzureADIDPAddedEvent
}

func NewOIDCIDPMigratedAzureADEvent(
	base *eventstore.BaseEvent,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	tenant string,
	isEmailVerified bool,
	options Options,
) *OIDCIDPMigratedAzureADEvent {
	return &OIDCIDPMigratedAzureADEvent{
		AzureADIDPAddedEvent: AzureADIDPAddedEvent{
			BaseEvent:       *base,
			ID:              id,
			Name:            name,
			ClientID:        clientID,
			ClientSecret:    clientSecret,
			Scopes:          scopes,
			Tenant:          tenant,
			IsEmailVerified: isEmailVerified,
			Options:         options,
		},
	}
}

func (e *OIDCIDPMigratedAzureADEvent) Data() interface{} {
	return e
}

func (e *OIDCIDPMigratedAzureADEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func OIDCIDPMigratedAzureADEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := AzureADIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OIDCIDPMigratedAzureADEvent{AzureADIDPAddedEvent: *e.(*AzureADIDPAddedEvent)}, nil
}

type OIDCIDPMigratedGoogleEvent struct {
	GoogleIDPAddedEvent
}

func NewOIDCIDPMigratedGoogleEvent(
	base *eventstore.BaseEvent,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	options Options,
) *OIDCIDPMigratedGoogleEvent {
	return &OIDCIDPMigratedGoogleEvent{
		GoogleIDPAddedEvent: GoogleIDPAddedEvent{
			BaseEvent:    *base,
			ID:           id,
			Name:         name,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       scopes,
			Options:      options,
		},
	}
}

func (e *OIDCIDPMigratedGoogleEvent) Data() interface{} {
	return e
}

func (e *OIDCIDPMigratedGoogleEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func OIDCIDPMigratedGoogleEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := GoogleIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OIDCIDPMigratedGoogleEvent{GoogleIDPAddedEvent: *e.(*GoogleIDPAddedEvent)}, nil
}
