package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetLabelPolicy(ctx context.Context, req *admin_pb.GetLabelPolicyRequest) (*admin_pb.GetLabelPolicyResponse, error) {
	policy, err := s.iam.GetDefaultLabelPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetLabelPolicyResponse{Policy: policy_grpc.ModelLabelPolicyToPb(policy)}, nil
}

func (s *Server) GetPreviewLabelPolicy(ctx context.Context, req *admin_pb.GetPreviewLabelPolicyRequest) (*admin_pb.GetPreviewLabelPolicyResponse, error) {
	policy, err := s.iam.GetDefaultPreviewLabelPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetPreviewLabelPolicyResponse{Policy: policy_grpc.ModelLabelPolicyToPb(policy)}, nil
}

func (s *Server) UpdateLabelPolicy(ctx context.Context, req *admin_pb.UpdateLabelPolicyRequest) (*admin_pb.UpdateLabelPolicyResponse, error) {
	policy, err := s.command.ChangeDefaultLabelPolicy(ctx, updateLabelPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateLabelPolicyResponse{
		Details: object.ChangeToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) ActivateLabelPolicy(ctx context.Context, req *admin_pb.ActivateLabelPolicyRequest) (*admin_pb.ActivateLabelPolicyResponse, error) {
	policy, err := s.command.ActivateDefaultLabelPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ActivateLabelPolicyResponse{
		Details: object.ChangeToDetailsPb(
			policy.Sequence,
			policy.EventDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveLabelPolicyLogo(ctx context.Context, req *admin_pb.RemoveLabelPolicyLogoRequest) (*admin_pb.RemoveLabelPolicyLogoResponse, error) {
	return nil, nil
}

func (s *Server) RemoveLabelPolicyLogoDark(ctx context.Context, req *admin_pb.RemoveLabelPolicyLogoDarkRequest) (*admin_pb.RemoveLabelPolicyLogoDarkResponse, error) {
	return nil, nil
}

func (s *Server) RemoveLabelPolicyIcon(ctx context.Context, req *admin_pb.RemoveLabelPolicyIconRequest) (*admin_pb.RemoveLabelPolicyIconResponse, error) {
	return nil, nil
}

func (s *Server) RemoveLabelPolicyIconDark(ctx context.Context, req *admin_pb.RemoveLabelPolicyIconDarkRequest) (*admin_pb.RemoveLabelPolicyIconDarkResponse, error) {
	return nil, nil
}

func (s *Server) RemoveLabelPolicyFont(ctx context.Context, req *admin_pb.RemoveLabelPolicyFontRequest) (*admin_pb.RemoveLabelPolicyFontResponse, error) {
	return nil, nil
}
