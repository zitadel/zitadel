package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/org"
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
		}
	}
}

func (wm *OrgPasswordLockoutPolicyWriteModel) Reduce() error {
	return wm.PasswordLockoutPolicyWriteModel.Reduce()
}

func (wm *OrgPasswordLockoutPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.PasswordLockoutPolicyWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}

func (wm *OrgPasswordLockoutPolicyWriteModel) NewChangedEvent(ctx context.Context, maxAttempts uint64, showLockoutFailure bool) (*org.PasswordLockoutPolicyChangedEvent, bool) {
	hasChanged := false
	changedEvent := org.NewPasswordLockoutPolicyChangedEvent(ctx)
	if wm.MaxAttempts != maxAttempts {
		hasChanged = true
		changedEvent.MaxAttempts = &maxAttempts
	}
	if wm.ShowLockOutFailures != showLockoutFailure {
		hasChanged = true
		changedEvent.ShowLockOutFailures = &showLockoutFailure
	}
	return changedEvent, hasChanged
}
