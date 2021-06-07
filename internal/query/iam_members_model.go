package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMMembersReadModel struct {
	MembersReadModel
}

func (rm *IAMMembersReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.MemberAddedEvent:
			rm.MembersReadModel.AppendEvents(&e.MemberAddedEvent)
		case *iam.MemberChangedEvent:
			rm.MembersReadModel.AppendEvents(&e.MemberChangedEvent)
		case *iam.MemberRemovedEvent:
			rm.MembersReadModel.AppendEvents(&e.MemberRemovedEvent)
		case *iam.MemberCascadeRemovedEvent:
			rm.MembersReadModel.AppendEvents(&e.MemberCascadeRemovedEvent)
		}
	}
}
