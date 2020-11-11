package member

import "github.com/caos/zitadel/internal/eventstore/v2"

type ReadModel struct {
	eventstore.ReadModel

	UserID string
	Roles  []string
}

func NewMemberReadModel(userID string) *ReadModel {
	return &ReadModel{
		ReadModel: *eventstore.NewReadModel(),
	}
}

func (rm *ReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *AddedEvent:
			rm.UserID = e.UserID
			rm.Roles = e.Roles
		case *ChangedEvent:
			rm.UserID = e.UserID
			rm.Roles = e.Roles
		}
	}
	return rm.ReadModel.Reduce()
}
