package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	machineSecretPrefix          = machineEventPrefix + "secret."
	MachineSecretSetType         = machineSecretPrefix + "set"
	MachineSecretHashUpdatedType = machineSecretPrefix + "updated"
	MachineSecretRemovedType     = machineSecretPrefix + "removed"
)

type MachineSecretSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	// New events only use EncodedHash. However, the ClientSecret field
	// is preserved to handle events older than the switch to Passwap.
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	HashedSecret string              `json:"hashedSecret,omitempty"`
}

func (e *MachineSecretSetEvent) Payload() interface{} {
	return e
}

func (e *MachineSecretSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewMachineSecretSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	hashedSecret string,
) *MachineSecretSetEvent {
	return &MachineSecretSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineSecretSetType,
		),
		HashedSecret: hashedSecret,
	}
}

func MachineSecretSetEventMapper(event eventstore.Event) (eventstore.Event, error) {
	credentialsSet := &MachineSecretSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(credentialsSet)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-lopbqu", "unable to unmarshal machine secret set")
	}

	return credentialsSet, nil
}

type MachineSecretRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *MachineSecretRemovedEvent) Payload() interface{} {
	return e
}

func (e *MachineSecretRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewMachineSecretRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *MachineSecretRemovedEvent {
	return &MachineSecretRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineSecretRemovedType,
		),
	}
}

func MachineSecretRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	credentialsRemoved := &MachineSecretRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(credentialsRemoved)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-quox9j2", "unable to unmarshal machine secret removed")
	}

	return credentialsRemoved, nil
}

type MachineSecretHashUpdatedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	HashedSecret          string `json:"hashedSecret,omitempty"`
}

func NewMachineSecretHashUpdatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	encoded string,
) *MachineSecretHashUpdatedEvent {
	return &MachineSecretHashUpdatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineSecretHashUpdatedType,
		),
		HashedSecret: encoded,
	}
}

func (e *MachineSecretHashUpdatedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *MachineSecretHashUpdatedEvent) Payload() interface{} {
	return e
}

func (e *MachineSecretHashUpdatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
