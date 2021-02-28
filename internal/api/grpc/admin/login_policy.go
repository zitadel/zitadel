package admin

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/api/grpc/idp"
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/api/grpc/policy"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	"github.com/caos/zitadel/internal/domain"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetLoginPolicy(ctx context.Context, _ *admin_pb.GetLoginPolicyRequest) (*admin_pb.GetLoginPolicyResponse, error) {
	policy, err := s.iam.GetDefaultLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetLoginPolicyResponse{Policy: policy_grpc.ModelLoginPolicyToPb(policy)}, nil
}

func (s *Server) UpdateLoginPolicy(ctx context.Context, p *admin_pb.UpdateLoginPolicyRequest) (*admin_pb.UpdateLoginPolicyResponse, error) {
	policy, err := s.command.ChangeDefaultLoginPolicy(ctx, updateLoginPolicyToDomain(p))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateLoginPolicyResponse{
		Details: object.ToDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) ListLoginPolicyIDPs(ctx context.Context, req *admin_pb.ListLoginPolicyIDPsRequest) (*admin_pb.ListLoginPolicyIDPsResponse, error) {
	res, err := s.iam.SearchDefaultIDPProviders(ctx, ListLoginPolicyIDPsRequestToModel(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListLoginPolicyIDPsResponse{
		Result:   idp.ExternalIDPViewsToLoginPolicyLinkPb(res.Result),
		MetaData: object.ToListDetails(res.TotalResult, res.Sequence, res.Timestamp),
	}, nil
}

func (s *Server) AddIDPToLoginPolicy(ctx context.Context, req *admin_pb.AddIDPToLoginPolicyRequest) (*admin_pb.AddIDPToLoginPolicyResponse, error) {
	idp, err := s.command.AddIDPProviderToDefaultLoginPolicy(ctx, &domain.IDPProvider{IDPConfigID: req.IdpId}) //TODO: old way was to also add type but this doesnt make sense in my point of view
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddIDPToLoginPolicyResponse{
		Details: object.ToDetailsPb(
			idp.Sequence,
			idp.CreationDate,
			idp.ChangeDate, idp.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveIDPFromLoginPolicy(ctx context.Context, req *admin_pb.RemoveIDPFromLoginPolicyRequest) (*admin_pb.RemoveIDPFromLoginPolicyResponse, error) {
	//TODO: dont understand current impelementation
	panic("implement me")
}

func (s *Server) ListLoginPolicySecondFactors(ctx context.Context, req *admin_pb.ListLoginPolicySecondFactorsRequest) (*admin_pb.ListLoginPolicySecondFactorsResponse, error) {
	result, err := s.iam.SearchDefaultSecondFactors(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListLoginPolicySecondFactorsResponse{
		//TODO: missing values from res
		MetaData: object.ToListDetails(result.TotalResult, 0, time.Time{}),
		Result:   policy.ModelSecondFactorTypesToPb(result.Result),
	}, nil
}

func (s *Server) AddSecondFactorToLoginPolicy(ctx context.Context, req *admin_pb.AddSecondFactorToLoginPolicyRequest) (*admin_pb.AddSecondFactorToLoginPolicyResponse, error) {
	result, err := s.command.AddSecondFactorToDefaultLoginPolicy(ctx, policy.SecondFactorTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	//TODO: details from be
	_ = result
	return &admin_pb.AddSecondFactorToLoginPolicyResponse{
		// Details: object.ToDetailsPb(
		// 	result.Sequence,
		// 	result.CreationDate,
		// 	result.ChangeDate,
		// 	result.ResourceOwner,
		// ),
	}, nil
}

func (s *Server) RemoveSecondFactorFromLoginPolicy(ctx context.Context, req *admin_pb.RemoveSecondFactorFromLoginPolicyRequest) (*admin_pb.RemoveSecondFactorFromLoginPolicyResponse, error) {
	err := s.command.RemoveSecondFactorFromDefaultLoginPolicy(ctx, policy.SecondFactorTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	//TODO: missing return value
	return &admin_pb.RemoveSecondFactorFromLoginPolicyResponse{}, nil
}

func (s *Server) ListLoginPolicyMultiFactors(ctx context.Context, req *admin_pb.ListLoginPolicyMultiFactorsRequest) (*admin_pb.ListLoginPolicyMultiFactorsResponse, error) {
	res, err := s.iam.SearchDefaultMultiFactors(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListLoginPolicyMultiFactorsResponse{
		//TODO: additional values
		MetaData: object.ToListDetails(res.TotalResult, 0, time.Time{}),
		Result:   policy.ModelMultiFactorTypesToPb(res.Result),
	}, nil
}

func (s *Server) AddMultiFactorToLoginPolicy(ctx context.Context, req *admin_pb.AddMultiFactorToLoginPolicyRequest) (*admin_pb.AddMultiFactorToLoginPolicyResponse, error) {
	result, err := s.command.AddMultiFactorToDefaultLoginPolicy(ctx, policy_grpc.MultiFactorTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	//TODO: return value
	_ = result
	return &admin_pb.AddMultiFactorToLoginPolicyResponse{
		// 	Details: object.ToDetailsPb(
		// 		result.Sequence,
		// 		result.CreationDate,
		// 		result.ChangeDate,
		// 		result.ResourceOwner,
		// 	),
	}, nil
}

func (s *Server) RemoveMultiFactorFromLoginPolicy(ctx context.Context, req *admin_pb.RemoveMultiFactorFromLoginPolicyRequest) (*admin_pb.RemoveMultiFactorFromLoginPolicyResponse, error) {
	err := s.command.RemoveMultiFactorFromDefaultLoginPolicy(ctx, policy.MultiFactorTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	//TODO: missing return value
	return &admin_pb.RemoveMultiFactorFromLoginPolicyResponse{}, nil
}
