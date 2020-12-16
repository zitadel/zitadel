package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type PolicyOrgIAMWriteModel struct {
	eventstore.WriteModel

	UserLoginMustBeDomain bool
	IsActive              bool
}

func (wm *PolicyOrgIAMWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.OrgIAMPolicyAddedEvent:
			wm.UserLoginMustBeDomain = e.UserLoginMustBeDomain
			wm.IsActive = true
		case *policy.OrgIAMPolicyChangedEvent:
			wm.UserLoginMustBeDomain = e.UserLoginMustBeDomain
		}
	}
	return wm.WriteModel.Reduce()
}
