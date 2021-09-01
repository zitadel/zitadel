package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	action_grpc "github.com/caos/zitadel/internal/api/grpc/action"
	obj_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) ListActions(ctx context.Context, req *mgmt_pb.ListActionsRequest) (*mgmt_pb.ListActionsResponse, error) {
	return nil, nil
}

func (s *Server) GetAction(ctx context.Context, req *mgmt_pb.GetActionRequest) (*mgmt_pb.GetActionResponse, error) {
	return nil, nil
}

func (s *Server) CreateAction(ctx context.Context, req *mgmt_pb.CreateActionRequest) (*mgmt_pb.CreateActionResponse, error) {
	id, details, err := s.command.AddAction(ctx, createActionRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.CreateActionResponse{
		Id: id,
		Details: obj_grpc.AddToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateAction(ctx context.Context, req *mgmt_pb.UpdateActionRequest) (*mgmt_pb.UpdateActionResponse, error) {
	details, err := s.command.ChangeAction(ctx, updateActionRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateActionResponse{
		Details: obj_grpc.AddToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner,
		),
	}, nil
}

func (s *Server) DeleteAction(ctx context.Context, req *mgmt_pb.DeleteActionRequest) (*mgmt_pb.DeleteActionResponse, error) {
	_, err := s.command.DeleteAction(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	return &mgmt_pb.DeleteActionResponse{}, err
}

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
	return nil, nil
}
