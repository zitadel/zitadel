package auth

import (
	"context"

	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) ListMyZitadelPermissions(ctx context.Context, _ *auth_pb.ListMyZitadelPermissionsRequest) (*auth_pb.ListMyZitadelPermissionsResponse, error) {
	perms, err := s.repo.SearchMyZitadelPermissions(ctx)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyZitadelPermissionsResponse{
		Result: perms,
	}, nil
}

func (s *Server) ListMyProjectPermissions(ctx context.Context, _ *auth_pb.ListMyProjectPermissionsRequest) (*auth_pb.ListMyProjectPermissionsResponse, error) {
	perms, err := s.repo.SearchMyProjectPermissions(ctx)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyProjectPermissionsResponse{
		Result: perms,
	}, nil
}
