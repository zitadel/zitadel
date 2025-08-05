package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/flow"
)

type FlowWriteModel struct {
	eventstore.WriteModel

	FlowType domain.FlowType
	State    domain.FlowState
	Triggers map[domain.TriggerType][]string
}

func NewFlowWriteModel(flowType domain.FlowType, resourceOwner string) *FlowWriteModel {
	return &FlowWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   resourceOwner,
			ResourceOwner: resourceOwner,
		},
		FlowType: flowType,
		Triggers: make(map[domain.TriggerType][]string),
	}
}

func (wm *FlowWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *flow.TriggerActionsSetEvent:
			if wm.Triggers == nil {
				wm.Triggers = make(map[domain.TriggerType][]string)
			}
			wm.Triggers[e.TriggerType] = e.ActionIDs
		case *flow.TriggerActionsCascadeRemovedEvent:
			remove(wm.Triggers[e.TriggerType], e.ActionID)
		case *flow.FlowClearedEvent:
			wm.Triggers = nil
		}
	}
	return wm.WriteModel.Reduce()
}

func remove(ids []string, id string) {
	for i := 0; i < len(ids); i++ {
		if ids[i] == id {
			ids = append(ids[:i], ids[i+1:]...)
			break
		}
	}
}
