package command

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/policy"
)

type InstanceLockoutPolicyWriteModel struct {
	LockoutPolicyWriteModel
}

func NewInstanceLockoutPolicyWriteModel(ctx context.Context) *InstanceLockoutPolicyWriteModel {
	return &InstanceLockoutPolicyWriteModel{
		LockoutPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   authz.GetInstance(ctx).InstanceID(),
				ResourceOwner: authz.GetInstance(ctx).InstanceID(),
			},
		},
	}
}

func (wm *InstanceLockoutPolicyWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.LockoutPolicyAddedEvent:
			wm.LockoutPolicyWriteModel.AppendEvents(&e.LockoutPolicyAddedEvent)
		case *instance.LockoutPolicyChangedEvent:
			wm.LockoutPolicyWriteModel.AppendEvents(&e.LockoutPolicyChangedEvent)
		}
	}
}

func (wm *InstanceLockoutPolicyWriteModel) Reduce() error {
	return wm.LockoutPolicyWriteModel.Reduce()
}

func (wm *InstanceLockoutPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.LockoutPolicyWriteModel.AggregateID).
		EventTypes(
			instance.LockoutPolicyAddedEventType,
			instance.LockoutPolicyChangedEventType).
		Builder()
}

func (wm *InstanceLockoutPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	maxAttempts uint64,
	showLockoutFailure bool) (*instance.LockoutPolicyChangedEvent, bool) {
	changes := make([]policy.LockoutPolicyChanges, 0)
	if wm.MaxPasswordAttempts != maxAttempts {
		changes = append(changes, policy.ChangeMaxAttempts(maxAttempts))
	}
	if wm.ShowLockOutFailures != showLockoutFailure {
		changes = append(changes, policy.ChangeShowLockOutFailures(showLockoutFailure))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := instance.NewLockoutPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
