package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type IAMOrgIAMPolicyWriteModel struct {
	PolicyOrgIAMWriteModel
}

func NewIAMOrgIAMPolicyWriteModel() *IAMOrgIAMPolicyWriteModel {
	return &IAMOrgIAMPolicyWriteModel{
		PolicyOrgIAMWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
		},
	}
}

func (wm *IAMOrgIAMPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.OrgIAMPolicyAddedEvent:
			wm.PolicyOrgIAMWriteModel.AppendEvents(&e.OrgIAMPolicyAddedEvent)
		case *iam.OrgIAMPolicyChangedEvent:
			wm.PolicyOrgIAMWriteModel.AppendEvents(&e.OrgIAMPolicyChangedEvent)
		}
	}
}

func (wm *IAMOrgIAMPolicyWriteModel) Reduce() error {
	return wm.PolicyOrgIAMWriteModel.Reduce()
}

func (wm *IAMOrgIAMPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.PolicyOrgIAMWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			iam.OrgIAMPolicyAddedEventType,
			iam.OrgIAMPolicyChangedEventType)
}

func (wm *IAMOrgIAMPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userLoginMustBeDomain bool) (*iam.OrgIAMPolicyChangedEvent, bool) {
	changes := make([]policy.OrgIAMPolicyChanges, 0)
	if wm.UserLoginMustBeDomain != userLoginMustBeDomain {
		changes = append(changes, policy.ChangeUserLoginMustBeDomain(userLoginMustBeDomain))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := iam.NewOrgIAMPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
