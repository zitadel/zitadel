package idp

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type RolesInfo struct {
	OrganizationID     string `json:"organizationId"`
	OrganizationDomain string `json:"organizationDomain"`
}

type ZitadelIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID                string              `json:"id"`
	Name              string              `json:"name"`
	Issuer            string              `json:"issuer"`
	ClientID          string              `json:"clientId"`
	ClientSecret      *crypto.CryptoValue `json:"clientSecret"`
	Scopes            []string            `json:"scopes,omitempty"`
	InstanceRolesInfo []RolesInfo         `json:"instanceRolesInfo,omitempty"`
	Options
}

func NewZitadelIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	name,
	issuer,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	options Options,
	instanceRolesInfo []RolesInfo,
) *ZitadelIDPAddedEvent {
	return &ZitadelIDPAddedEvent{
		BaseEvent:         *base,
		ID:                id,
		Name:              name,
		Issuer:            issuer,
		ClientID:          clientID,
		ClientSecret:      clientSecret,
		Scopes:            scopes,
		Options:           options,
		InstanceRolesInfo: instanceRolesInfo,
	}
}

func (e *ZitadelIDPAddedEvent) Payload() interface{} {
	return e
}

func (e *ZitadelIDPAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type ZitadelIDPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID                string              `json:"id"`
	Name              *string             `json:"name,omitempty"`
	Issuer            *string             `json:"issuer,omitempty"`
	ClientID          *string             `json:"clientId,omitempty"`
	ClientSecret      *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Scopes            *[]string            `json:"scopes,omitempty"`
	InstanceRolesInfo *[]RolesInfo         `json:"instanceRolesInfo,omitempty"`
	OptionChanges
}

func NewZitadelIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []ZitadelIDPChanges,
) *ZitadelIDPChangedEvent {
	e := &ZitadelIDPChangedEvent{
		BaseEvent: *base,
		ID:        id,
	}
	for _, change := range changes {
		change(e)
	}
	return e
}

type ZitadelIDPChanges func(*ZitadelIDPChangedEvent)

func ChangeZitadelIDPName(name string) ZitadelIDPChanges {
	return func(e *ZitadelIDPChangedEvent) {
		e.Name = &name
	}
}

func ChangeZitadelIDPIssuer(issuer string) ZitadelIDPChanges {
	return func(e *ZitadelIDPChangedEvent) {
		e.Issuer = &issuer
	}
}

func ChangeZitadelIDPClientID(clientID string) ZitadelIDPChanges {
	return func(e *ZitadelIDPChangedEvent) {
		e.ClientID = &clientID
	}
}

func ChangeZitadelIDPClientSecret(clientSecret *crypto.CryptoValue) ZitadelIDPChanges {
	return func(e *ZitadelIDPChangedEvent) {
		e.ClientSecret = clientSecret
	}
}

func ChangeZitadelIDPScopes(scopes []string) ZitadelIDPChanges {
	return func(e *ZitadelIDPChangedEvent) {
		// explicitly set them to empty in case the scopes are unset
		if scopes == nil {
			scopes = make([]string, 0)
		}
		e.Scopes = &scopes
	}
}

func ChangeZitadelIDPInstanceRolesInfo(instanceRolesInfo []RolesInfo) ZitadelIDPChanges {
	return func(e *ZitadelIDPChangedEvent) {
		// explicitly set them to empty in case the instance roles are unset
		if instanceRolesInfo == nil {
			instanceRolesInfo = make([]RolesInfo, 0)
		}
		e.InstanceRolesInfo = &instanceRolesInfo
	}
}

func ChangeZitadelIDPOptions(options OptionChanges) ZitadelIDPChanges {
	return func(e *ZitadelIDPChangedEvent) {
		e.OptionChanges = options
	}
}

func (e *ZitadelIDPChangedEvent) Payload() interface{} {
	return e
}

func (e *ZitadelIDPChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
