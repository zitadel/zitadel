package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) AddCustomLabelPolicy(ctx context.Context, req *mgmt_pb.AddCustomLabelPolicyRequest) (*mgmt_pb.AddCustomLabelPolicyResponse, error) {
	policy, err := s.command.AddLabelPolicy(ctx, authz.GetCtxData(ctx).OrgID, addLabelPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddCustomLabelPolicyResponse{
		Details: object.AddToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateCustomLabelPolicy(ctx context.Context, req *mgmt_pb.UpdateCustomLabelPolicyRequest) (*mgmt_pb.UpdateCustomLabelPolicyResponse, error) {
	policy, err := s.command.ChangeLabelPolicy(ctx, authz.GetCtxData(ctx).OrgID, updateLabelPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateCustomLabelPolicyResponse{
		Details: object.ChangeToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) ActivateCustomLabelPolicy(ctx context.Context, req *mgmt_pb.ActivateCustomLabelPolicyRequest) (*mgmt_pb.ActivateCustomLabelPolicyResponse, error) {
	policy, err := s.command.ActivateLabelPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ActivateCustomLabelPolicyResponse{
		Details: object.ChangeToDetailsPb(
			policy.Sequence,
			policy.EventDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetLabelPolicyToDefault(ctx context.Context, req *mgmt_pb.ResetLabelPolicyToDefaultRequest) (*mgmt_pb.ResetLabelPolicyToDefaultResponse, error) {
	objectDetails, err := s.command.RemoveLabelPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetLabelPolicyToDefaultResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveCustomLabelPolicyLogo(ctx context.Context, req *mgmt_pb.RemoveCustomLabelPolicyLogoRequest) (*mgmt_pb.RemoveCustomLabelPolicyLogoResponse, error) {
	policy, err := s.command.RemoveLogoLabelPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveCustomLabelPolicyLogoResponse{
		Details: object.ChangeToDetailsPb(
			policy.Sequence,
			policy.EventDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveCustomLabelPolicyLogoDark(ctx context.Context, req *mgmt_pb.RemoveCustomLabelPolicyLogoDarkRequest) (*mgmt_pb.RemoveCustomLabelPolicyLogoDarkResponse, error) {
	policy, err := s.command.RemoveLogoDarkLabelPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveCustomLabelPolicyLogoDarkResponse{
		Details: object.ChangeToDetailsPb(
			policy.Sequence,
			policy.EventDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveCustomLabelPolicyIcon(ctx context.Context, req *mgmt_pb.RemoveCustomLabelPolicyIconRequest) (*mgmt_pb.RemoveCustomLabelPolicyIconResponse, error) {
	policy, err := s.command.RemoveIconLabelPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveCustomLabelPolicyIconResponse{
		Details: object.ChangeToDetailsPb(
			policy.Sequence,
			policy.EventDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveCustomLabelPolicyIconDark(ctx context.Context, req *mgmt_pb.RemoveCustomLabelPolicyIconDarkRequest) (*mgmt_pb.RemoveCustomLabelPolicyIconDarkResponse, error) {
	policy, err := s.command.RemoveIconDarkLabelPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveCustomLabelPolicyIconDarkResponse{
		Details: object.ChangeToDetailsPb(
			policy.Sequence,
			policy.EventDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveCustomLabelPolicyFont(ctx context.Context, req *mgmt_pb.RemoveCustomLabelPolicyFontRequest) (*mgmt_pb.RemoveCustomLabelPolicyFontResponse, error) {
	policy, err := s.command.RemoveFontLabelPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveCustomLabelPolicyFontResponse{
		Details: object.ChangeToDetailsPb(
			policy.Sequence,
			policy.EventDate,
			policy.ResourceOwner,
		),
	}, nil
}
