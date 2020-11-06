package member

import "github.com/caos/zitadel/internal/eventstore/v2"

type MemberReadModel struct {
	eventstore.ReadModel

	UserID string
	Roles  []string
}

func NewMemberReadModel(userID string) *MemberReadModel {
	return &MemberReadModel{
		ReadModel: *eventstore.NewReadModel(),
	}
}

func (rm *MemberReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *MemberAddedEvent:
			rm.UserID = e.UserID
			rm.Roles = e.Roles
		case *MemberChangedEvent:
			rm.UserID = e.UserID
			rm.Roles = e.Roles
		}
	}
	return rm.ReadModel.Reduce()
}

type MembersReadModel struct {
	eventstore.ReadModel
	Members []*MemberReadModel
}

func NewMembersReadModel() *MembersReadModel {
	return &MembersReadModel{
		ReadModel: *eventstore.NewReadModel(),
		Members:   []*MemberReadModel{},
	}
}

func (rm *MembersReadModel) MemberByUserID(id string) (idx int, member *MemberReadModel) {
	for idx, member = range rm.Members {
		if member.UserID == id {
			return idx, member
		}
	}
	return -1, nil
}

func (rm *MembersReadModel) AppendEvents(events ...eventstore.EventReader) (err error) {
	for _, event := range events {
		switch e := event.(type) {
		case *MemberAddedEvent:
			member := NewMemberReadModel(e.UserID)
			rm.Members = append(rm.Members, member)
			member.AppendEvents(e)
		case *MemberChangedEvent:
			_, member := rm.MemberByUserID(e.UserID)
			member.AppendEvents(e)
		case *MemberRemovedEvent:
			idx, _ := rm.MemberByUserID(e.UserID)
			if idx < 0 {
				return nil
			}
			copy(rm.Members[idx:], rm.Members[idx+1:])
			rm.Members[len(rm.Members)-1] = nil
			rm.Members = rm.Members[:len(rm.Members)-1]
		}
	}
	return err
}

func (rm *MembersReadModel) Reduce() (err error) {
	for _, member := range rm.Members {
		err = member.Reduce()
		if err != nil {
			return err
		}
	}
	return rm.ReadModel.Reduce()
}
