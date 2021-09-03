package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

type OrgFlowWriteModel struct {
	FlowWriteModel
}

func NewOrgFlowWriteModel(flowType domain.FlowType, resourceOwner string) *OrgFlowWriteModel {
	return &OrgFlowWriteModel{
		FlowWriteModel: *NewFlowWriteModel(flowType, resourceOwner),
	}
}

func (wm *OrgFlowWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.TriggerActionsSetEvent:
			wm.FlowWriteModel.AppendEvents(&e.TriggerActionsSetEvent)
		}
	}
}

func (wm *OrgFlowWriteModel) Reduce() error {
	return wm.FlowWriteModel.Reduce()
}

func (wm *OrgFlowWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		EventTypes(org.TriggerActionsSetEventType).
		Builder()
}
