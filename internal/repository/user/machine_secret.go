package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	machineSecretPrefix             = machineEventPrefix + "secret."
	MachineSecretSetType            = machineSecretPrefix + "set"
	MachineSecretRemovedType        = machineSecretPrefix + "removed"
	MachineSecretCheckSucceededType = machineSecretPrefix + "check.succeeded"
	MachineSecretCheckFailedType    = machineSecretPrefix + "check.failed"
)

type MachineSecretSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
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
	clientSecret *crypto.CryptoValue,
) *MachineSecretSetEvent {
	return &MachineSecretSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineSecretSetType,
		),
		ClientSecret: clientSecret,
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

type MachineSecretCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *MachineSecretCheckSucceededEvent) Payload() interface{} {
	return e
}

func (e *MachineSecretCheckSucceededEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewMachineSecretCheckSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *MachineSecretCheckSucceededEvent {
	return &MachineSecretCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineSecretCheckSucceededType,
		),
	}
}

func MachineSecretCheckSucceededEventMapper(event eventstore.Event) (eventstore.Event, error) {
	check := &MachineSecretCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(check)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-x002n1p", "unable to unmarshal machine secret check succeeded")
	}

	return check, nil
}

type MachineSecretCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *MachineSecretCheckFailedEvent) Payload() interface{} {
	return e
}

func (e *MachineSecretCheckFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewMachineSecretCheckFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *MachineSecretCheckFailedEvent {
	return &MachineSecretCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineSecretCheckFailedType,
		),
	}
}

func MachineSecretCheckFailedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	check := &MachineSecretCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(check)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-x7901b1l", "unable to unmarshal machine secret check failed")
	}

	return check, nil
}
