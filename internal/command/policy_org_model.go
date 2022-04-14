package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/policy"
)

type PolicyDomainWriteModel struct {
	eventstore.WriteModel

	UserLoginMustBeDomain bool
	ValidateOrgDomains    bool
	State                 domain.PolicyState
}

func (wm *PolicyDomainWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.DomainPolicyAddedEvent:
			wm.UserLoginMustBeDomain = e.UserLoginMustBeDomain
			wm.ValidateOrgDomains = e.ValidateOrgDomains
			wm.State = domain.PolicyStateActive
		case *policy.DomainPolicyChangedEvent:
			if e.UserLoginMustBeDomain != nil {
				wm.UserLoginMustBeDomain = *e.UserLoginMustBeDomain
			}
			if e.ValidateOrgDomains != nil {
				wm.ValidateOrgDomains = *e.ValidateOrgDomains
			}
		}
	}
	return wm.WriteModel.Reduce()
}
