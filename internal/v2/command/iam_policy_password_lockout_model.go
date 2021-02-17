package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type IAMPasswordLockoutPolicyWriteModel struct {
	PasswordLockoutPolicyWriteModel
}

func NewIAMPasswordLockoutPolicyWriteModel() *IAMPasswordLockoutPolicyWriteModel {
	return &IAMPasswordLockoutPolicyWriteModel{
		PasswordLockoutPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
		},
	}
}

func (wm *IAMPasswordLockoutPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.PasswordLockoutPolicyAddedEvent:
			wm.PasswordLockoutPolicyWriteModel.AppendEvents(&e.PasswordLockoutPolicyAddedEvent)
		case *iam.PasswordLockoutPolicyChangedEvent:
			wm.PasswordLockoutPolicyWriteModel.AppendEvents(&e.PasswordLockoutPolicyChangedEvent)
		}
	}
}

func (wm *IAMPasswordLockoutPolicyWriteModel) Reduce() error {
	return wm.PasswordLockoutPolicyWriteModel.Reduce()
}

func (wm *IAMPasswordLockoutPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.PasswordLockoutPolicyWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}

func (wm *IAMPasswordLockoutPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	maxAttempts uint64,
	showLockoutFailure bool) (*iam.PasswordLockoutPolicyChangedEvent, bool) {
	changes := make([]policy.PasswordLockoutPolicyChanges, 0)
	if wm.MaxAttempts != maxAttempts {
		changes = append(changes, policy.ChangeMaxAttempts(maxAttempts))
	}
	if wm.ShowLockOutFailures != showLockoutFailure {
		changes = append(changes, policy.ChangeShowLockOutFailures(showLockoutFailure))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := iam.NewPasswordLockoutPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
