package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceCustomLoginTextReadModel struct {
	CustomLoginTextReadModel
}

func NewInstanceCustomLoginTextReadModel(ctx context.Context, lang language.Tag) *InstanceCustomLoginTextReadModel {
	return &InstanceCustomLoginTextReadModel{
		CustomLoginTextReadModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   authz.GetInstance(ctx).InstanceID(),
				ResourceOwner: authz.GetInstance(ctx).InstanceID(),
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
