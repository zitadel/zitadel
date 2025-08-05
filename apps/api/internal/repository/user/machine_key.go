package user

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	machineKeyEventPrefix      = machineEventPrefix + "key."
	MachineKeyAddedEventType   = machineKeyEventPrefix + "added"
	MachineKeyRemovedEventType = machineKeyEventPrefix + "removed"
)

type MachineKeyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	KeyID          string              `json:"keyId,omitempty"`
	KeyType        domain.AuthNKeyType `json:"type,omitempty"`
	ExpirationDate time.Time           `json:"expirationDate,omitempty"`
	PublicKey      []byte              `json:"publicKey,omitempty"`
}

func (e *MachineKeyAddedEvent) Payload() interface{} {
	return e
}

func (e *MachineKeyAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewMachineKeyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	keyID string,
	keyType domain.AuthNKeyType,
	expirationDate time.Time,
	publicKey []byte,
) *MachineKeyAddedEvent {
	return &MachineKeyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineKeyAddedEventType,
		),
		KeyID:          keyID,
		KeyType:        keyType,
		ExpirationDate: expirationDate,
		PublicKey:      publicKey,
	}
}

func MachineKeyAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	machineKeyAdded := &MachineKeyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(machineKeyAdded)
	if err != nil {
		// first events had wrong payload.
		// the keys were removed later, that's why we ignore them here.
		//nolint:errorlint
		if unwrapErr, ok := err.(*json.UnmarshalTypeError); ok && unwrapErr.Field == "publicKey" {
			return machineKeyAdded, nil
		}
		return nil, zerrors.ThrowInternal(err, "USER-p0ovS", "unable to unmarshal machine key added")
	}

	return machineKeyAdded, nil
}

type MachineKeyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	KeyID string `json:"keyId,omitempty"`
}

func (e *MachineKeyRemovedEvent) Payload() interface{} {
	return e
}

func (e *MachineKeyRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewMachineKeyRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	keyID string,
) *MachineKeyRemovedEvent {
	return &MachineKeyRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MachineKeyRemovedEventType,
		),
		KeyID: keyID,
	}
}

func MachineKeyRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	machineRemoved := &MachineKeyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(machineRemoved)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-5Gm9s", "unable to unmarshal machine key removed")
	}

	return machineRemoved, nil
}
