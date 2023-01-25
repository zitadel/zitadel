package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/zitadel/zitadel/internal/api/grpc/policy"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Server) GetNotificationPolicy(ctx context.Context, _ *mgmt_pb.GetNotificationPolicyRequest) (*mgmt_pb.GetNotificationPolicyResponse, error) {
	policy, err := s.query.NotificationPolicyByOrg(ctx, true, authz.GetCtxData(ctx).OrgID, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetNotificationPolicyResponse{Policy: policy_grpc.ModelNotificationPolicyToPb(policy)}, nil
}

func (s *Server) GetDefaultNotificationPolicy(ctx context.Context, _ *mgmt_pb.GetDefaultNotificationPolicyRequest) (*mgmt_pb.GetDefaultNotificationPolicyResponse, error) {
	policy, err := s.query.DefaultNotificationPolicy(ctx, true)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultNotificationPolicyResponse{Policy: policy_grpc.ModelNotificationPolicyToPb(policy)}, nil
}

func (s *Server) AddCustomNotificationPolicy(ctx context.Context, req *mgmt_pb.AddCustomNotificationPolicyRequest) (*mgmt_pb.AddCustomNotificationPolicyResponse, error) {
	result, err := s.command.AddNotificationPolicy(ctx, authz.GetCtxData(ctx).OrgID, req.GetPasswordChange())
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddCustomNotificationPolicyResponse{
		Details: object.AddToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateCustomNotificationPolicy(ctx context.Context, req *mgmt_pb.UpdateCustomNotificationPolicyRequest) (*mgmt_pb.UpdateCustomNotificationPolicyResponse, error) {
	result, err := s.command.ChangeNotificationPolicy(ctx, authz.GetCtxData(ctx).OrgID, req.GetPasswordChange())
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateCustomNotificationPolicyResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.EventDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetNotificationPolicyToDefault(ctx context.Context, _ *mgmt_pb.ResetNotificationPolicyToDefaultRequest) (*mgmt_pb.ResetNotificationPolicyToDefaultResponse, error) {
	objectDetails, err := s.command.RemoveNotificationPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetNotificationPolicyToDefaultResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}
