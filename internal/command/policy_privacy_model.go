package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/policy"
)

type PrivacyPolicyWriteModel struct {
	eventstore.WriteModel

	TOSLink     string
	PrivacyLink string
	State       domain.PolicyState
}

func (wm *PrivacyPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.PrivacyPolicyAddedEvent:
			wm.TOSLink = e.TOSLink
			wm.PrivacyLink = e.PrivacyLink
			wm.State = domain.PolicyStateActive
		case *policy.PrivacyPolicyChangedEvent:
			if e.PrivacyLink != nil {
				wm.PrivacyLink = *e.PrivacyLink
			}
			if e.TOSLink != nil {
				wm.TOSLink = *e.TOSLink
			}
		case *policy.PrivacyPolicyRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
