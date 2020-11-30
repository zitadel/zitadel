package iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
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

func PasswordAgePolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.PasswordAgePolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordAgePolicyAddedEvent{PasswordAgePolicyAddedEvent: *e.(*policy.PasswordAgePolicyAddedEvent)}, nil
}

type PasswordAgePolicyChangedEvent struct {
	policy.PasswordAgePolicyChangedEvent
}

func PasswordAgePolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.PasswordAgePolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordAgePolicyChangedEvent{PasswordAgePolicyChangedEvent: *e.(*policy.PasswordAgePolicyChangedEvent)}, nil
}
