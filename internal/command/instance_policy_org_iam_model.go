package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/policy"
)

type InstanceOrgIAMPolicyWriteModel struct {
	PolicyOrgIAMWriteModel
}

func NewInstanceOrgIAMPolicyWriteModel() *InstanceOrgIAMPolicyWriteModel {
	return &InstanceOrgIAMPolicyWriteModel{
		PolicyOrgIAMWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
		},
	}
}

func (wm *InstanceOrgIAMPolicyWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.OrgIAMPolicyAddedEvent:
			wm.PolicyOrgIAMWriteModel.AppendEvents(&e.OrgIAMPolicyAddedEvent)
		case *instance.OrgIAMPolicyChangedEvent:
			wm.PolicyOrgIAMWriteModel.AppendEvents(&e.OrgIAMPolicyChangedEvent)
		}
	}
}

func (wm *InstanceOrgIAMPolicyWriteModel) Reduce() error {
	return wm.PolicyOrgIAMWriteModel.Reduce()
}

func (wm *InstanceOrgIAMPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.PolicyOrgIAMWriteModel.AggregateID).
		EventTypes(
			instance.OrgIAMPolicyAddedEventType,
			instance.OrgIAMPolicyChangedEventType).
		Builder()
}

func (wm *InstanceOrgIAMPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userLoginMustBeDomain bool) (*instance.OrgIAMPolicyChangedEvent, bool) {
	changes := make([]policy.OrgIAMPolicyChanges, 0)
	if wm.UserLoginMustBeDomain != userLoginMustBeDomain {
		changes = append(changes, policy.ChangeUserLoginMustBeDomain(userLoginMustBeDomain))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := instance.NewOrgIAMPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
