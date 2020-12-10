package second_factors

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/login/second_factors"
)

const (
	AggregateType = "iam"
)

type SecondFactorWriteModel struct {
	eventstore.WriteModel
	SecondFactor second_factors.SecondFactoryWriteModel

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
