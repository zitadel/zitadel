package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/flow"
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
		case *flow.FlowClearedEvent:
			wm.Triggers = nil
		}
	}
	return wm.WriteModel.Reduce()
}
