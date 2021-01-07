package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/member"
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
