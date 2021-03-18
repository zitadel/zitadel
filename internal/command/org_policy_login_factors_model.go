package command

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
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
		case *org.LoginPolicySecondFactorAddedEvent:
			wm.WriteModel.AppendEvents(&e.SecondFactorAddedEvent)
		case *org.LoginPolicySecondFactorRemovedEvent:
			wm.WriteModel.AppendEvents(&e.SecondFactorRemovedEvent)
		}
	}
}

func (wm *OrgSecondFactorWriteModel) Reduce() error {
	return wm.WriteModel.Reduce()
}

func (wm *OrgSecondFactorWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.WriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			org.LoginPolicySecondFactorAddedEventType,
			org.LoginPolicySecondFactorRemovedEventType)
}

type OrgMultiFactorWriteModel struct {
	MultiFactorWriteModel
}

func NewOrgMultiFactorWriteModel(orgID string) *OrgMultiFactorWriteModel {
	return &OrgMultiFactorWriteModel{
		MultiFactorWriteModel{
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
		case *org.LoginPolicyMultiFactorAddedEvent:
			wm.WriteModel.AppendEvents(&e.MultiFactorAddedEvent)
		case *org.LoginPolicyMultiFactorRemovedEvent:
			wm.WriteModel.AppendEvents(&e.MultiFactorRemovedEvent)
		}
	}
}

func (wm *OrgMultiFactorWriteModel) Reduce() error {
	return wm.WriteModel.Reduce()
}

func (wm *OrgMultiFactorWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.WriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			org.LoginPolicyMultiFactorAddedEventType,
			org.LoginPolicyMultiFactorRemovedEventType)
}
