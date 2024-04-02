package command

import (
	"context"
	"slices"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceSecurityPolicyWriteModel struct {
	eventstore.WriteModel
	SecurityPolicy
}

func NewInstanceSecurityPolicyWriteModel(ctx context.Context) *InstanceSecurityPolicyWriteModel {
	return &InstanceSecurityPolicyWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   authz.GetInstance(ctx).InstanceID(),
			ResourceOwner: authz.GetInstance(ctx).InstanceID(),
		},
	}
}

func (wm *InstanceSecurityPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		if e, ok := event.(*instance.SecurityPolicySetEvent); ok {

			if e.EnableIframeEmbedding != nil {
				wm.EnableIframeEmbedding = *e.EnableIframeEmbedding
			} else if e.Enabled != nil {
				wm.EnableIframeEmbedding = *e.Enabled
			}
			if e.AllowedOrigins != nil {
				wm.AllowedOrigins = *e.AllowedOrigins
			}
			if e.EnableImpersonation != nil {
				wm.EnableImpersonation = *e.EnableImpersonation
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *InstanceSecurityPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.SecurityPolicySetEventType).
		Builder()
}

func (wm *InstanceSecurityPolicyWriteModel) NewSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	policy *SecurityPolicy,
) (*instance.SecurityPolicySetEvent, error) {
	changes := make([]instance.SecurityPolicyChanges, 0, 2)
	var err error

	if wm.EnableIframeEmbedding != policy.EnableIframeEmbedding {
		changes = append(changes, instance.ChangeSecurityPolicyEnableIframeEmbedding(policy.EnableIframeEmbedding))
	}
	if !slices.Equal(wm.AllowedOrigins, policy.AllowedOrigins) {
		changes = append(changes, instance.ChangeSecurityPolicyAllowedOrigins(policy.AllowedOrigins))
	}
	if wm.EnableImpersonation != policy.EnableImpersonation {
		changes = append(changes, instance.ChangeSecurityPolicyEnableImpersonation(policy.EnableImpersonation))
	}
	changeEvent, err := instance.NewSecurityPolicySetEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}
