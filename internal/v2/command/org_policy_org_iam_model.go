package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

type ORGOrgIAMPolicyWriteModel struct {
	PolicyOrgIAMWriteModel
}

func NewORGOrgIAMPolicyWriteModel(orgID string) *ORGOrgIAMPolicyWriteModel {
	return &ORGOrgIAMPolicyWriteModel{
		PolicyOrgIAMWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID: orgID,
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
		}
	}
}

func (wm *ORGOrgIAMPolicyWriteModel) Reduce() error {
	return wm.PolicyOrgIAMWriteModel.Reduce()
}

func (wm *ORGOrgIAMPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.PolicyOrgIAMWriteModel.AggregateID)
}

func (wm *ORGOrgIAMPolicyWriteModel) NewChangedEvent(ctx context.Context, userLoginMustBeDomain bool) (*org.OrgIAMPolicyChangedEvent, bool) {
	hasChanged := false
	changedEvent := org.NewOrgIAMPolicyChangedEvent(ctx)
	if wm.UserLoginMustBeDomain != userLoginMustBeDomain {
		hasChanged = true
		changedEvent.UserLoginMustBeDomain = userLoginMustBeDomain
	}
	return changedEvent, hasChanged
}
