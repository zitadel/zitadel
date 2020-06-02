package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
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

func (s *Server) GetMyZitadelPermissions(ctx context.Context, _ *empty.Empty) (*MyPermissions, error) {
	perms, err := s.repo.SearchMyZitadelPermissions(ctx)
	if err != nil {
		return nil, err
	}
	return &MyPermissions{Permissions: perms}, nil
}
