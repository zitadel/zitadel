package management

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/internal/api/authz"
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) CreateProject(ctx context.Context, in *management.ProjectCreateRequest) (*management.Project, error) {
	project, err := s.project.CreateProject(ctx, in.Name)
	if err != nil {
		return nil, err
	}
	return projectFromModel(project), nil
}
func (s *Server) UpdateProject(ctx context.Context, in *management.ProjectUpdateRequest) (*management.Project, error) {
	project, err := s.project.UpdateProject(ctx, projectUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return projectFromModel(project), nil
}
func (s *Server) DeactivateProject(ctx context.Context, in *management.ProjectID) (*management.Project, error) {
	project, err := s.project.DeactivateProject(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return projectFromModel(project), nil
}
func (s *Server) ReactivateProject(ctx context.Context, in *management.ProjectID) (*management.Project, error) {
	project, err := s.project.ReactivateProject(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return projectFromModel(project), nil
}

func (s *Server) SearchProjects(ctx context.Context, in *management.ProjectSearchRequest) (*management.ProjectSearchResponse, error) {
	request := projectSearchRequestsToModel(in)
	request.AppendMyResourceOwnerQuery(grpc_util.GetHeader(ctx, http.ZitadelOrgID))
	response, err := s.project.SearchProjects(ctx, request)
	if err != nil {
		return nil, err
	}
	return projectSearchResponseFromModel(response), nil
}

func (s *Server) ProjectByID(ctx context.Context, id *management.ProjectID) (*management.ProjectView, error) {
	project, err := s.project.ProjectByID(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return projectViewFromModel(project), nil
}

func (s *Server) SearchGrantedProjects(ctx context.Context, in *management.GrantedProjectSearchRequest) (*management.ProjectGrantSearchResponse, error) {
	request := grantedProjectSearchRequestsToModel(in)
	request.AppendMyOrgQuery(grpc_util.GetHeader(ctx, http.ZitadelOrgID))
	response, err := s.project.SearchProjectGrants(ctx, request)
	if err != nil {
		return nil, err
	}
	return projectGrantSearchResponseFromModel(response), nil
}

func (s *Server) GetGrantedProjectByID(ctx context.Context, in *management.ProjectGrantID) (*management.ProjectGrantView, error) {
	project, err := s.project.ProjectGrantViewByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return projectGrantFromGrantedProjectModel(project), nil
}

func (s *Server) AddProjectRole(ctx context.Context, in *management.ProjectRoleAdd) (*management.ProjectRole, error) {
	role, err := s.project.AddProjectRole(ctx, projectRoleAddToModel(in))
	if err != nil {
		return nil, err
	}
	return projectRoleFromModel(role), nil
}

func (s *Server) BulkAddProjectRole(ctx context.Context, in *management.ProjectRoleAddBulk) (*empty.Empty, error) {
	err := s.project.BulkAddProjectRole(ctx, projectRoleAddBulkToModel(in))
	return &empty.Empty{}, err
}

func (s *Server) ChangeProjectRole(ctx context.Context, in *management.ProjectRoleChange) (*management.ProjectRole, error) {
	role, err := s.project.ChangeProjectRole(ctx, projectRoleChangeToModel(in))
	if err != nil {
		return nil, err
	}
	return projectRoleFromModel(role), nil
}

func (s *Server) RemoveProjectRole(ctx context.Context, in *management.ProjectRoleRemove) (*empty.Empty, error) {
	err := s.project.RemoveProjectRole(ctx, in.Id, in.Key)
	return &empty.Empty{}, err
}

func (s *Server) SearchProjectRoles(ctx context.Context, in *management.ProjectRoleSearchRequest) (*management.ProjectRoleSearchResponse, error) {
	request := projectRoleSearchRequestsToModel(in)
	request.AppendMyOrgQuery(authz.GetCtxData(ctx).OrgID)
	request.AppendProjectQuery(in.ProjectId)
	response, err := s.project.SearchProjectRoles(ctx, request)
	if err != nil {
		return nil, err
	}
	return projectRoleSearchResponseFromModel(response), nil
}

func (s *Server) ProjectChanges(ctx context.Context, changesRequest *management.ChangeRequest) (*management.Changes, error) {
	response, err := s.project.ProjectChanges(ctx, changesRequest.Id, changesRequest.SequenceOffset, changesRequest.Limit, changesRequest.Asc)
	if err != nil {
		return nil, err
	}
	return projectChangesToResponse(response, changesRequest.GetSequenceOffset(), changesRequest.GetLimit()), nil
}
