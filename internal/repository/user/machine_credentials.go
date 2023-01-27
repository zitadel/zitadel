package user

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	machineCredentialsPrefix             = machineEventPrefix + "credentials."
	MachineCredentialsSetType            = machineCredentialsPrefix + "set"
	MachineCredentialsRemovedType        = machineCredentialsPrefix + "removed"
	MachineCredentialsCheckSucceededType = machineCredentialsPrefix + "check.succeeded"
	MachineCredentialsCheckFailedType    = machineCredentialsPrefix + "check.failed"
)

type MachineCredentialsSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
}

func (e *MachineCredentialsSetEvent) Data() interface{} {
	return e
}

func (e *MachineCredentialsSetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewMachineCredentialsSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	clientSecret *crypto.CryptoValue,
) *MachineCredentialsSetEvent {
	return &MachineCredentialsSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineCredentialsSetType,
		),
		ClientSecret: clientSecret,
	}
}

func MachineCredentialsSetEventMapper(event *repository.Event) (eventstore.Event, error) {
	credentialsSet := &MachineCredentialsSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, credentialsSet)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-lrv2di", "unable to unmarshal machine added")
	}

	return credentialsSet, nil
}

type MachineCredentialsRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *MachineCredentialsRemovedEvent) Data() interface{} {
	return e
}

func (e *MachineCredentialsRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewMachineCredentialsRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *MachineCredentialsRemovedEvent {
	return &MachineCredentialsRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineCredentialsRemovedType,
		),
	}
}

func MachineCredentialsRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	credentialsRemoved := &MachineCredentialsRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, credentialsRemoved)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-lrv2ei", "unable to unmarshal machine added")
	}

	return credentialsRemoved, nil
}

type MachineCredentialsCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *MachineCredentialsCheckSucceededEvent) Data() interface{} {
	return e
}

func (e *MachineCredentialsCheckSucceededEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewMachineCredentialsCheckSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *MachineCredentialsCheckSucceededEvent {
	return &MachineCredentialsCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineCredentialsCheckSucceededType,
		),
	}
}

func MachineCredentialsCheckSucceededEventMapper(event *repository.Event) (eventstore.Event, error) {
	check := &MachineCredentialsCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, check)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-x9000ja", "unable to unmarshal machine added")
	}

	return check, nil
}

type MachineCredentialsCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *MachineCredentialsCheckFailedEvent) Data() interface{} {
	return e
}

func (e *MachineCredentialsCheckFailedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewMachineCredentialsCheckFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *MachineCredentialsCheckFailedEvent {
	return &MachineCredentialsCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineCredentialsCheckFailedType,
		),
	}
}

func MachineCredentialsCheckFailedEventMapper(event *repository.Event) (eventstore.Event, error) {
	check := &MachineCredentialsCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, check)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-x9000ja", "unable to unmarshal machine added")
	}

	return check, nil
}
