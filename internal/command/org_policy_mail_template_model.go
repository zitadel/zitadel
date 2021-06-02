package command

import (
	"context"
	"reflect"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
)

type OrgMailTemplateWriteModel struct {
	MailTemplateWriteModel
}

func NewOrgMailTemplateWriteModel(orgID string) *OrgMailTemplateWriteModel {
	return &OrgMailTemplateWriteModel{
		MailTemplateWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
		},
	}
}

func (wm *OrgMailTemplateWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.MailTemplateAddedEvent:
			wm.MailTemplateWriteModel.AppendEvents(&e.MailTemplateAddedEvent)
		case *org.MailTemplateChangedEvent:
			wm.MailTemplateWriteModel.AppendEvents(&e.MailTemplateChangedEvent)
		case *org.MailTemplateRemovedEvent:
			wm.MailTemplateWriteModel.AppendEvents(&e.MailTemplateRemovedEvent)
		}
	}
}

func (wm *OrgMailTemplateWriteModel) Reduce() error {
	return wm.MailTemplateWriteModel.Reduce()
}

func (wm *OrgMailTemplateWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.MailTemplateWriteModel.AggregateID).
		EventTypes(
			org.MailTemplateAddedEventType,
			org.MailTemplateChangedEventType,
			org.MailTemplateRemovedEventType).
		SearchQueryBuilder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func (wm *OrgMailTemplateWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	template []byte,
) (*org.MailTemplateChangedEvent, bool) {
	changes := make([]policy.MailTemplateChanges, 0)
	if !reflect.DeepEqual(wm.Template, template) {
		changes = append(changes, policy.ChangeTemplate(template))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := org.NewMailTemplateChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
