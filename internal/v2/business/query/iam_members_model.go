package query

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/iam"
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
		}
	}
}
