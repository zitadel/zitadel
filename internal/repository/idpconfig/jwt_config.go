package idpconfig

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	JWTConfigAddedEventType   eventstore.EventType = "jwt.config.added"
	JWTConfigChangedEventType eventstore.EventType = "jwt.config.changed"
)

type JWTConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID  string `json:"idpConfigId"`
	JWTEndpoint  string `json:"jwtEndpoint,omitempty"`
	Issuer       string `json:"issuer,omitempty"`
	KeysEndpoint string `json:"keysEndpoint,omitempty"`
}

func (e *JWTConfigAddedEvent) Data() interface{} {
	return e
}

func (e *JWTConfigAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewJWTConfigAddedEvent(
	base *eventstore.BaseEvent,
	idpConfigID,
	jwtEndpoint,
	issuer,
	keysEndpoint string,
) *JWTConfigAddedEvent {
	return &JWTConfigAddedEvent{
		BaseEvent:    *base,
		IDPConfigID:  idpConfigID,
		JWTEndpoint:  jwtEndpoint,
		Issuer:       issuer,
		KeysEndpoint: keysEndpoint,
	}
}

func JWTConfigAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &JWTConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "JWT-m0fwf", "unable to unmarshal event")
	}

	return e, nil
}

type JWTConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID string `json:"idpConfigId"`

	JWTEndpoint  *string `json:"jwtEndpoint,omitempty"`
	Issuer       *string `json:"issuer,omitempty"`
	KeysEndpoint *string `json:"keysEndpoint,omitempty"`
}

func (e *JWTConfigChangedEvent) Data() interface{} {
	return e
}

func (e *JWTConfigChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewJWTConfigChangedEvent(
	base *eventstore.BaseEvent,
	idpConfigID string,
	changes []JWTConfigChanges,
) (*JWTConfigChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDPCONFIG-fn93s", "Errors.NoChangesFound")
	}
	changeEvent := &JWTConfigChangedEvent{
		BaseEvent:   *base,
		IDPConfigID: idpConfigID,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type JWTConfigChanges func(*JWTConfigChangedEvent)

func ChangeJWTEndpoint(jwtEndpoint string) func(*JWTConfigChangedEvent) {
	return func(e *JWTConfigChangedEvent) {
		e.JWTEndpoint = &jwtEndpoint
	}
}

func ChangeJWTIssuer(issuer string) func(*JWTConfigChangedEvent) {
	return func(e *JWTConfigChangedEvent) {
		e.Issuer = &issuer
	}
}

func ChangeKeysEndpoint(keysEndpoint string) func(*JWTConfigChangedEvent) {
	return func(e *JWTConfigChangedEvent) {
		e.KeysEndpoint = &keysEndpoint
	}
}

func JWTConfigChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &JWTConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "JWT-fk3fs", "unable to unmarshal event")
	}

	return e, nil
}
