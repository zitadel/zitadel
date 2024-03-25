package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceCustomMessageTextWriteModel struct {
	CustomMessageTextReadModel
}

func NewInstanceCustomMessageTextWriteModel(ctx context.Context, messageTextType string, lang language.Tag) *InstanceCustomMessageTextWriteModel {
	return &InstanceCustomMessageTextWriteModel{
		CustomMessageTextReadModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   authz.GetInstance(ctx).InstanceID(),
				ResourceOwner: authz.GetInstance(ctx).InstanceID(),
			},
			MessageTextType: messageTextType,
			Language:        lang,
		},
	}
}

func (wm *InstanceCustomMessageTextWriteModel) AppendEvents(events ...eventstore.Event) {
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

func (wm *InstanceCustomMessageTextWriteModel) Reduce() error {
	return wm.CustomMessageTextReadModel.Reduce()
}

func (wm *InstanceCustomMessageTextWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.CustomMessageTextReadModel.AggregateID).
		EventTypes(instance.CustomTextSetEventType, instance.CustomTextRemovedEventType, instance.CustomTextTemplateRemovedEventType).
		Builder()
}
