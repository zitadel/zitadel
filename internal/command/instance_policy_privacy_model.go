package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/policy"
)

type InstancePrivacyPolicyWriteModel struct {
	PrivacyPolicyWriteModel
}

func NewInstancePrivacyPolicyWriteModel(instanceID string) *InstancePrivacyPolicyWriteModel {
	return &InstancePrivacyPolicyWriteModel{
		PrivacyPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
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
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := instance.NewPrivacyPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
