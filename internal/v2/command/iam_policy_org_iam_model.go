package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
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
		ResourceOwner(wm.ResourceOwner)
}

func (wm *IAMOrgIAMPolicyWriteModel) NewChangedEvent(ctx context.Context, userLoginMustBeDomain bool) (*iam.OrgIAMPolicyChangedEvent, bool) {
	hasChanged := false
	changedEvent := iam.NewOrgIAMPolicyChangedEvent(ctx)
	if wm.UserLoginMustBeDomain != userLoginMustBeDomain {
		hasChanged = true
		changedEvent.UserLoginMustBeDomain = &userLoginMustBeDomain
	}
	return changedEvent, hasChanged
}
