package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
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

func (wm *OrgPrivacyPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
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
	privacyLink string,
) (*org.PrivacyPolicyChangedEvent, bool) {

	changes := make([]policy.PrivacyPolicyChanges, 0)
	if wm.TOSLink != tosLink {
		changes = append(changes, policy.ChangeTOSLink(tosLink))
	}
	if wm.PrivacyLink != privacyLink {
		changes = append(changes, policy.ChangePrivacyLink(privacyLink))
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
