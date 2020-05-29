package grpc

import (
	"context"
)

func (s *Server) SearchUserGrant(ctx context.Context, in *UserGrantSearchRequest) (*UserGrantSearchResponse, error) {
	response, err := s.repo.SearchUserGrants(ctx, userGrantSearchRequestsToModel(in))
	if err != nil {
		return nil, err
	}
	return userGrantSearchResponseFromModel(response), nil
}

func (s *Server) SearchMyProjectOrgs(ctx context.Context, in *MyProjectOrgSearchRequest) (*MyProjectOrgSearchResponse, error) {
	response, err := s.repo.SearchMyProjectOrgs(ctx, myProjectOrgSearchRequestRequestsToModel(in))
	if err != nil {
		return nil, err
	}
	return projectOrgSearchResponseFromModel(response), nil
}
