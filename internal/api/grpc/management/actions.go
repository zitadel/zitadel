package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	action_grpc "github.com/caos/zitadel/internal/api/grpc/action"
	obj_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) ListActions(ctx context.Context, req *mgmt_pb.ListActionsRequest) (*mgmt_pb.ListActionsResponse, error) {
	query, _ := listActionsToQuery(authz.GetCtxData(ctx).OrgID, req)
	actions, err := s.query.SearchActions(ctx, query)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListActionsResponse{
		Result: action_grpc.ActionsToPb(actions),
	}, nil
}

func (s *Server) GetAction(ctx context.Context, req *mgmt_pb.GetActionRequest) (*mgmt_pb.GetActionResponse, error) {
	action, err := s.query.GetAction(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetActionResponse{
		Action: action_grpc.ActionToPb(action),
	}, nil
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

func (s *Server) DeactivateAction(ctx context.Context, req *mgmt_pb.DeactivateActionRequest) (*mgmt_pb.DeactivateActionResponse, error) {
	details, err := s.command.DeactivateAction(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	return &mgmt_pb.DeactivateActionResponse{
		Details: obj_grpc.AddToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner,
		),
	}, err
}

func (s *Server) ReactivateAction(ctx context.Context, req *mgmt_pb.ReactivateActionRequest) (*mgmt_pb.ReactivateActionResponse, error) {
	details, err := s.command.ReactivateAction(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateActionResponse{
		Details: obj_grpc.AddToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner,
		),
	}, nil
}

func (s *Server) DeleteAction(ctx context.Context, req *mgmt_pb.DeleteActionRequest) (*mgmt_pb.DeleteActionResponse, error) {
	flowTypes, err := s.query.GetFlowTypesOfActionID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	_, err = s.command.DeleteAction(ctx, req.Id, authz.GetCtxData(ctx).OrgID, flowTypes...)
	return &mgmt_pb.DeleteActionResponse{}, err
}
