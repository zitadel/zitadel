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

	userID        string
	aggregateType eventstore.AggregateType
	aggregateID   string
}

func NewWriteModel(
	userID string,
	aggregateType eventstore.AggregateType,
	aggregateID string,
) *WriteModel {

	return &WriteModel{
		WriteModel:    *eventstore.NewWriteModel(),
		userID:        userID,
		aggregateType: aggregateType,
		aggregateID:   aggregateID,
	}
}

//Reduce extends eventstore.ReadModel
func (wm *WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *AddedEvent:
			if e.UserID != wm.userID {
				continue
			}
			wm.UserID = e.UserID
			wm.Roles = e.Roles
		case *ChangedEvent:
			if e.UserID != wm.userID {
				continue
			}
			wm.UserID = e.UserID
			wm.Roles = e.Roles
		case *RemovedEvent:
			if e.UserID != wm.userID {
				continue
			}
			wm.Roles = nil
			wm.IsRemoved = true
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *WriteModel) Query() *eventstore.SearchQueryFactory {
	return eventstore.NewSearchQueryFactory(eventstore.ColumnsEvent, wm.aggregateType).
		AggregateIDs(wm.aggregateID)
}
