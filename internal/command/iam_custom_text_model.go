package command

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/iam"
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

func (wm *IAMCustomTextWriteModel) AppendEvents(events ...eventstore.Event) {
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
