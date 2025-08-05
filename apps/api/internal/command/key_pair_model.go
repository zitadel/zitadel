package command

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/keypair"
)

type KeyPairWriteModel struct {
	eventstore.WriteModel

	Usage       crypto.KeyUsage
	Algorithm   string
	PrivateKey  *domain.Key
	PublicKey   *domain.Key
	Certificate *domain.Key
}

func NewKeyPairWriteModel(aggregateID, resourceOwner string) *KeyPairWriteModel {
	return &KeyPairWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   aggregateID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *KeyPairWriteModel) AppendEvents(events ...eventstore.Event) {
	wm.WriteModel.AppendEvents(events...)
}

func (wm *KeyPairWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *keypair.AddedEvent:
			wm.Usage = e.Usage
			wm.Algorithm = e.Algorithm
			wm.PrivateKey = &domain.Key{
				Key:    e.PrivateKey.Key,
				Expiry: e.PrivateKey.Expiry,
			}
			wm.PublicKey = &domain.Key{
				Key:    e.PublicKey.Key,
				Expiry: e.PublicKey.Expiry,
			}
		case *keypair.AddedCertificateEvent:
			wm.Certificate = &domain.Key{
				Key:    e.Certificate.Key,
				Expiry: e.Certificate.Expiry,
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *KeyPairWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(keypair.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(keypair.AddedEventType, keypair.AddedCertificateEventType).
		Builder()
}

func KeyPairAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, keypair.AggregateType, keypair.AggregateVersion)
}
