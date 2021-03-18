package user

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	machineEventPrefix      = userEventTypePrefix + "machine."
	MachineAddedEventType   = machineEventPrefix + "added"
	MachineChangedEventType = machineEventPrefix + "changed"
)

type MachineAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserName              string `json:"userName"`
	userLoginMustBeDomain bool   `json:"-"`

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (e *MachineAddedEvent) Data() interface{} {
	return e
}

func (e *MachineAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddUsernameUniqueConstraint(e.UserName, e.Aggregate().ResourceOwner, e.userLoginMustBeDomain)}
}

func NewMachineAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userName,
	name,
	description string,
	userLoginMustBeDomain bool,
) *MachineAddedEvent {
	return &MachineAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineAddedEventType,
		),
		UserName:              userName,
		Name:                  name,
		Description:           description,
		userLoginMustBeDomain: userLoginMustBeDomain,
	}
}

func MachineAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	machineAdded := &MachineAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, machineAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-tMv9s", "unable to unmarshal machine added")
	}

	return machineAdded, nil
}

type MachineChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserName string `json:"userName"`

	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (e *MachineChangedEvent) Data() interface{} {
	return e
}

func (e *MachineChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewMachineChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []MachineChanges,
) (*MachineChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "USER-3M9fs", "Errors.NoChangesFound")
	}
	changeEvent := &MachineChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineChangedEventType,
		),
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type MachineChanges func(event *MachineChangedEvent)

func ChangeName(name string) func(event *MachineChangedEvent) {
	return func(e *MachineChangedEvent) {
		e.Name = &name
	}
}

func ChangeDescription(description string) func(event *MachineChangedEvent) {
	return func(e *MachineChangedEvent) {
		e.Description = &description
	}
}

func MachineChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	machineChanged := &MachineChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, machineChanged)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-4M9ds", "unable to unmarshal machine changed")
	}

	return machineChanged, nil
}
