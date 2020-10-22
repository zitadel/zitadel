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

func (s *Server) GetDefaultLoginPolicySoftwareMFAs(ctx context.Context, _ *empty.Empty) (*admin.SoftwareMFAResult, error) {
	result, err := s.iam.SearchDefaultSoftwareMFAs(ctx)
	if err != nil {
		return nil, err
	}
	return softwareMFAResultFromModel(result), nil
}

func (s *Server) AddSoftwareMFAToDefaultLoginPolicy(ctx context.Context, mfa *admin.SoftwareMFA) (*admin.SoftwareMFA, error) {
	result, err := s.iam.AddSoftwareMFAToLoginPolicy(ctx, softwareMFATypeToModel(mfa))
	if err != nil {
		return nil, err
	}
	return softwareMFAFromModel(result), nil
}

func (s *Server) RemoveSoftwareMFAFromDefaultLoginPolicy(ctx context.Context, mfa *admin.SoftwareMFA) (*empty.Empty, error) {
	err := s.iam.RemoveSoftwareMFAFromLoginPolicy(ctx, softwareMFATypeToModel(mfa))
	return &empty.Empty{}, err
}

func (s *Server) GetDefaultLoginPolicyHardwareMFAs(ctx context.Context, _ *empty.Empty) (*admin.HardwareMFAResult, error) {
	result, err := s.iam.SearchDefaultHardwareMFAs(ctx)
	if err != nil {
		return nil, err
	}
	return hardwareMFAResultFromModel(result), nil
}

func (s *Server) AddHardwareMFAToDefaultLoginPolicy(ctx context.Context, mfa *admin.HardwareMFA) (*admin.HardwareMFA, error) {
	result, err := s.iam.AddHardwareMFAToLoginPolicy(ctx, hardwareMFATypeToModel(mfa))
	if err != nil {
		return nil, err
	}
	return hardwareMFAFromModel(result), nil
}

func (s *Server) RemoveHardwareMFAFromDefaultLoginPolicy(ctx context.Context, mfa *admin.HardwareMFA) (*empty.Empty, error) {
	err := s.iam.RemoveHardwareMFAFromLoginPolicy(ctx, hardwareMFATypeToModel(mfa))
	return &empty.Empty{}, err
}
