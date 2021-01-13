package management

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) CreateOrg(ctx context.Context, request *management.OrgCreateRequest) (_ *management.Org, err error) {
	ctxData := authz.GetCtxData(ctx)
	org, err := s.command.AddOrg(ctx, request.Name, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return orgFromDomain(org), err
}

func (s *Server) GetMyOrg(ctx context.Context, _ *empty.Empty) (*management.OrgView, error) {
	org, err := s.org.OrgByID(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return orgViewFromModel(org), nil
}

func (s *Server) GetOrgByDomainGlobal(ctx context.Context, in *management.Domain) (*management.OrgView, error) {
	org, err := s.org.OrgByDomainGlobal(ctx, in.Domain)
	if err != nil {
		return nil, err
	}
	return orgViewFromModel(org), nil
}

func (s *Server) DeactivateMyOrg(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.command.DeactivateOrg(ctx, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) ReactivateMyOrg(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.command.ReactivateOrg(ctx, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) SearchMyOrgDomains(ctx context.Context, in *management.OrgDomainSearchRequest) (*management.OrgDomainSearchResponse, error) {
	domains, err := s.org.SearchMyOrgDomains(ctx, orgDomainSearchRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return orgDomainSearchResponseFromModel(domains), nil
}

func (s *Server) AddMyOrgDomain(ctx context.Context, in *management.AddOrgDomainRequest) (*management.OrgDomain, error) {
	domain, err := s.command.AddOrgDomain(ctx, addOrgDomainToDomain(ctx, in))
	if err != nil {
		return nil, err
	}
	return orgDomainFromDomain(domain), nil
}

func (s *Server) GenerateMyOrgDomainValidation(ctx context.Context, in *management.OrgDomainValidationRequest) (*management.OrgDomainValidationResponse, error) {
	token, url, err := s.command.GenerateOrgDomainValidation(ctx, orgDomainValidationToDomain(ctx, in))
	if err != nil {
		return nil, err
	}
	return &management.OrgDomainValidationResponse{
		Token: token,
		Url:   url,
	}, nil
}

func (s *Server) ValidateMyOrgDomain(ctx context.Context, in *management.ValidateOrgDomainRequest) (*empty.Empty, error) {
	err := s.command.ValidateOrgDomain(ctx, validateOrgDomainToDomain(ctx, in))
	return &empty.Empty{}, err
}
func (s *Server) SetMyPrimaryOrgDomain(ctx context.Context, in *management.PrimaryOrgDomainRequest) (*empty.Empty, error) {
	err := s.command.SetPrimaryOrgDomain(ctx, primaryOrgDomainToDomain(ctx, in))
	return &empty.Empty{}, err
}

func (s *Server) RemoveMyOrgDomain(ctx context.Context, in *management.RemoveOrgDomainRequest) (*empty.Empty, error) {
	err := s.command.RemoveOrgDomain(ctx, removeOrgDomainToDomain(ctx, in))
	return &empty.Empty{}, err
}

func (s *Server) OrgChanges(ctx context.Context, changesRequest *management.ChangeRequest) (*management.Changes, error) {
	response, err := s.org.OrgChanges(ctx, changesRequest.Id, changesRequest.SequenceOffset, changesRequest.Limit, changesRequest.Asc)
	if err != nil {
		return nil, err
	}
	return orgChangesToResponse(response, changesRequest.GetSequenceOffset(), changesRequest.GetLimit()), nil
}

func (s *Server) GetMyOrgIamPolicy(ctx context.Context, _ *empty.Empty) (_ *management.OrgIamPolicyView, err error) {
	policy, err := s.org.GetMyOrgIamPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return orgIamPolicyViewFromModel(policy), err
}
