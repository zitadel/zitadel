package user

import (
	"context"
	"encoding/json"
	"time"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	machineTokenEventPrefix = machineEventPrefix + "token."
	MachineTokenAddedType   = machineTokenEventPrefix + "added"
	MachineTokenRemovedType = machineTokenEventPrefix + "removed"
)

type MachineTokenAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TokenID    string    `json:"tokenId"`
	Expiration time.Time `json:"expiration"`
	Scopes     []string  `json:"scopes"`
}

func (e *MachineTokenAddedEvent) Data() interface{} {
	return e
}

func (e *MachineTokenAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewMachineTokenAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tokenID string,
	expiration time.Time,
	scopes []string,
) *MachineTokenAddedEvent {
	return &MachineTokenAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineTokenAddedType,
		),
		TokenID:    tokenID,
		Expiration: expiration,
		Scopes:     scopes,
	}
}

func MachineTokenAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	tokenAdded := &MachineTokenAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, tokenAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Dbges", "unable to unmarshal token added")
	}

	return tokenAdded, nil
}

type MachineTokenRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TokenID string `json:"tokenId"`
}

func (e *MachineTokenRemovedEvent) Data() interface{} {
	return e
}

func (e *MachineTokenRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewMachineTokenRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tokenID string,
) *MachineTokenRemovedEvent {
	return &MachineTokenRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineTokenRemovedType,
		),
		TokenID: tokenID,
	}
}

func MachineTokenRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	tokenRemoved := &MachineTokenRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, tokenRemoved)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Dbneg", "unable to unmarshal token added")
	}

	return tokenRemoved, nil
}
