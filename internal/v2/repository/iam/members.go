package iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/members"
)

type MembersReadModel struct {
	members.ReadModel
}

func (rm *MembersReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *MemberAddedEvent:
			rm.ReadModel.AppendEvents(&e.AddedEvent)
		case *MemberChangedEvent:
			rm.ReadModel.AppendEvents(&e.ChangedEvent)
		case *MemberRemovedEvent:
			rm.ReadModel.AppendEvents(&e.RemovedEvent)
		}
	}
}
