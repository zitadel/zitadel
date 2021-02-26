package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/api/grpc/policy"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	"github.com/caos/zitadel/internal/domain"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetDefaultLoginPolicy(ctx context.Context, _ *admin_pb.GetDefaultLoginPolicyRequest) (*admin_pb.GetDefaultLoginPolicyResponse, error) {
	policy, err := s.iam.GetDefaultLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultLoginPolicyResponse{Policy: policy_grpc.ModelLoginPolicyToPb(policy)}, nil
}

func (s *Server) UpdateDefaultLoginPolicy(ctx context.Context, p *admin_pb.UpdateDefaultLoginPolicyRequest) (*admin_pb.UpdateDefaultLoginPolicyResponse, error) {
	policy, err := s.command.ChangeDefaultLoginPolicy(ctx, updateDefaultLoginPolicyToDomain(p))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateDefaultLoginPolicyResponse{
		Details: object.ToDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) AddIDPToDefaultLoginPolicy(ctx context.Context, req *admin_pb.AddIDPToDefaultLoginPolicyRequest) (*admin_pb.AddIDPToDefaultLoginPolicyResponse, error) {
	idp, err := s.command.AddIDPProviderToDefaultLoginPolicy(ctx, &domain.IDPProvider{IDPConfigID: req.IdpId}) //TODO: old way was to also add type but this doesnt make sense in my point of view
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddIDPToDefaultLoginPolicyResponse{
		Details: object.ToDetailsPb(
			idp.Sequence,
			idp.CreationDate,
			idp.ChangeDate, idp.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveIDPFromDefaultLoginPolicy(ctx context.Context, req *admin_pb.RemoveIDPFromDefaultLoginPolicyRequest) (*admin_pb.RemoveIDPFromDefaultLoginPolicyResponse, error) {
	//TODO: dont understand current impelementation
	return nil, nil
}

func (s *Server) AddSecondFactorToDefaultLoginPolicy(ctx context.Context, req *admin_pb.AddSecondFactorToDefaultLoginPolicyRequest) (*admin_pb.AddSecondFactorToDefaultLoginPolicyResponse, error) {
	result, err := s.command.AddSecondFactorToDefaultLoginPolicy(ctx, policy.SecondFactorTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	//TODO: return value
	_ = result
	// return &admin_pb.AddSecondFactorToDefaultLoginPolicyResponse{
	// 	Details: object.ToDetailsPb(
	// 		result.Sequence,
	// 		result.CreationDate,
	// 		result.ChangeDate,
	// 		result.ResourceOwner,
	// 	),
	// }, nil
	return nil, nil
}

func (s *Server) RemoveSecondFactorFromDefaultLoginPolicy(ctx context.Context, req *admin_pb.RemoveSecondFactorFromDefaultLoginPolicyRequest) (*admin_pb.RemoveSecondFactorFromDefaultLoginPolicyResponse, error) {
	err := s.command.RemoveSecondFactorFromDefaultLoginPolicy(ctx, policy.SecondFactorTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	//TODO: missing return value
	return &admin_pb.RemoveSecondFactorFromDefaultLoginPolicyResponse{}, nil
}

func (s *Server) AddMultiFactorToDefaultLoginPolicy(ctx context.Context, req *admin_pb.AddMultiFactorToDefaultLoginPolicyRequest) (*admin_pb.AddMultiFactorToDefaultLoginPolicyResponse, error) {
	result, err := s.command.AddMultiFactorToDefaultLoginPolicy(ctx, policy_grpc.MultiFactorTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	//TODO: return value
	_ = result
	// return &admin_pb.AddMultiFactorToDefaultLoginPolicyResponse{
	// 	Details: object.ToDetailsPb(
	// 		result.Sequence,
	// 		result.CreationDate,
	// 		result.ChangeDate,
	// 		result.ResourceOwner,
	// 	),
	// }, nil
	return nil, nil
}

func (s *Server) RemoveMultiFactorFromDefaultLoginPolicy(ctx context.Context, req *admin_pb.RemoveMultiFactorFromDefaultLoginPolicyRequest) (*admin_pb.RemoveMultiFactorFromDefaultLoginPolicyResponse, error) {
	err := s.command.RemoveMultiFactorFromDefaultLoginPolicy(ctx, policy.MultiFactorTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	//TODO: missing return value
	return &admin_pb.RemoveMultiFactorFromDefaultLoginPolicyResponse{}, nil
}
