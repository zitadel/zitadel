package command

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
)

type InstanceCustomLoginTextReadModel struct {
	CustomLoginTextReadModel
}

func NewInstanceCustomLoginTextReadModel(lang language.Tag) *InstanceCustomLoginTextReadModel {
	return &InstanceCustomLoginTextReadModel{
		CustomLoginTextReadModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
			Language: lang,
		},
	}
}

func (wm *InstanceCustomLoginTextReadModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.CustomTextSetEvent:
			wm.CustomLoginTextReadModel.AppendEvents(&e.CustomTextSetEvent)
		case *instance.CustomTextRemovedEvent:
			wm.CustomLoginTextReadModel.AppendEvents(&e.CustomTextRemovedEvent)
		case *instance.CustomTextTemplateRemovedEvent:
			wm.CustomLoginTextReadModel.AppendEvents(&e.CustomTextTemplateRemovedEvent)
		}
	}
}

func (wm *InstanceCustomLoginTextReadModel) Reduce() error {
	return wm.CustomLoginTextReadModel.Reduce()
}

func (wm *InstanceCustomLoginTextReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.CustomLoginTextReadModel.AggregateID).
		AggregateTypes(instance.AggregateType).
		EventTypes(
			instance.CustomTextSetEventType,
			instance.CustomTextRemovedEventType,
			instance.CustomTextTemplateRemovedEventType).
		Builder()
}
