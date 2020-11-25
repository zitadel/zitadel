package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	LoginPolicyAddedEventType   = orgEventTypePrefix + policy.LoginPolicyAddedEventType
	LoginPolicyChangedEventType = orgEventTypePrefix + policy.LoginPolicyChangedEventType
)

type LoginPolicyReadModel struct{ policy.LoginPolicyReadModel }

func (rm *LoginPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *LoginPolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.LoginPolicyAddedEvent)
		case *LoginPolicyChangedEvent:
			rm.ReadModel.AppendEvents(&e.LoginPolicyChangedEvent)
		case *policy.LoginPolicyAddedEvent, *policy.LoginPolicyChangedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}

type LoginPolicyAddedEvent struct {
	policy.LoginPolicyAddedEvent
}

type LoginPolicyChangedEvent struct {
	policy.LoginPolicyChangedEvent
}
