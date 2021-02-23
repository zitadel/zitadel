package command

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type PolicyOrgIAMWriteModel struct {
	eventstore.WriteModel

	UserLoginMustBeDomain bool
	State                 domain.PolicyState
}

func (wm *PolicyOrgIAMWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.OrgIAMPolicyAddedEvent:
			wm.UserLoginMustBeDomain = e.UserLoginMustBeDomain
			wm.State = domain.PolicyStateActive
		case *policy.OrgIAMPolicyChangedEvent:
			if e.UserLoginMustBeDomain != nil {
				wm.UserLoginMustBeDomain = *e.UserLoginMustBeDomain
			}
		}
	}
	return wm.WriteModel.Reduce()
}
