package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

type IAMOrgIAMPolicyWriteModel struct {
	PolicyOrgIAMWriteModel
}

func NewIAMOrgIAMPolicyWriteModel(iamID string) *IAMOrgIAMPolicyWriteModel {
	return &IAMOrgIAMPolicyWriteModel{
		PolicyOrgIAMWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID: iamID,
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
		AggregateIDs(wm.PolicyOrgIAMWriteModel.AggregateID)
}

func (wm *IAMOrgIAMPolicyWriteModel) NewChangedEvent(userLoginMustBeDomain bool) (*iam.OrgIAMPolicyChangedEvent, bool) {
	hasChanged := false
	changedEvent := &iam.OrgIAMPolicyChangedEvent{}
	if wm.UserLoginMustBeDomain != userLoginMustBeDomain {
		hasChanged = true
		changedEvent.UserLoginMustBeDomain = userLoginMustBeDomain
	}
	return changedEvent, hasChanged
}
