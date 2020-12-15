package idpprovider

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/login/idpprovider"
)

const (
	AggregateType = "iam"
)

type WriteModel struct {
	idpprovider.WriteModel
	IsRemoved bool
}

func NewWriteModel(iamID, idpConfigID string) *WriteModel {
	return &WriteModel{
		WriteModel: idpprovider.WriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID: iamID,
			},
			IDPConfigID: idpConfigID,
		},
		IsRemoved: false,
	}
}

func (wm *WriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *AddedEvent:
			if e.IDPConfigID != wm.IDPConfigID {
				continue
			}
			wm.WriteModel.AppendEvents(&e.AddedEvent)
		}
	}
}

func (wm *WriteModel) Reduce() error {
	return wm.WriteModel.Reduce()
}

func (wm *WriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.AggregateID)
}
