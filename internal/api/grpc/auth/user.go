package auth

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/change"
	"github.com/caos/zitadel/internal/api/grpc/metadata"
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

func (s *Server) ListMyMetadata(ctx context.Context, req *auth_pb.ListMyMetadataRequest) (*auth_pb.ListMyMetadataResponse, error) {
	res, err := s.repo.SearchMyMetadata(ctx, ListUserMetadataToDomain(req))
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyMetadataResponse{
		Result: metadata.MetadataListToPb(res.Result),
		Details: obj_grpc.ToListDetails(
			res.TotalResult,
			res.Sequence,
			res.Timestamp,
		),
	}, nil
}

func (s *Server) GetMyMetadata(ctx context.Context, req *auth_pb.GetMyMetadataRequest) (*auth_pb.GetMyMetadataResponse, error) {
	data, err := s.repo.GetMyMetadataByKey(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyMetadataResponse{
		Metadata: metadata.DomainMetadataToPb(data),
	}, nil
}

func (s *Server) SetMyMetadata(ctx context.Context, req *auth_pb.SetMyMetadataRequest) (*auth_pb.SetMyMetadataResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	result, err := s.command.SetUserMetadata(ctx, &domain.Metadata{Key: req.Key, Value: req.Value}, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.SetMyMetadataResponse{
		Details: obj_grpc.AddToDetailsPb(
			result.Sequence,
			result.ChangeDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) BulkSetMyMetadata(ctx context.Context, req *auth_pb.BulkSetMyMetadataRequest) (*auth_pb.BulkSetMyMetadataResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	result, err := s.command.BulkSetUserMetadata(ctx, ctxData.UserID, ctxData.ResourceOwner, BulkSetMetadataToDomain(req)...)
	if err != nil {
		return nil, err
	}
	return &auth_pb.BulkSetMyMetadataResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(result),
	}, nil
}

func (s *Server) RemoveMyMetadata(ctx context.Context, req *auth_pb.RemoveMyMetadataRequest) (*auth_pb.RemoveMyMetadataResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	result, err := s.command.RemoveUserMetadata(ctx, req.Key, ctxData.UserID, ctxData.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &auth_pb.RemoveMyMetadataResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(result),
	}, nil
}

func (s *Server) BulkRemoveMyMetadata(ctx context.Context, req *auth_pb.BulkRemoveMyMetadataRequest) (*auth_pb.BulkRemoveMyMetadataResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	result, err := s.command.BulkRemoveUserMetadata(ctx, ctxData.UserID, ctxData.ResourceOwner, req.Keys...)
	if err != nil {
		return nil, err
	}
	return &auth_pb.BulkRemoveMyMetadataResponse{
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
