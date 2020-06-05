package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/internal/errors"
)

func (s *Server) GetOrgByID(ctx context.Context, orgID *OrgID) (*Org, error) {
	org, err := s.org.OrgByID(ctx, orgID.Id)
	if err != nil {
		return nil, err
	}
	return orgFromModel(org), nil
}

func (s *Server) GetOrgByDomainGlobal(ctx context.Context, in *OrgDomain) (*Org, error) {
	org, err := s.org.OrgByDomainGlobal(ctx, in.Domain)
	if err != nil {
		return nil, err
	}
	return orgFromView(org), nil
}

func (s *Server) DeactivateOrg(ctx context.Context, in *OrgID) (*Org, error) {
	org, err := s.org.DeactivateOrg(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return orgFromModel(org), nil
}

func (s *Server) ReactivateOrg(ctx context.Context, in *OrgID) (*Org, error) {
	org, err := s.org.ReactivateOrg(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return orgFromModel(org), nil
}

func (s *Server) SearchMyOrgDomains(ctx context.Context, in *OrgDomainSearchRequest) (*OrgDomainSearchResponse, error) {
	domains, err := s.org.SearchMyOrgDomains(ctx, orgDomainSearchRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return orgDomainSearchResponseFromModel(domains), nil
}

func (s *Server) AddMyOrgDomain(ctx context.Context, in *AddOrgDomainRequest) (*OrgDomain, error) {
	domain, err := s.org.AddMyOrgDomain(ctx, addOrgDomainToModel(in))
	if err != nil {
		return nil, err
	}
	return orgDomainFromModel(domain), nil
}

func (s *Server) RemoveMyOrgDomain(ctx context.Context, in *RemoveOrgDomainRequest) (*empty.Empty, error) {
	err := s.org.RemoveMyOrgDomain(ctx, in.Domain)
	return &empty.Empty{}, err
}

func (s *Server) OrgChanges(ctx context.Context, changesRequest *ChangeRequest) (*Changes, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-DNiIq", "unimplemented")
}
