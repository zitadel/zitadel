package auth

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/change"
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/api/grpc/org"
	user_grpc "github.com/caos/zitadel/internal/api/grpc/user"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) GetMyUser(ctx context.Context, _ *auth_pb.GetMyUserRequest) (*auth_pb.GetMyUserResponse, error) {
	user, err := s.repo.MyUser(ctx)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyUserResponse{User: user_grpc.UserToPb(user)}, nil
}

func (s *Server) ListMyUserChanges(ctx context.Context, req *auth_pb.ListMyUserChangesRequest) (*auth_pb.ListMyUserChangesResponse, error) {
	changes, err := s.repo.MyUserChanges(ctx, req.Query.Offset, uint64(req.Query.Limit), req.Query.Asc)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyUserChangesResponse{
		Result: change.UserChangesToPb(changes.Changes),
	}, nil
}

func (s *Server) ListMyUserSessions(ctx context.Context, req *auth_pb.ListMyUserSessionsRequest) (*auth_pb.ListMyUserSessionsResponse, error) {
	userSessions, err := s.repo.GetMyUserSessions(ctx)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyUserSessionsResponse{
		Result: user_grpc.UserSessionsToPb(userSessions),
	}, nil
}

func (s *Server) UpdateMyUserName(ctx context.Context, req *auth_pb.UpdateMyUserNameRequest) (*auth_pb.UpdateMyUserNameResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	err := s.command.ChangeUsername(ctx, ctxData.ResourceOwner, ctxData.UserID, req.UserName)
	if err != nil {
		return nil, err
	}
	//TODO: return values
	return &auth_pb.UpdateMyUserNameResponse{}, nil
}

func ctxToObjectRoot(ctx context.Context) models.ObjectRoot {
	ctxData := authz.GetCtxData(ctx)
	return models.ObjectRoot{
		AggregateID:   ctxData.UserID,
		ResourceOwner: ctxData.ResourceOwner,
	}
}

func (s *Server) ListMyUserGrants(ctx context.Context, req *auth_pb.ListMyUserGrantsRequest) (*auth_pb.ListMyUserGrantsResponse, error) {
	res, err := s.repo.SearchMyUserGrants(ctx, ListMyUserGrantsRequestToModel(req))
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyUserGrantsResponse{
		Result: UserGrantsToPb(res.Result),
		Details: object.ToListDetails(
			res.TotalResult,
			res.Sequence,
			res.Timestamp,
		),
	}, nil
}
func (s *Server) ListMyProjectOrgs(ctx context.Context, req *auth_pb.ListMyProjectOrgsRequest) (*auth_pb.ListMyProjectOrgsResponse, error) {
	res, err := s.repo.SearchMyProjectOrgs(ctx, ListMyProjectOrgsRequestToModel(req))
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyProjectOrgsResponse{
		//TODO: not all details
		Details: object.ToListDetails(res.TotalResult, 0, time.Time{}),
		Result:  org.OrgsToPb(res.Result),
	}, nil
}

func ListMyProjectOrgsRequestToModel(req *auth_pb.ListMyProjectOrgsRequest) *grant_model.UserGrantSearchRequest {
	return &grant_model.UserGrantSearchRequest{
		Offset: req.Query.Offset,
		Limit:  uint64(req.Query.Limit),
		Asc:    req.Query.Asc,
		// Queries: queries,//TODO:user grant queries missing in proto
	}
}
