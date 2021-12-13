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

	userQuery, err := query.NewMembershipUserIDQuery(authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	memberships, err := s.query.Memberships(ctx, &query.MembershipSearchQuery{
		Queries: []query.SearchQuery{userQuery},
	})
	// if err != nil {
	// 	return nil, err
	// }
	details, err := s.command.RemoveUser(ctx, ctxData.UserID, ctxData.ResourceOwner, membershipToDomain(memberships.Memberships), userGrantsToIDs(grants.Result)...)
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
	} else {
		memberships, err := s.myOrgsQuery(ctx, ctxData)
		if err != nil {
			return nil, err
		}

		if !isIAMAdmin(memberships.Memberships) {
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

func (s *Server) myOrgsQuery(ctx context.Context, ctxData authz.CtxData) (*query.Memberships, error) {
	userQuery, err := query.NewMembershipUserIDQuery(ctxData.UserID)
	if err != nil {
		return nil, err
	}
	return s.query.Memberships(ctx, &query.MembershipSearchQuery{
		Queries: []query.SearchQuery{userQuery},
	})
}

func isIAMAdmin(memberships []*query.Membership) bool {
	for _, m := range memberships {
		if m.IAM != nil {
			return true
		}
	}
	return false
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

func membershipToDomain(memberships []*query.Membership) []*domain.UserMembership {
	result := make([]*domain.UserMembership, len(memberships))
	for i, membership := range memberships {
		typ, aggID, objID := MemberTypeToDomain(membership)
		result[i] = &domain.UserMembership{
			UserID:        membership.UserID,
			MemberType:    typ,
			AggregateID:   aggID,
			ObjectID:      objID,
			Roles:         membership.Roles,
			DisplayName:   membership.DisplayName,
			CreationDate:  membership.CreationDate,
			ChangeDate:    membership.ChangeDate,
			ResourceOwner: membership.ResourceOwner,
			//TODO: implement
			// ResourceOwnerName: membership.ResourceOwnerName,
			Sequence: membership.Sequence,
		}
	}
	return result
}

func MemberTypeToDomain(m *query.Membership) (_ domain.MemberType, aggID, objID string) {
	if m.Org != nil {
		return domain.MemberTypeOrganisation, m.Org.OrgID, ""
	} else if m.IAM != nil {
		return domain.MemberTypeIam, m.IAM.IAMID, ""
	} else if m.Project != nil {
		return domain.MemberTypeProject, m.Project.ProjectID, ""
	} else if m.ProjectGrant != nil {
		return domain.MemberTypeProjectGrant, m.ProjectGrant.ProjectID, m.ProjectGrant.GrantID
	}
	return domain.MemberTypeUnspecified, "", ""
}

func userGrantsToIDs(userGrants []*grant_model.UserGrantView) []string {
	converted := make([]string, len(userGrants))
	for i, grant := range userGrants {
		converted[i] = grant.ID
	}
	return converted
}
