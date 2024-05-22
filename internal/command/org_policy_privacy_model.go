package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

type OrgPrivacyPolicyWriteModel struct {
	PrivacyPolicyWriteModel
}

func NewOrgPrivacyPolicyWriteModel(orgID string) *OrgPrivacyPolicyWriteModel {
	return &OrgPrivacyPolicyWriteModel{
		PrivacyPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
		},
	}
}

func (wm *OrgPrivacyPolicyWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.PrivacyPolicyAddedEvent:
			wm.PrivacyPolicyWriteModel.AppendEvents(&e.PrivacyPolicyAddedEvent)
		case *org.PrivacyPolicyChangedEvent:
			wm.PrivacyPolicyWriteModel.AppendEvents(&e.PrivacyPolicyChangedEvent)
		case *org.PrivacyPolicyRemovedEvent:
			wm.PrivacyPolicyWriteModel.AppendEvents(&e.PrivacyPolicyRemovedEvent)
		}
	}
}

func (wm *OrgPrivacyPolicyWriteModel) Reduce() error {
	return wm.PrivacyPolicyWriteModel.Reduce()
}

func (wm *OrgPrivacyPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.PrivacyPolicyWriteModel.AggregateID).
		AggregateTypes(org.AggregateType).
		EventTypes(org.PrivacyPolicyAddedEventType,
			org.PrivacyPolicyChangedEventType,
			org.PrivacyPolicyRemovedEventType).
		Builder()
}

func (wm *OrgPrivacyPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tosLink,
	privacyLink,
	helpLink string,
	supportEmail domain.EmailAddress,
	docsLink, customLink, customLinkText string,
) (*org.PrivacyPolicyChangedEvent, bool) {

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
	changedEvent, err := org.NewPrivacyPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
