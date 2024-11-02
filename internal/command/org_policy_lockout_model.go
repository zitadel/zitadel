package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

type OrgLockoutPolicyWriteModel struct {
	LockoutPolicyWriteModel
}

func NewOrgLockoutPolicyWriteModel(orgID string) *OrgLockoutPolicyWriteModel {
	return &OrgLockoutPolicyWriteModel{
		LockoutPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
		},
	}
}

func (wm *OrgLockoutPolicyWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.LockoutPolicyAddedEvent:
			wm.LockoutPolicyWriteModel.AppendEvents(&e.LockoutPolicyAddedEvent)
		case *org.LockoutPolicyChangedEvent:
			wm.LockoutPolicyWriteModel.AppendEvents(&e.LockoutPolicyChangedEvent)
		case *org.LockoutPolicyRemovedEvent:
			wm.LockoutPolicyWriteModel.AppendEvents(&e.LockoutPolicyRemovedEvent)
		}
	}
}

func (wm *OrgLockoutPolicyWriteModel) Reduce() error {
	return wm.LockoutPolicyWriteModel.Reduce()
}

func (wm *OrgLockoutPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.LockoutPolicyWriteModel.AggregateID).
		EventTypes(org.LockoutPolicyAddedEventType,
			org.LockoutPolicyChangedEventType,
			org.LockoutPolicyRemovedEventType).
		Builder()
}

func (wm *OrgLockoutPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	maxPasswordAttempts,
	maxOTPAttempts uint64,
	showLockoutFailure bool) (*org.LockoutPolicyChangedEvent, bool) {
	changes := make([]policy.LockoutPolicyChanges, 0)
	if wm.MaxPasswordAttempts != maxPasswordAttempts {
		changes = append(changes, policy.ChangeMaxPasswordAttempts(maxPasswordAttempts))
	}
	if wm.MaxOTPAttempts != maxOTPAttempts {
		changes = append(changes, policy.ChangeMaxOTPAttempts(maxOTPAttempts))
	}
	if wm.ShowLockOutFailures != showLockoutFailure {
		changes = append(changes, policy.ChangeShowLockOutFailures(showLockoutFailure))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := org.NewLockoutPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
