package auth

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) SearchMyUserGrant(ctx context.Context, in *auth.UserGrantSearchRequest) (*auth.UserGrantSearchResponse, error) {
	response, err := s.repo.SearchMyUserGrants(ctx, userGrantSearchRequestsToModel(in))
	if err != nil {
		return nil, err
	}
	return userGrantSearchResponseFromModel(response), nil
}

func (s *Server) SearchMyProjectOrgs(ctx context.Context, in *auth.MyProjectOrgSearchRequest) (*auth.MyProjectOrgSearchResponse, error) {
	response, err := s.repo.SearchMyProjectOrgs(ctx, myProjectOrgSearchRequestRequestsToModel(in))
	if err != nil {
		return nil, err
	}
	return projectOrgSearchResponseFromModel(response), nil
}

func (s *Server) GetMyZitadelPermissions(ctx context.Context, _ *empty.Empty) (*auth.MyPermissions, error) {
	perms, err := s.repo.SearchMyZitadelPermissions(ctx)
	if err != nil {
		return nil, err
	}
	return &auth.MyPermissions{Permissions: perms}, nil
}

func (s *Server) GetMyProjectPermissions(ctx context.Context, _ *empty.Empty) (*auth.MyPermissions, error) {
	perms, err := s.repo.SearchMyProjectPermissions(ctx)
	if err != nil {
		return nil, err
	}
	return &auth.MyPermissions{Permissions: perms}, nil
}
