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
		case *org.CustomTextTemplateRemovedEvent:
			wm.CustomMessageTextReadModel.AppendEvents(&e.CustomTextTemplateRemovedEvent)
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
			org.CustomTextTemplateRemovedEventType)
}

type OrgCustomMessageTemplatesReadModel struct {
	CustomMessageTemplatesReadModel
}

func NewOrgCustomMessageTextsWriteModel(orgID string) *OrgCustomMessageTemplatesReadModel {
	return &OrgCustomMessageTemplatesReadModel{
		CustomMessageTemplatesReadModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			CustomMessageTemplate: make(map[string]*CustomText),
		},
	}
}

func (wm *OrgCustomMessageTemplatesReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.CustomTextSetEvent:
			wm.CustomMessageTemplatesReadModel.AppendEvents(&e.CustomTextSetEvent)
		case *org.CustomTextRemovedEvent:
			wm.CustomMessageTemplatesReadModel.AppendEvents(&e.CustomTextRemovedEvent)
		case *org.CustomTextTemplateRemovedEvent:
			wm.CustomMessageTemplatesReadModel.AppendEvents(&e.CustomTextTemplateRemovedEvent)
		}
	}
}

func (wm *OrgCustomMessageTemplatesReadModel) Reduce() error {
	return wm.CustomMessageTemplatesReadModel.Reduce()
}

func (wm *OrgCustomMessageTemplatesReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.CustomMessageTemplatesReadModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			org.CustomTextSetEventType,
			org.CustomTextRemovedEventType,
			org.CustomTextTemplateRemovedEventType)
}
