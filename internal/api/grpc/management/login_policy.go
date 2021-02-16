package management

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetLoginPolicy(ctx context.Context, _ *empty.Empty) (*management.LoginPolicyView, error) {
	result, err := s.org.GetLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return loginPolicyViewFromModel(result), nil
}

func (s *Server) GetDefaultLoginPolicy(ctx context.Context, _ *empty.Empty) (*management.LoginPolicyView, error) {
	result, err := s.org.GetDefaultLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return loginPolicyViewFromModel(result), nil
}

func (s *Server) CreateLoginPolicy(ctx context.Context, policy *management.LoginPolicyRequest) (*management.LoginPolicy, error) {
	result, err := s.command.AddLoginPolicy(ctx, authz.GetCtxData(ctx).OrgID, loginPolicyRequestToDomain(ctx, policy))
	if err != nil {
		return nil, err
	}
	return loginPolicyFromDomain(result), nil
}

func (s *Server) UpdateLoginPolicy(ctx context.Context, policy *management.LoginPolicyRequest) (*management.LoginPolicy, error) {
	result, err := s.command.ChangeLoginPolicy(ctx, authz.GetCtxData(ctx).OrgID, loginPolicyRequestToDomain(ctx, policy))
	if err != nil {
		return nil, err
	}
	return loginPolicyFromDomain(result), nil
}

func (s *Server) RemoveLoginPolicy(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.command.RemoveLoginPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) GetLoginPolicyIdpProviders(ctx context.Context, request *management.IdpProviderSearchRequest) (*management.IdpProviderSearchResponse, error) {
	result, err := s.org.SearchIDPProviders(ctx, idpProviderSearchRequestToModel(request))
	if err != nil {
		return nil, err
	}
	return idpProviderSearchResponseFromModel(result), nil
}

func (s *Server) AddIdpProviderToLoginPolicy(ctx context.Context, provider *management.IdpProviderAdd) (*management.IdpProvider, error) {
	result, err := s.command.AddIDPProviderToLoginPolicy(ctx, authz.GetCtxData(ctx).OrgID, idpProviderAddToDomain(ctx, provider))
	if err != nil {
		return nil, err
	}
	return idpProviderFromDomain(result), nil
}

func (s *Server) RemoveIdpProviderFromLoginPolicy(ctx context.Context, provider *management.IdpProviderID) (*empty.Empty, error) {
	externalIDPs, err := s.user.ExternalIDPsByIDPConfigIDAndResourceOwner(ctx, provider.IdpConfigId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return &empty.Empty{}, err
	}
	err = s.command.RemoveIDPProviderFromLoginPolicy(ctx, authz.GetCtxData(ctx).OrgID, idpProviderIDToDomain(ctx, provider), externalIDPViewsToDomain(externalIDPs)...)
	return &empty.Empty{}, err
}

func (s *Server) GetLoginPolicySecondFactors(ctx context.Context, _ *empty.Empty) (*management.SecondFactorsResult, error) {
	result, err := s.org.SearchSecondFactors(ctx)
	if err != nil {
		return nil, err
	}
	return secondFactorResultFromModel(result), nil
}

func (s *Server) AddSecondFactorToLoginPolicy(ctx context.Context, mfa *management.SecondFactor) (*management.SecondFactor, error) {
	result, err := s.command.AddSecondFactorToLoginPolicy(ctx, secondFactorTypeToDomain(mfa), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return secondFactorFromDomain(result), nil
}

func (s *Server) RemoveSecondFactorFromLoginPolicy(ctx context.Context, mfa *management.SecondFactor) (*empty.Empty, error) {
	err := s.command.RemoveSecondFactorFromLoginPolicy(ctx, secondFactorTypeToDomain(mfa), authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) GetLoginPolicyMultiFactors(ctx context.Context, _ *empty.Empty) (*management.MultiFactorsResult, error) {
	result, err := s.org.SearchMultiFactors(ctx)
	if err != nil {
		return nil, err
	}
	return multiFactorResultFromModel(result), nil
}

func (s *Server) AddMultiFactorToLoginPolicy(ctx context.Context, mfa *management.MultiFactor) (*management.MultiFactor, error) {
	result, err := s.command.AddMultiFactorToLoginPolicy(ctx, multiFactorTypeToDomain(mfa), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return multiFactorFromDomain(result), nil
}

func (s *Server) RemoveMultiFactorFromLoginPolicy(ctx context.Context, mfa *management.MultiFactor) (*empty.Empty, error) {
	err := s.command.RemoveMultiFactorFromLoginPolicy(ctx, multiFactorTypeToDomain(mfa), authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}
