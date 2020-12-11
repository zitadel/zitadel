package idpprovider

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/login/idpprovider"
)

const (
	AggregateType = "iam"
)

type WriteModel struct {
	eventstore.WriteModel
	Provider idpprovider.WriteModel

	idpConfigID string
	iamID       string

	IsRemoved bool
}

func NewWriteModel(iamID, idpConfigID string) *WriteModel {
	return &WriteModel{
		iamID:       iamID,
		idpConfigID: idpConfigID,
	}
}

func (wm *WriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *AddedEvent:
			if e.IDPConfigID != wm.idpConfigID {
				continue
			}
			wm.Provider.AppendEvents(&e.AddedEvent)
		}
	}
}

func (wm *WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *AddedEvent:
			if e.IDPConfigID != wm.idpConfigID {
				continue
			}
			wm.IsRemoved = false
		case *RemovedEvent:
			if e.IDPConfigID != wm.idpConfigID {
				continue
			}
			wm.IsRemoved = true
		}
	}
	if err := wm.Provider.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *WriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.iamID)
}
