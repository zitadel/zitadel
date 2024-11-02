package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

type InstancePrivacyPolicyWriteModel struct {
	PrivacyPolicyWriteModel
}

func NewInstancePrivacyPolicyWriteModel(ctx context.Context) *InstancePrivacyPolicyWriteModel {
	return &InstancePrivacyPolicyWriteModel{
		PrivacyPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   authz.GetInstance(ctx).InstanceID(),
				ResourceOwner: authz.GetInstance(ctx).InstanceID(),
			},
		},
	}
}

func (wm *InstancePrivacyPolicyWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.PrivacyPolicyAddedEvent:
			wm.PrivacyPolicyWriteModel.AppendEvents(&e.PrivacyPolicyAddedEvent)
		case *instance.PrivacyPolicyChangedEvent:
			wm.PrivacyPolicyWriteModel.AppendEvents(&e.PrivacyPolicyChangedEvent)
		}
	}
}

func (wm *InstancePrivacyPolicyWriteModel) Reduce() error {
	return wm.PrivacyPolicyWriteModel.Reduce()
}

func (wm *InstancePrivacyPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.PrivacyPolicyWriteModel.AggregateID).
		EventTypes(
			instance.PrivacyPolicyAddedEventType,
			instance.PrivacyPolicyChangedEventType).
		Builder()
}

func (wm *InstancePrivacyPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tosLink,
	privacyLink,
	helpLink string,
	supportEmail domain.EmailAddress,
	docsLink, customLink, customLinkText string,
) (*instance.PrivacyPolicyChangedEvent, bool) {

	changes := make([]policy.PrivacyPolicyChanges, 0)
	if wm.TOSLink != tosLink {
		changes = append(changes, policy.ChangeTOSLink(tosLink))
	}
	if wm.PrivacyLink != privacyLink {
		changes = append(changes, policy.ChangePrivacyLink(privacyLink))
	}
	if wm.HelpLink != helpLink {
		changes = append(changes, policy.ChangeHelpLink(helpLink))
	}
	if wm.SupportEmail != supportEmail {
		changes = append(changes, policy.ChangeSupportEmail(supportEmail))
	}
	if wm.DocsLink != docsLink {
		changes = append(changes, policy.ChangeDocsLink(docsLink))
	}
	if wm.CustomLink != customLink {
		changes = append(changes, policy.ChangeCustomLink(customLink))
	}
	if wm.CustomLinkText != customLinkText {
		changes = append(changes, policy.ChangeCustomLinkText(customLinkText))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := instance.NewPrivacyPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
