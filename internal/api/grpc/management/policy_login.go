package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/idp"
	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	"github.com/caos/zitadel/internal/api/grpc/user"
	"github.com/caos/zitadel/internal/domain"
	"time"

	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetLoginPolicy(ctx context.Context, req *mgmt_pb.GetLoginPolicyRequest) (*mgmt_pb.GetLoginPolicyResponse, error) {
	policy, err := s.org.GetLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetLoginPolicyResponse{Policy: policy_grpc.ModelLoginPolicyToPb(policy)}, nil
}

func (s *Server) GetDefaultLoginPolicy(ctx context.Context, req *mgmt_pb.GetDefaultLoginPolicyRequest) (*mgmt_pb.GetDefaultLoginPolicyResponse, error) {
	policy, err := s.org.GetDefaultLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultLoginPolicyResponse{Policy: policy_grpc.ModelLoginPolicyToPb(policy)}, nil
}

func (s *Server) AddCustomLoginPolicy(ctx context.Context, req *mgmt_pb.AddCustomLoginPolicyRequest) (*mgmt_pb.AddCustomLoginPolicyResponse, error) {
	policy, err := s.command.AddLoginPolicy(ctx, authz.GetCtxData(ctx).OrgID, addLoginPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddCustomLoginPolicyResponse{
		Details: object.ToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateCustomLoginPolicy(ctx context.Context, req *mgmt_pb.UpdateCustomLoginPolicyRequest) (*mgmt_pb.UpdateCustomLoginPolicyResponse, error) {
	policy, err := s.command.ChangeLoginPolicy(ctx, authz.GetCtxData(ctx).OrgID, updateLoginPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateCustomLoginPolicyResponse{
		Details: object.ToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetLoginPolicyToDefault(ctx context.Context, req *mgmt_pb.ResetLoginPolicyToDefaultRequest) (*mgmt_pb.ResetLoginPolicyToDefaultResponse, error) {
	objectDetails, err := s.command.RemoveLoginPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetLoginPolicyToDefaultResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ListLoginPolicyIDPs(ctx context.Context, req *mgmt_pb.ListLoginPolicyIDPsRequest) (*mgmt_pb.ListLoginPolicyIDPsResponse, error) {
	res, err := s.org.SearchIDPProviders(ctx, ListLoginPolicyIDPsRequestToModel(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListLoginPolicyIDPsResponse{
		Result:  idp.ExternalIDPViewsToLoginPolicyLinkPb(res.Result),
		Details: object.ToListDetails(res.TotalResult, res.Sequence, res.Timestamp),
	}, nil
}

func (s *Server) AddIDPToLoginPolicy(ctx context.Context, req *mgmt_pb.AddIDPToLoginPolicyRequest) (*mgmt_pb.AddIDPToLoginPolicyResponse, error) {
	idp, err := s.command.AddIDPProviderToLoginPolicy(ctx, authz.GetCtxData(ctx).OrgID, &domain.IDPProvider{IDPConfigID: req.IdpId}) //TODO: old way was to also add type but this doesnt make sense in my point of view
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddIDPToLoginPolicyResponse{
		Details: object.ToDetailsPb(
			idp.Sequence,
			idp.ChangeDate,
			idp.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveIDPFromLoginPolicy(ctx context.Context, req *mgmt_pb.RemoveIDPFromLoginPolicyRequest) (*mgmt_pb.RemoveIDPFromLoginPolicyResponse, error) {
	externalIDPs, err := s.user.ExternalIDPsByIDPConfigID(ctx, req.IdpId)
	if err != nil {
		return nil, err
	}
	objectDetails, err := s.command.RemoveIDPProviderFromLoginPolicy(ctx, authz.GetCtxData(ctx).OrgID, &domain.IDPProvider{IDPConfigID: req.IdpId}, user.ExternalIDPViewsToExternalIDPs(externalIDPs)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveIDPFromLoginPolicyResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ListLoginPolicySecondFactors(ctx context.Context, req *mgmt_pb.ListLoginPolicySecondFactorsRequest) (*mgmt_pb.ListLoginPolicySecondFactorsResponse, error) {
	result, err := s.org.SearchSecondFactors(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListLoginPolicySecondFactorsResponse{
		//TODO: missing values from res
		Details: object.ToListDetails(result.TotalResult, 0, time.Time{}),
		Result:  policy_grpc.ModelSecondFactorTypesToPb(result.Result),
	}, nil
}

func (s *Server) AddSecondFactorToLoginPolicy(ctx context.Context, req *mgmt_pb.AddSecondFactorToLoginPolicyRequest) (*mgmt_pb.AddSecondFactorToLoginPolicyResponse, error) {
	_, objectDetails, err := s.command.AddSecondFactorToDefaultLoginPolicy(ctx, policy_grpc.SecondFactorTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddSecondFactorToLoginPolicyResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveSecondFactorFromLoginPolicy(ctx context.Context, req *mgmt_pb.RemoveSecondFactorFromLoginPolicyRequest) (*mgmt_pb.RemoveSecondFactorFromLoginPolicyResponse, error) {
	objectDetails, err := s.command.RemoveSecondFactorFromDefaultLoginPolicy(ctx, policy_grpc.SecondFactorTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveSecondFactorFromLoginPolicyResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ListLoginPolicyMultiFactors(ctx context.Context, req *mgmt_pb.ListLoginPolicyMultiFactorsRequest) (*mgmt_pb.ListLoginPolicyMultiFactorsResponse, error) {
	res, err := s.org.SearchMultiFactors(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListLoginPolicyMultiFactorsResponse{
		//TODO: additional values
		Details: object.ToListDetails(res.TotalResult, 0, time.Time{}),
		Result:  policy_grpc.ModelMultiFactorTypesToPb(res.Result),
	}, nil
}

func (s *Server) AddMultiFactorToLoginPolicy(ctx context.Context, req *mgmt_pb.AddMultiFactorToLoginPolicyRequest) (*mgmt_pb.AddMultiFactorToLoginPolicyResponse, error) {
	_, objectDetails, err := s.command.AddMultiFactorToDefaultLoginPolicy(ctx, policy_grpc.MultiFactorTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddMultiFactorToLoginPolicyResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveMultiFactorFromLoginPolicy(ctx context.Context, req *mgmt_pb.RemoveMultiFactorFromLoginPolicyRequest) (*mgmt_pb.RemoveMultiFactorFromLoginPolicyResponse, error) {
	objectDetails, err := s.command.RemoveMultiFactorFromDefaultLoginPolicy(ctx, policy_grpc.MultiFactorTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveMultiFactorFromLoginPolicyResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}
