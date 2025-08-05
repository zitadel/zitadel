package idpconfig

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	HeaderName   string `json:"headerName,omitempty"`
}

func (e *JWTConfigAddedEvent) Payload() interface{} {
	return e
}

func (e *JWTConfigAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewJWTConfigAddedEvent(
	base *eventstore.BaseEvent,
	idpConfigID,
	jwtEndpoint,
	issuer,
	keysEndpoint,
	headerName string,
) *JWTConfigAddedEvent {
	return &JWTConfigAddedEvent{
		BaseEvent:    *base,
		IDPConfigID:  idpConfigID,
		JWTEndpoint:  jwtEndpoint,
		Issuer:       issuer,
		KeysEndpoint: keysEndpoint,
		HeaderName:   headerName,
	}
}

func JWTConfigAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &JWTConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "JWT-m0fwf", "unable to unmarshal event")
	}

	return e, nil
}

type JWTConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID string `json:"idpConfigId"`

	JWTEndpoint  *string `json:"jwtEndpoint,omitempty"`
	Issuer       *string `json:"issuer,omitempty"`
	KeysEndpoint *string `json:"keysEndpoint,omitempty"`
	HeaderName   *string `json:"headerName,omitempty"`
}

func (e *JWTConfigChangedEvent) Payload() interface{} {
	return e
}

func (e *JWTConfigChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewJWTConfigChangedEvent(
	base *eventstore.BaseEvent,
	idpConfigID string,
	changes []JWTConfigChanges,
) (*JWTConfigChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IDPCONFIG-fn93s", "Errors.NoChangesFound")
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

func ChangeHeaderName(headerName string) func(*JWTConfigChangedEvent) {
	return func(e *JWTConfigChangedEvent) {
		e.HeaderName = &headerName
	}
}

func JWTConfigChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &JWTConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "JWT-fk3fs", "unable to unmarshal event")
	}

	return e, nil
}
