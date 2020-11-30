package iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	OrgIAMPolicyAddedEventType = iamEventTypePrefix + policy.OrgIAMPolicyAddedEventType
)

type OrgIAMPolicyReadModel struct{ policy.OrgIAMPolicyReadModel }

func (rm *OrgIAMPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *OrgIAMPolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.OrgIAMPolicyAddedEvent)
		case *policy.OrgIAMPolicyAddedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}

type OrgIAMPolicyAddedEvent struct {
	policy.OrgIAMPolicyAddedEvent
}

func OrgIAMPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.OrgIAMPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OrgIAMPolicyAddedEvent{OrgIAMPolicyAddedEvent: *e.(*policy.OrgIAMPolicyAddedEvent)}, nil
}
