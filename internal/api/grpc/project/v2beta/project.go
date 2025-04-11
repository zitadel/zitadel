package project

import (
	"context"

	"github.com/muhlemmer/gu"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	filter "github.com/zitadel/zitadel/internal/api/grpc/filter/v2beta"
	object_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	project_grpc "github.com/zitadel/zitadel/internal/api/grpc/project"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	project_pb "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func (s *Server) CreateProject(ctx context.Context, req *project_pb.CreateProjectRequest) (*project_pb.CreateProjectResponse, error) {
	add := projectCreateToCommand(req)
	project, err := s.command.AddProject(ctx, add)
	if err != nil {
		return nil, err
	}
	var creationDate *timestamppb.Timestamp
	if !project.EventDate.IsZero() {
		creationDate = timestamppb.New(project.EventDate)
	}
	return &project_pb.CreateProjectResponse{
		Id:           add.AggregateID,
		CreationDate: creationDate,
	}, nil
}

func projectCreateToCommand(req *project_pb.CreateProjectRequest) *command.AddProject {
	var aggregateID string
	if req.Id != nil {
		aggregateID = *req.Id
	}
	return &command.AddProject{
		ObjectRoot: models.ObjectRoot{
			ResourceOwner: req.OrganizationId,
			AggregateID:   aggregateID,
		},
		Name:                   req.Name,
		ProjectRoleAssertion:   req.ProjectRoleAssertion,
		ProjectRoleCheck:       req.ProjectRoleCheck,
		HasProjectCheck:        req.HasProjectCheck,
		PrivateLabelingSetting: privateLabelingSettingToDomain(req.PrivateLabelingSetting),
	}
}

func privateLabelingSettingToDomain(setting project_pb.PrivateLabelingSetting) domain.PrivateLabelingSetting {
	switch setting {
	case project_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY:
		return domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy
	case project_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY:
		return domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy
	default:
		return domain.PrivateLabelingSettingUnspecified
	}
}

func (s *Server) UpdateProject(ctx context.Context, req *project_pb.UpdateProjectRequest) (*project_pb.UpdateProjectResponse, error) {
	project, err := s.command.ChangeProject(ctx, projectUpdateToCommand(req))
	if err != nil {
		return nil, err
	}
	var changeDate *timestamppb.Timestamp
	if !project.EventDate.IsZero() {
		changeDate = timestamppb.New(project.EventDate)
	}
	return &project_pb.UpdateProjectResponse{
		ChangeDate: changeDate,
	}, nil
}

func projectUpdateToCommand(req *project_pb.UpdateProjectRequest) *command.ChangeProject {
	var labeling *domain.PrivateLabelingSetting
	if req.PrivateLabelingSetting != nil {
		labeling = gu.Ptr(privateLabelingSettingToDomain(*req.PrivateLabelingSetting))
	}
	return &command.ChangeProject{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.Id,
		},
		Name:                   req.Name,
		ProjectRoleAssertion:   req.ProjectRoleAssertion,
		ProjectRoleCheck:       req.ProjectRoleCheck,
		HasProjectCheck:        req.HasProjectCheck,
		PrivateLabelingSetting: labeling,
	}
}

func (s *Server) DeleteProject(ctx context.Context, req *project_pb.DeleteProjectRequest) (*project_pb.DeleteProjectResponse, error) {
	projectQuery, err := query.NewUserGrantProjectIDSearchQuery(req.Id)
	if err != nil {
		return nil, err
	}
	grants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: []query.SearchQuery{projectQuery},
	}, true)
	if err != nil {
		return nil, err
	}

	deletedAt, err := s.command.DeleteProject(ctx, req.Id, "", userGrantsToIDs(grants.UserGrants)...)
	if err != nil {
		return nil, err
	}
	var deletionDate *timestamppb.Timestamp
	if !deletedAt.IsZero() {
		deletionDate = timestamppb.New(deletedAt)
	}
	return &project_pb.DeleteProjectResponse{
		DeletionDate: deletionDate,
	}, nil
}

func (s *Server) DeactivateProject(ctx context.Context, req *project_pb.DeactivateProjectRequest) (*project_pb.DeactivateProjectResponse, error) {
	details, err := s.command.DeactivateProject(ctx, req.Id, "")
	if err != nil {
		return nil, err
	}
	var changeDate *timestamppb.Timestamp
	if !details.EventDate.IsZero() {
		changeDate = timestamppb.New(details.EventDate)
	}
	return &project_pb.DeactivateProjectResponse{
		ChangeDate: changeDate,
	}, nil
}

func (s *Server) ActivateProject(ctx context.Context, req *project_pb.ActivateProjectRequest) (*project_pb.ActivateProjectResponse, error) {
	details, err := s.command.ReactivateProject(ctx, req.Id, "")
	if err != nil {
		return nil, err
	}
	var changeDate *timestamppb.Timestamp
	if !details.EventDate.IsZero() {
		changeDate = timestamppb.New(details.EventDate)
	}
	return &project_pb.ActivateProjectResponse{
		ChangeDate: changeDate,
	}, nil
}

func (s *Server) RemoveProject(ctx context.Context, req *mgmt_pb.RemoveProjectRequest) (*mgmt_pb.RemoveProjectResponse, error) {
	projectQuery, err := query.NewUserGrantProjectIDSearchQuery(req.Id)
	if err != nil {
		return nil, err
	}
	grants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: []query.SearchQuery{projectQuery},
	}, true)
	if err != nil {
		return nil, err
	}
	details, err := s.command.RemoveProject(ctx, req.Id, authz.GetCtxData(ctx).OrgID, userGrantsToIDs(grants.UserGrants)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveProjectResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func userGrantsToIDs(userGrants []*query.UserGrant) []string {
	converted := make([]string, len(userGrants))
	for i, grant := range userGrants {
		converted[i] = grant.ID
	}
	return converted
}

func (s *Server) GetProjectByID(ctx context.Context, req *mgmt_pb.GetProjectByIDRequest) (*mgmt_pb.GetProjectByIDResponse, error) {
	project, err := s.query.ProjectByID(ctx, true, req.Id)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetProjectByIDResponse{
		Project: project_grpc.ProjectViewToPb(project),
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
	resp, err := s.query.SearchProjects(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &project_pb.ListProjectsResponse{
		Result:     projectsToPb(resp.Projects),
		Pagination: filter.QueryToPaginationPb(queries.SearchRequest, resp.SearchResponse),
	}, nil
}

func (s *Server) listProjectRequestToModel(req *project_pb.ListProjectsRequest) (*query.ProjectSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(s.systemDefaults, req.Pagination)
	queries, err := projectFiltersToQuery(req.Filters)
	if err != nil {
		return nil, err
	}
	return &query.ProjectSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
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
	case *project_pb.ProjectSearchFilter_ProjectsOnlyFilter:
		return projectOnlyProjectsFilterToQuery(q.ProjectsOnlyFilter)
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

func projectOnlyProjectsFilterToQuery(q *project_pb.ProjectsOnlyFilter) (query.SearchQuery, error) {
	// TODO
	return nil, nil
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
		HasProjectCheck:        project.HasProjectCheck,
		ProjectRoleAssertion:   project.ProjectRoleAssertion,
		ProjectRoleCheck:       project.ProjectRoleCheck,
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
