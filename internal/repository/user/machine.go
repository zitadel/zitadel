package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	machineEventPrefix      = userEventTypePrefix + "machine."
	MachineAddedEventType   = machineEventPrefix + "added"
	MachineChangedEventType = machineEventPrefix + "changed"
)

type MachineAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserName              string `json:"userName"`
	userLoginMustBeDomain bool

	Name            string               `json:"name,omitempty"`
	Description     string               `json:"description,omitempty"`
	AccessTokenType domain.OIDCTokenType `json:"accessTokenType,omitempty"`
}

func (e *MachineAddedEvent) Payload() interface{} {
	return e
}

func (e *MachineAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddUsernameUniqueConstraint(e.UserName, e.Aggregate().ResourceOwner, e.userLoginMustBeDomain)}
}

func NewMachineAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userName,
	name,
	description string,
	userLoginMustBeDomain bool,
	accessTokenType domain.OIDCTokenType,
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
		AccessTokenType:       accessTokenType,
	}
}

func MachineAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	machineAdded := &MachineAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(machineAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-tMv9s", "unable to unmarshal machine added")
	}

	return machineAdded, nil
}

type MachineChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name            *string               `json:"name,omitempty"`
	Description     *string               `json:"description,omitempty"`
	AccessTokenType *domain.OIDCTokenType `json:"accessTokenType,omitempty"`
}

func (e *MachineChangedEvent) Payload() interface{} {
	return e
}

func (e *MachineChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func ChangeAccessTokenType(accessTokenType domain.OIDCTokenType) func(event *MachineChangedEvent) {
	return func(e *MachineChangedEvent) {
		e.AccessTokenType = &accessTokenType
	}
}

func MachineChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	machineChanged := &MachineChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(machineChanged)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-4M9ds", "unable to unmarshal machine changed")
	}

	return machineChanged, nil
}
