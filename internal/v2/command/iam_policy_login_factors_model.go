package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

type IAMSecondFactorWriteModel struct {
	SecondFactorWriteModel
}

func NewIAMSecondFactorWriteModel() *IAMSecondFactorWriteModel {
	return &IAMSecondFactorWriteModel{
		SecondFactorWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
		},
	}
}

func (wm *IAMSecondFactorWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.LoginPolicySecondFactorAddedEvent:
			wm.WriteModel.AppendEvents(&e.SecondFactorAddedEvent)
		}
	}
}

func (wm *IAMSecondFactorWriteModel) Reduce() error {
	return wm.WriteModel.Reduce()
}

func (wm *IAMSecondFactorWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.WriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}

type IAMMultiFactorWriteModel struct {
	MultiFactoryWriteModel
}

func NewIAMMultiFactorWriteModel() *IAMMultiFactorWriteModel {
	return &IAMMultiFactorWriteModel{
		MultiFactoryWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
		},
	}
}

func (wm *IAMMultiFactorWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.LoginPolicyMultiFactorAddedEvent:
			wm.WriteModel.AppendEvents(&e.MultiFactorAddedEvent)
		}
	}
}

func (wm *IAMMultiFactorWriteModel) Reduce() error {
	return wm.WriteModel.Reduce()
}

func (wm *IAMMultiFactorWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.WriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}
