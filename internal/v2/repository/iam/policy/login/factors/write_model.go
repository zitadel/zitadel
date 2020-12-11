package factors

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/login/factors"
)

const (
	AggregateType = "iam"
)

type SecondFactorWriteModel struct {
	eventstore.WriteModel
	SecondFactor factors.SecondFactoryWriteModel

	iamID string
}

func NewSecondFactorWriteModel(iamID string) *SecondFactorWriteModel {
	return &SecondFactorWriteModel{
		iamID: iamID,
	}
}

func (wm *SecondFactorWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *LoginPolicySecondFactorAddedEvent:
			wm.SecondFactor.AppendEvents(&e.SecondFactorAddedEvent)
		}
	}
}

func (wm *SecondFactorWriteModel) Reduce() error {
	if err := wm.SecondFactor.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *SecondFactorWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.iamID)
}

type MultiFactorWriteModel struct {
	eventstore.WriteModel
	MultiFactor factors.MultiFactoryWriteModel

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
