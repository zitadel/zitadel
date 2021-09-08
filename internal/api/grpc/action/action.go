package action

import (
	"github.com/caos/zitadel/internal/domain"
	action_pb "github.com/caos/zitadel/pkg/grpc/action"
)

func FlowTypeToDomain(flowType action_pb.FlowType) domain.FlowType {
	switch flowType { //TODO: converter
	case action_pb.FlowType_FLOW_TYPE_EXTERNAL_AUTHENTICATION:
		return domain.FlowTypeExternalAuthentication
	case action_pb.FlowType_FLOW_TYPE_EXTERNAL_REGISTRATION:
		return domain.FlowTypeExternalRegistration
	default:
		return domain.FlowTypeUnspecified
	}
}

func TriggerTypeToDomain(flowType action_pb.TriggerType) domain.TriggerType {
	switch flowType { //TODO: converter
	case action_pb.TriggerType_TRIGGER_TYPE_POST_AUTHENTICATION:
		return domain.TriggerTypePostAuthentication
	default:
		return domain.TriggerTypeUnspecified
	}
}
