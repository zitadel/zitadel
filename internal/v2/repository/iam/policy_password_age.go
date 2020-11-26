package iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	PasswordAgePolicyAddedEventType   = iamEventTypePrefix + policy.PasswordAgePolicyAddedEventType
	PasswordAgePolicyChangedEventType = iamEventTypePrefix + policy.PasswordAgePolicyChangedEventType
)

type PasswordAgePolicyReadModel struct {
	policy.PasswordAgePolicyReadModel
}

func (rm *PasswordAgePolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *PasswordAgePolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.PasswordAgePolicyAddedEvent)
		case *PasswordAgePolicyChangedEvent:
			rm.ReadModel.AppendEvents(&e.PasswordAgePolicyChangedEvent)
		case *policy.PasswordAgePolicyAddedEvent,
			*policy.PasswordAgePolicyChangedEvent:

			rm.ReadModel.AppendEvents(e)
		}
	}
}

type PasswordAgePolicyAddedEvent struct {
	policy.PasswordAgePolicyAddedEvent
}

type PasswordAgePolicyChangedEvent struct {
	policy.PasswordAgePolicyChangedEvent
}
