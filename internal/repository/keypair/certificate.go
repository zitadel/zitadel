package keypair

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	AddedCertificateEventType = eventTypePrefix + "certificate.added"
)

type AddedCertificateEvent struct {
	eventstore.BaseEvent `json:"-"`

	Usage       domain.KeyUsage `json:"usage"`
	Algorithm   string          `json:"algorithm"`
	Certificate *Key            `json:"certificate"`
}

func (e *AddedCertificateEvent) Data() interface{} {
	return e
}

func (e *AddedCertificateEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewAddedCertificateEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	usage domain.KeyUsage,
	algorithm string,
	certificateCrypto *crypto.CryptoValue,
	certificateExpiration time.Time) *AddedCertificateEvent {
	return &AddedCertificateEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AddedCertificateEventType,
		),
		Usage:     usage,
		Algorithm: algorithm,
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
