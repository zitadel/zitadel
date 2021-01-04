package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

type MemberWriteModel struct {
	eventstore.WriteModel

	UserID   string
	Roles    []string
	IsActive bool
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
			wm.IsActive = true
		case *member.MemberChangedEvent:
			wm.Roles = e.Roles
		case *member.MemberRemovedEvent:
			wm.Roles = nil
			wm.IsActive = false
		}
	}
	return wm.WriteModel.Reduce()
}
