package iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	PasswordLockoutPolicyAddedEventType   = IamEventTypePrefix + policy.PasswordLockoutPolicyAddedEventType
	PasswordLockoutPolicyChangedEventType = IamEventTypePrefix + policy.PasswordLockoutPolicyChangedEventType
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

func PasswordLockoutPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.PasswordLockoutPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordLockoutPolicyAddedEvent{PasswordLockoutPolicyAddedEvent: *e.(*policy.PasswordLockoutPolicyAddedEvent)}, nil
}

type PasswordLockoutPolicyChangedEvent struct {
	policy.PasswordLockoutPolicyChangedEvent
}

func PasswordLockoutPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.PasswordLockoutPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordLockoutPolicyChangedEvent{PasswordLockoutPolicyChangedEvent: *e.(*policy.PasswordLockoutPolicyChangedEvent)}, nil
}
