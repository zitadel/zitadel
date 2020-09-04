package admin

import (
	"context"

	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes/empty"
)

// ToDo Michi
func (s *Server) GetDefaultLabelPolicy(ctx context.Context, _ *empty.Empty) (*admin.DefaultLabelPolicyView, error) {
	result, err := s.iam.GetDefaultLabelPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return labelPolicyViewFromModel(result), nil
}

func (s *Server) UpdateDefaultLabelPolicy(ctx context.Context, policy *admin.DefaultLabelPolicy) (*admin.DefaultLabelPolicy, error) {
	result, err := s.iam.ChangeDefaultLabelPolicy(ctx, labelPolicyToModel(policy))
	if err != nil {
		return nil, err
	}
	return labelPolicyFromModel(result), nil
}

// func (s *Server) GetDefaultLabelPolicyIdpProviders(ctx context.Context, request *admin.IdpProviderSearchRequest) (*admin.IdpProviderSearchResponse, error) {
// 	result, err := s.iam.SearchDefaultIDPProviders(ctx, idpProviderSearchRequestToModel(request))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return idpProviderSearchResponseFromModel(result), nil
// }

// func (s *Server) AddIdpProviderToDefaultLabelPolicy(ctx context.Context, provider *admin.IdpProviderID) (*admin.IdpProviderID, error) {
// 	result, err := s.iam.AddIDPProviderToLabelPolicy(ctx, idpProviderToModel(provider))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return idpProviderFromModel(result), nil
// }

// func (s *Server) RemoveIdpProviderFromDefaultLabelPolicy(ctx context.Context, provider *admin.IdpProviderID) (*empty.Empty, error) {
// 	err := s.iam.RemoveIDPProviderFromLabelPolicy(ctx, idpProviderToModel(provider))
// 	return &empty.Empty{}, err
// }
