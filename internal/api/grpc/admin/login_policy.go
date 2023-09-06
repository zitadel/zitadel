package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/idp"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/zitadel/zitadel/internal/api/grpc/policy"
	"github.com/zitadel/zitadel/internal/api/grpc/user"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) GetLoginPolicy(ctx context.Context, _ *admin_pb.GetLoginPolicyRequest) (*admin_pb.GetLoginPolicyResponse, error) {
	policy, err := s.query.DefaultLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetLoginPolicyResponse{Policy: policy_grpc.ModelLoginPolicyToPb(policy)}, nil
}

func (s *Server) UpdateLoginPolicy(ctx context.Context, p *admin_pb.UpdateLoginPolicyRequest) (*admin_pb.UpdateLoginPolicyResponse, error) {
	policy, err := s.command.ChangeDefaultLoginPolicy(ctx, updateLoginPolicyToCommand(p))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateLoginPolicyResponse{
		Details: object.ChangeToDetailsPb(
			policy.Sequence,
			policy.EventDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) ListLoginPolicyIDPs(ctx context.Context, req *admin_pb.ListLoginPolicyIDPsRequest) (*admin_pb.ListLoginPolicyIDPsResponse, error) {
	res, err := s.query.IDPLoginPolicyLinks(ctx, authz.GetInstance(ctx).InstanceID(), ListLoginPolicyIDPsRequestToQuery(req), false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListLoginPolicyIDPsResponse{
		Result:  idp.IDPLoginPolicyLinksToPb(res.Links),
		Details: object.ToListDetails(res.Count, res.Sequence, res.LastRun),
	}, nil
}

func (s *Server) AddIDPToLoginPolicy(ctx context.Context, req *admin_pb.AddIDPToLoginPolicyRequest) (*admin_pb.AddIDPToLoginPolicyResponse, error) {
	idp, err := s.command.AddIDPProviderToDefaultLoginPolicy(ctx, &domain.IDPProvider{IDPConfigID: req.IdpId})
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddIDPToLoginPolicyResponse{
		Details: object.AddToDetailsPb(
			idp.Sequence,
			idp.ChangeDate,
			idp.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveIDPFromLoginPolicy(ctx context.Context, req *admin_pb.RemoveIDPFromLoginPolicyRequest) (*admin_pb.RemoveIDPFromLoginPolicyResponse, error) {
	idpQuery, err := query.NewIDPUserLinkIDPIDSearchQuery(req.IdpId)
	if err != nil {
		return nil, err
	}
	idps, err := s.query.IDPUserLinks(ctx, &query.IDPUserLinksSearchQuery{
		Queries: []query.SearchQuery{idpQuery},
	}, true)

	if err != nil {
		return nil, err
	}

	objectDetails, err := s.command.RemoveIDPProviderFromDefaultLoginPolicy(ctx, &domain.IDPProvider{IDPConfigID: req.IdpId}, user.ExternalIDPViewsToExternalIDPs(idps.Links)...)
	if err != nil {
		return nil, err
	}
	return &admin_pb.RemoveIDPFromLoginPolicyResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ListLoginPolicySecondFactors(ctx context.Context, req *admin_pb.ListLoginPolicySecondFactorsRequest) (*admin_pb.ListLoginPolicySecondFactorsResponse, error) {
	result, err := s.query.DefaultSecondFactors(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListLoginPolicySecondFactorsResponse{
		Details: object.ToListDetails(result.Count, result.Sequence, result.LastRun),
		Result:  policy_grpc.ModelSecondFactorTypesToPb(result.Factors),
	}, nil
}

func (s *Server) AddSecondFactorToLoginPolicy(ctx context.Context, req *admin_pb.AddSecondFactorToLoginPolicyRequest) (*admin_pb.AddSecondFactorToLoginPolicyResponse, error) {
	objectDetails, err := s.command.AddSecondFactorToDefaultLoginPolicy(ctx, policy_grpc.SecondFactorTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddSecondFactorToLoginPolicyResponse{
		Details: object.DomainToAddDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveSecondFactorFromLoginPolicy(ctx context.Context, req *admin_pb.RemoveSecondFactorFromLoginPolicyRequest) (*admin_pb.RemoveSecondFactorFromLoginPolicyResponse, error) {
	objectDetails, err := s.command.RemoveSecondFactorFromDefaultLoginPolicy(ctx, policy_grpc.SecondFactorTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	return &admin_pb.RemoveSecondFactorFromLoginPolicyResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ListLoginPolicyMultiFactors(ctx context.Context, req *admin_pb.ListLoginPolicyMultiFactorsRequest) (*admin_pb.ListLoginPolicyMultiFactorsResponse, error) {
	res, err := s.query.DefaultMultiFactors(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListLoginPolicyMultiFactorsResponse{
		Details: object.ToListDetails(res.Count, res.Sequence, res.LastRun),
		Result:  policy_grpc.ModelMultiFactorTypesToPb(res.Factors),
	}, nil
}

func (s *Server) AddMultiFactorToLoginPolicy(ctx context.Context, req *admin_pb.AddMultiFactorToLoginPolicyRequest) (*admin_pb.AddMultiFactorToLoginPolicyResponse, error) {
	objectDetails, err := s.command.AddMultiFactorToDefaultLoginPolicy(ctx, policy_grpc.MultiFactorTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddMultiFactorToLoginPolicyResponse{
		Details: object.DomainToAddDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveMultiFactorFromLoginPolicy(ctx context.Context, req *admin_pb.RemoveMultiFactorFromLoginPolicyRequest) (*admin_pb.RemoveMultiFactorFromLoginPolicyResponse, error) {
	objectDetails, err := s.command.RemoveMultiFactorFromDefaultLoginPolicy(ctx, policy_grpc.MultiFactorTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	return &admin_pb.RemoveMultiFactorFromLoginPolicyResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}
