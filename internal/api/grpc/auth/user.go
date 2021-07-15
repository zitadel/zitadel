package auth

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/change"
	"github.com/caos/zitadel/internal/api/grpc/object"
	obj_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/api/grpc/org"
	user_grpc "github.com/caos/zitadel/internal/api/grpc/user"
	"github.com/caos/zitadel/internal/domain"
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
	sequence, limit, asc := change.ChangeQueryToModel(req.Query)
	features, err := s.repo.GetOrgFeatures(ctx, authz.GetCtxData(ctx).ResourceOwner)
	if err != nil {
		return nil, err
	}
	changes, err := s.repo.MyUserChanges(ctx, sequence, limit, asc, features.AuditLogRetention)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyUserChangesResponse{
		Result: change.UserChangesToPb(changes.Changes),
	}, nil
}

func (s *Server) ListMyMetaData(ctx context.Context, req *auth_pb.ListMyMetaDataRequest) (*auth_pb.ListMyMetaDataResponse, error) {
	return nil, nil
}

func (s *Server) GetMyMetaData(ctx context.Context, req *auth_pb.GetMyMetaDataRequest) (*auth_pb.GetMyMetaDataResponse, error) {
	return nil, nil
}

func (s *Server) SetMyMetaData(ctx context.Context, req *auth_pb.SetMyMetaDataRequest) (*auth_pb.SetMyMetaDataResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	result, err := s.command.SetUserMetaData(ctx, &domain.MetaData{Key: req.Key, Value: req.Value}, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.SetMyMetaDataResponse{
		Details: obj_grpc.AddToDetailsPb(
			result.Sequence,
			result.ChangeDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) BulkSetMyMetaData(ctx context.Context, req *auth_pb.BulkSetMyMetaDataRequest) (*auth_pb.BulkSetMyMetaDataResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	result, err := s.command.BulkSetUserMetaData(ctx, ctxData.UserID, ctxData.ResourceOwner, BulkSetMetaDataToDomain(req)...)
	if err != nil {
		return nil, err
	}
	return &auth_pb.BulkSetMyMetaDataResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(result),
	}, nil
}

func (s *Server) RemoveMyMetaData(ctx context.Context, req *auth_pb.RemoveMyMetaDataRequest) (*auth_pb.RemoveMyMetaDataResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	result, err := s.command.RemoveUserMetaData(ctx, req.Key, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.RemoveMyMetaDataResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(result),
	}, nil
}

func (s *Server) BulkRemoveMyMetaData(ctx context.Context, req *auth_pb.BulkRemoveMyMetaDataRequest) (*auth_pb.BulkRemoveMyMetaDataResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	result, err := s.command.BulkRemoveUserMetaData(ctx, ctxData.UserID, ctxData.ResourceOwner, req.Keys...)
	if err != nil {
		return nil, err
	}
	return &auth_pb.BulkRemoveMyMetaDataResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(result),
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
	objectDetails, err := s.command.ChangeUsername(ctx, ctxData.ResourceOwner, ctxData.UserID, req.UserName)
	if err != nil {
		return nil, err
	}
	return &auth_pb.UpdateMyUserNameResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
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
	r, err := ListMyProjectOrgsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	res, err := s.repo.SearchMyProjectOrgs(ctx, r)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyProjectOrgsResponse{
		//TODO: not all details
		Details: object.ToListDetails(res.TotalResult, 0, time.Time{}),
		Result:  org.OrgsToPb(res.Result),
	}, nil
}

func ListMyProjectOrgsRequestToModel(req *auth_pb.ListMyProjectOrgsRequest) (*grant_model.UserGrantSearchRequest, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := org.OrgQueriesToUserGrantModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &grant_model.UserGrantSearchRequest{
		Offset:  offset,
		Limit:   limit,
		Asc:     asc,
		Queries: queries,
	}, nil
}
