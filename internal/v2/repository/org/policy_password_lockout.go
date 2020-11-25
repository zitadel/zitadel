package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	PasswordLockoutPolicyAddedEventType   = orgEventTypePrefix + policy.PasswordLockoutPolicyAddedEventType
	PasswordLockoutPolicyChangedEventType = orgEventTypePrefix + policy.PasswordLockoutPolicyChangedEventType
)

type PasswordLockoutPolicyReadModel struct {
	policy.PasswordLockoutPolicyReadModel
}

func (rm *PasswordLockoutPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *PasswordLockoutPolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.PasswordLockoutPolicyAddedEvent)
		case *PasswordLockoutPolicyChangedEvent:
			rm.ReadModel.AppendEvents(&e.PasswordLockoutPolicyChangedEvent)
		case *policy.PasswordLockoutPolicyAddedEvent, *policy.PasswordLockoutPolicyChangedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}

type PasswordLockoutPolicyAddedEvent struct {
	policy.PasswordLockoutPolicyAddedEvent
}

type PasswordLockoutPolicyChangedEvent struct {
	policy.PasswordLockoutPolicyChangedEvent
}
