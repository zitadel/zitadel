package idp

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SAMLIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID                string              `json:"id"`
	Name              string              `json:"name,omitempty"`
	Metadata          []byte              `json:"metadata,omitempty"`
	Key               *crypto.CryptoValue `json:"key,omitempty"`
	Certificate       []byte              `json:"certificate,omitempty"`
	Binding           string              `json:"binding,omitempty"`
	WithSignedRequest bool                `json:"withSignedRequest,omitempty"`
	Options
}

func NewSAMLIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	name string,
	metadata []byte,
	key *crypto.CryptoValue,
	certificate []byte,
	binding string,
	withSignedRequest bool,
	options Options,
) *SAMLIDPAddedEvent {
	return &SAMLIDPAddedEvent{
		BaseEvent:         *base,
		ID:                id,
		Name:              name,
		Metadata:          metadata,
		Key:               key,
		Certificate:       certificate,
		Binding:           binding,
		WithSignedRequest: withSignedRequest,
		Options:           options,
	}
}

func (e *SAMLIDPAddedEvent) Payload() interface{} {
	return e
}

func (e *SAMLIDPAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SAMLIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &SAMLIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IDP-v9uajo3k71", "unable to unmarshal event")
	}

	return e, nil
}

type SAMLIDPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID                string              `json:"id"`
	Name              *string             `json:"name,omitempty"`
	Metadata          []byte              `json:"metadata,omitempty"`
	Key               *crypto.CryptoValue `json:"key,omitempty"`
	Certificate       []byte              `json:"certificate,omitempty"`
	Binding           *string             `json:"binding,omitempty"`
	WithSignedRequest *bool               `json:"withSignedRequest,omitempty"`
	OptionChanges
}

func NewSAMLIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []SAMLIDPChanges,
) (*SAMLIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IDP-cz6mnf860t", "Errors.NoChangesFound")
	}
	changedEvent := &SAMLIDPChangedEvent{
		BaseEvent: *base,
		ID:        id,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent, nil
}

type SAMLIDPChanges func(*SAMLIDPChangedEvent)

func ChangeSAMLName(name string) func(*SAMLIDPChangedEvent) {
	return func(e *SAMLIDPChangedEvent) {
		e.Name = &name
	}
}

func ChangeSAMLMetadata(metadata []byte) func(*SAMLIDPChangedEvent) {
	return func(e *SAMLIDPChangedEvent) {
		e.Metadata = metadata
	}
}

func ChangeSAMLKey(key *crypto.CryptoValue) func(*SAMLIDPChangedEvent) {
	return func(e *SAMLIDPChangedEvent) {
		e.Key = key
	}
}

func ChangeSAMLCertificate(certificate []byte) func(*SAMLIDPChangedEvent) {
	return func(e *SAMLIDPChangedEvent) {
		e.Certificate = certificate
	}
}

func ChangeSAMLBinding(binding string) func(*SAMLIDPChangedEvent) {
	return func(e *SAMLIDPChangedEvent) {
		e.Binding = &binding
	}
}

func ChangeSAMLWithSignedRequest(withSignedRequest bool) func(*SAMLIDPChangedEvent) {
	return func(e *SAMLIDPChangedEvent) {
		e.WithSignedRequest = &withSignedRequest
	}
}

func ChangeSAMLOptions(options OptionChanges) func(*SAMLIDPChangedEvent) {
	return func(e *SAMLIDPChangedEvent) {
		e.OptionChanges = options
	}
}

func (e *SAMLIDPChangedEvent) Payload() interface{} {
	return e
}

func (e *SAMLIDPChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SAMLIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &SAMLIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IDP-w1t1824tw5", "unable to unmarshal event")
	}

	return e, nil
}
