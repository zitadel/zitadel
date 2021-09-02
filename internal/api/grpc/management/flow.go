package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	action_grpc "github.com/caos/zitadel/internal/api/grpc/action"
	obj_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) ListFlows(ctx context.Context, req *mgmt_pb.ListFlowsRequest) (*mgmt_pb.ListFlowsResponse, error) {
	return nil, nil
}

func (s *Server) CreateFlow(ctx context.Context, req *mgmt_pb.CreateFlowRequest) (*mgmt_pb.CreateFlowResponse, error) {
	return nil, nil
}

func (s *Server) GetFlow(ctx context.Context, req *mgmt_pb.GetFlowRequest) (*mgmt_pb.GetFlowResponse, error) {
	return nil, nil
}

func (s *Server) DeleteFlow(ctx context.Context, req *mgmt_pb.DeleteFlowRequest) (*mgmt_pb.DeleteFlowResponse, error) {
	_, err := s.command.DeleteFlow(ctx, action_grpc.FlowTypeToDomain(req.Type), authz.GetCtxData(ctx).OrgID)
	return &mgmt_pb.DeleteFlowResponse{}, err
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
