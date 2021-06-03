package admin

import (
	"context"

	features_grpc "github.com/caos/zitadel/internal/api/grpc/features"
	object_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetDefaultFeatures(ctx context.Context, _ *admin_pb.GetDefaultFeaturesRequest) (*admin_pb.GetDefaultFeaturesResponse, error) {
	features, err := s.features.GetDefaultFeatures(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultFeaturesResponse{
		Features: features_grpc.FeaturesFromModel(features),
	}, nil
}

func (s *Server) SetDefaultFeatures(ctx context.Context, req *admin_pb.SetDefaultFeaturesRequest) (*admin_pb.SetDefaultFeaturesResponse, error) {
	details, err := s.command.SetDefaultFeatures(ctx, setDefaultFeaturesRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetDefaultFeaturesResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) GetOrgFeatures(ctx context.Context, req *admin_pb.GetOrgFeaturesRequest) (*admin_pb.GetOrgFeaturesResponse, error) {
	features, err := s.features.GetOrgFeatures(ctx, req.OrgId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetOrgFeaturesResponse{
		Features: features_grpc.FeaturesFromModel(features),
	}, nil
}

func (s *Server) SetOrgFeatures(ctx context.Context, req *admin_pb.SetOrgFeaturesRequest) (*admin_pb.SetOrgFeaturesResponse, error) {
	details, err := s.command.SetOrgFeatures(ctx, req.OrgId, setOrgFeaturesRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetOrgFeaturesResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ResetOrgFeatures(ctx context.Context, req *admin_pb.ResetOrgFeaturesRequest) (*admin_pb.ResetOrgFeaturesResponse, error) {
	details, err := s.command.RemoveOrgFeatures(ctx, req.OrgId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ResetOrgFeaturesResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func setDefaultFeaturesRequestToDomain(req *admin_pb.SetDefaultFeaturesRequest) *domain.Features {
	return &domain.Features{
		TierName:                 req.TierName,
		TierDescription:          req.Description,
		AuditLogRetention:        req.AuditLogRetention.AsDuration(),
		LoginPolicyFactors:       req.LoginPolicyFactors,
		LoginPolicyIDP:           req.LoginPolicyIdp,
		LoginPolicyPasswordless:  req.LoginPolicyPasswordless,
		LoginPolicyRegistration:  req.LoginPolicyRegistration,
		LoginPolicyUsernameLogin: req.LoginPolicyUsernameLogin,
		LoginPolicyPasswordReset: req.LoginPolicyPasswordReset,
		PasswordComplexityPolicy: req.PasswordComplexityPolicy,
		LabelPolicy:              req.LabelPolicy,
		CustomDomain:             req.CustomDomain,
	}
}

func setOrgFeaturesRequestToDomain(req *admin_pb.SetOrgFeaturesRequest) *domain.Features {
	return &domain.Features{
		TierName:                 req.TierName,
		TierDescription:          req.Description,
		State:                    features_grpc.FeaturesStateToDomain(req.State),
		StateDescription:         req.StateDescription,
		AuditLogRetention:        req.AuditLogRetention.AsDuration(),
		LoginPolicyFactors:       req.LoginPolicyFactors,
		LoginPolicyIDP:           req.LoginPolicyIdp,
		LoginPolicyPasswordless:  req.LoginPolicyPasswordless,
		LoginPolicyRegistration:  req.LoginPolicyRegistration,
		LoginPolicyUsernameLogin: req.LoginPolicyUsernameLogin,
		LoginPolicyPasswordReset: req.LoginPolicyPasswordReset,
		PasswordComplexityPolicy: req.PasswordComplexityPolicy,
		LabelPolicy:              req.LabelPolicy,
		CustomDomain:             req.CustomDomain,
	}
}
