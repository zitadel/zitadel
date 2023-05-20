package keypair

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
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

func (e *AddedCertificateEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func AddedCertificateEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &AddedCertificateEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "KEY-4n9vs", "unable to unmarshal certificate added")
	}

	return e, nil
}
