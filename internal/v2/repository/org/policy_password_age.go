package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/business/query"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	PasswordAgePolicyAddedEventType   = orgEventTypePrefix + policy.PasswordAgePolicyAddedEventType
	PasswordAgePolicyChangedEventType = orgEventTypePrefix + policy.PasswordAgePolicyChangedEventType
)

type PasswordAgePolicyReadModel struct {
	query.PasswordAgePolicyReadModel
}

func (rm *PasswordAgePolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *PasswordAgePolicyAddedEvent:
			rm.PasswordAgePolicyReadModel.AppendEvents(&e.PasswordAgePolicyAddedEvent)
		case *PasswordAgePolicyChangedEvent:
			rm.PasswordAgePolicyReadModel.AppendEvents(&e.PasswordAgePolicyChangedEvent)
		case *policy.PasswordAgePolicyAddedEvent, *policy.PasswordAgePolicyChangedEvent:
			rm.PasswordAgePolicyReadModel.AppendEvents(e)
		}
	}
}

type PasswordAgePolicyAddedEvent struct {
	policy.PasswordAgePolicyAddedEvent
}

type PasswordAgePolicyChangedEvent struct {
	policy.PasswordAgePolicyChangedEvent
}
