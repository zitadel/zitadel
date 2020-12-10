package multi_factors

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/login/multi_factors"
)

const (
	AggregateType = "iam"
)

type MultiFactorWriteModel struct {
	eventstore.WriteModel
	MultiFactor multi_factors.MultiFactoryWriteModel

	iamID string
}

func NewMultiFactorWriteModel(iamID string) *MultiFactorWriteModel {
	return &MultiFactorWriteModel{
		iamID: iamID,
	}
}

func (wm *MultiFactorWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *LoginPolicyMultiFactorAddedEvent:
			wm.MultiFactor.AppendEvents(&e.MultiFactorAddedEvent)
		}
	}
}

func (wm *MultiFactorWriteModel) Reduce() error {
	if err := wm.MultiFactor.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *MultiFactorWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.iamID)
}
