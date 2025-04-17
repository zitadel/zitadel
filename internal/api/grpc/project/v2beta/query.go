package project

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	filter "github.com/zitadel/zitadel/internal/api/grpc/filter/v2beta"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	project_pb "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func (s *Server) GetProject(ctx context.Context, req *project_pb.GetProjectRequest) (*project_pb.GetProjectResponse, error) {
	project, err := s.query.ProjectByID(ctx, true, req.Id)
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
	err = queries.AppendPermissionQueries(authz.GetRequestPermissionsFromCtx(ctx))
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchProjects(ctx, queries, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return &project_pb.ListProjectsResponse{
		Projects:   projectsToPb(resp.Projects),
		Pagination: filter.QueryToPaginationPb(queries.SearchRequest, resp.SearchResponse),
	}, nil
}

func (s *Server) listProjectRequestToModel(req *project_pb.ListProjectsRequest) (*query.ProjectSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(s.systemDefaults, req.Pagination)
	queriesProjects, queriesGrants, err := projectFiltersToQuery(req.Filters)
	if err != nil {
		return nil, err
	}
	return &query.ProjectSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries:      queriesProjects,
		GrantQueries: queriesGrants,
	}, nil
}

func projectFiltersToQuery(queries []*project_pb.ProjectSearchFilter) (_ []query.SearchQuery, _ []query.SearchQuery, err error) {
	q1 := make([]query.SearchQuery, len(queries))
	for i, qry := range queries {
		q1[i], err = projectFilterToModel(qry)
		if err != nil {
			return nil, nil, err
		}
	}
	q2 := make([]query.SearchQuery, len(queries))
	for i, qry := range queries {
		q2[i], err = projectGrantFilterToModel(qry)
		if err != nil {
			return nil, nil, err
		}
	}
	return q1, q2, nil
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
		return nil, nil
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-vR9nC", "List.Query.Invalid")
	}
}

func projectGrantFilterToModel(filter *project_pb.ProjectSearchFilter) (query.SearchQuery, error) {
	switch q := filter.Filter.(type) {
	case *project_pb.ProjectSearchFilter_ProjectGrantResourceOwnerFilter:
		return projectGrantResourceOwnerFilterToQuery(q.ProjectGrantResourceOwnerFilter)
	case *project_pb.ProjectSearchFilter_ProjectOrganizationIdFilter:
		return projectGrantOrganizationIDFilterToQuery(q.ProjectOrganizationIdFilter)
	case *project_pb.ProjectSearchFilter_ProjectNameFilter,
		*project_pb.ProjectSearchFilter_InProjectIdsFilter,
		*project_pb.ProjectSearchFilter_ProjectResourceOwnerFilter:
		return nil, nil
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-vR9nC", "List.Query.Invalid")
	}
}

func projectNameFilterToQuery(q *project_pb.ProjectNameFilter) (query.SearchQuery, error) {
	return query.NewProjectNameSearchQuery(filter.TextMethodPbToQuery(q.Method), q.GetProjectName())
}

func projectInIDsFilterToQuery(q *project_pb.InProjectIDsFilter) (query.SearchQuery, error) {
	return query.NewProjectIDSearchQuery(q.ProjectIds)
}

func projectResourceOwnerFilterToQuery(q *project_pb.ProjectResourceOwnerFilter) (query.SearchQuery, error) {
	return query.NewProjectResourceOwnerSearchQuery(q.ProjectResourceOwner)
}

func projectOrganizationIDFilterToQuery(q *project_pb.ProjectOrganizationIDFilter) (query.SearchQuery, error) {
	return query.NewProjectResourceOwnerSearchQuery(q.ProjectOrganizationId)
}

func projectGrantResourceOwnerFilterToQuery(q *project_pb.ProjectGrantResourceOwnerFilter) (query.SearchQuery, error) {
	return query.NewProjectResourceOwnerSearchQuery(q.ProjectGrantResourceOwner)
}

func projectGrantOrganizationIDFilterToQuery(q *project_pb.ProjectOrganizationIDFilter) (query.SearchQuery, error) {
	return query.NewProjectGrantGrantedOrgIDSearchQuery(q.ProjectOrganizationId)
}

func projectsToPb(projects []*query.Project) []*project_pb.Project {
	o := make([]*project_pb.Project, len(projects))
	for i, org := range projects {
		o[i] = projectToPb(org)
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

func projectStateToPb(state domain.ProjectState) project_pb.ProjectState {
	switch state {
	case domain.ProjectStateActive:
		return project_pb.ProjectState_PROJECT_STATE_ACTIVE
	case domain.ProjectStateInactive:
		return project_pb.ProjectState_PROJECT_STATE_INACTIVE
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
	default:
		return project_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_UNSPECIFIED
	}
}
