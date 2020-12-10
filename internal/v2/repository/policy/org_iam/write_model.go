package org_iam

import "github.com/caos/zitadel/internal/eventstore/v2"

type WriteModel struct {
	eventstore.WriteModel

	UserLoginMustBeDomain bool
}

func (wm *WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *AddedEvent:
			wm.UserLoginMustBeDomain = e.UserLoginMustBeDomain
		case *ChangedEvent:
			wm.UserLoginMustBeDomain = e.UserLoginMustBeDomain
		}
	}
	return wm.WriteModel.Reduce()
}
