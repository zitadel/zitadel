package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	OrgIAMPolicyAddedEventType = orgEventTypePrefix + policy.OrgIAMPolicyAddedEventType
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
