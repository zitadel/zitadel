package idpconfig

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/domain"
)

const (
	OIDCConfigAddedEventType eventstore.EventType = "oidc.config.added"
	ConfigChangedEventType   eventstore.EventType = "oidc.config.changed"
)

type OIDCConfigAddedEvent struct {
	eventstore.BaseEvent

	IDPConfigID  string              `json:"idpConfigId"`
	ClientID     string              `json:"clientId,omitempty"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Issuer       string              `json:"issuer,omitempty"`
	Scopes       []string            `json:"scpoes,omitempty"`

	IDPDisplayNameMapping domain.OIDCMappingField `json:"idpDisplayNameMapping,omitempty"`
	UserNameMapping       domain.OIDCMappingField `json:"usernameMapping,omitempty"`
}

func (e *OIDCConfigAddedEvent) Data() interface{} {
	return e
}

func (e *OIDCConfigAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewOIDCConfigAddedEvent(
	base *eventstore.BaseEvent,
	clientID,
	idpConfigID,
	issuer string,
	clientSecret *crypto.CryptoValue,
	idpDisplayNameMapping,
	userNameMapping domain.OIDCMappingField,
	scopes ...string,
) *OIDCConfigAddedEvent {

	return &OIDCConfigAddedEvent{
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

func OIDCConfigAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &OIDCConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}

type OIDCConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID string `json:"idpConfigId"`

	ClientID     string              `json:"clientId,omitempty"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	Issuer       string              `json:"issuer,omitempty"`
	Scopes       []string            `json:"scpoes,omitempty"`

	IDPDisplayNameMapping domain.OIDCMappingField `json:"idpDisplayNameMapping,omitempty"`
	UserNameMapping       domain.OIDCMappingField `json:"usernameMapping,omitempty"`
}

func (e *OIDCConfigChangedEvent) Data() interface{} {
	return e
}

func (e *OIDCConfigChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewOIDCConfigChangedEvent(
	base *eventstore.BaseEvent,
) *OIDCConfigChangedEvent {
	return &OIDCConfigChangedEvent{
		BaseEvent: *base,
	}
}

func OIDCConfigChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &OIDCConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}
