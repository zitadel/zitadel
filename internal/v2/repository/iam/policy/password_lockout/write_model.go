package password_lockout

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_lockout"
)

const (
	AggregateType = "iam"
)

type PasswordLockoutPolicyWriteModel struct {
	eventstore.WriteModel
	Policy password_lockout.PasswordLockoutPolicyWriteModel

	iamID string
}

func NewPasswordLockoutPolicyWriteModel(iamID string) *PasswordLockoutPolicyWriteModel {
	return &PasswordLockoutPolicyWriteModel{
		iamID: iamID,
	}
}

func (wm *PasswordLockoutPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *PasswordLockoutPolicyAddedEvent:
			wm.Policy.AppendEvents(&e.PasswordLockoutPolicyAddedEvent)
		case *PasswordLockoutPolicyChangedEvent:
			wm.Policy.AppendEvents(&e.PasswordLockoutPolicyChangedEvent)
		}
	}
}

func (wm *PasswordLockoutPolicyWriteModel) Reduce() error {
	if err := wm.Policy.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *PasswordLockoutPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.iamID)
}
