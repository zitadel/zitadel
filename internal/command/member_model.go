package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/member"
)

type MemberWriteModel struct {
	eventstore.WriteModel

	UserID string
	Roles  []string

	State domain.MemberState
}

func NewMemberWriteModel(userID string) *MemberWriteModel {
	return &MemberWriteModel{
		UserID: userID,
	}
}

func (wm *MemberWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *member.MemberAddedEvent:
			wm.UserID = e.UserID
			wm.Roles = e.Roles
			wm.State = domain.MemberStateActive
		case *member.MemberChangedEvent:
			wm.Roles = e.Roles
		case *member.MemberRemovedEvent:
			wm.Roles = nil
			wm.State = domain.MemberStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func isMemberExisting(state domain.MemberState) bool {
	return !hasMemberState(state, domain.MemberStateUnspecified, domain.MemberStateRemoved)
}

func hasMemberState(check domain.MemberState, states ...domain.MemberState) bool {
	for _, state := range states {
		if check == state {
			return true
		}
	}
	return false
}
