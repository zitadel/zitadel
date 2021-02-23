package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type OrgIAMPolicyReadModel struct {
	eventstore.ReadModel

	UserLoginMustBeDomain bool
}

func (rm *OrgIAMPolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *policy.OrgIAMPolicyAddedEvent:
			rm.UserLoginMustBeDomain = e.UserLoginMustBeDomain
		case *policy.OrgIAMPolicyChangedEvent:
			if e.UserLoginMustBeDomain != nil {
				rm.UserLoginMustBeDomain = *e.UserLoginMustBeDomain
			}
		}
	}
	return rm.ReadModel.Reduce()
}
