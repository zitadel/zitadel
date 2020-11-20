package member

import "github.com/caos/zitadel/internal/eventstore/v2"

//WriteModel is used to create events
// It has no computed fields and represents the data
// which can be changed
type WriteModel struct {
	eventstore.WriteModel

	UserID string
	Roles  []string
}

//Reduce extends eventstore.ReadModel
func (rm *WriteModel) Reduce() error {
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
