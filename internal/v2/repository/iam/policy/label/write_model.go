package label

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/label"
)

const (
	AggregateType = "iam"
)

type WriteModel struct {
	eventstore.WriteModel
	Policy label.WriteModel

	iamID string
}

func NewWriteModel(iamID string) *WriteModel {
	return &WriteModel{
		iamID: iamID,
	}
}

func (wm *WriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
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
	if err := wm.Policy.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *WriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.iamID)
}
