package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_lockout"
)

var (
	PasswordLockoutPolicyAddedEventType   = orgEventTypePrefix + password_lockout.PasswordLockoutPolicyAddedEventType
	PasswordLockoutPolicyChangedEventType = orgEventTypePrefix + password_lockout.PasswordLockoutPolicyChangedEventType
)

type PasswordLockoutPolicyReadModel struct {
	password_lockout.ReadModel
}

func (rm *PasswordLockoutPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *PasswordLockoutPolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.AddedEvent)
		case *PasswordLockoutPolicyChangedEvent:
			rm.ReadModel.AppendEvents(&e.ChangedEvent)
		case *password_lockout.AddedEvent, *password_lockout.ChangedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}

type PasswordLockoutPolicyAddedEvent struct {
	password_lockout.AddedEvent
}

type PasswordLockoutPolicyChangedEvent struct {
	password_lockout.ChangedEvent
}
