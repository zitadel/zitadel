package auth

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/change"
	"github.com/zitadel/zitadel/internal/api/grpc/metadata"
	obj_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/api/grpc/org"
	user_grpc "github.com/zitadel/zitadel/internal/api/grpc/user"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/user"
	auth_pb "github.com/zitadel/zitadel/pkg/grpc/auth"
)

func (s *Server) GetMyUser(ctx context.Context, _ *auth_pb.GetMyUserRequest) (*auth_pb.GetMyUserResponse, error) {
	user, err := s.query.GetUserByID(ctx, true, authz.GetCtxData(ctx).UserID, false)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyUserResponse{User: user_grpc.UserToPb(user, s.assetsAPIDomain(ctx))}, nil
}

func (s *Server) RemoveMyUser(ctx context.Context, _ *auth_pb.RemoveMyUserRequest) (*auth_pb.RemoveMyUserResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	userGrantUserID, err := query.NewUserGrantUserIDSearchQuery(ctxData.UserID)
	if err != nil {
		return nil, err
	}
	queries := &query.UserGrantsQueries{Queries: []query.SearchQuery{userGrantUserID}}
	grants, err := s.query.UserGrants(ctx, queries, true, false)
	if err != nil {
		return nil, err
	}

	userQuery, err := query.NewMembershipUserIDQuery(authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	memberships, err := s.query.Memberships(ctx, &query.MembershipSearchQuery{
		Queries: []query.SearchQuery{userQuery},
	}, false)
	if err != nil {
		return nil, err
	}
	details, err := s.command.RemoveUser(ctx, ctxData.UserID, ctxData.ResourceOwner, cascadingMemberships(memberships.Memberships), userGrantsToIDs(grants.UserGrants)...)
	if err != nil {
		return nil, err
	}
	return &auth_pb.RemoveMyUserResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ListMyUserChanges(ctx context.Context, req *auth_pb.ListMyUserChangesRequest) (*auth_pb.ListMyUserChangesResponse, error) {
	var (
		limit    uint64
		sequence uint64
		asc      bool
	)
	if req.Query != nil {
		limit = uint64(req.Query.Limit)
		sequence = req.Query.Sequence
		asc = req.Query.Asc
	}

	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AllowTimeTravel().
		Limit(limit).
		OrderDesc().
		ResourceOwner(authz.GetCtxData(ctx).ResourceOwner).
		AddQuery().
		SequenceGreater(sequence).
		AggregateTypes(user.AggregateType).
		AggregateIDs(authz.GetCtxData(ctx).UserID).
		Builder()
	if asc {
		query.OrderAsc()
	}

	changes, err := s.query.SearchEvents(ctx, query, s.auditLogRetention)
	if err != nil {
		return nil, err
	}

	return &auth_pb.ListMyUserChangesResponse{
		Result: change.EventsToChangesPb(changes, s.assetsAPIDomain(ctx)),
	}, nil
}

func (s *Server) ListMyMetadata(ctx context.Context, req *auth_pb.ListMyMetadataRequest) (*auth_pb.ListMyMetadataResponse, error) {
	queries, err := ListUserMetadataToQuery(req)
	if err != nil {
		return nil, err
	}
	res, err := s.query.SearchUserMetadata(ctx, true, authz.GetCtxData(ctx).UserID, queries, false)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyMetadataResponse{
		Result:  metadata.UserMetadataListToPb(res.Metadata),
		Details: obj_grpc.ToListDetails(res.Count, res.Sequence, res.Timestamp),
	}, nil
}

func (s *Server) GetMyMetadata(ctx context.Context, req *auth_pb.GetMyMetadataRequest) (*auth_pb.GetMyMetadataResponse, error) {
	data, err := s.query.GetUserMetadataByKey(ctx, true, authz.GetCtxData(ctx).UserID, req.Key, false)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyMetadataResponse{
		Metadata: metadata.UserMetadataToPb(data),
	}, nil
}

func (s *Server) ListMyUserSessions(ctx context.Context, req *auth_pb.ListMyUserSessionsRequest) (*auth_pb.ListMyUserSessionsResponse, error) {
	userSessions, err := s.repo.GetMyUserSessions(ctx)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyUserSessionsResponse{
		Result: user_grpc.UserSessionsToPb(userSessions, s.assetsAPIDomain(ctx)),
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
	queries, err := ListMyUserGrantsRequestToQuery(ctx, req)
	if err != nil {
		return nil, err
	}
	res, err := s.query.UserGrants(ctx, queries, false, false)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyUserGrantsResponse{
		Result:  UserGrantsToPb(res.UserGrants),
		Details: obj_grpc.ToListDetails(res.Count, res.Sequence, res.Timestamp),
	}, nil
}

func (s *Server) ListMyProjectOrgs(ctx context.Context, req *auth_pb.ListMyProjectOrgsRequest) (*auth_pb.ListMyProjectOrgsResponse, error) {
	queries, err := ListMyProjectOrgsRequestToQuery(req)
	if err != nil {
		return nil, err
	}

	ctxData := authz.GetCtxData(ctx)

	//client of user is not in project of ZITADEL
	if ctxData.ProjectID != authz.GetInstance(ctx).ProjectID() {
		userGrantProjectID, err := query.NewUserGrantProjectIDSearchQuery(ctxData.ProjectID)
		if err != nil {
			return nil, err
		}
		userGrantUserID, err := query.NewUserGrantUserIDSearchQuery(ctxData.UserID)
		if err != nil {
			return nil, err
		}

		grants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{Queries: []query.SearchQuery{userGrantProjectID, userGrantUserID}}, false, false)
		if err != nil {
			return nil, err
		}

		ids := make([]string, 0, len(grants.UserGrants))
		for _, grant := range grants.UserGrants {
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
			for _, membership := range memberships.Memberships {
				orgID := membership.ResourceOwner
				if membership.ProjectGrant != nil && membership.ProjectGrant.GrantedOrgID != "" {
					orgID = membership.ProjectGrant.GrantedOrgID
				}
				ids = appendIfNotExists(ids, orgID)
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
	}, false)
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
		typ, displayName, aggID, objID := MemberTypeToDomain(membership)
		result[i] = &domain.UserMembership{
			UserID:        membership.UserID,
			MemberType:    typ,
			AggregateID:   aggID,
			ObjectID:      objID,
			Roles:         membership.Roles,
			DisplayName:   displayName,
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

func MemberTypeToDomain(m *query.Membership) (_ domain.MemberType, displayName, aggID, objID string) {
	if m.Org != nil {
		return domain.MemberTypeOrganisation, m.Org.Name, m.Org.OrgID, ""
	} else if m.IAM != nil {
		return domain.MemberTypeIam, m.IAM.Name, m.IAM.IAMID, ""
	} else if m.Project != nil {
		return domain.MemberTypeProject, m.Project.Name, m.Project.ProjectID, ""
	} else if m.ProjectGrant != nil {
		return domain.MemberTypeProjectGrant, m.ProjectGrant.ProjectName, m.ProjectGrant.ProjectID, m.ProjectGrant.GrantID
	}
	return domain.MemberTypeUnspecified, "", "", ""
}

func cascadingMemberships(memberships []*query.Membership) []*command.CascadingMembership {
	cascades := make([]*command.CascadingMembership, len(memberships))
	for i, membership := range memberships {
		cascades[i] = &command.CascadingMembership{
			UserID:        membership.UserID,
			ResourceOwner: membership.ResourceOwner,
			IAM:           cascadingIAMMembership(membership.IAM),
			Org:           cascadingOrgMembership(membership.Org),
			Project:       cascadingProjectMembership(membership.Project),
			ProjectGrant:  cascadingProjectGrantMembership(membership.ProjectGrant),
		}
	}
	return cascades
}

func cascadingIAMMembership(membership *query.IAMMembership) *command.CascadingIAMMembership {
	if membership == nil {
		return nil
	}
	return &command.CascadingIAMMembership{IAMID: membership.IAMID}
}
func cascadingOrgMembership(membership *query.OrgMembership) *command.CascadingOrgMembership {
	if membership == nil {
		return nil
	}
	return &command.CascadingOrgMembership{OrgID: membership.OrgID}
}
func cascadingProjectMembership(membership *query.ProjectMembership) *command.CascadingProjectMembership {
	if membership == nil {
		return nil
	}
	return &command.CascadingProjectMembership{ProjectID: membership.ProjectID}
}
func cascadingProjectGrantMembership(membership *query.ProjectGrantMembership) *command.CascadingProjectGrantMembership {
	if membership == nil {
		return nil
	}
	return &command.CascadingProjectGrantMembership{ProjectID: membership.ProjectID, GrantID: membership.GrantID}
}

func userGrantsToIDs(userGrants []*query.UserGrant) []string {
	converted := make([]string, len(userGrants))
	for i, grant := range userGrants {
		converted[i] = grant.ID
	}
	return converted
}
