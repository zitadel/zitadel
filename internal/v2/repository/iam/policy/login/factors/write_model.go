package factors

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/login/factors"
)

const (
	AggregateType = "iam"
)

type SecondFactorWriteModel struct {
	SecondFactor factors.SecondFactorWriteModel
}

func NewSecondFactorWriteModel(iamID string) *SecondFactorWriteModel {
	return &SecondFactorWriteModel{
		SecondFactor: factors.SecondFactorWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID: iamID,
			},
		},
	}
}

func (wm *SecondFactorWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *LoginPolicySecondFactorAddedEvent:
			wm.SecondFactor.AppendEvents(&e.SecondFactorAddedEvent)
		}
	}
}

func (wm *SecondFactorWriteModel) Reduce() error {
	return wm.SecondFactor.Reduce()
}

func (wm *SecondFactorWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.SecondFactor.AggregateID)
}

type MultiFactorWriteModel struct {
	MultiFactor factors.MultiFactoryWriteModel
}

func NewMultiFactorWriteModel(iamID string) *MultiFactorWriteModel {
	return &MultiFactorWriteModel{
		MultiFactor: factors.MultiFactoryWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID: iamID,
			},
		},
	}
}

func (wm *MultiFactorWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *LoginPolicyMultiFactorAddedEvent:
			wm.MultiFactor.AppendEvents(&e.MultiFactorAddedEvent)
		}
	}
}

func (wm *MultiFactorWriteModel) Reduce() error {
	return wm.MultiFactor.Reduce()
}

func (wm *MultiFactorWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.MultiFactor.AggregateID)
}
