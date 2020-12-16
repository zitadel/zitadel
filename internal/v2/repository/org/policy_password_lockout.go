package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/business/query"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	PasswordLockoutPolicyAddedEventType   = orgEventTypePrefix + policy.PasswordLockoutPolicyAddedEventType
	PasswordLockoutPolicyChangedEventType = orgEventTypePrefix + policy.PasswordLockoutPolicyChangedEventType
)

type PasswordLockoutPolicyReadModel struct {
	query.PasswordLockoutPolicyReadModel
}

func (rm *PasswordLockoutPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *PasswordLockoutPolicyAddedEvent:
			rm.PasswordLockoutPolicyReadModel.AppendEvents(&e.PasswordLockoutPolicyAddedEvent)
		case *PasswordLockoutPolicyChangedEvent:
			rm.PasswordLockoutPolicyReadModel.AppendEvents(&e.PasswordLockoutPolicyChangedEvent)
		case *policy.PasswordLockoutPolicyAddedEvent, *policy.PasswordLockoutPolicyChangedEvent:
			rm.PasswordLockoutPolicyReadModel.AppendEvents(e)
		}
	}
}

type PasswordLockoutPolicyAddedEvent struct {
	policy.PasswordLockoutPolicyAddedEvent
}

type PasswordLockoutPolicyChangedEvent struct {
	policy.PasswordLockoutPolicyChangedEvent
}
