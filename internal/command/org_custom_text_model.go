package command

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

type OrgCustomTextWriteModel struct {
	CustomTextWriteModel
}

func NewOrgCustomTextWriteModel(orgID, key string, language language.Tag) *OrgCustomTextWriteModel {
	return &OrgCustomTextWriteModel{
		CustomTextWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			Key:      key,
			Language: language,
		},
	}
}

func (wm *OrgCustomTextWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.CustomTextSetEvent:
			wm.CustomTextWriteModel.AppendEvents(&e.CustomTextSetEvent)
		case *org.CustomTextRemovedEvent:
			wm.CustomTextWriteModel.AppendEvents(&e.CustomTextRemovedEvent)
		}
	}
}

func (wm *OrgCustomTextWriteModel) Reduce() error {
	return wm.CustomTextWriteModel.Reduce()
}

func (wm *OrgCustomTextWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.CustomTextWriteModel.AggregateID).
		EventTypes(org.CustomTextSetEventType,
			org.CustomTextRemovedEventType)
	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}
