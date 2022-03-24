package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
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
	policy, err := s.command.AddOrgDomainPolicy(ctx, req.OrgId, domainPolicyToDomain(req.UserLoginMustBeDomain))
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

func domainPolicyToDomain(userLoginMustBeDomain bool) *domain.DomainPolicy {
	return &domain.DomainPolicy{
		UserLoginMustBeDomain: userLoginMustBeDomain,
	}
}

func updateDomainPolicyToDomain(req *admin_pb.UpdateDomainPolicyRequest) *domain.DomainPolicy {
	return &domain.DomainPolicy{
		// ObjectRoot: models.ObjectRoot{
		// 	// AggreagateID: //TODO: there should only be ONE default
		// },
		UserLoginMustBeDomain: req.UserLoginMustBeDomain,
	}
}

func updateCustomDomainPolicyToDomain(req *admin_pb.UpdateCustomDomainPolicyRequest) *domain.DomainPolicy {
	return &domain.DomainPolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.OrgId,
		},
		UserLoginMustBeDomain: req.UserLoginMustBeDomain,
	}
}
