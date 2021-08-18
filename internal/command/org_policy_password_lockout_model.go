package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
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

func (wm *OrgLockoutPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
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
	maxAttempts uint64,
	showLockoutFailure bool) (*org.LockoutPolicyChangedEvent, bool) {
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
	changedEvent, err := org.NewLockoutPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
