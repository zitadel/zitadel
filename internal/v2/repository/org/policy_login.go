package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/business/query"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	LoginPolicyAddedEventType   = orgEventTypePrefix + policy.LoginPolicyAddedEventType
	LoginPolicyChangedEventType = orgEventTypePrefix + policy.LoginPolicyChangedEventType
)

type LoginPolicyReadModel struct{ query.LoginPolicyReadModel }

func (rm *LoginPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *LoginPolicyAddedEvent:
			rm.LoginPolicyReadModel.AppendEvents(&e.LoginPolicyAddedEvent)
		case *LoginPolicyChangedEvent:
			rm.LoginPolicyReadModel.AppendEvents(&e.LoginPolicyChangedEvent)
		case *policy.LoginPolicyAddedEvent, *policy.LoginPolicyChangedEvent:
			rm.LoginPolicyReadModel.AppendEvents(e)
		}
	}
}

type LoginPolicyAddedEvent struct {
	policy.LoginPolicyAddedEvent
}

type LoginPolicyChangedEvent struct {
	policy.LoginPolicyChangedEvent
}
