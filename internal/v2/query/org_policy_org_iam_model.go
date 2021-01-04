package query

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/org"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type OrgOrgIAMPolicyReadModel struct{ OrgIAMPolicyReadModel }

func (rm *OrgOrgIAMPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.OrgIAMPolicyAddedEvent:
			rm.OrgIAMPolicyReadModel.AppendEvents(&e.OrgIAMPolicyAddedEvent)
		case *policy.OrgIAMPolicyAddedEvent:
			rm.OrgIAMPolicyReadModel.AppendEvents(e)
		}
	}
}
