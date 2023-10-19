package keypair

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	AddedCertificateEventType = eventTypePrefix + "certificate.added"
)

type AddedCertificateEvent struct {
	eventstore.BaseEvent `json:"-"`

	Certificate *Key `json:"certificate"`
}

func (e *AddedCertificateEvent) Payload() interface{} {
	return e
}

func (e *AddedCertificateEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewAddedCertificateEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	certificateCrypto *crypto.CryptoValue,
	certificateExpiration time.Time) *AddedCertificateEvent {
	return &AddedCertificateEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AddedCertificateEventType,
		),
		Certificate: &Key{
			Key:    certificateCrypto,
			Expiry: certificateExpiration,
		},
	}
}

func AddedCertificateEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &AddedCertificateEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "KEY-4n9vs", "unable to unmarshal certificate added")
	}

	return e, nil
}
