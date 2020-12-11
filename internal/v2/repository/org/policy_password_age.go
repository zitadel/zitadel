package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_age"
)

var (
	PasswordAgePolicyAddedEventType   = orgEventTypePrefix + password_age.PasswordAgePolicyAddedEventType
	PasswordAgePolicyChangedEventType = orgEventTypePrefix + password_age.PasswordAgePolicyChangedEventType
)

type PasswordAgePolicyReadModel struct {
	password_age.ReadModel
}

func (rm *PasswordAgePolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *PasswordAgePolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.AddedEvent)
		case *PasswordAgePolicyChangedEvent:
			rm.ReadModel.AppendEvents(&e.ChangedEvent)
		case *password_age.AddedEvent, *password_age.ChangedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}

type PasswordAgePolicyAddedEvent struct {
	password_age.AddedEvent
}

type PasswordAgePolicyChangedEvent struct {
	password_age.ChangedEvent
}
