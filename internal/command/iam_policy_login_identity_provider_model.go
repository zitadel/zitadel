package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMIdentityProviderWriteModel struct {
	IdentityProviderWriteModel
}

func NewIAMIdentityProviderWriteModel(idpConfigID string) *IAMIdentityProviderWriteModel {
	return &IAMIdentityProviderWriteModel{
		IdentityProviderWriteModel: IdentityProviderWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
			IDPConfigID: idpConfigID,
		},
	}
}

func (wm *IAMIdentityProviderWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.IdentityProviderAddedEvent:
			if e.IDPConfigID != wm.IDPConfigID {
				continue
			}
			wm.IdentityProviderWriteModel.AppendEvents(&e.IdentityProviderAddedEvent)
		case *iam.IdentityProviderRemovedEvent:
			if e.IDPConfigID != wm.IDPConfigID {
				continue
			}
			wm.IdentityProviderWriteModel.AppendEvents(&e.IdentityProviderRemovedEvent)
		}
	}
}

func (wm *IAMIdentityProviderWriteModel) Reduce() error {
	return wm.IdentityProviderWriteModel.Reduce()
}

func (wm *IAMIdentityProviderWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(wm.AggregateID).
		SearchQueryBuilder()
}
