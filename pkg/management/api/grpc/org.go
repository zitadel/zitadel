package grpc

import (
	"context"
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

func (s *Server) OrgChanges(ctx context.Context, changesRequest *ChangeRequest) (*Changes, error) {
	response, err := s.org.OrgChanges(ctx, changesRequest.Id, 0, 0)
	if err != nil {
		return nil, err
	}
	return changesToResponse(response, changesRequest.GetSequenceOffset(), changesRequest.GetLimit()), nil
}
