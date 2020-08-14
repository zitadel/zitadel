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

func (s *Server) UpdateDefaultLoginPolicy(ctx context.Context, policy *admin.DefaultLoginPolicy) (*admin.DefaultLoginPolicy, error) {
	result, err := s.iam.ChangeDefaultLoginPolicy(ctx, loginPolicyToModel(policy))
	if err != nil {
		return nil, err
	}
	return loginPolicyFromModel(result), nil
}

func (s *Server) GetDefaultLoginPolicyIdpProviders(ctx context.Context, request *admin.IdpProviderSearchRequest) (*admin.IdpProviderSearchResponse, error) {
	result, err := s.iam.SearchDefaultIdpProviders(ctx, idpProviderSearchRequestToModel(request))
	if err != nil {
		return nil, err
	}
	return idpProviderSearchResponseFromModel(result), nil
}

func (s *Server) AddIdpProviderToDefaultLoginPolicy(ctx context.Context, provider *admin.IdpProviderID) (*admin.IdpProviderID, error) {
	result, err := s.iam.AddIdpProviderToLoginPolicy(ctx, idpProviderToModel(provider))
	if err != nil {
		return nil, err
	}
	return idpProviderFromModel(result), nil
}

func (s *Server) RemoveIdpProviderFromDefaultLoginPolicy(ctx context.Context, provider *admin.IdpProviderID) (*empty.Empty, error) {
	err := s.iam.RemoveIdpProviderFromLoginPolicy(ctx, idpProviderToModel(provider))
	return &empty.Empty{}, err
}
