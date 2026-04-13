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
	// todo (@grvijayan): will be added along with UpdateZitadelProvider changes
}
