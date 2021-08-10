package auth

import (
	"context"

	obj_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	user_grpc "github.com/caos/zitadel/internal/api/grpc/user"
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

func (s *Server) ListMyMemberships(ctx context.Context, req *auth_pb.ListMyMembershipsRequest) (*auth_pb.ListMyMembershipsResponse, error) {
	request, err := ListMyMembershipsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	response, err := s.repo.SearchMyUserMemberships(ctx, request)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyMembershipsResponse{
		Result: user_grpc.MembershipsToMembershipsPb(response.Result),
		Details: obj_grpc.ToListDetails(
			response.TotalResult,
			response.Sequence,
			response.Timestamp,
		),
	}, nil
}
