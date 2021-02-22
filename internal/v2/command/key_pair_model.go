package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	keypair "github.com/caos/zitadel/internal/v2/repository/keypair"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

type KeyPairWriteModel struct {
	eventstore.WriteModel

	Usage      domain.KeyUsage
	Algorithm  string
	PrivateKey *domain.Key
	PublicKey  *domain.Key
}

func NewKeyPairWriteModel(aggregateID, resourceOwner string) *KeyPairWriteModel {
	return &KeyPairWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   aggregateID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *KeyPairWriteModel) AppendEvents(events ...eventstore.EventReader) {
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
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *KeyPairWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, project.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(keypair.AddedEventType)
}

func KeyPairAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, keypair.AggregateType, keypair.AggregateVersion)

}
