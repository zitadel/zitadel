package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_complexity"
)

var (
	PasswordComplexityPolicyAddedEventType   = orgEventTypePrefix + password_complexity.PasswordComplexityPolicyAddedEventType
	PasswordComplexityPolicyChangedEventType = orgEventTypePrefix + password_complexity.PasswordComplexityPolicyChangedEventType
)

type PasswordComplexityPolicyReadModel struct {
	password_complexity.ReadModel
}

func (rm *PasswordComplexityPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *PasswordComplexityPolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.AddedEvent)
		case *PasswordComplexityPolicyChangedEvent:
			rm.ReadModel.AppendEvents(&e.ChangedEvent)
		case *password_complexity.AddedEvent, *password_complexity.ChangedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}

type PasswordComplexityPolicyAddedEvent struct {
	password_complexity.AddedEvent
}

type PasswordComplexityPolicyChangedEvent struct {
	password_complexity.ChangedEvent
}
