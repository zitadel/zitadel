package command

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

type OrgCustomLoginTextReadModel struct {
	CustomLoginTextReadModel
}

func NewOrgCustomLoginTextReadModel(orgID string, lang language.Tag) *OrgCustomLoginTextReadModel {
	return &OrgCustomLoginTextReadModel{
		CustomLoginTextReadModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			Language: lang,
		},
	}
}

func (wm *OrgCustomLoginTextReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.CustomTextSetEvent:
			wm.CustomLoginTextReadModel.AppendEvents(&e.CustomTextSetEvent)
		case *org.CustomTextRemovedEvent:
			wm.CustomLoginTextReadModel.AppendEvents(&e.CustomTextRemovedEvent)
		case *org.CustomTextTemplateRemovedEvent:
			wm.CustomLoginTextReadModel.AppendEvents(&e.CustomTextTemplateRemovedEvent)
		}
	}
}

func (wm *OrgCustomLoginTextReadModel) Reduce() error {
	return wm.CustomLoginTextReadModel.Reduce()
}

func (wm *OrgCustomLoginTextReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.CustomLoginTextReadModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			org.CustomTextSetEventType,
			org.CustomTextRemovedEventType,
			org.CustomTextTemplateRemovedEventType)
}
