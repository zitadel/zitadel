package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

type IAMPasswordLockoutPolicyWriteModel struct {
	PasswordLockoutPolicyWriteModel
}

func NewIAMPasswordLockoutPolicyWriteModel(iamID string) *IAMPasswordLockoutPolicyWriteModel {
	return &IAMPasswordLockoutPolicyWriteModel{
		PasswordLockoutPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID: iamID,
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
		AggregateIDs(wm.PasswordLockoutPolicyWriteModel.AggregateID)
}

func (wm *IAMPasswordLockoutPolicyWriteModel) NewChangedEvent(maxAttempts uint64, showLockoutFailure bool) (*iam.PasswordLockoutPolicyChangedEvent, bool) {
	hasChanged := false
	changedEvent := &iam.PasswordLockoutPolicyChangedEvent{}
	if wm.MaxAttempts == maxAttempts {
		hasChanged = true
		changedEvent.MaxAttempts = maxAttempts
	}
	if wm.ShowLockOutFailures == showLockoutFailure {
		hasChanged = true
		changedEvent.ShowLockOutFailures = showLockoutFailure
	}
	return changedEvent, hasChanged
}
