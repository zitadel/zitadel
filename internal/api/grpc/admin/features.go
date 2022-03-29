package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	features_grpc "github.com/caos/zitadel/internal/api/grpc/features"
	object_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetDefaultFeatures(ctx context.Context, _ *admin_pb.GetDefaultFeaturesRequest) (*admin_pb.GetDefaultFeaturesResponse, error) {
	features, err := s.query.DefaultFeatures(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultFeaturesResponse{
		Features: features_grpc.ModelFeaturesToPb(features),
	}, nil
}

func (s *Server) SetDefaultFeatures(ctx context.Context, req *admin_pb.SetDefaultFeaturesRequest) (*admin_pb.SetDefaultFeaturesResponse, error) {
	details, err := s.command.SetDefaultFeatures(ctx, authz.GetInstance(ctx).ID, setDefaultFeaturesRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetDefaultFeaturesResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) GetOrgFeatures(ctx context.Context, req *admin_pb.GetOrgFeaturesRequest) (*admin_pb.GetOrgFeaturesResponse, error) {
	features, err := s.query.FeaturesByOrgID(ctx, req.OrgId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetOrgFeaturesResponse{
		Features: features_grpc.ModelFeaturesToPb(features),
	}, nil
}

func (s *Server) SetOrgFeatures(ctx context.Context, req *admin_pb.SetOrgFeaturesRequest) (*admin_pb.SetOrgFeaturesResponse, error) {
	details, err := s.command.SetOrgFeatures(ctx, authz.GetInstance(ctx).ID, req.OrgId, setOrgFeaturesRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetOrgFeaturesResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ResetOrgFeatures(ctx context.Context, req *admin_pb.ResetOrgFeaturesRequest) (*admin_pb.ResetOrgFeaturesResponse, error) {
	details, err := s.command.RemoveOrgFeatures(ctx, authz.GetInstance(ctx).ID, req.OrgId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ResetOrgFeaturesResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func setDefaultFeaturesRequestToDomain(req *admin_pb.SetDefaultFeaturesRequest) *domain.Features {
	actionsAllowed := features_grpc.ActionsAllowedToDomain(req.ActionsAllowed)
	if req.Actions {
		actionsAllowed = domain.ActionsAllowedUnlimited
	}
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
		LabelPolicyPrivateLabel:  req.LabelPolicy || req.LabelPolicyPrivateLabel,
		LabelPolicyWatermark:     req.LabelPolicyWatermark,
		CustomDomain:             req.CustomDomain,
		PrivacyPolicy:            req.PrivacyPolicy,
		MetadataUser:             req.MetadataUser,
		CustomTextLogin:          req.CustomTextLogin || req.CustomText,
		CustomTextMessage:        req.CustomTextMessage,
		LockoutPolicy:            req.LockoutPolicy,
		ActionsAllowed:           actionsAllowed,
		MaxActions:               int(req.MaxActions),
	}
}

func setOrgFeaturesRequestToDomain(req *admin_pb.SetOrgFeaturesRequest) *domain.Features {
	actionsAllowed := features_grpc.ActionsAllowedToDomain(req.ActionsAllowed)
	if req.Actions {
		actionsAllowed = domain.ActionsAllowedUnlimited
	}
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
		LabelPolicyPrivateLabel:  req.LabelPolicy || req.LabelPolicyPrivateLabel,
		LabelPolicyWatermark:     req.LabelPolicyWatermark,
		CustomDomain:             req.CustomDomain,
		PrivacyPolicy:            req.PrivacyPolicy,
		MetadataUser:             req.MetadataUser,
		CustomTextLogin:          req.CustomTextLogin || req.CustomText,
		CustomTextMessage:        req.CustomTextMessage,
		LockoutPolicy:            req.LockoutPolicy,
		ActionsAllowed:           actionsAllowed,
		MaxActions:               int(req.MaxActions),
	}
}
