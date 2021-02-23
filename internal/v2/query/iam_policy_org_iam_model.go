package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/policy"
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
