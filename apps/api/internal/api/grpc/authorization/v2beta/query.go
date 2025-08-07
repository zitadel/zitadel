package authorization

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	filter "github.com/zitadel/zitadel/internal/api/grpc/filter/v2beta"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	authorization "github.com/zitadel/zitadel/pkg/grpc/authorization/v2beta"
	filter_pb "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
)

func (s *Server) ListAuthorizations(ctx context.Context, req *connect.Request[authorization.ListAuthorizationsRequest]) (*connect.Response[authorization.ListAuthorizationsResponse], error) {
	queries, err := s.listAuthorizationsRequestToModel(req.Msg)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.UserGrants(ctx, queries, false, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&authorization.ListAuthorizationsResponse{
		Authorizations: userGrantsToPb(resp.UserGrants),
		Pagination:     filter.QueryToPaginationPb(queries.SearchRequest, resp.SearchResponse),
	}), nil
}

func (s *Server) listAuthorizationsRequestToModel(req *authorization.ListAuthorizationsRequest) (*query.UserGrantsQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(s.systemDefaults, req.Pagination)
	if err != nil {
		return nil, err
	}
	queries, err := AuthorizationQueriesToQuery(req.Filters)
	if err != nil {
		return nil, err
	}
	return &query.UserGrantsQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: authorizationFieldNameToSortingColumn(req.GetSortingColumn()),
		},
		Queries: queries,
	}, nil
}

func authorizationFieldNameToSortingColumn(field authorization.AuthorizationFieldName) query.Column {
	switch field {
	case authorization.AuthorizationFieldName_AUTHORIZATION_FIELD_NAME_UNSPECIFIED:
		return query.UserGrantCreationDate
	case authorization.AuthorizationFieldName_AUTHORIZATION_FIELD_NAME_CREATED_DATE:
		return query.UserGrantCreationDate
	case authorization.AuthorizationFieldName_AUTHORIZATION_FIELD_NAME_CHANGED_DATE:
		return query.UserGrantChangeDate
	case authorization.AuthorizationFieldName_AUTHORIZATION_FIELD_NAME_ID:
		return query.UserGrantID
	case authorization.AuthorizationFieldName_AUTHORIZATION_FIELD_NAME_USER_ID:
		return query.UserGrantUserID
	case authorization.AuthorizationFieldName_AUTHORIZATION_FIELD_NAME_PROJECT_ID:
		return query.UserGrantProjectID
	case authorization.AuthorizationFieldName_AUTHORIZATION_FIELD_NAME_ORGANIZATION_ID:
		return query.UserGrantResourceOwner
	case authorization.AuthorizationFieldName_AUTHORIZATION_FIELD_NAME_USER_ORGANIZATION_ID:
		return query.UserResourceOwnerCol
	default:
		return query.UserGrantCreationDate
	}
}

