package member

import "github.com/caos/zitadel/internal/eventstore/v2"

//WriteModel is used to create events
// It has no computed fields and represents the data
// which can be changed
type WriteModel struct {
	eventstore.WriteModel

	UserID    string
	Roles     []string
	IsRemoved bool
}

func NewWriteModel(userID string) *WriteModel {
	return &WriteModel{
		UserID: userID,
	}
}

//Reduce extends eventstore.ReadModel
func (wm *WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *AddedEvent:
			wm.UserID = e.UserID
			wm.Roles = e.Roles
		case *ChangedEvent:
			wm.Roles = e.Roles
		case *RemovedEvent:
			wm.Roles = nil
			wm.IsRemoved = true
		}
	}
	return wm.WriteModel.Reduce()
}

// func (wm *WriteModel) Query() *eventstore.SearchQueryFactory {
// 	return eventstore.NewSearchQueryFactory(eventstore.ColumnsEvent, wm.aggregateType).
// 		AggregateIDs(wm.aggregateID)
// }
