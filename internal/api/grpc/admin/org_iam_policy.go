package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetOrgIAMPolicy(ctx context.Context, _ *admin_pb.GetOrgIAMPolicyRequest) (*admin_pb.GetOrgIAMPolicyResponse, error) {
	policy, err := s.iam.GetDefaultOrgIAMPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetOrgIAMPolicyResponse{Policy: policy_grpc.OrgIAMPolicyToPb(policy)}, nil
}

func (s *Server) GetCustomOrgIAMPolicy(ctx context.Context, req *admin_pb.GetCustomOrgIAMPolicyRequest) (*admin_pb.GetCustomOrgIAMPolicyResponse, error) {
	policy, err := s.org.GetOrgIAMPolicyByID(ctx, req.OrgId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetCustomOrgIAMPolicyResponse{Policy: policy_grpc.OrgIAMPolicyToPb(policy)}, nil
}

func (s *Server) AddCustomOrgIAMPolicy(ctx context.Context, req *admin_pb.AddCustomOrgIAMPolicyRequest) (*admin_pb.AddCustomOrgIAMPolicyResponse, error) {
	policy, err := s.command.AddOrgIAMPolicy(ctx, req.OrgId, toDomainOrgIAMPolicy(req.UserLoginMustBeDomain))
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
	config, err := s.command.ChangeDefaultOrgIAMPolicy(ctx, updateOrgIAMPolicyToDomain(req))
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
	config, err := s.command.ChangeOrgIAMPolicy(ctx, req.OrgId, updateCustomOrgIAMPolicyToDomain(req))
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

func (s *Server) ResetCustomOrgIAMPolicyTo(ctx context.Context, req *admin_pb.ResetCustomOrgIAMPolicyToDefaultRequest) (*admin_pb.ResetCustomOrgIAMPolicyToDefaultResponse, error) {
	err := s.command.RemoveOrgIAMPolicy(ctx, req.OrgId)
	if err != nil {
		return nil, err
	}
	return nil, nil //TOOD: return data
}

func toDomainOrgIAMPolicy(userLoginMustBeDomain bool) *domain.OrgIAMPolicy {
	return &domain.OrgIAMPolicy{
		UserLoginMustBeDomain: userLoginMustBeDomain,
	}
}

func updateOrgIAMPolicyToDomain(req *admin_pb.UpdateOrgIAMPolicyRequest) *domain.OrgIAMPolicy {
	return &domain.OrgIAMPolicy{
		// ObjectRoot: models.ObjectRoot{
		// 	// AggreagateID: //TODO: there should only be ONE default
		// },
		UserLoginMustBeDomain: req.UserLoginMustBeDomain,
	}
}

func updateCustomOrgIAMPolicyToDomain(req *admin_pb.UpdateCustomOrgIAMPolicyRequest) *domain.OrgIAMPolicy {
	return &domain.OrgIAMPolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.OrgId,
		},
		UserLoginMustBeDomain: req.UserLoginMustBeDomain,
	}
}
