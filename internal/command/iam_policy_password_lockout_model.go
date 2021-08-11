package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/policy"
)

type IAMLockoutPolicyWriteModel struct {
	LockoutPolicyWriteModel
}

func NewIAMLockoutPolicyWriteModel() *IAMLockoutPolicyWriteModel {
	return &IAMLockoutPolicyWriteModel{
		LockoutPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
		},
	}
}

func (wm *IAMLockoutPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.LockoutPolicyAddedEvent:
			wm.LockoutPolicyWriteModel.AppendEvents(&e.LockoutPolicyAddedEvent)
		case *iam.LockoutPolicyChangedEvent:
			wm.LockoutPolicyWriteModel.AppendEvents(&e.LockoutPolicyChangedEvent)
		}
	}
}

func (wm *IAMLockoutPolicyWriteModel) Reduce() error {
	return wm.LockoutPolicyWriteModel.Reduce()
}

func (wm *IAMLockoutPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(wm.LockoutPolicyWriteModel.AggregateID).
		EventTypes(
			iam.LockoutPolicyAddedEventType,
			iam.LockoutPolicyChangedEventType).
		Builder()
}

func (wm *IAMLockoutPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	maxAttempts uint64,
	showLockoutFailure bool) (*iam.LockoutPolicyChangedEvent, bool) {
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
	changedEvent, err := iam.NewLockoutPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