func AuthorizationQueriesToQuery(queries []*authorization.AuthorizationsSearchFilter) (q []query.SearchQuery, err error) {
	q = make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = AuthorizationSearchFilterToQuery(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func AuthorizationSearchFilterToQuery(query *authorization.AuthorizationsSearchFilter) (query.SearchQuery, error) {
	switch q := query.Filter.(type) {
	case *authorization.AuthorizationsSearchFilter_AuthorizationIds:
		return AuthorizationIDQueryToModel(q.AuthorizationIds)
	case *authorization.AuthorizationsSearchFilter_OrganizationId:
		return AuthorizationOrganizationIDQueryToModel(q.OrganizationId)
	case *authorization.AuthorizationsSearchFilter_State:
		return AuthorizationStateQueryToModel(q.State)
	case *authorization.AuthorizationsSearchFilter_UserId:
		return AuthorizationUserUserIDQueryToModel(q.UserId)
	case *authorization.AuthorizationsSearchFilter_UserOrganizationId:
		return AuthorizationUserOrganizationIDQueryToModel(q.UserOrganizationId)
	case *authorization.AuthorizationsSearchFilter_UserPreferredLoginName:
		return AuthorizationUserNameQueryToModel(q.UserPreferredLoginName)
	case *authorization.AuthorizationsSearchFilter_UserDisplayName:
		return AuthorizationDisplayNameQueryToModel(q.UserDisplayName)
	case *authorization.AuthorizationsSearchFilter_ProjectId:
		return AuthorizationProjectIDQueryToModel(q.ProjectId)
	case *authorization.AuthorizationsSearchFilter_ProjectName:
		return AuthorizationProjectNameQueryToModel(q.ProjectName)
	case *authorization.AuthorizationsSearchFilter_RoleKey:
		return AuthorizationRoleKeyQueryToModel(q.RoleKey)
	case *authorization.AuthorizationsSearchFilter_ProjectGrantId:
		return AuthorizationProjectGrantIDQueryToModel(q.ProjectGrantId)
	default:
		return nil, errors.New("invalid query")
	}
}

func AuthorizationIDQueryToModel(q *filter_pb.InIDsFilter) (query.SearchQuery, error) {
	return query.NewUserGrantInIDsSearchQuery(q.Ids)
}

func AuthorizationDisplayNameQueryToModel(q *authorization.UserDisplayNameQuery) (query.SearchQuery, error) {
	return query.NewUserGrantDisplayNameQuery(q.DisplayName, filter.TextMethodPbToQuery(q.Method))
}

func AuthorizationOrganizationIDQueryToModel(q *filter_pb.IDFilter) (query.SearchQuery, error) {
	return query.NewUserGrantResourceOwnerSearchQuery(q.Id)
}

func AuthorizationProjectIDQueryToModel(q *filter_pb.IDFilter) (query.SearchQuery, error) {
	return query.NewUserGrantProjectIDSearchQuery(q.Id)
}

func AuthorizationProjectNameQueryToModel(q *authorization.ProjectNameQuery) (query.SearchQuery, error) {
	return query.NewUserGrantProjectNameQuery(q.Name, filter.TextMethodPbToQuery(q.Method))
}

func AuthorizationProjectGrantIDQueryToModel(q *filter_pb.IDFilter) (query.SearchQuery, error) {
	return query.NewUserGrantGrantIDSearchQuery(q.Id)
}

func AuthorizationRoleKeyQueryToModel(q *authorization.RoleKeyQuery) (query.SearchQuery, error) {
	return query.NewUserGrantRoleQuery(q.Key)
}

func AuthorizationUserNameQueryToModel(q *authorization.UserPreferredLoginNameQuery) (query.SearchQuery, error) {
	return query.NewUserGrantUsernameQuery(q.LoginName, filter.TextMethodPbToQuery(q.Method))
}

func AuthorizationUserUserIDQueryToModel(q *filter_pb.IDFilter) (query.SearchQuery, error) {
	return query.NewUserGrantUserIDSearchQuery(q.Id)
}

func AuthorizationUserOrganizationIDQueryToModel(q *filter_pb.IDFilter) (query.SearchQuery, error) {
	return query.NewUserGrantUserResourceOwnerSearchQuery(q.Id)
}

func AuthorizationStateQueryToModel(q *authorization.StateQuery) (query.SearchQuery, error) {
	return query.NewUserGrantStateQuery(domain.UserGrantState(q.State))
}

func userGrantsToPb(userGrants []*query.UserGrant) []*authorization.Authorization {
	o := make([]*authorization.Authorization, len(userGrants))
	for i, grant := range userGrants {
		o[i] = userGrantToPb(grant)
	}
	return o
}

func userGrantToPb(userGrant *query.UserGrant) *authorization.Authorization {
	var grantID, grantedOrgID *string
	if userGrant.GrantID != "" {
		grantID = &userGrant.GrantID
	}
	if userGrant.GrantedOrgID != "" {
		grantedOrgID = &userGrant.GrantedOrgID
	}
	return &authorization.Authorization{
		Id:                    userGrant.ID,
		ProjectId:             userGrant.ProjectID,
		ProjectName:           userGrant.ProjectName,
		ProjectOrganizationId: userGrant.ProjectResourceOwner,
		ProjectGrantId:        grantID,
		GrantedOrganizationId: grantedOrgID,
		OrganizationId:        userGrant.ResourceOwner,
		CreationDate:          timestamppb.New(userGrant.CreationDate),
		ChangeDate:            timestamppb.New(userGrant.ChangeDate),
		State:                 userGrantStateToPb(userGrant.State),
		User: &authorization.User{
			Id:                 userGrant.UserID,
			PreferredLoginName: userGrant.PreferredLoginName,
			DisplayName:        userGrant.DisplayName,
			AvatarUrl:          userGrant.AvatarURL,
			OrganizationId:     userGrant.UserResourceOwner,
		},
		Roles: userGrant.Roles,
	}
}

func userGrantStateToPb(state domain.UserGrantState) authorization.State {
	switch state {
	case domain.UserGrantStateActive:
		return authorization.State_STATE_ACTIVE
	case domain.UserGrantStateInactive:
		return authorization.State_STATE_INACTIVE
	case domain.UserGrantStateUnspecified, domain.UserGrantStateRemoved:
		return authorization.State_STATE_UNSPECIFIED
	default:
		return authorization.State_STATE_UNSPECIFIED
	}
}
