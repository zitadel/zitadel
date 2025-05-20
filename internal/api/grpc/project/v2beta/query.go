package project

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	filter "github.com/zitadel/zitadel/internal/api/grpc/filter/v2beta"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	project_pb "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func (s *Server) GetProject(ctx context.Context, req *project_pb.GetProjectRequest) (*project_pb.GetProjectResponse, error) {
	project, err := s.query.GetProjectByIDWithPermission(ctx, true, req.Id, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return &project_pb.GetProjectResponse{
		Project: projectToPb(project),
	}, nil
}

func (s *Server) ListProjects(ctx context.Context, req *project_pb.ListProjectsRequest) (*project_pb.ListProjectsResponse, error) {
	queries, err := s.listProjectRequestToModel(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchGrantedProjects(ctx, queries, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return &project_pb.ListProjectsResponse{
		Projects:   grantedProjectsToPb(resp.GrantedProjects),
		Pagination: filter.QueryToPaginationPb(queries.SearchRequest, resp.SearchResponse),
	}, nil
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

func projectInIDsFilterToQuery(q *project_pb.InProjectIDsFilter) (query.SearchQuery, error) {
	return query.NewGrantedProjectIDSearchQuery(q.ProjectIds)
}

func projectResourceOwnerFilterToQuery(q *project_pb.ProjectResourceOwnerFilter) (query.SearchQuery, error) {
	return query.NewGrantedProjectResourceOwnerSearchQuery(q.ProjectResourceOwner)
}

func projectOrganizationIDFilterToQuery(q *project_pb.ProjectOrganizationIDFilter) (query.SearchQuery, error) {
	return query.NewGrantedProjectOrganizationIDSearchQuery(q.ProjectOrganizationId)
}

func projectGrantResourceOwnerFilterToQuery(q *project_pb.ProjectGrantResourceOwnerFilter) (query.SearchQuery, error) {
	return query.NewGrantedProjectGrantResourceOwnerSearchQuery(q.ProjectGrantResourceOwner)
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
	return &project_pb.Project{
		Id:                     project.ProjectID,
		OrganizationId:         project.ResourceOwner,
		CreationDate:           timestamppb.New(project.CreationDate),
		ChangeDate:             timestamppb.New(project.ChangeDate),
		State:                  projectStateToPb(project.ProjectState),
		Name:                   project.ProjectName,
		PrivateLabelingSetting: privateLabelingSettingToPb(project.PrivateLabelingSetting),
		ProjectAccessRequired:  project.HasProjectCheck,
		ProjectRoleAssertion:   project.ProjectRoleAssertion,
		AuthorizationRequired:  project.ProjectRoleCheck,
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
