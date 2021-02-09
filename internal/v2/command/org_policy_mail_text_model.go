package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/org"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type OrgMailTextWriteModel struct {
	MailTextWriteModel
}

func NewOrgMailTextWriteModel(orgID string) *OrgMailTextWriteModel {
	return &OrgMailTextWriteModel{
		MailTextWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
		},
	}
}

func (wm *OrgMailTextWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.MailTextAddedEvent:
			wm.MailTextWriteModel.AppendEvents(&e.MailTextAddedEvent)
		case *org.MailTextChangedEvent:
			wm.MailTextWriteModel.AppendEvents(&e.MailTextChangedEvent)
		}
	}
}

func (wm *OrgMailTextWriteModel) Reduce() error {
	return wm.MailTextWriteModel.Reduce()
}

func (wm *OrgMailTextWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.MailTextWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}

func (wm *OrgMailTextWriteModel) NewChangedEvent(
	ctx context.Context,
	mailTextType,
	language,
	title,
	preHeader,
	subject,
	greeting,
	text,
	buttonText string,
) (*org.MailTextChangedEvent, bool) {
	changes := make([]policy.MailTextChanges, 0)
	if wm.Title != title {
		changes = append(changes, policy.ChangeTitle(title))
	}
	if wm.PreHeader != preHeader {
		changes = append(changes, policy.ChangePreHeader(preHeader))
	}
	if wm.Subject != subject {
		changes = append(changes, policy.ChangeSubject(subject))
	}
	if wm.Greeting != greeting {
		changes = append(changes, policy.ChangeGreeting(greeting))
	}
	if wm.Text != text {
		changes = append(changes, policy.ChangeText(text))
	}
	if wm.ButtonText != buttonText {
		changes = append(changes, policy.ChangeButtonText(buttonText))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := org.NewMailTextChangedEvent(ctx, mailTextType, language, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
