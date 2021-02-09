package command

import (
	"context"
	"reflect"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/org"
	"github.com/caos/zitadel/internal/v2/repository/policy"
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
		}
	}
}

func (wm *OrgMailTemplateWriteModel) Reduce() error {
	return wm.MailTemplateWriteModel.Reduce()
}

func (wm *OrgMailTemplateWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.MailTemplateWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}

func (wm *OrgMailTemplateWriteModel) NewChangedEvent(
	ctx context.Context,
	template []byte,
) (*org.MailTemplateChangedEvent, bool) {
	changes := make([]policy.MailTemplateChanges, 0)
	if !reflect.DeepEqual(wm.Template, template) {
		changes = append(changes, policy.ChangeTemplate(template))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := org.NewMailTemplateChangedEvent(ctx, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
