package user

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"time"
)

const (
	machineKeyEventPrefix      = machineEventPrefix + "key."
	MachineKeyAddedEventType   = machineEventPrefix + "added"
	MachineKeyRemovedEventType = machineEventPrefix + "removed"
)

type MachineKeyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	KeyID          string         `json:"keyId,omitempty"`
	Type           MachineKeyType `json:"type,omitempty"`
	ExpirationDate time.Time      `json:"expirationDate,omitempty"`
	PublicKey      []byte         `json:"publicKey,omitempty"`
}

func (e *MachineKeyAddedEvent) CheckPrevious() bool {
	return false
}

func (e *MachineKeyAddedEvent) Data() interface{} {
	return e
}

func NewMachineKeyAddedEvent(
	ctx context.Context,
	keyID string,
	keyType MachineKeyType,
	expirationDate time.Time,
	publicKey []byte) *MachineKeyAddedEvent {
	return &MachineKeyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			MachineKeyAddedEventType,
		),
		KeyID:          keyID,
		Type:           keyType,
		ExpirationDate: expirationDate,
		PublicKey:      publicKey,
	}
}

func MachineKeyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	machineAdded := &MachineKeyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, machineAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-rEs8f", "unable to unmarshal machine key added")
	}

	return machineAdded, nil
}

type MachineKeyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	KeyID string `json:"keyId,omitempty"`
}

func (e *MachineKeyRemovedEvent) CheckPrevious() bool {
	return false
}

func (e *MachineKeyRemovedEvent) Data() interface{} {
	return e
}

func NewMachineKeyRemovedEvent(
	ctx context.Context,
	keyID string) *MachineKeyRemovedEvent {
	return &MachineKeyRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			MachineKeyRemovedEventType,
		),
		KeyID: keyID,
	}
}

func MachineKeyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	machineRemoved := &MachineKeyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, machineRemoved)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-5Gm9s", "unable to unmarshal machine key removed")
	}

	return machineRemoved, nil
}
