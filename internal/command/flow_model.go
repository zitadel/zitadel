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
}

func NewFlowWriteModel(flowType domain.FlowType, resourceOwner string) *FlowWriteModel {
	return &FlowWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   resourceOwner,
			ResourceOwner: resourceOwner,
		},
		FlowType: flowType,
	}
}

func (wm *FlowWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *flow.TriggerActionsSetEvent:
			_ = e
		}
	}
	return wm.WriteModel.Reduce()
}
