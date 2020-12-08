package org_iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/org_iam"
)

type OrgIAMPolicyReadModel struct{ org_iam.OrgIAMPolicyReadModel }

func (rm *OrgIAMPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *OrgIAMPolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.OrgIAMPolicyAddedEvent)
		case *org_iam.OrgIAMPolicyAddedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}
