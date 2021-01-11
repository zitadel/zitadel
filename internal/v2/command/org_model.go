package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

type OrgWriteModel struct {
	eventstore.WriteModel

	Name  string
	State domain.OrgState
}

func NewOrgWriteModel(orgID string) *OrgWriteModel {
	return &OrgWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   orgID,
			ResourceOwner: orgID,
		},
	}
}

func (wm *OrgWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *org.OrgAddedEvent,
			*iam.LabelPolicyChangedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *org.OrgAddedEvent:
			wm.Name = e.Name
			wm.State = domain.OrgStateActive
		case *org.OrgChangedEvent:
			wm.Name = e.Name
		}
	}
	return nil
}

func (wm *OrgWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.AggregateID)
}

func OrgAggregateFromWriteModel(wm *eventstore.WriteModel) *org.Aggregate {
	return &org.Aggregate{
		Aggregate: *eventstore.AggregateFromWriteModel(wm, org.AggregateType, org.AggregateVersion),
	}
}
