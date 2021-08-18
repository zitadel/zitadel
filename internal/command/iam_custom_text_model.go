package command

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMCustomTextWriteModel struct {
	CustomTextWriteModel
}

func NewIAMCustomTextWriteModel(key string, language language.Tag) *IAMCustomTextWriteModel {
	return &IAMCustomTextWriteModel{
		CustomTextWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
			Key:      key,
			Language: language,
		},
	}
}

func (wm *IAMCustomTextWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.CustomTextSetEvent:
			wm.CustomTextWriteModel.AppendEvents(&e.CustomTextSetEvent)
		}
	}
}

func (wm *IAMCustomTextWriteModel) Reduce() error {
	return wm.CustomTextWriteModel.Reduce()
}

func (wm *IAMCustomTextWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.CustomTextWriteModel.AggregateID).
		AggregateTypes(iam.AggregateType).
		EventTypes(
			iam.CustomTextSetEventType).
		Builder()
}
