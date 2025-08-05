package policy

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	loginPolicyIDPProviderPrevix             = loginPolicyPrefix + "idpprovider."
	LoginPolicyIDPProviderAddedType          = loginPolicyIDPProviderPrevix + "added"
	LoginPolicyIDPProviderRemovedType        = loginPolicyIDPProviderPrevix + "removed"
	LoginPolicyIDPProviderCascadeRemovedType = loginPolicyIDPProviderPrevix + "cascade.removed"
)

type IdentityProviderAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID     string                      `json:"idpConfigId,omitempty"`
	IDPProviderType domain.IdentityProviderType `json:"idpProviderType,omitempty"`
}

func (e *IdentityProviderAddedEvent) Payload() interface{} {
	return e
}

func (e *IdentityProviderAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewIdentityProviderAddedEvent(
	base *eventstore.BaseEvent,
	idpConfigID string,
	idpProviderType domain.IdentityProviderType,
) *IdentityProviderAddedEvent {

	return &IdentityProviderAddedEvent{
		*base,
		idpConfigID,
		idpProviderType,
	}
}

func IdentityProviderAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &IdentityProviderAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROVI-bfNnp", "Errors.Internal")
	}

	return e, nil
}

type IdentityProviderRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID string `json:"idpConfigId"`
}

func (e *IdentityProviderRemovedEvent) Payload() interface{} {
	return e
}

func (e *IdentityProviderRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewIdentityProviderRemovedEvent(
	base *eventstore.BaseEvent,
	idpConfigID string,
) *IdentityProviderRemovedEvent {
	return &IdentityProviderRemovedEvent{
		BaseEvent:   *base,
		IDPConfigID: idpConfigID,
	}
}

func IdentityProviderRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &IdentityProviderRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROVI-6H0KQ", "Errors.Internal")
	}

	return e, nil
}

type IdentityProviderCascadeRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPConfigID string `json:"idpConfigId"`
}

func (e *IdentityProviderCascadeRemovedEvent) Payload() interface{} {
	return e
}

func (e *IdentityProviderCascadeRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewIdentityProviderCascadeRemovedEvent(
	base *eventstore.BaseEvent,
	idpConfigID string,
) *IdentityProviderCascadeRemovedEvent {
	return &IdentityProviderCascadeRemovedEvent{
		BaseEvent:   *base,
		IDPConfigID: idpConfigID,
	}
}

func IdentityProviderCascadeRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &IdentityProviderCascadeRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROVI-7M9fs", "Errors.Internal")
	}

	return e, nil
}
