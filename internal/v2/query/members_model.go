package query

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

type MembersReadModel struct {
	eventstore.ReadModel

	Members []*MemberReadModel
}

func (rm *MembersReadModel) MemberByUserID(id string) (idx int, member *MemberReadModel) {
	for idx, member = range rm.Members {
		if member.UserID == id {
			return idx, member
		}
	}
	return -1, nil
}

func (rm *MembersReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *member.MemberAddedEvent:
			m := NewMemberReadModel(e.UserID)
			rm.Members = append(rm.Members, m)
			m.AppendEvents(e)
		case *member.MemberChangedEvent:
			_, m := rm.MemberByUserID(e.UserID)
			m.AppendEvents(e)
		case *member.MemberRemovedEvent:
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

func (rm *MembersReadModel) Reduce() (err error) {
	for _, m := range rm.Members {
		err = m.Reduce()
		if err != nil {
			return err
		}
	}
	return rm.ReadModel.Reduce()
}
