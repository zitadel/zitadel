package management

import (
	"context"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes/empty"
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
	result, err := s.org.AddLoginPolicy(ctx, loginPolicyRequestToModel(policy))
	if err != nil {
		return nil, err
	}
	return loginPolicyFromModel(result), nil
}

func (s *Server) UpdateLoginPolicy(ctx context.Context, policy *management.LoginPolicyRequest) (*management.LoginPolicy, error) {
	result, err := s.org.ChangeLoginPolicy(ctx, loginPolicyRequestToModel(policy))
	if err != nil {
		return nil, err
	}
	return loginPolicyFromModel(result), nil
}

func (s *Server) RemoveLoginPolicy(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.org.RemoveLoginPolicy(ctx)
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
	result, err := s.org.AddIDPProviderToLoginPolicy(ctx, idpProviderAddToModel(provider))
	if err != nil {
		return nil, err
	}
	return idpProviderFromModel(result), nil
}

func (s *Server) RemoveIdpProviderFromLoginPolicy(ctx context.Context, provider *management.IdpProviderID) (*empty.Empty, error) {
	err := s.org.RemoveIDPProviderFromLoginPolicy(ctx, idpProviderToModel(provider))
	return &empty.Empty{}, err
}

func (s *Server) GetLoginPolicySoftwareMFAs(ctx context.Context, _ *empty.Empty) (*management.SoftwareMFAResult, error) {
	result, err := s.org.SearchSoftwareMFAs(ctx)
	if err != nil {
		return nil, err
	}
	return softwareMFAResultFromModel(result), nil
}

func (s *Server) AddSoftwareMFAToLoginPolicy(ctx context.Context, mfa *management.SoftwareMFA) (*management.SoftwareMFA, error) {
	result, err := s.org.AddSoftwareMFAToLoginPolicy(ctx, softwareMFATypeToModel(mfa))
	if err != nil {
		return nil, err
	}
	return softwareMFAFromModel(result), nil
}

func (s *Server) RemoveSoftwareMFAFromLoginPolicy(ctx context.Context, mfa *management.SoftwareMFA) (*empty.Empty, error) {
	err := s.org.RemoveSoftwareMFAFromLoginPolicy(ctx, softwareMFATypeToModel(mfa))
	return &empty.Empty{}, err
}

func (s *Server) GetLoginPolicyHardwareMFAs(ctx context.Context, _ *empty.Empty) (*management.HardwareMFAResult, error) {
	result, err := s.org.SearchHardwareMFAs(ctx)
	if err != nil {
		return nil, err
	}
	return hardwareMFAResultFromModel(result), nil
}

func (s *Server) AddHardwareMFAToLoginPolicy(ctx context.Context, mfa *management.HardwareMFA) (*management.HardwareMFA, error) {
	result, err := s.org.AddHardwareMFAToLoginPolicy(ctx, hardwareMFATypeToModel(mfa))
	if err != nil {
		return nil, err
	}
	return hardwareMFAFromModel(result), nil
}

func (s *Server) RemoveHardwareMFAFromLoginPolicy(ctx context.Context, mfa *management.HardwareMFA) (*empty.Empty, error) {
	err := s.org.RemoveHardwareMFAFromLoginPolicy(ctx, hardwareMFATypeToModel(mfa))
	return &empty.Empty{}, err
}
