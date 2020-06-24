package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/pkg/management/grpc"
)

func (s *Server) GetMyOrg(ctx context.Context, _ *empty.Empty) (*grpc.OrgView, error) {
	org, err := s.org.OrgByID(ctx, auth.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return orgViewFromModel(org), nil
}

func (s *Server) GetOrgByDomainGlobal(ctx context.Context, in *grpc.Domain) (*grpc.OrgView, error) {
	org, err := s.org.OrgByDomainGlobal(ctx, in.Domain)
	if err != nil {
		return nil, err
	}
	return orgViewFromModel(org), nil
}

func (s *Server) DeactivateMyOrg(ctx context.Context, _ *empty.Empty) (*grpc.Org, error) {
	org, err := s.org.DeactivateOrg(ctx, auth.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return orgFromModel(org), nil
}

func (s *Server) ReactivateMyOrg(ctx context.Context, _ *empty.Empty) (*grpc.Org, error) {
	org, err := s.org.ReactivateOrg(ctx, auth.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return orgFromModel(org), nil
}

func (s *Server) SearchMyOrgDomains(ctx context.Context, in *grpc.OrgDomainSearchRequest) (*grpc.OrgDomainSearchResponse, error) {
	domains, err := s.org.SearchMyOrgDomains(ctx, orgDomainSearchRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return orgDomainSearchResponseFromModel(domains), nil
}

func (s *Server) AddMyOrgDomain(ctx context.Context, in *grpc.AddOrgDomainRequest) (*grpc.OrgDomain, error) {
	domain, err := s.org.AddMyOrgDomain(ctx, addOrgDomainToModel(in))
	if err != nil {
		return nil, err
	}
	return orgDomainFromModel(domain), nil
}

func (s *Server) RemoveMyOrgDomain(ctx context.Context, in *grpc.RemoveOrgDomainRequest) (*empty.Empty, error) {
	err := s.org.RemoveMyOrgDomain(ctx, in.Domain)
	return &empty.Empty{}, err
}

func (s *Server) OrgChanges(ctx context.Context, changesRequest *grpc.ChangeRequest) (*grpc.Changes, error) {
	response, err := s.org.OrgChanges(ctx, changesRequest.Id, 0, 0)
	if err != nil {
		return nil, err
	}
	return orgChangesToResponse(response, changesRequest.GetSequenceOffset(), changesRequest.GetLimit()), nil
}

func (s *Server) GetMyOrgIamPolicy(ctx context.Context, _ *empty.Empty) (_ *grpc.OrgIamPolicy, err error) {
	policy, err := s.org.GetMyOrgIamPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return orgIamPolicyFromModel(policy), err
}
