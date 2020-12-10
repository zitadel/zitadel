package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/login"
)

var (
	LoginPolicyAddedEventType   = orgEventTypePrefix + login.LoginPolicyAddedEventType
	LoginPolicyChangedEventType = orgEventTypePrefix + login.LoginPolicyChangedEventType
)

type LoginPolicyReadModel struct{ login.ReadModel }

func (rm *LoginPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *LoginPolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.AddedEvent)
		case *LoginPolicyChangedEvent:
			rm.ReadModel.AppendEvents(&e.ChangedEvent)
		case *login.AddedEvent, *login.ChangedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}

type LoginPolicyAddedEvent struct {
	login.AddedEvent
}

type LoginPolicyChangedEvent struct {
	login.ChangedEvent
}
