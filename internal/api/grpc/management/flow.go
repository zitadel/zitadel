package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	action_grpc "github.com/caos/zitadel/internal/api/grpc/action"
	obj_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetFlow(ctx context.Context, req *mgmt_pb.GetFlowRequest) (*mgmt_pb.GetFlowResponse, error) {
	flow, err := s.query.GetFlow(ctx, action_grpc.FlowTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetFlowResponse{
		Flow: action_grpc.FlowToPb(flow),
	}, nil
}

func (s *Server) ClearFlow(ctx context.Context, req *mgmt_pb.ClearFlowRequest) (*mgmt_pb.ClearFlowResponse, error) {
	details, err := s.command.ClearFlow(ctx, action_grpc.FlowTypeToDomain(req.Type), authz.GetCtxData(ctx).OrgID)
	return &mgmt_pb.ClearFlowResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(details),
	}, err
}

func (s *Server) SetTriggerActions(ctx context.Context, req *mgmt_pb.SetTriggerActionsRequest) (*mgmt_pb.SetTriggerActionsResponse, error) {
	details, err := s.command.SetTriggerActions(
		ctx,
		action_grpc.FlowTypeToDomain(req.FlowType),
		action_grpc.TriggerTypeToDomain(req.TriggerType),
		req.ActionIds,
		authz.GetCtxData(ctx).OrgID,
	)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetTriggerActionsResponse{
		Details: obj_grpc.AddToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner,
		),
	}, nil
}
