package idp

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type JWTIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string `json:"id"`
	Name         string `json:"name,omitempty"`
	Issuer       string `json:"issuer,omitempty"`
	JWTEndpoint  string `json:"jwtEndpoint,omitempty"`
	KeysEndpoint string `json:"keysEndpoint,omitempty"`
	HeaderName   string `json:"headerName,omitempty"`
	Options
}

func NewJWTIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	name,
	issuer,
	jwtEndpoint,
	keysEndpoint,
	headerName string,
	options Options,
) *JWTIDPAddedEvent {
	return &JWTIDPAddedEvent{
		BaseEvent:    *base,
		ID:           id,
		Name:         name,
		Issuer:       issuer,
		JWTEndpoint:  jwtEndpoint,
		KeysEndpoint: keysEndpoint,
		HeaderName:   headerName,
		Options:      options,
	}
}

func (e *JWTIDPAddedEvent) Payload() interface{} {
	return e
}

func (e *JWTIDPAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func JWTIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &JWTIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IDP-Et1dq", "unable to unmarshal event")
	}

	return e, nil
}

type JWTIDPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string  `json:"id"`
	Name         *string `json:"name,omitempty"`
	Issuer       *string `json:"issuer,omitempty"`
	JWTEndpoint  *string `json:"jwtEndpoint,omitempty"`
	KeysEndpoint *string `json:"keysEndpoint,omitempty"`
	HeaderName   *string `json:"headerName,omitempty"`
	OptionChanges
}

func NewJWTIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []JWTIDPChanges,
) (*JWTIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IDP-BH3dl", "Errors.NoChangesFound")
	}
	changedEvent := &JWTIDPChangedEvent{
		BaseEvent: *base,
		ID:        id,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent, nil
}

type JWTIDPChanges func(*JWTIDPChangedEvent)

func ChangeJWTName(name string) func(*JWTIDPChangedEvent) {
	return func(e *JWTIDPChangedEvent) {
		e.Name = &name
	}
}

func ChangeJWTIssuer(issuer string) func(*JWTIDPChangedEvent) {
	return func(e *JWTIDPChangedEvent) {
		e.Issuer = &issuer
	}
}

func ChangeJWTEndpoint(jwtEndpoint string) func(*JWTIDPChangedEvent) {
	return func(e *JWTIDPChangedEvent) {
		e.JWTEndpoint = &jwtEndpoint
	}
}

func ChangeJWTKeysEndpoint(keysEndpoint string) func(*JWTIDPChangedEvent) {
	return func(e *JWTIDPChangedEvent) {
		e.KeysEndpoint = &keysEndpoint
	}
}

func ChangeJWTHeaderName(headerName string) func(*JWTIDPChangedEvent) {
	return func(e *JWTIDPChangedEvent) {
		e.HeaderName = &headerName
	}
}

func ChangeJWTOptions(options OptionChanges) func(*JWTIDPChangedEvent) {
	return func(e *JWTIDPChangedEvent) {
		e.OptionChanges = options
	}
}

func (e *JWTIDPChangedEvent) Payload() interface{} {
	return e
}

func (e *JWTIDPChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func JWTIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &JWTIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IDP-D3gjzh", "unable to unmarshal event")
	}

	return e, nil
}
