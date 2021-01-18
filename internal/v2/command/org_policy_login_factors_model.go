package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

type OrgSecondFactorWriteModel struct {
	SecondFactorWriteModel
}

func NewOrgSecondFactorWriteModel(orgID string) *OrgSecondFactorWriteModel {
	return &OrgSecondFactorWriteModel{
		SecondFactorWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
		},
	}
}

func (wm *OrgSecondFactorWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.LoginPolicySecondFactorAddedEvent:
			wm.WriteModel.AppendEvents(&e.SecondFactorAddedEvent)
		}
	}
}

func (wm *OrgSecondFactorWriteModel) Reduce() error {
	return wm.WriteModel.Reduce()
}

func (wm *OrgSecondFactorWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.WriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}

type OrgMultiFactorWriteModel struct {
	MultiFactoryWriteModel
}

func NewOrgMultiFactorWriteModel(orgID string) *OrgMultiFactorWriteModel {
	return &OrgMultiFactorWriteModel{
		MultiFactoryWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
		},
	}
}

func (wm *OrgMultiFactorWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.LoginPolicyMultiFactorAddedEvent:
			wm.WriteModel.AppendEvents(&e.MultiFactorAddedEvent)
		}
	}
}

func (wm *OrgMultiFactorWriteModel) Reduce() error {
	return wm.WriteModel.Reduce()
}

func (wm *OrgMultiFactorWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.WriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}
