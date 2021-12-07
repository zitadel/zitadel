package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/change"
	"github.com/caos/zitadel/internal/api/grpc/metadata"
	obj_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/api/grpc/org"
	user_grpc "github.com/caos/zitadel/internal/api/grpc/user"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/query"
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
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
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
		Details: obj_grpc.ToListDetails(
			res.TotalResult,
			res.Sequence,
			res.Timestamp,
		),
	}, nil
}

func (s *Server) ListMyProjectOrgs(ctx context.Context, req *auth_pb.ListMyProjectOrgsRequest) (*auth_pb.ListMyProjectOrgsResponse, error) {
	queries, err := ListMyProjectOrgsRequestToQuery(req)
	if err != nil {
		return nil, err
	}

	iam, err := s.query.IAMByID(ctx, domain.IAMID)
	if err != nil {
		return nil, err
	}
	ctxData := authz.GetCtxData(ctx)

	//client of user is not in project of ZITADEL
	if ctxData.ProjectID != iam.IAMProjectID {
		grants, err := s.repo.UserGrantsByProjectAndUserID(ctxData.ProjectID, ctxData.UserID)
		if err != nil {
			return nil, err
		}

		ids := make([]string, 0, len(grants))
		for _, grant := range grants {
			ids = appendIfNotExists(ids, grant.ResourceOwner)
		}

		idsQuery, err := query.NewOrgIDsSearchQuery(ids...)
		if err != nil {
			return nil, err
		}
		queries.Queries = append(queries.Queries, idsQuery)
	} else if authz.HasGlobalExplicitPermission(authz.GetAllPermissionsFromCtx(ctx), "iam.read") {
		//user is allowed to read all organisation
		//no additional query required
	} else {
		// all orgs of my meberships
		userQuery, err := query.NewMembershipUserIDQuery(ctxData.UserID)
		if err != nil {
			return nil, err
		}
		memberships, err := s.query.Memberships(ctx, &query.MembershipSearchQuery{
			Queries: []query.SearchQuery{userQuery},
		})
		if err != nil {
			return nil, err
		}

		ids := make([]string, 0, len(memberships.Memberships))
		for _, grant := range memberships.Memberships {
			ids = appendIfNotExists(ids, grant.ResourceOwner)
		}

		idsQuery, err := query.NewOrgIDsSearchQuery(ids...)
		if err != nil {
			return nil, err
		}
		queries.Queries = append(queries.Queries, idsQuery)
	}

	orgs, err := s.query.SearchOrgs(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyProjectOrgsResponse{
		Details: obj_grpc.ToListDetails(orgs.Count, orgs.Sequence, orgs.Timestamp),
		Result:  org.OrgsToPb(orgs.Orgs),
	}, nil
}

func appendIfNotExists(array []string, value string) []string {
	for _, a := range array {
		if a == value {
			return array
		}
	}
	return append(array, value)
}

func ListMyProjectOrgsRequestToQuery(req *auth_pb.ListMyProjectOrgsRequest) (*query.OrgSearchQueries, error) {
	offset, limit, asc := obj_grpc.ListQueryToModel(req.Query)
	queries, err := org.OrgQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.OrgSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}
