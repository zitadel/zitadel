package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/policy"
)

type IAMOrgIAMPolicyReadModel struct{ OrgIAMPolicyReadModel }

func (rm *IAMOrgIAMPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.OrgIAMPolicyAddedEvent:
			rm.OrgIAMPolicyReadModel.AppendEvents(&e.OrgIAMPolicyAddedEvent)
		case *policy.OrgIAMPolicyAddedEvent:
			rm.OrgIAMPolicyReadModel.AppendEvents(e)
		}
	}
}
