package command

import (
	"context"
	"reflect"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"

	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

type InstanceMailTemplateWriteModel struct {
	MailTemplateWriteModel
}

func NewInstanceMailTemplateWriteModel(ctx context.Context) *InstanceMailTemplateWriteModel {
	return &InstanceMailTemplateWriteModel{
		MailTemplateWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   authz.GetInstance(ctx).InstanceID(),
				ResourceOwner: authz.GetInstance(ctx).InstanceID(),
				InstanceID:    authz.GetInstance(ctx).InstanceID(),
			},
		},
	}
}

func (wm *InstanceMailTemplateWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.MailTemplateAddedEvent:
			wm.MailTemplateWriteModel.AppendEvents(&e.MailTemplateAddedEvent)
		case *instance.MailTemplateChangedEvent:
			wm.MailTemplateWriteModel.AppendEvents(&e.MailTemplateChangedEvent)
		}
	}
}

func (wm *InstanceMailTemplateWriteModel) Reduce() error {
	return wm.MailTemplateWriteModel.Reduce()
}

func (wm *InstanceMailTemplateWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.MailTemplateWriteModel.AggregateID).
		EventTypes(
			instance.MailTemplateAddedEventType,
			instance.MailTemplateChangedEventType).
		Builder()
}

func (wm *InstanceMailTemplateWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	template []byte,
) (*instance.MailTemplateChangedEvent, bool) {
	changes := make([]policy.MailTemplateChanges, 0)
	if !reflect.DeepEqual(wm.Template, template) {
		changes = append(changes, policy.ChangeTemplate(template))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := instance.NewMailTemplateChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
