package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/zitadel/zitadel/internal/api/grpc/policy"
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
	policy, err := s.query.DomainPolicyByOrg(ctx, true, req.OrgId, false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomDomainPolicyResponse{Policy: policy_grpc.DomainPolicyToPb(policy)}, nil
}

func (s *Server) AddCustomDomainPolicy(ctx context.Context, req *admin_pb.AddCustomDomainPolicyRequest) (*admin_pb.AddCustomDomainPolicyResponse, error) {
	details, err := s.command.AddOrgDomainPolicy(ctx, req.OrgId, req.UserLoginMustBeDomain, req.ValidateOrgDomains, req.SmtpSenderAddressMatchesInstanceDomain)
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddCustomDomainPolicyResponse{
		Details: object.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateDomainPolicy(ctx context.Context, req *admin_pb.UpdateDomainPolicyRequest) (*admin_pb.UpdateDomainPolicyResponse, error) {
	details, err := s.command.ChangeDefaultDomainPolicy(ctx, req.UserLoginMustBeDomain, req.ValidateOrgDomains, req.SmtpSenderAddressMatchesInstanceDomain)
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateDomainPolicyResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) UpdateCustomDomainPolicy(ctx context.Context, req *admin_pb.UpdateCustomDomainPolicyRequest) (*admin_pb.UpdateCustomDomainPolicyResponse, error) {
	details, err := s.command.ChangeOrgDomainPolicy(ctx, req.OrgId, req.UserLoginMustBeDomain, req.ValidateOrgDomains, req.SmtpSenderAddressMatchesInstanceDomain)
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateCustomDomainPolicyResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ResetCustomDomainPolicyToDefault(ctx context.Context, req *admin_pb.ResetCustomDomainPolicyToDefaultRequest) (*admin_pb.ResetCustomDomainPolicyToDefaultResponse, error) {
	details, err := s.command.RemoveOrgDomainPolicy(ctx, req.OrgId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ResetCustomDomainPolicyToDefaultResponse{Details: object.DomainToChangeDetailsPb(details)}, nil
}

// the following requests only exist for backwards compatibility
// OrgIAMPolicy has been replaced by DomainPolicy, which also extends it with validateOrgDomains and smtpSenderAddressMatchesInstanceDomain
// Add and Update requests will therefore set the previous default (true)

func (s *Server) AddCustomOrgIAMPolicy(ctx context.Context, req *admin_pb.AddCustomOrgIAMPolicyRequest) (*admin_pb.AddCustomOrgIAMPolicyResponse, error) {
	details, err := s.command.AddOrgDomainPolicy(ctx, req.OrgId, req.UserLoginMustBeDomain, true, true)
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddCustomOrgIAMPolicyResponse{
		Details: object.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateOrgIAMPolicy(ctx context.Context, req *admin_pb.UpdateOrgIAMPolicyRequest) (*admin_pb.UpdateOrgIAMPolicyResponse, error) {
	details, err := s.command.ChangeDefaultDomainPolicy(ctx, req.UserLoginMustBeDomain, true, true)
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateOrgIAMPolicyResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) UpdateCustomOrgIAMPolicy(ctx context.Context, req *admin_pb.UpdateCustomOrgIAMPolicyRequest) (*admin_pb.UpdateCustomOrgIAMPolicyResponse, error) {
	details, err := s.command.ChangeOrgDomainPolicy(ctx, req.OrgId, req.UserLoginMustBeDomain, true, true)
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateCustomOrgIAMPolicyResponse{
		Details: object.DomainToChangeDetailsPb(details),
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
	policy, err := s.query.DomainPolicyByOrg(ctx, true, req.OrgId, false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomOrgIAMPolicyResponse{Policy: policy_grpc.DomainPolicyToOrgIAMPb(policy)}, nil
}

func (s *Server) ResetCustomOrgIAMPolicyToDefault(ctx context.Context, req *admin_pb.ResetCustomOrgIAMPolicyToDefaultRequest) (*admin_pb.ResetCustomOrgIAMPolicyToDefaultResponse, error) {
	details, err := s.command.RemoveOrgDomainPolicy(ctx, req.OrgId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ResetCustomOrgIAMPolicyToDefaultResponse{Details: object.DomainToChangeDetailsPb(details)}, nil
}
