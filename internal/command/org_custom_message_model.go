package command

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

type OrgCustomMessageTextReadModel struct {
	CustomMessageTextReadModel
}

func NewOrgCustomMessageTextWriteModel(orgID, messageTextType string, lang language.Tag) *OrgCustomMessageTextReadModel {
	return &OrgCustomMessageTextReadModel{
		CustomMessageTextReadModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			MessageTextType: messageTextType,
			Language:        lang,
		},
	}
}

func (wm *OrgCustomMessageTextReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.CustomTextSetEvent:
			wm.CustomMessageTextReadModel.AppendEvents(&e.CustomTextSetEvent)
		case *org.CustomTextRemovedEvent:
			wm.CustomMessageTextReadModel.AppendEvents(&e.CustomTextRemovedEvent)
		case *org.CustomTextMessageRemovedEvent:
			wm.CustomMessageTextReadModel.AppendEvents(&e.CustomTextMessageRemovedEvent)
		}
	}
}

func (wm *OrgCustomMessageTextReadModel) Reduce() error {
	return wm.CustomMessageTextReadModel.Reduce()
}

func (wm *OrgCustomMessageTextReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.CustomMessageTextReadModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			org.CustomTextSetEventType,
			org.CustomTextRemovedEventType,
			org.CustomTextMessageRemovedEventType)
}
