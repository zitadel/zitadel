package project

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	filter "github.com/zitadel/zitadel/internal/api/grpc/filter/v2beta"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	filter_pb "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	project_pb "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func (s *Server) GetProject(ctx context.Context, req *connect.Request[project_pb.GetProjectRequest]) (*connect.Response[project_pb.GetProjectResponse], error) {
	project, err := s.query.GetProjectByIDWithPermission(ctx, true, req.Msg.GetId(), s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&project_pb.GetProjectResponse{
		Project: projectToPb(project),
	}), nil
}

func (s *Server) ListProjects(ctx context.Context, req *connect.Request[project_pb.ListProjectsRequest]) (*connect.Response[project_pb.ListProjectsResponse], error) {
	queries, err := s.listProjectRequestToModel(req.Msg)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchGrantedProjects(ctx, queries, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&project_pb.ListProjectsResponse{
		Projects:   grantedProjectsToPb(resp.GrantedProjects),
		Pagination: filter.QueryToPaginationPb(queries.SearchRequest, resp.SearchResponse),
	}), nil
}

func (s *Server) listProjectRequestToModel(req *project_pb.ListProjectsRequest) (*query.ProjectAndGrantedProjectSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(s.systemDefaults, req.Pagination)
	if err != nil {
		return nil, err
	}
	queries, err := projectFiltersToQuery(req.Filters)
	if err != nil {
		return nil, err
	}
	return &query.ProjectAndGrantedProjectSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: grantedProjectFieldNameToSortingColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func grantedProjectFieldNameToSortingColumn(field *project_pb.ProjectFieldName) query.Column {
	if field == nil {
		return query.GrantedProjectColumnCreationDate
	}
	switch *field {
	case project_pb.ProjectFieldName_PROJECT_FIELD_NAME_CREATION_DATE:
		return query.GrantedProjectColumnCreationDate
	case project_pb.ProjectFieldName_PROJECT_FIELD_NAME_ID:
		return query.GrantedProjectColumnID
	case project_pb.ProjectFieldName_PROJECT_FIELD_NAME_NAME:
		return query.GrantedProjectColumnName
	case project_pb.ProjectFieldName_PROJECT_FIELD_NAME_CHANGE_DATE:
		return query.GrantedProjectColumnChangeDate
	case project_pb.ProjectFieldName_PROJECT_FIELD_NAME_UNSPECIFIED:
		return query.GrantedProjectColumnCreationDate
	default:
		return query.GrantedProjectColumnCreationDate
	}
}

func projectFiltersToQuery(queries []*project_pb.ProjectSearchFilter) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, qry := range queries {
		q[i], err = projectFilterToModel(qry)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func projectFilterToModel(filter *project_pb.ProjectSearchFilter) (query.SearchQuery, error) {
	switch q := filter.Filter.(type) {
	case *project_pb.ProjectSearchFilter_ProjectNameFilter:
		return projectNameFilterToQuery(q.ProjectNameFilter)
	case *project_pb.ProjectSearchFilter_InProjectIdsFilter:
		return projectInIDsFilterToQuery(q.InProjectIdsFilter)
	case *project_pb.ProjectSearchFilter_ProjectResourceOwnerFilter:
		return projectResourceOwnerFilterToQuery(q.ProjectResourceOwnerFilter)
	case *project_pb.ProjectSearchFilter_ProjectOrganizationIdFilter:
		return projectOrganizationIDFilterToQuery(q.ProjectOrganizationIdFilter)
	case *project_pb.ProjectSearchFilter_ProjectGrantResourceOwnerFilter:
		return projectGrantResourceOwnerFilterToQuery(q.ProjectGrantResourceOwnerFilter)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-vR9nC", "List.Query.Invalid")
	}
}

func projectNameFilterToQuery(q *project_pb.ProjectNameFilter) (query.SearchQuery, error) {
	return query.NewGrantedProjectNameSearchQuery(filter.TextMethodPbToQuery(q.Method), q.GetProjectName())
}

func projectInIDsFilterToQuery(q *filter_pb.InIDsFilter) (query.SearchQuery, error) {
	return query.NewGrantedProjectIDSearchQuery(q.Ids)
}

func projectResourceOwnerFilterToQuery(q *filter_pb.IDFilter) (query.SearchQuery, error) {
	return query.NewGrantedProjectResourceOwnerSearchQuery(q.Id)
}

func projectOrganizationIDFilterToQuery(q *filter_pb.IDFilter) (query.SearchQuery, error) {
	return query.NewGrantedProjectOrganizationIDSearchQuery(q.Id)
}

func projectGrantResourceOwnerFilterToQuery(q *filter_pb.IDFilter) (query.SearchQuery, error) {
	return query.NewGrantedProjectGrantResourceOwnerSearchQuery(q.Id)
}

func grantedProjectsToPb(projects []*query.GrantedProject) []*project_pb.Project {
	o := make([]*project_pb.Project, len(projects))
	for i, org := range projects {
		o[i] = grantedProjectToPb(org)
	}
	return o
}

func projectToPb(project *query.Project) *project_pb.Project {
	return &project_pb.Project{
		Id:                     project.ID,
		OrganizationId:         project.ResourceOwner,
		CreationDate:           timestamppb.New(project.CreationDate),
		ChangeDate:             timestamppb.New(project.ChangeDate),
		State:                  projectStateToPb(project.State),
		Name:                   project.Name,
		PrivateLabelingSetting: privateLabelingSettingToPb(project.PrivateLabelingSetting),
		ProjectAccessRequired:  project.HasProjectCheck,
		ProjectRoleAssertion:   project.ProjectRoleAssertion,
		AuthorizationRequired:  project.ProjectRoleCheck,
	}
}

func grantedProjectToPb(project *query.GrantedProject) *project_pb.Project {
	var grantedOrganizationID, grantedOrganizationName *string
	if project.GrantedOrgID != "" {
		grantedOrganizationID = &project.GrantedOrgID
	}
	if project.OrgName != "" {
		grantedOrganizationName = &project.OrgName
	}

	return &project_pb.Project{
		Id:                      project.ProjectID,
		OrganizationId:          project.ResourceOwner,
		CreationDate:            timestamppb.New(project.CreationDate),
		ChangeDate:              timestamppb.New(project.ChangeDate),
		State:                   projectStateToPb(project.ProjectState),
		Name:                    project.ProjectName,
		PrivateLabelingSetting:  privateLabelingSettingToPb(project.PrivateLabelingSetting),
		ProjectAccessRequired:   project.HasProjectCheck,
		ProjectRoleAssertion:    project.ProjectRoleAssertion,
		AuthorizationRequired:   project.ProjectRoleCheck,
		GrantedOrganizationId:   grantedOrganizationID,
		GrantedOrganizationName: grantedOrganizationName,
		GrantedState:            grantedProjectStateToPb(project.ProjectGrantState),
	}
}

func projectStateToPb(state domain.ProjectState) project_pb.ProjectState {
	switch state {
	case domain.ProjectStateActive:
		return project_pb.ProjectState_PROJECT_STATE_ACTIVE
	case domain.ProjectStateInactive:
		return project_pb.ProjectState_PROJECT_STATE_INACTIVE
	case domain.ProjectStateUnspecified, domain.ProjectStateRemoved:
		return project_pb.ProjectState_PROJECT_STATE_UNSPECIFIED
	default:
		return project_pb.ProjectState_PROJECT_STATE_UNSPECIFIED
	}
}
func grantedProjectStateToPb(state domain.ProjectGrantState) project_pb.GrantedProjectState {
	switch state {
	case domain.ProjectGrantStateActive:
		return project_pb.GrantedProjectState_GRANTED_PROJECT_STATE_ACTIVE
	case domain.ProjectGrantStateInactive:
		return project_pb.GrantedProjectState_GRANTED_PROJECT_STATE_INACTIVE
	case domain.ProjectGrantStateUnspecified, domain.ProjectGrantStateRemoved:
		return project_pb.GrantedProjectState_GRANTED_PROJECT_STATE_UNSPECIFIED
	default:
		return project_pb.GrantedProjectState_GRANTED_PROJECT_STATE_UNSPECIFIED
	}
}

func privateLabelingSettingToPb(setting domain.PrivateLabelingSetting) project_pb.PrivateLabelingSetting {
	switch setting {
	case domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy:
		return project_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY
	case domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy:
		return project_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY
	case domain.PrivateLabelingSettingUnspecified:
		return project_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_UNSPECIFIED
	default:
		return project_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_UNSPECIFIED
	}
}

func (s *Server) ListProjectGrants(ctx context.Context, req *connect.Request[project_pb.ListProjectGrantsRequest]) (*connect.Response[project_pb.ListProjectGrantsResponse], error) {
	queries, err := s.listProjectGrantsRequestToModel(req.Msg)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchProjectGrants(ctx, queries, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&project_pb.ListProjectGrantsResponse{
		ProjectGrants: projectGrantsToPb(resp.ProjectGrants),
		Pagination:    filter.QueryToPaginationPb(queries.SearchRequest, resp.SearchResponse),
	}), nil
}

func (s *Server) listProjectGrantsRequestToModel(req *project_pb.ListProjectGrantsRequest) (*query.ProjectGrantSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(s.systemDefaults, req.Pagination)
	if err != nil {
		return nil, err
	}
	queries, err := projectGrantFiltersToModel(req.Filters)
	if err != nil {
		return nil, err
	}
	return &query.ProjectGrantSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: projectGrantFieldNameToSortingColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func projectGrantFieldNameToSortingColumn(field *project_pb.ProjectGrantFieldName) query.Column {
	if field == nil {
		return query.ProjectGrantColumnCreationDate
	}
	switch *field {
	case project_pb.ProjectGrantFieldName_PROJECT_GRANT_FIELD_NAME_PROJECT_ID:
		return query.ProjectGrantColumnProjectID
	case project_pb.ProjectGrantFieldName_PROJECT_GRANT_FIELD_NAME_CREATION_DATE:
		return query.ProjectGrantColumnCreationDate
	case project_pb.ProjectGrantFieldName_PROJECT_GRANT_FIELD_NAME_CHANGE_DATE:
		return query.ProjectGrantColumnChangeDate
	case project_pb.ProjectGrantFieldName_PROJECT_GRANT_FIELD_NAME_UNSPECIFIED:
		return query.ProjectGrantColumnCreationDate
	default:
		return query.ProjectGrantColumnCreationDate
	}
}

func projectGrantFiltersToModel(queries []*project_pb.ProjectGrantSearchFilter) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, qry := range queries {
		q[i], err = projectGrantFilterToModel(qry)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func projectGrantFilterToModel(filter *project_pb.ProjectGrantSearchFilter) (query.SearchQuery, error) {
	switch q := filter.Filter.(type) {
	case *project_pb.ProjectGrantSearchFilter_ProjectNameFilter:
		return projectNameFilterToQuery(q.ProjectNameFilter)
	case *project_pb.ProjectGrantSearchFilter_RoleKeyFilter:
		return query.NewProjectGrantRoleKeySearchQuery(q.RoleKeyFilter.Key)
	case *project_pb.ProjectGrantSearchFilter_InProjectIdsFilter:
		return query.NewProjectGrantProjectIDsSearchQuery(q.InProjectIdsFilter.Ids)
	case *project_pb.ProjectGrantSearchFilter_ProjectResourceOwnerFilter:
		return query.NewProjectGrantResourceOwnerSearchQuery(q.ProjectResourceOwnerFilter.Id)
	case *project_pb.ProjectGrantSearchFilter_ProjectGrantResourceOwnerFilter:
		return query.NewProjectGrantGrantedOrgIDSearchQuery(q.ProjectGrantResourceOwnerFilter.Id)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-M099f", "List.Query.Invalid")
	}
}

func projectGrantsToPb(projects []*query.ProjectGrant) []*project_pb.ProjectGrant {
	p := make([]*project_pb.ProjectGrant, len(projects))
	for i, project := range projects {
		p[i] = projectGrantToPb(project)
	}
	return p
}

func projectGrantToPb(project *query.ProjectGrant) *project_pb.ProjectGrant {
	return &project_pb.ProjectGrant{
		OrganizationId:          project.ResourceOwner,
		CreationDate:            timestamppb.New(project.CreationDate),
		ChangeDate:              timestamppb.New(project.ChangeDate),
		GrantedOrganizationId:   project.GrantedOrgID,
		GrantedOrganizationName: project.OrgName,
		GrantedRoleKeys:         project.GrantedRoleKeys,
		ProjectId:               project.ProjectID,
		ProjectName:             project.ProjectName,
		State:                   projectGrantStateToPb(project.State),
	}
}

func projectGrantStateToPb(state domain.ProjectGrantState) project_pb.ProjectGrantState {
	switch state {
	case domain.ProjectGrantStateActive:
		return project_pb.ProjectGrantState_PROJECT_GRANT_STATE_ACTIVE
	case domain.ProjectGrantStateInactive:
		return project_pb.ProjectGrantState_PROJECT_GRANT_STATE_INACTIVE
	case domain.ProjectGrantStateUnspecified, domain.ProjectGrantStateRemoved:
		return project_pb.ProjectGrantState_PROJECT_GRANT_STATE_UNSPECIFIED
	default:
		return project_pb.ProjectGrantState_PROJECT_GRANT_STATE_UNSPECIFIED
	}
}

func (s *Server) ListProjectRoles(ctx context.Context, req *connect.Request[project_pb.ListProjectRolesRequest]) (*connect.Response[project_pb.ListProjectRolesResponse], error) {
	queries, err := s.listProjectRolesRequestToModel(req.Msg)
	if err != nil {
		return nil, err
	}
	err = queries.AppendProjectIDQuery(req.Msg.GetProjectId())
	if err != nil {
		return nil, err
	}
	roles, err := s.query.SearchProjectRoles(ctx, true, queries, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&project_pb.ListProjectRolesResponse{
		ProjectRoles: roleViewsToPb(roles.ProjectRoles),
		Pagination:   filter.QueryToPaginationPb(queries.SearchRequest, roles.SearchResponse),
	}), nil
}

func (s *Server) listProjectRolesRequestToModel(req *project_pb.ListProjectRolesRequest) (*query.ProjectRoleSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(s.systemDefaults, req.Pagination)
	if err != nil {
		return nil, err
	}
	queries, err := roleQueriesToModel(req.Filters)
	if err != nil {
		return nil, err
	}
	return &query.ProjectRoleSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: projectRoleFieldNameToSortingColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func projectRoleFieldNameToSortingColumn(field *project_pb.ProjectRoleFieldName) query.Column {
	if field == nil {
		return query.ProjectRoleColumnCreationDate
	}
	switch *field {
	case project_pb.ProjectRoleFieldName_PROJECT_ROLE_FIELD_NAME_KEY:
		return query.ProjectRoleColumnKey
	case project_pb.ProjectRoleFieldName_PROJECT_ROLE_FIELD_NAME_CREATION_DATE:
		return query.ProjectRoleColumnCreationDate
	case project_pb.ProjectRoleFieldName_PROJECT_ROLE_FIELD_NAME_CHANGE_DATE:
		return query.ProjectRoleColumnChangeDate
	case project_pb.ProjectRoleFieldName_PROJECT_ROLE_FIELD_NAME_UNSPECIFIED:
		return query.ProjectRoleColumnCreationDate
	default:
		return query.ProjectRoleColumnCreationDate
	}
}

func roleQueriesToModel(queries []*project_pb.ProjectRoleSearchFilter) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = roleQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func roleQueryToModel(apiQuery *project_pb.ProjectRoleSearchFilter) (query.SearchQuery, error) {
	switch q := apiQuery.Filter.(type) {
	case *project_pb.ProjectRoleSearchFilter_RoleKeyFilter:
		return query.NewProjectRoleKeySearchQuery(filter.TextMethodPbToQuery(q.RoleKeyFilter.Method), q.RoleKeyFilter.Key)
	case *project_pb.ProjectRoleSearchFilter_DisplayNameFilter:
		return query.NewProjectRoleDisplayNameSearchQuery(filter.TextMethodPbToQuery(q.DisplayNameFilter.Method), q.DisplayNameFilter.DisplayName)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-fms0e", "List.Query.Invalid")
	}
}

func roleViewsToPb(roles []*query.ProjectRole) []*project_pb.ProjectRole {
	o := make([]*project_pb.ProjectRole, len(roles))
	for i, org := range roles {
		o[i] = roleViewToPb(org)
	}
	return o
}

func roleViewToPb(role *query.ProjectRole) *project_pb.ProjectRole {
	return &project_pb.ProjectRole{
		ProjectId:    role.ProjectID,
		Key:          role.Key,
		CreationDate: timestamppb.New(role.CreationDate),
		ChangeDate:   timestamppb.New(role.ChangeDate),
		DisplayName:  role.DisplayName,
		Group:        role.Group,
	}
}
