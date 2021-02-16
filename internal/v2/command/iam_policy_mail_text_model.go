package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type IAMMailTextWriteModel struct {
	MailTextWriteModel
}

func NewIAMMailTextWriteModel(mailTextType, language string) *IAMMailTextWriteModel {
	return &IAMMailTextWriteModel{
		MailTextWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
			MailTextType: mailTextType,
			Language:     language,
		},
	}
}

func (wm *IAMMailTextWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.MailTextAddedEvent:
			wm.MailTextWriteModel.AppendEvents(&e.MailTextAddedEvent)
		case *iam.MailTextChangedEvent:
			wm.MailTextWriteModel.AppendEvents(&e.MailTextChangedEvent)
		}
	}
}

func (wm *IAMMailTextWriteModel) Reduce() error {
	return wm.MailTextWriteModel.Reduce()
}

func (wm *IAMMailTextWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.MailTextWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}

func (wm *IAMMailTextWriteModel) NewChangedEvent(
	ctx context.Context,
	mailTextType,
	language,
	title,
	preHeader,
	subject,
	greeting,
	text,
	buttonText string,
) (*iam.MailTextChangedEvent, bool) {
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
	changedEvent, err := iam.NewMailTextChangedEvent(ctx, mailTextType, language, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
