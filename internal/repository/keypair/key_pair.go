package usergrant

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"time"
)

const (
	eventTypePrefix = eventstore.EventType("key_pair.")
	AddedEventType  = eventTypePrefix + "added"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Usage      domain.KeyUsage `json:"usage"`
	Algorithm  string          `json:"algorithm"`
	PrivateKey *Key            `json:"privateKey"`
	PublicKey  *Key            `json:"publicKey"`
}

type Key struct {
	Key    *crypto.CryptoValue `json:"key"`
	Expiry time.Time           `json:"expiry"`
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	usage domain.KeyUsage,
	algorithm string,
	privateCrypto,
	publicCrypto *crypto.CryptoValue,
	privateKeyExpiration,
	publicKeyExpiration time.Time) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AddedEventType,
		),
		Usage:     usage,
		Algorithm: algorithm,
		PrivateKey: &Key{
			Key:    privateCrypto,
			Expiry: privateKeyExpiration,
		},
		PublicKey: &Key{
			Key:    publicCrypto,
			Expiry: publicKeyExpiration,
		},
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "KEY-4n8vs", "unable to unmarshal key pair added")
	}

	return e, nil
}
