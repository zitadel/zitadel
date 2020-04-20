package grpc

import (
	"context"
)

func (s *Server) GetOrgByID(ctx context.Context, orgID *OrgID) (_ *Org, err error) {
	org, err := s.org.OrgByID(ctx, orgID.Id)
	if err != nil {
		return nil, err
	}
	return orgFromModel(org), nil
}

func (s *Server) SearchOrgs(ctx context.Context, request *OrgSearchRequest) (_ *OrgSearchResponse, err error) {
	orgs, err := s.org.SearchOrgs(ctx)
	if err != nil {
		return nil, err
	}
	return &OrgSearchResponse{Result: orgsFromModel(orgs),
		Limit:  request.Limit,
		Offset: request.Offset,
		// TotalResult: , TODO: total result from search
	}, nil
}

func (s *Server) IsOrgUnique(ctx context.Context, request *UniqueOrgRequest) (org *UniqueOrgResponse, err error) {
	isUnique, err := s.org.IsOrgUnique(ctx, request.Name, request.Domain)

	return &UniqueOrgResponse{IsUnique: isUnique}, err
}

func (s *Server) SetUpOrg(ctx context.Context, orgSetUp *OrgSetUpRequest) (_ *OrgSetUpResponse, err error) {
	setUp, err := s.org.SetUpOrg(ctx, setUpRequestToModel(orgSetUp))
	if err != nil {
		return nil, err
	}
	return setUpOrgResponseFromModel(setUp), err
}
