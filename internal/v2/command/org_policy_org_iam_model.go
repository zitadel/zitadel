package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/org"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type ORGOrgIAMPolicyWriteModel struct {
	PolicyOrgIAMWriteModel
}

func NewORGOrgIAMPolicyWriteModel(orgID string) *ORGOrgIAMPolicyWriteModel {
	return &ORGOrgIAMPolicyWriteModel{
		PolicyOrgIAMWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
		},
	}
}

func (wm *ORGOrgIAMPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.OrgIAMPolicyAddedEvent:
			wm.PolicyOrgIAMWriteModel.AppendEvents(&e.OrgIAMPolicyAddedEvent)
		case *org.OrgIAMPolicyChangedEvent:
			wm.PolicyOrgIAMWriteModel.AppendEvents(&e.OrgIAMPolicyChangedEvent)
		case *org.OrgIAMPolicyRemovedEvent:
			wm.PolicyOrgIAMWriteModel.AppendEvents(&e.OrgIAMPolicyRemovedEvent)
		}
	}
}

func (wm *ORGOrgIAMPolicyWriteModel) Reduce() error {
	return wm.PolicyOrgIAMWriteModel.Reduce()
}

func (wm *ORGOrgIAMPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.PolicyOrgIAMWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(org.OrgIAMPolicyAddedEventType,
			org.OrgIAMPolicyChangedEventType,
			org.OrgIAMPolicyRemovedEventType)
}

func (wm *ORGOrgIAMPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userLoginMustBeDomain bool) (*org.OrgIAMPolicyChangedEvent, bool) {
	changes := make([]policy.OrgIAMPolicyChanges, 0)
	if wm.UserLoginMustBeDomain != userLoginMustBeDomain {
		changes = append(changes, policy.ChangeUserLoginMustBeDomain(userLoginMustBeDomain))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := org.NewOrgIAMPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
