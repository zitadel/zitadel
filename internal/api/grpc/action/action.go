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
