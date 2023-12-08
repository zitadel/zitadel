package keypair

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
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

func (e *AddedEvent) Payload() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func AddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "KEY-4n8vs", "unable to unmarshal key pair added")
	}

	return e, nil
}
