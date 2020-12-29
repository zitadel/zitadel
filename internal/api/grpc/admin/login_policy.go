package admin

import (
	"context"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetDefaultLoginPolicy(ctx context.Context, _ *empty.Empty) (*admin.DefaultLoginPolicyView, error) {
	result, err := s.iam.GetDefaultLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return loginPolicyViewFromModel(result), nil
}

func (s *Server) UpdateDefaultLoginPolicy(ctx context.Context, policy *admin.DefaultLoginPolicyRequest) (*admin.DefaultLoginPolicy, error) {
	result, err := s.iam.ChangeDefaultLoginPolicy(ctx, loginPolicyToModel(policy))
	if err != nil {
		return nil, err
	}
	return loginPolicyFromModel(result), nil
}

func (s *Server) GetDefaultLoginPolicyIdpProviders(ctx context.Context, request *admin.IdpProviderSearchRequest) (*admin.IdpProviderSearchResponse, error) {
	result, err := s.iam.SearchDefaultIDPProviders(ctx, idpProviderSearchRequestToModel(request))
	if err != nil {
		return nil, err
	}
	return idpProviderSearchResponseFromModel(result), nil
}

func (s *Server) AddIdpProviderToDefaultLoginPolicy(ctx context.Context, provider *admin.IdpProviderID) (*admin.IdpProviderID, error) {
	result, err := s.iam.AddIDPProviderToLoginPolicy(ctx, idpProviderToModel(provider))
	if err != nil {
		return nil, err
	}
	return idpProviderFromModel(result), nil
}

func (s *Server) RemoveIdpProviderFromDefaultLoginPolicy(ctx context.Context, provider *admin.IdpProviderID) (*empty.Empty, error) {
	err := s.iam.RemoveIDPProviderFromLoginPolicy(ctx, idpProviderToModel(provider))
	return &empty.Empty{}, err
}

func (s *Server) GetDefaultLoginPolicySecondFactors(ctx context.Context, _ *empty.Empty) (*admin.SecondFactorsResult, error) {
	result, err := s.iam.SearchDefaultSecondFactors(ctx)
	if err != nil {
		return nil, err
	}
	return secondFactorsResultFromModel(result), nil
}

func (s *Server) AddSecondFactorToDefaultLoginPolicy(ctx context.Context, mfa *admin.SecondFactor) (*admin.SecondFactor, error) {
	result, err := s.iam.AddSecondFactorToLoginPolicy(ctx, secondFactorTypeToModel(mfa))
	if err != nil {
		return nil, err
	}
	return secondFactorFromModel(result), nil
}

func (s *Server) RemoveSecondFactorFromDefaultLoginPolicy(ctx context.Context, mfa *admin.SecondFactor) (*empty.Empty, error) {
	err := s.iam.RemoveSecondFactorFromLoginPolicy(ctx, secondFactorTypeToModel(mfa))
	return &empty.Empty{}, err
}

func (s *Server) GetDefaultLoginPolicyMultiFactors(ctx context.Context, _ *empty.Empty) (*admin.MultiFactorsResult, error) {
	result, err := s.iam.SearchDefaultMultiFactors(ctx)
	if err != nil {
		return nil, err
	}
	return multiFactorResultFromModel(result), nil
}

func (s *Server) AddMultiFactorToDefaultLoginPolicy(ctx context.Context, mfa *admin.MultiFactor) (*admin.MultiFactor, error) {
	result, err := s.iam.AddMultiFactorToLoginPolicy(ctx, multiFactorTypeToModel(mfa))
	if err != nil {
		return nil, err
	}
	return multiFactorFromModel(result), nil
}

func (s *Server) RemoveMultiFactorFromDefaultLoginPolicy(ctx context.Context, mfa *admin.MultiFactor) (*empty.Empty, error) {
	err := s.iam.RemoveMultiFactorFromLoginPolicy(ctx, multiFactorTypeToModel(mfa))
	return &empty.Empty{}, err
}
