package oidc

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	ConfigAddedEventType eventstore.EventType = "oidc.config.added"
)

type ConfigAddedEvent struct {
	eventstore.BaseEvent

	IDPConfigID  string              `json:"idpConfigId"`
	ClientID     string              `json:"clientId,omitempty"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Issuer       string              `json:"issuer,omitempty"`
	Scopes       []string            `json:"scpoes,omitempty"`

	IDPDisplayNameMapping MappingField `json:"idpDisplayNameMapping,omitempty"`
	UserNameMapping       MappingField `json:"usernameMapping,omitempty"`
}

func (e *ConfigAddedEvent) Data() interface{} {
	return e
}

func NewConfigAddedEvent(
	base *eventstore.BaseEvent,
	clientID,
	idpConfigID,
	issuer string,
	clientSecret *crypto.CryptoValue,
	idpDisplayNameMapping,
	userNameMapping MappingField,
	scopes ...string,
) *ConfigAddedEvent {

	return &ConfigAddedEvent{
		BaseEvent:             *base,
		IDPConfigID:           idpConfigID,
		ClientID:              clientID,
		ClientSecret:          clientSecret,
		Issuer:                issuer,
		Scopes:                scopes,
		IDPDisplayNameMapping: idpDisplayNameMapping,
		UserNameMapping:       userNameMapping,
	}
}

func ConfigAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}
