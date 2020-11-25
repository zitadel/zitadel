package iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	PasswordComplexityPolicyAddedEventType   = iamEventTypePrefix + policy.PasswordComplexityPolicyAddedEventType
	PasswordComplexityPolicyChangedEventType = iamEventTypePrefix + policy.PasswordComplexityPolicyChangedEventType
)

type PasswordComplexityPolicyReadModel struct {
	policy.PasswordComplexityPolicyReadModel
}

func (rm *PasswordComplexityPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *PasswordComplexityPolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.PasswordComplexityPolicyAddedEvent)
		case *PasswordComplexityPolicyChangedEvent:
			rm.ReadModel.AppendEvents(&e.PasswordComplexityPolicyChangedEvent)
		case *policy.PasswordComplexityPolicyAddedEvent, *policy.PasswordComplexityPolicyChangedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}

type PasswordComplexityPolicyAddedEvent struct {
	policy.PasswordComplexityPolicyAddedEvent
}

type PasswordComplexityPolicyChangedEvent struct {
	policy.PasswordComplexityPolicyChangedEvent
}
