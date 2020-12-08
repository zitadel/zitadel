package iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	PasswordComplexityPolicyAddedEventType   = IamEventTypePrefix + policy.PasswordComplexityPolicyAddedEventType
	PasswordComplexityPolicyChangedEventType = IamEventTypePrefix + policy.PasswordComplexityPolicyChangedEventType
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
		case *policy.PasswordComplexityPolicyAddedEvent,
			*policy.PasswordComplexityPolicyChangedEvent:

			rm.ReadModel.AppendEvents(e)
		}
	}
}

type PasswordComplexityPolicyAddedEvent struct {
	policy.PasswordComplexityPolicyAddedEvent
}

func PasswordComplexityPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.PasswordComplexityPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordComplexityPolicyAddedEvent{PasswordComplexityPolicyAddedEvent: *e.(*policy.PasswordComplexityPolicyAddedEvent)}, nil
}

type PasswordComplexityPolicyChangedEvent struct {
	policy.PasswordComplexityPolicyChangedEvent
}

func PasswordComplexityPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.PasswordComplexityPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordComplexityPolicyChangedEvent{PasswordComplexityPolicyChangedEvent: *e.(*policy.PasswordComplexityPolicyChangedEvent)}, nil
}
