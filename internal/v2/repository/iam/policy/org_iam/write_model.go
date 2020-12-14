package org_iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/org_iam"
)

const (
	AggregateType = "iam"
)

type WriteModel struct {
	org_iam.WriteModel
}

func NewWriteModel(iamID string) *WriteModel {
	return &WriteModel{
		org_iam.WriteModel{
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
			wm.WriteModel.AppendEvents(&e.AddedEvent)
		case *ChangedEvent:
			wm.WriteModel.AppendEvents(&e.ChangedEvent)
		}
	}
}

func (wm *WriteModel) Reduce() error {
	return wm.WriteModel.Reduce()
}

func (wm *WriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.WriteModel.AggregateID)
}
