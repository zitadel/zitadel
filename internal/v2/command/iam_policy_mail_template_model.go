package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
	"reflect"

	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type IAMMailTemplateWriteModel struct {
	MailTemplateWriteModel
}

func NewIAMMailTemplateWriteModel() *IAMMailTemplateWriteModel {
	return &IAMMailTemplateWriteModel{
		MailTemplateWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
		},
	}
}

func (wm *IAMMailTemplateWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.MailTemplateAddedEvent:
			wm.MailTemplateWriteModel.AppendEvents(&e.MailTemplateAddedEvent)
		case *iam.MailTemplateChangedEvent:
			wm.MailTemplateWriteModel.AppendEvents(&e.MailTemplateChangedEvent)
		}
	}
}

func (wm *IAMMailTemplateWriteModel) Reduce() error {
	return wm.MailTemplateWriteModel.Reduce()
}

func (wm *IAMMailTemplateWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.MailTemplateWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			iam.MailTemplateAddedEventType,
			iam.MailTemplateChangedEventType)
}

func (wm *IAMMailTemplateWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	template []byte,
) (*iam.MailTemplateChangedEvent, bool) {
	changes := make([]policy.MailTemplateChanges, 0)
	if !reflect.DeepEqual(wm.Template, template) {
		changes = append(changes, policy.ChangeTemplate(template))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := iam.NewMailTemplateChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
