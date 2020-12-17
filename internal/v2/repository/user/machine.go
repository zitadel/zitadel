package user

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	machineEventPrefix      = userEventTypePrefix + "machine."
	MachineAddedEventType   = machineEventPrefix + "added"
	MachineChangedEventType = machineEventPrefix + "changed"
)

type MachineAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserName string `json:"userName"`

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (e *MachineAddedEvent) Data() interface{} {
	return e
}

func NewMachineAddedEvent(
	ctx context.Context,
	userName,
	name,
	description string,
) *MachineAddedEvent {
	return &MachineAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			MachineAddedEventType,
		),
		UserName:    userName,
		Name:        name,
		Description: description,
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

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (e *MachineChangedEvent) Data() interface{} {
	return e
}

func NewMachineChangedEvent(
	ctx context.Context,
) *MachineChangedEvent {
	return &MachineChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			MachineChangedEventType,
		),
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
