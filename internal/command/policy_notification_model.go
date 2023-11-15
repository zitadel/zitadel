package command

import (
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
	"github.com/zitadel/zitadel/v2/internal/repository/policy"
)

type NotificationPolicyWriteModel struct {
	eventstore.WriteModel

	PasswordChange bool
	State          domain.PolicyState
}

func (wm *NotificationPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.NotificationPolicyAddedEvent:
			wm.PasswordChange = e.PasswordChange
			wm.State = domain.PolicyStateActive
		case *policy.NotificationPolicyChangedEvent:
			if e.PasswordChange != nil {
				wm.PasswordChange = *e.PasswordChange
			}
		case *policy.NotificationPolicyRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
