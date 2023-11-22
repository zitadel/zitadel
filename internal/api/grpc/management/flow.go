package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	action_grpc "github.com/zitadel/zitadel/internal/api/grpc/action"
	obj_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	action_pb "github.com/zitadel/zitadel/pkg/grpc/action"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Server) ListFlowTypes(ctx context.Context, _ *mgmt_pb.ListFlowTypesRequest) (*mgmt_pb.ListFlowTypesResponse, error) {
	return &mgmt_pb.ListFlowTypesResponse{
		Result: []*action_pb.FlowType{
			action_grpc.FlowTypeToPb(domain.FlowTypeExternalAuthentication),
			action_grpc.FlowTypeToPb(domain.FlowTypeCustomiseToken),
			action_grpc.FlowTypeToPb(domain.FlowTypeInternalAuthentication),
			action_grpc.FlowTypeToPb(domain.FlowTypeCustomizeSAMLResponse),
		},
	}, nil
}

func (s *Server) ListFlowTriggerTypes(ctx context.Context, req *mgmt_pb.ListFlowTriggerTypesRequest) (*mgmt_pb.ListFlowTriggerTypesResponse, error) {
	triggerTypes := action_grpc.FlowTypeToDomain(req.Type).TriggerTypes()
	if len(triggerTypes) == 0 {
		return nil, errors.ThrowNotFound(nil, "MANAG-P2OBk", "Errors.NotFound")
	}
	return &mgmt_pb.ListFlowTriggerTypesResponse{
		Result: action_grpc.TriggerTypesToPb(triggerTypes),
	}, nil
}

func (s *Server) GetFlow(ctx context.Context, req *mgmt_pb.GetFlowRequest) (*mgmt_pb.GetFlowResponse, error) {
	flow, err := s.query.GetFlow(ctx, action_grpc.FlowTypeToDomain(req.Type), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetFlowResponse{
		Flow: action_grpc.FlowToPb(flow),
	}, nil
}

func (s *Server) ClearFlow(ctx context.Context, req *mgmt_pb.ClearFlowRequest) (*mgmt_pb.ClearFlowResponse, error) {
	details, err := s.command.ClearFlow(ctx, action_grpc.FlowTypeToDomain(req.Type), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
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
