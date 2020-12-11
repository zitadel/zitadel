package org_iam

import "github.com/caos/zitadel/internal/eventstore/v2"

type ReadModel struct {
	eventstore.ReadModel

	UserLoginMustBeDomain bool
}

func (rm *ReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *AddedEvent:
			rm.UserLoginMustBeDomain = e.UserLoginMustBeDomain
		case *ChangedEvent:
			rm.UserLoginMustBeDomain = e.UserLoginMustBeDomain
		}
	}
	return rm.ReadModel.Reduce()
}
