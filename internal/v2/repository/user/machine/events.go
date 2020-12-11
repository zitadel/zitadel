package machine

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	machineEventPrefix      = eventstore.EventType("user.machine.")
	MachineAddedEventType   = machineEventPrefix + "added"
	MachineChangedEventType = machineEventPrefix + "changed"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserName string `json:"userName"`

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func NewAddedEvent(
	ctx context.Context,
	userName,
	name,
	description string,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			MachineAddedEventType,
		),
		UserName:    userName,
		Name:        name,
		Description: description,
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	machineAdded := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, machineAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-tMv9s", "unable to unmarshal machine added")
	}

	return machineAdded, nil
}

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserName string `json:"userName"`

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (e *ChangedEvent) Data() interface{} {
	return e
}

func NewChangedEvent(
	ctx context.Context,
	userName,
	name,
	description string,
) *ChangedEvent {
	return &ChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			MachineChangedEventType,
		),
		UserName:    userName,
		Name:        name,
		Description: description,
	}
}

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	machineChanged := &ChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, machineChanged)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-4M9ds", "unable to unmarshal machine changed")
	}

	return machineChanged, nil
}
