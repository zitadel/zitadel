package keys

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"time"
)

const (
	machineKeyEventPrefix      = eventstore.EventType("user.machine.key.")
	MachineKeyAddedEventType   = machineKeyEventPrefix + "added"
	MachineKeyRemovedEventType = machineKeyEventPrefix + "removed"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	KeyID          string         `json:"keyId,omitempty"`
	Type           MachineKeyType `json:"type,omitempty"`
	ExpirationDate time.Time      `json:"expirationDate,omitempty"`
	PublicKey      []byte         `json:"publicKey,omitempty"`
}

func (e *AddedEvent) CheckPrevious() bool {
	return false
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func NewAddedEvent(
	ctx context.Context,
	keyID string,
	keyType MachineKeyType,
	expirationDate time.Time,
	publicKey []byte,
) *AddedEvent {
	return &AddedEvent{
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

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	machineAdded := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, machineAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-rEs8f", "unable to unmarshal machine key added")
	}

	return machineAdded, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	KeyID string `json:"keyId,omitempty"`
}

func (e *RemovedEvent) CheckPrevious() bool {
	return false
}

func (e *RemovedEvent) Data() interface{} {
	return e
}

func NewRemovedEvent(
	ctx context.Context,
	keyID string,
) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			MachineKeyRemovedEventType,
		),
		KeyID: keyID,
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	machineRemoved := &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, machineRemoved)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-5Gm9s", "unable to unmarshal machine key removed")
	}

	return machineRemoved, nil
}
