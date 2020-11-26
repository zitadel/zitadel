package member

import "github.com/caos/zitadel/internal/eventstore/v2"

//ReadModel represenets the default member view.
// It's computed from events.
type ReadModel struct {
	eventstore.ReadModel

	UserID string
	Roles  []string
}

//NewMemberReadModel is the default constructor of ReadModel
func NewMemberReadModel(userID string) *ReadModel {
	return &ReadModel{
		UserID: userID,
	}
}

//Reduce extends eventstore.ReadModel
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
