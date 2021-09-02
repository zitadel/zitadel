package action

import (
	"github.com/caos/zitadel/internal/domain"
	action_pb "github.com/caos/zitadel/pkg/grpc/action"
)

func FlowTypeToDomain(flowType action_pb.FlowType) domain.FlowType {
	switch flowType { //TODO: converter
	case action_pb.FlowType_FLOW_TYPE_REGISTER:
		return domain.FlowTypeRegister
	default:
		return domain.FlowTypeUnspecified
	}
}

func TriggerTypeToDomain(flowType action_pb.TriggerType) domain.TriggerType {
	switch flowType { //TODO: converter
	case action_pb.TriggerType_TRIGGER_TYPE_PRE_REGISTER:
		return domain.TriggerTypePreRegister
	default:
		return domain.TriggerTypeUnspecified
	}
}
