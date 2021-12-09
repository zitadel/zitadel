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
	user_model "github.com/caos/zitadel/internal/user/model"
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

func (s *Server) RemoveMyUser(ctx context.Context, _ *auth_pb.RemoveMyUserRequest) (*auth_pb.RemoveMyUserResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	grants, err := s.repo.SearchMyUserGrants(ctx, &grant_model.UserGrantSearchRequest{Queries: []*grant_model.UserGrantSearchQuery{}})
	if err != nil {
		return nil, err
	}
	membersShips, err := s.repo.SearchMyUserMemberships(ctx, &user_model.UserMembershipSearchRequest{Queries: []*user_model.UserMembershipSearchQuery{}})
	if err != nil {
		return nil, err
	}
	details, err := s.command.RemoveUser(ctx, ctxData.UserID, ctxData.ResourceOwner, UserMembershipViewsToDomain(membersShips.Result), userGrantsToIDs(grants.Result)...)
	if err != nil {
		return nil, err
	}
	return &auth_pb.RemoveMyUserResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ListMyUserChanges(ctx context.Context, req *auth_pb.ListMyUserChangesRequest) (*auth_pb.ListMyUserChangesResponse, error) {
	sequence, limit, asc := change.ChangeQueryToModel(req.Query)
	features, err := s.query.FeaturesByOrgID(ctx, authz.GetCtxData(ctx).ResourceOwner)
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

func UserMembershipViewsToDomain(memberships []*user_model.UserMembershipView) []*domain.UserMembership {
	result := make([]*domain.UserMembership, len(memberships))
	for i, membership := range memberships {
		result[i] = &domain.UserMembership{
			UserID:            membership.UserID,
			MemberType:        MemberTypeToDomain(membership.MemberType),
			AggregateID:       membership.AggregateID,
			ObjectID:          membership.ObjectID,
			Roles:             membership.Roles,
			DisplayName:       membership.DisplayName,
			CreationDate:      membership.CreationDate,
			ChangeDate:        membership.ChangeDate,
			ResourceOwner:     membership.ResourceOwner,
			ResourceOwnerName: membership.ResourceOwnerName,
			Sequence:          membership.Sequence,
		}
	}
	return result
}

func MemberTypeToDomain(mType user_model.MemberType) domain.MemberType {
	switch mType {
	case user_model.MemberTypeIam:
		return domain.MemberTypeIam
	case user_model.MemberTypeOrganisation:
		return domain.MemberTypeOrganisation
	case user_model.MemberTypeProject:
		return domain.MemberTypeProject
	case user_model.MemberTypeProjectGrant:
		return domain.MemberTypeProjectGrant
	default:
		return domain.MemberTypeUnspecified
	}
}

func userGrantsToIDs(userGrants []*grant_model.UserGrantView) []string {
	converted := make([]string, len(userGrants))
	for i, grant := range userGrants {
		converted[i] = grant.ID
	}
	return converted
}
