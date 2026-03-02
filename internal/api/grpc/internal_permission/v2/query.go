package internal_permission

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	filter_pb "github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	"github.com/zitadel/zitadel/pkg/grpc/internal_permission/v2"
)

func (s *Server) ListAdministrators(ctx context.Context, req *connect.Request[internal_permission.ListAdministratorsRequest]) (*connect.Response[internal_permission.ListAdministratorsResponse], error) {
	queries, err := s.listAdministratorsRequestToModel(req.Msg)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchAdministrators(ctx, queries, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&internal_permission.ListAdministratorsResponse{
		Administrators: administratorsToPb(resp.Administrators),
		Pagination:     filter.QueryToPaginationPb(queries.SearchRequest, resp.SearchResponse),
	}), nil
}

func (s *Server) listAdministratorsRequestToModel(req *internal_permission.ListAdministratorsRequest) (*query.MembershipSearchQuery, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(s.systemDefaults, req.Pagination)
	if err != nil {
		return nil, err
	}
	queries, err := administratorSearchFiltersToQuery(req.Filters, 0) // start at level 0
	if err != nil {
		return nil, err
	}
	return &query.MembershipSearchQuery{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: administratorFieldNameToSortingColumn(req.GetSortingColumn()),
		},
		Queries: queries,
	}, nil
}

func administratorFieldNameToSortingColumn(field internal_permission.AdministratorFieldName) query.Column {
	switch field {
	case internal_permission.AdministratorFieldName_ADMINISTRATOR_FIELD_NAME_CREATION_DATE:
		return query.MembershipCreationDate
	case internal_permission.AdministratorFieldName_ADMINISTRATOR_FIELD_NAME_USER_ID:
		return query.MembershipUserID
	case internal_permission.AdministratorFieldName_ADMINISTRATOR_FIELD_NAME_CHANGE_DATE:
		return query.MembershipChangeDate
	case internal_permission.AdministratorFieldName_ADMINISTRATOR_FIELD_NAME_UNSPECIFIED:
		return query.MembershipCreationDate
	default:
		return query.MembershipCreationDate
	}
}

func administratorSearchFiltersToQuery(queries []*internal_permission.AdministratorSearchFilter, level uint8) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, qry := range queries {
		q[i], err = administratorFilterToQuery(qry, level)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func administratorFilterToQuery(filter *internal_permission.AdministratorSearchFilter, level uint8) (query.SearchQuery, error) {
	if level > 20 {
		// can't go deeper than 20 levels of nesting.
		return nil, zerrors.ThrowInvalidArgument(nil, "PERM-zsQ97", "Errors.Query.TooManyNestingLevels")
	}
	switch q := filter.Filter.(type) {
	case *internal_permission.AdministratorSearchFilter_InUserIdsFilter:
		return inUserIDsFilterToQuery(q.InUserIdsFilter)
	case *internal_permission.AdministratorSearchFilter_CreationDate:
		return creationDateFilterToQuery(q.CreationDate)
	case *internal_permission.AdministratorSearchFilter_ChangeDate:
		return changeDateFilterToQuery(q.ChangeDate)
	case *internal_permission.AdministratorSearchFilter_UserOrganizationId:
		return userResourceOwnerFilterToQuery(q.UserOrganizationId)
	case *internal_permission.AdministratorSearchFilter_UserPreferredLoginName:
		return userLoginNameFilterToQuery(q.UserPreferredLoginName)
	case *internal_permission.AdministratorSearchFilter_UserDisplayName:
		return userDisplayNameFilterToQuery(q.UserDisplayName)
	case *internal_permission.AdministratorSearchFilter_Resource:
		return resourceFilterToQuery(q.Resource)
	case *internal_permission.AdministratorSearchFilter_Role:
		return roleFilterToQuery(q.Role)
	case *internal_permission.AdministratorSearchFilter_And:
		return andQueryToQuery(q.And, level)
	case *internal_permission.AdministratorSearchFilter_Or:
		return orQueryToQuery(q.Or, level)
	case *internal_permission.AdministratorSearchFilter_Not:
		return notQueryToQuery(q.Not, level)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "PERM-vR9nC", "List.Query.Invalid")
	}
}

func inUserIDsFilterToQuery(q *filter_pb.InIDsFilter) (query.SearchQuery, error) {
	return query.NewMemberInUserIDsSearchQuery(q.GetIds())
}

func userResourceOwnerFilterToQuery(q *filter_pb.IDFilter) (query.SearchQuery, error) {
	return query.NewAdministratorUserResourceOwnerSearchQuery(q.GetId())
}

func userLoginNameFilterToQuery(q *internal_permission.UserPreferredLoginNameFilter) (query.SearchQuery, error) {
	return query.NewAdministratorUserLoginNameSearchQuery(q.GetPreferredLoginName())
}

