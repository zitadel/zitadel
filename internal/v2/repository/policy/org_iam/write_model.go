package org_iam

import "github.com/caos/zitadel/internal/eventstore/v2"

type OrgIAMPolicyWriteModel struct {
	eventstore.WriteModel

	UserLoginMustBeDomain bool
}

func (wm *OrgIAMPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *OrgIAMPolicyAddedEvent:
			wm.UserLoginMustBeDomain = e.UserLoginMustBeDomain
		case *OrgIAMPolicyChangedEvent:
			wm.UserLoginMustBeDomain = e.UserLoginMustBeDomain
		}
	}
	return wm.WriteModel.Reduce()
}
