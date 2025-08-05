package auth

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	policy_grpc "github.com/zitadel/zitadel/internal/api/grpc/policy"
	auth_pb "github.com/zitadel/zitadel/pkg/grpc/auth"
)

func (s *Server) GetMyLabelPolicy(ctx context.Context, _ *auth_pb.GetMyLabelPolicyRequest) (*auth_pb.GetMyLabelPolicyResponse, error) {
	policy, err := s.query.ActiveLabelPolicyByOrg(ctx, authz.GetCtxData(ctx).OrgID, false)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyLabelPolicyResponse{
		Policy: policy_grpc.ModelLabelPolicyToPb(policy, s.assetsAPIDomain(ctx)),
	}, nil
}

func (s *Server) GetMyPrivacyPolicy(ctx context.Context, _ *auth_pb.GetMyPrivacyPolicyRequest) (*auth_pb.GetMyPrivacyPolicyResponse, error) {
	policy, err := s.query.PrivacyPolicyByOrg(ctx, true, authz.GetCtxData(ctx).OrgID, false)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyPrivacyPolicyResponse{
		Policy: policy_grpc.ModelPrivacyPolicyToPb(policy),
	}, nil
}

func (s *Server) GetMyLoginPolicy(ctx context.Context, _ *auth_pb.GetMyLoginPolicyRequest) (*auth_pb.GetMyLoginPolicyResponse, error) {
	policy, err := s.query.LoginPolicyByID(ctx, true, authz.GetCtxData(ctx).OrgID, false)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyLoginPolicyResponse{
		Policy: policy_grpc.ModelLoginPolicyToPb(policy),
	}, nil
}
