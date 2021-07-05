package command

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMCustomLoginTextReadModel struct {
	CustomLoginTextReadModel
}

func NewIAMCustomLoginTextReadModel(lang language.Tag) *IAMCustomLoginTextReadModel {
	return &IAMCustomLoginTextReadModel{
		CustomLoginTextReadModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
			Language: lang,
		},
	}
}

func (wm *IAMCustomLoginTextReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.CustomTextSetEvent:
			wm.CustomLoginTextReadModel.AppendEvents(&e.CustomTextSetEvent)
		case *iam.CustomTextRemovedEvent:
			wm.CustomLoginTextReadModel.AppendEvents(&e.CustomTextRemovedEvent)
		}
	}
}

func (wm *IAMCustomLoginTextReadModel) Reduce() error {
	return wm.CustomLoginTextReadModel.Reduce()
}

func (wm *IAMCustomLoginTextReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.CustomLoginTextReadModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			iam.CustomTextSetEventType,
			iam.CustomTextRemovedEventType)
}
