package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/org"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type OrgPasswordLockoutPolicyWriteModel struct {
	PasswordLockoutPolicyWriteModel
}

func NewOrgPasswordLockoutPolicyWriteModel(orgID string) *OrgPasswordLockoutPolicyWriteModel {
	return &OrgPasswordLockoutPolicyWriteModel{
		PasswordLockoutPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
		},
	}
}

func (wm *OrgPasswordLockoutPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.PasswordLockoutPolicyAddedEvent:
			wm.PasswordLockoutPolicyWriteModel.AppendEvents(&e.PasswordLockoutPolicyAddedEvent)
		case *org.PasswordLockoutPolicyChangedEvent:
			wm.PasswordLockoutPolicyWriteModel.AppendEvents(&e.PasswordLockoutPolicyChangedEvent)
		case *org.PasswordLockoutPolicyRemovedEvent:
			wm.PasswordLockoutPolicyWriteModel.AppendEvents(&e.PasswordLockoutPolicyRemovedEvent)
		}
	}
}

func (wm *OrgPasswordLockoutPolicyWriteModel) Reduce() error {
	return wm.PasswordLockoutPolicyWriteModel.Reduce()
}

func (wm *OrgPasswordLockoutPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.PasswordLockoutPolicyWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(org.PasswordLockoutPolicyAddedEventType,
			org.PasswordLockoutPolicyChangedEventType,
			org.PasswordLockoutPolicyRemovedEventType)
}

func (wm *OrgPasswordLockoutPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	maxAttempts uint64,
	showLockoutFailure bool) (*org.PasswordLockoutPolicyChangedEvent, bool) {
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
	changedEvent, err := org.NewPasswordLockoutPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
