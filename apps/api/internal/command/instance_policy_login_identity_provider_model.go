package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceIdentityProviderWriteModel struct {
	IdentityProviderWriteModel
}

func NewInstanceIdentityProviderWriteModel(ctx context.Context, idpConfigID string) *InstanceIdentityProviderWriteModel {
	return &InstanceIdentityProviderWriteModel{
		IdentityProviderWriteModel: IdentityProviderWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   authz.GetInstance(ctx).InstanceID(),
				ResourceOwner: authz.GetInstance(ctx).InstanceID(),
				InstanceID:    authz.GetInstance(ctx).InstanceID(),
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
