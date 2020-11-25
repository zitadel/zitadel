package members

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

type ReadModel struct {
	eventstore.ReadModel

	Members []*member.ReadModel
}

func NewMembersReadModel() *ReadModel {
	return &ReadModel{
		ReadModel: *eventstore.NewReadModel(),
		Members:   []*member.ReadModel{},
	}
}

func (rm *ReadModel) MemberByUserID(id string) (idx int, member *member.ReadModel) {
	for idx, member = range rm.Members {
		if member.UserID == id {
			return idx, member
		}
	}
	return -1, nil
}

func (rm *ReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *member.AddedEvent:
			member := member.NewMemberReadModel(e.UserID)
			rm.Members = append(rm.Members, member)
			member.AppendEvents(e)
		case *member.ChangedEvent:
			_, member := rm.MemberByUserID(e.UserID)
			member.AppendEvents(e)
		case *member.RemovedEvent:
			idx, _ := rm.MemberByUserID(e.UserID)
			if idx < 0 {
				continue
			}
			copy(rm.Members[idx:], rm.Members[idx+1:])
			rm.Members[len(rm.Members)-1] = nil
			rm.Members = rm.Members[:len(rm.Members)-1]
		}
	}
}

func (rm *ReadModel) Reduce() (err error) {
	for _, member := range rm.Members {
		err = member.Reduce()
		if err != nil {
			return err
		}
	}
	return rm.ReadModel.Reduce()
}
