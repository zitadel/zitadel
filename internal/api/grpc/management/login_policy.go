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

func (s *Server) CreateLoginPolicy(ctx context.Context, policy *management.LoginPolicyAdd) (*management.LoginPolicy, error) {
	result, err := s.org.AddLoginPolicy(ctx, loginPolicyAddToModel(policy))
	if err != nil {
		return nil, err
	}
	return loginPolicyFromModel(result), nil
}

func (s *Server) UpdateLoginPolicy(ctx context.Context, policy *management.LoginPolicy) (*management.LoginPolicy, error) {
	result, err := s.org.ChangeLoginPolicy(ctx, loginPolicyToModel(policy))
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
	result, err := s.org.SearchIdpProviders(ctx, idpProviderSearchRequestToModel(request))
	if err != nil {
		return nil, err
	}
	return idpProviderSearchResponseFromModel(result), nil
}

func (s *Server) AddIdpProviderToLoginPolicy(ctx context.Context, provider *management.IdpProviderAdd) (*management.IdpProvider, error) {
	result, err := s.org.AddIdpProviderToLoginPolicy(ctx, idpProviderAddToModel(provider))
	if err != nil {
		return nil, err
	}
	return idpProviderFromModel(result), nil
}

func (s *Server) RemoveIdpProviderFromLoginPolicy(ctx context.Context, provider *management.IdpProviderID) (*empty.Empty, error) {
	err := s.org.RemoveIdpProviderFromLoginPolicy(ctx, idpProviderToModel(provider))
	return &empty.Empty{}, err
}
