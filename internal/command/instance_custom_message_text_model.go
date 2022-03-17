package command

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
)

type InstanceCustomMessageTextReadModel struct {
	CustomMessageTextReadModel
}

func NewInstanceCustomMessageTextWriteModel(messageTextType string, lang language.Tag) *InstanceCustomMessageTextReadModel {
	return &InstanceCustomMessageTextReadModel{
		CustomMessageTextReadModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
			MessageTextType: messageTextType,
			Language:        lang,
		},
	}
}

func (wm *InstanceCustomMessageTextReadModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.CustomTextSetEvent:
			wm.CustomMessageTextReadModel.AppendEvents(&e.CustomTextSetEvent)
		case *instance.CustomTextRemovedEvent:
			wm.CustomMessageTextReadModel.AppendEvents(&e.CustomTextRemovedEvent)
		case *instance.CustomTextTemplateRemovedEvent:
			wm.CustomMessageTextReadModel.AppendEvents(&e.CustomTextTemplateRemovedEvent)
		}
	}
}

func (wm *InstanceCustomMessageTextReadModel) Reduce() error {
	return wm.CustomMessageTextReadModel.Reduce()
}

func (wm *InstanceCustomMessageTextReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.CustomMessageTextReadModel.AggregateID).
		EventTypes(instance.CustomTextSetEventType, instance.CustomTextRemovedEventType, instance.CustomTextTemplateRemovedEventType).
		Builder()
}
