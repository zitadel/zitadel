package org_iam

import "github.com/caos/zitadel/internal/eventstore/v2"

type OrgIAMPolicyReadModel struct {
	eventstore.ReadModel

	UserLoginMustBeDomain bool
}

func (rm *OrgIAMPolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *OrgIAMPolicyAddedEvent:
			rm.UserLoginMustBeDomain = e.UserLoginMustBeDomain
		case *OrgIAMPolicyChangedEvent:
			rm.UserLoginMustBeDomain = e.UserLoginMustBeDomain
		}
	}
	return rm.ReadModel.Reduce()
}
