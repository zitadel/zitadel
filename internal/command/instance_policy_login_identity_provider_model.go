package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
)

type InstanceIdentityProviderWriteModel struct {
	IdentityProviderWriteModel
}

func NewInstanceIdentityProviderWriteModel(idpConfigID string) *InstanceIdentityProviderWriteModel {
	return &InstanceIdentityProviderWriteModel{
		IdentityProviderWriteModel: IdentityProviderWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
			IDPConfigID: idpConfigID,
		},
	}
}

func (wm *InstanceIdentityProviderWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.IdentityProviderAddedEvent:
			if e.IDPConfigID != wm.IDPConfigID {
				continue
			}
			wm.IdentityProviderWriteModel.AppendEvents(&e.IdentityProviderAddedEvent)
		case *instance.IdentityProviderRemovedEvent:
			if e.IDPConfigID != wm.IDPConfigID {
				continue
			}
			wm.IdentityProviderWriteModel.AppendEvents(&e.IdentityProviderRemovedEvent)
		}
	}
}

func (wm *InstanceIdentityProviderWriteModel) Reduce() error {
	return wm.IdentityProviderWriteModel.Reduce()
}

func (wm *InstanceIdentityProviderWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		Builder()
}
