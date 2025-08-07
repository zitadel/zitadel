package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type OrgFlowWriteModel struct {
	FlowWriteModel
}

func NewOrgFlowWriteModel(flowType domain.FlowType, resourceOwner string) *OrgFlowWriteModel {
	return &OrgFlowWriteModel{
		FlowWriteModel: *NewFlowWriteModel(flowType, resourceOwner),
	}
}

func (wm *OrgFlowWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.TriggerActionsSetEvent:
			if e.FlowType != wm.FlowType {
				continue
			}
			wm.FlowWriteModel.AppendEvents(&e.TriggerActionsSetEvent)
		case *org.TriggerActionsCascadeRemovedEvent:
			if e.FlowType != wm.FlowType {
				continue
			}
			wm.FlowWriteModel.AppendEvents(&e.TriggerActionsCascadeRemovedEvent)
		case *org.FlowClearedEvent:
			if e.FlowType != wm.FlowType {
				continue
			}
			wm.FlowWriteModel.AppendEvents(&e.FlowClearedEvent)
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
		EventTypes(org.TriggerActionsSetEventType,
			org.TriggerActionsCascadeRemovedEventType,
			org.FlowClearedEventType).
		Builder()
}
