package auth

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	user_grpc "github.com/zitadel/zitadel/internal/api/grpc/user"
	"github.com/zitadel/zitadel/internal/query"
	auth_pb "github.com/zitadel/zitadel/pkg/grpc/auth"
)

func (s *Server) ListMyZitadelPermissions(ctx context.Context, _ *auth_pb.ListMyZitadelPermissionsRequest) (*auth_pb.ListMyZitadelPermissionsResponse, error) {
	perms, err := s.query.MyZitadelPermissions(ctx, authz.GetCtxData(ctx).OrgID, authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyZitadelPermissionsResponse{
		Result: perms.Permissions,
	}, nil
}

func (s *Server) ListMyProjectPermissions(ctx context.Context, _ *auth_pb.ListMyProjectPermissionsRequest) (*auth_pb.ListMyProjectPermissionsResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	userGrantOrgID, err := query.NewUserGrantResourceOwnerSearchQuery(ctxData.OrgID)
	if err != nil {
		return nil, err
	}
	userGrantProjectID, err := query.NewUserGrantProjectIDSearchQuery(ctxData.ProjectID)
	if err != nil {
		return nil, err
	}
	userGrantUserID, err := query.NewUserGrantUserIDSearchQuery(ctxData.UserID)
	if err != nil {
		return nil, err
	}
	userGrant, err := s.query.UserGrant(ctx, true, false, userGrantOrgID, userGrantProjectID, userGrantUserID)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyProjectPermissionsResponse{
		Result: userGrant.Roles,
	}, nil
}

func (s *Server) ListMyMemberships(ctx context.Context, req *auth_pb.ListMyMembershipsRequest) (*auth_pb.ListMyMembershipsResponse, error) {
	request, err := ListMyMembershipsRequestToModel(ctx, req)
	if err != nil {
		return nil, err
	}
	response, err := s.query.Memberships(ctx, request, false, false)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyMembershipsResponse{
		Result:  user_grpc.MembershipsToMembershipsPb(response.Memberships),
		Details: object.ToListDetails(response.Count, response.Sequence, response.LastRun),
	}, nil
}
