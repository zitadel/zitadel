package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

type OrgMembersReadModel struct {
	MembersReadModel
}

func (rm *OrgMembersReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.MemberAddedEvent:
			rm.MembersReadModel.AppendEvents(&e.MemberAddedEvent)
		case *org.MemberChangedEvent:
			rm.MembersReadModel.AppendEvents(&e.MemberChangedEvent)
		case *org.MemberRemovedEvent:
			rm.MembersReadModel.AppendEvents(&e.MemberRemovedEvent)
		}
	}
}

type OrgMemberReadModel MemberReadModel

func (rm *OrgMemberReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.MemberAddedEvent:
			rm.ReadModel.AppendEvents(&e.MemberAddedEvent)
		case *org.MemberChangedEvent:
			rm.ReadModel.AppendEvents(&e.MemberChangedEvent)
		}
	}
}
