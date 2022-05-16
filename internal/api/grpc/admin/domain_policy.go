package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/zitadel/zitadel/internal/api/grpc/policy"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) GetDomainPolicy(ctx context.Context, _ *admin_pb.GetDomainPolicyRequest) (*admin_pb.GetDomainPolicyResponse, error) {
	policy, err := s.query.DefaultDomainPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDomainPolicyResponse{Policy: policy_grpc.DomainPolicyToPb(policy)}, nil
}

func (s *Server) GetCustomDomainPolicy(ctx context.Context, req *admin_pb.GetCustomDomainPolicyRequest) (*admin_pb.GetCustomDomainPolicyResponse, error) {
	policy, err := s.query.DomainPolicyByOrg(ctx, req.OrgId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomDomainPolicyResponse{Policy: policy_grpc.DomainPolicyToPb(policy)}, nil
}

func (s *Server) AddCustomDomainPolicy(ctx context.Context, req *admin_pb.AddCustomDomainPolicyRequest) (*admin_pb.AddCustomDomainPolicyResponse, error) {
	policy, err := s.command.AddOrgDomainPolicy(ctx, req.OrgId, domainPolicyToDomain(req.UserLoginMustBeDomain, req.ValidateOrgDomains, req.SmtpSenderAddressMatchesInstanceDomain))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddCustomDomainPolicyResponse{
		Details: object.AddToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateDomainPolicy(ctx context.Context, req *admin_pb.UpdateDomainPolicyRequest) (*admin_pb.UpdateDomainPolicyResponse, error) {
	config, err := s.command.ChangeDefaultDomainPolicy(ctx, updateDomainPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateDomainPolicyResponse{
		Details: object.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateCustomDomainPolicy(ctx context.Context, req *admin_pb.UpdateCustomDomainPolicyRequest) (*admin_pb.UpdateCustomDomainPolicyResponse, error) {
	config, err := s.command.ChangeOrgDomainPolicy(ctx, req.OrgId, updateCustomDomainPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateCustomDomainPolicyResponse{
		Details: object.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetCustomDomainPolicyTo(ctx context.Context, req *admin_pb.ResetCustomDomainPolicyToDefaultRequest) (*admin_pb.ResetCustomDomainPolicyToDefaultResponse, error) {
	err := s.command.RemoveOrgDomainPolicy(ctx, req.OrgId)
	if err != nil {
		return nil, err
	}
	return nil, nil //TOOD: return data
}

func domainPolicyToDomain(userLoginMustBeDomain, validateOrgDomains, smtpSenderAddressMatchesInstanceDomain bool) *domain.DomainPolicy {
	return &domain.DomainPolicy{
		UserLoginMustBeDomain:                  userLoginMustBeDomain,
		ValidateOrgDomains:                     validateOrgDomains,
		SMTPSenderAddressMatchesInstanceDomain: smtpSenderAddressMatchesInstanceDomain,
	}
}

func updateDomainPolicyToDomain(req *admin_pb.UpdateDomainPolicyRequest) *domain.DomainPolicy {
	return &domain.DomainPolicy{
		UserLoginMustBeDomain:                  req.UserLoginMustBeDomain,
		ValidateOrgDomains:                     req.ValidateOrgDomains,
		SMTPSenderAddressMatchesInstanceDomain: req.SmtpSenderAddressMatchesInstanceDomain,
	}
}

func updateCustomDomainPolicyToDomain(req *admin_pb.UpdateCustomDomainPolicyRequest) *domain.DomainPolicy {
	return &domain.DomainPolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.OrgId,
		},
		UserLoginMustBeDomain:                  req.UserLoginMustBeDomain,
		ValidateOrgDomains:                     req.ValidateOrgDomains,
		SMTPSenderAddressMatchesInstanceDomain: req.SmtpSenderAddressMatchesInstanceDomain,
	}
}

func (s *Server) AddCustomOrgIAMPolicy(ctx context.Context, req *admin_pb.AddCustomOrgIAMPolicyRequest) (*admin_pb.AddCustomOrgIAMPolicyResponse, error) {
	policy, err := s.command.AddOrgDomainPolicy(ctx, req.OrgId, domainPolicyToDomain(req.UserLoginMustBeDomain, true, true))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddCustomOrgIAMPolicyResponse{
		Details: object.AddToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateOrgIAMPolicy(ctx context.Context, req *admin_pb.UpdateOrgIAMPolicyRequest) (*admin_pb.UpdateOrgIAMPolicyResponse, error) {
	config, err := s.command.ChangeDefaultDomainPolicy(ctx, updateOrgIAMPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateOrgIAMPolicyResponse{
		Details: object.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateCustomOrgIAMPolicy(ctx context.Context, req *admin_pb.UpdateCustomOrgIAMPolicyRequest) (*admin_pb.UpdateCustomOrgIAMPolicyResponse, error) {
	config, err := s.command.ChangeOrgDomainPolicy(ctx, req.OrgId, updateCustomOrgIAMPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateCustomOrgIAMPolicyResponse{
		Details: object.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetOrgIAMPolicy(ctx context.Context, _ *admin_pb.GetOrgIAMPolicyRequest) (*admin_pb.GetOrgIAMPolicyResponse, error) {
	policy, err := s.query.DefaultDomainPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetOrgIAMPolicyResponse{Policy: policy_grpc.DomainPolicyToOrgIAMPb(policy)}, nil
}

func (s *Server) GetCustomOrgIAMPolicy(ctx context.Context, req *admin_pb.GetCustomOrgIAMPolicyRequest) (*admin_pb.GetCustomOrgIAMPolicyResponse, error) {
	policy, err := s.query.DomainPolicyByOrg(ctx, req.OrgId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomOrgIAMPolicyResponse{Policy: policy_grpc.DomainPolicyToOrgIAMPb(policy)}, nil
}

func updateOrgIAMPolicyToDomain(req *admin_pb.UpdateOrgIAMPolicyRequest) *domain.DomainPolicy {
	return &domain.DomainPolicy{
		UserLoginMustBeDomain:                  req.UserLoginMustBeDomain,
		ValidateOrgDomains:                     true,
		SMTPSenderAddressMatchesInstanceDomain: true,
	}
}

func updateCustomOrgIAMPolicyToDomain(req *admin_pb.UpdateCustomOrgIAMPolicyRequest) *domain.DomainPolicy {
	return &domain.DomainPolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.OrgId,
		},
		UserLoginMustBeDomain:                  req.UserLoginMustBeDomain,
		ValidateOrgDomains:                     true,
		SMTPSenderAddressMatchesInstanceDomain: true,
	}
}
