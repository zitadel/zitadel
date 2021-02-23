package policy

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	loginPolicyIDPProviderPrevix             = loginPolicyPrefix + "idpprovider."
	LoginPolicyIDPProviderAddedType          = loginPolicyIDPProviderPrevix + "added"
	LoginPolicyIDPProviderRemovedType        = loginPolicyIDPProviderPrevix + "removed"
	LoginPolicyIDPProviderCascadeRemovedType = loginPolicyIDPProviderPrevix + "cascade.removed"
)

type IdentityProviderAddedEvent struct {
	eventstore.BaseEvent

	IDPConfigID     string                      `json:"idpConfigId,omitempty"`
	IDPProviderType domain.IdentityProviderType `json:"idpProviderType,omitempty"`
}

func (e *IdentityProviderAddedEvent) Data() interface{} {
	return e
}

func (e *IdentityProviderAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func IdentityProviderAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &IdentityProviderAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROVI-bfNnp", "Errors.Internal")
	}

	return e, nil
}

type IdentityProviderRemovedEvent struct {
	eventstore.BaseEvent

	IDPConfigID string `json:"idpConfigId"`
}

func (e *IdentityProviderRemovedEvent) Data() interface{} {
	return e
}

func (e *IdentityProviderRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func IdentityProviderRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &IdentityProviderRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROVI-6H0KQ", "Errors.Internal")
	}

	return e, nil
}

type IdentityProviderCascadeRemovedEvent struct {
	eventstore.BaseEvent

	IDPConfigID string `json:"idpConfigId"`
}

func (e *IdentityProviderCascadeRemovedEvent) Data() interface{} {
	return e
}

func (e *IdentityProviderCascadeRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func IdentityProviderCascadeRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &IdentityProviderCascadeRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROVI-7M9fs", "Errors.Internal")
	}

	return e, nil
}
