package login

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/login"
)

const (
	AggregateType = "iam"
)

type WriteModel struct {
	Policy login.WriteModel
}

func NewWriteModel(iamID string) *WriteModel {
	return &WriteModel{
		Policy: login.WriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID: iamID,
			},
		},
	}
}

func (wm *WriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *AddedEvent:
			wm.Policy.AppendEvents(&e.AddedEvent)
		case *ChangedEvent:
			wm.Policy.AppendEvents(&e.ChangedEvent)
		}
	}
}

func (wm *WriteModel) Reduce() error {
	return wm.Policy.Reduce()
}

func (wm *WriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.Policy.AggregateID)
}
