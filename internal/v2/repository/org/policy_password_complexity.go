package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/business/query"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	PasswordComplexityPolicyAddedEventType   = orgEventTypePrefix + policy.PasswordComplexityPolicyAddedEventType
	PasswordComplexityPolicyChangedEventType = orgEventTypePrefix + policy.PasswordComplexityPolicyChangedEventType
)

type PasswordComplexityPolicyReadModel struct {
	query.PasswordComplexityPolicyReadModel
}

func (rm *PasswordComplexityPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *PasswordComplexityPolicyAddedEvent:
			rm.PasswordComplexityPolicyReadModel.AppendEvents(&e.PasswordComplexityPolicyAddedEvent)
		case *PasswordComplexityPolicyChangedEvent:
			rm.PasswordComplexityPolicyReadModel.AppendEvents(&e.PasswordComplexityPolicyChangedEvent)
		case *policy.PasswordComplexityPolicyAddedEvent, *policy.PasswordComplexityPolicyChangedEvent:
			rm.PasswordComplexityPolicyReadModel.AppendEvents(e)
		}
	}
}

type PasswordComplexityPolicyAddedEvent struct {
	policy.PasswordComplexityPolicyAddedEvent
}

type PasswordComplexityPolicyChangedEvent struct {
	policy.PasswordComplexityPolicyChangedEvent
}