func userDisplayNameFilterToQuery(q *internal_permission.UserDisplayNameFilter) (query.SearchQuery, error) {
	return query.NewAdministratorUserDisplayNameSearchQuery(q.GetDisplayName())
}

func creationDateFilterToQuery(q *filter_pb.TimestampFilter) (query.SearchQuery, error) {
	return query.NewMembershipCreationDateQuery(q.GetTimestamp().AsTime(), filter.TimestampMethodPbToQuery(q.Method))
}

func changeDateFilterToQuery(q *filter_pb.TimestampFilter) (query.SearchQuery, error) {
	return query.NewMembershipChangeDateQuery(q.GetTimestamp().AsTime(), filter.TimestampMethodPbToQuery(q.Method))
}

func resourceFilterToQuery(q *internal_permission.ResourceFilter) (query.SearchQuery, error) {
	switch r := q.GetResource().(type) {
	case *internal_permission.ResourceFilter_Instance:
		if q.GetInstance() {
			return query.NewMembershipIsIAMQuery()
		}
	case *internal_permission.ResourceFilter_OrganizationId:
		return query.NewMembershipOrgIDQuery(r.OrganizationId)
	case *internal_permission.ResourceFilter_ProjectId:
		projectIDQuery, err := query.NewMembershipProjectIDQuery(r.ProjectId)
		if err != nil {
			return nil, err
		}
		notGrantedQuery, err := query.NewMembershipNotGrantedSearchQuery()
		if err != nil {
			return nil, err
		}
		return query.NewAndQuery(projectIDQuery, notGrantedQuery)
	case *internal_permission.ResourceFilter_ProjectGrant_:
		projectIDQuery, err := query.NewMembershipProjectIDQuery(r.ProjectGrant.GetProjectId())
		if err != nil {
			return nil, err
		}
		grantedOrganizationID, err := query.NewMembershipGrantedOrgIDSearchQuery(r.ProjectGrant.GetOrganizationId())
		if err != nil {
			return nil, err
		}
		return query.NewAndQuery(projectIDQuery, grantedOrganizationID)
	}
	return nil, nil
}

func roleFilterToQuery(q *internal_permission.RoleFilter) (query.SearchQuery, error) {
	return query.NewMembershipRoleQuery(q.GetRoleKey())
}

func orQueryToQuery(q *internal_permission.OrFilter, level uint8) (query.SearchQuery, error) {
	mappedQueries, err := administratorSearchFiltersToQuery(q.Queries, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewOrQuery(mappedQueries...)
}
func andQueryToQuery(q *internal_permission.AndFilter, level uint8) (query.SearchQuery, error) {
	mappedQueries, err := administratorSearchFiltersToQuery(q.Queries, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewAndQuery(mappedQueries...)
}
func notQueryToQuery(q *internal_permission.NotFilter, level uint8) (query.SearchQuery, error) {
	mappedQuery, err := administratorFilterToQuery(q.Query, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewNotQuery(mappedQuery)
}

func administratorsToPb(administrators []*query.Administrator) []*internal_permission.Administrator {
	a := make([]*internal_permission.Administrator, len(administrators))
	for i, admin := range administrators {
		a[i] = administratorToPb(admin)
	}
	return a
}

func administratorToPb(admin *query.Administrator) *internal_permission.Administrator {
	var resource internal_permission.Resource
	if admin.Instance != nil {
		resource = &internal_permission.Administrator_Instance{Instance: true}
	}
	if admin.Org != nil {
		resource = &internal_permission.Administrator_Organization{
			Organization: &internal_permission.Organization{
				Id:   admin.Org.OrgID,
				Name: admin.Org.Name,
			},
		}
	}
	if admin.Project != nil {
		resource = &internal_permission.Administrator_Project{
			Project: &internal_permission.Project{
				Id:             admin.Project.ProjectID,
				Name:           admin.Project.Name,
				OrganizationId: admin.Project.ResourceOwner,
			},
		}
	}
	if admin.ProjectGrant != nil {
		resource = &internal_permission.Administrator_ProjectGrant{
			ProjectGrant: &internal_permission.ProjectGrant{
				Id:                    admin.ProjectGrant.GrantID,
				ProjectId:             admin.ProjectGrant.ProjectID,
				ProjectName:           admin.ProjectGrant.ProjectName,
				OrganizationId:        admin.ProjectGrant.ResourceOwner,
				GrantedOrganizationId: admin.ProjectGrant.GrantedOrgID,
			},
		}
	}

	return &internal_permission.Administrator{
		CreationDate: timestamppb.New(admin.CreationDate),
		ChangeDate:   timestamppb.New(admin.ChangeDate),
		User: &internal_permission.User{
			Id:                 admin.User.UserID,
			PreferredLoginName: admin.User.LoginName,
			DisplayName:        admin.User.DisplayName,
			OrganizationId:     admin.User.ResourceOwner,
		},
		Resource: resource,
		Roles:    admin.Roles,
	}
}
