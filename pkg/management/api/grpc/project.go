package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) CreateProject(ctx context.Context, in *ProjectCreateRequest) (*Project, error) {
	project, err := s.project.CreateProject(ctx, in.Name)
	if err != nil {
		return nil, err
	}
	return projectFromModel(project), nil
}
func (s *Server) UpdateProject(ctx context.Context, in *ProjectUpdateRequest) (*Project, error) {
	project, err := s.project.UpdateProject(ctx, projectUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return projectFromModel(project), nil
}
func (s *Server) DeactivateProject(ctx context.Context, in *ProjectID) (*Project, error) {
	project, err := s.project.DeactivateProject(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return projectFromModel(project), nil
}
func (s *Server) ReactivateProject(ctx context.Context, in *ProjectID) (*Project, error) {
	project, err := s.project.ReactivateProject(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return projectFromModel(project), nil
}

func (s *Server) SearchGrantedProjects(ctx context.Context, in *GrantedProjectSearchRequest) (*GrantedProjectSearchResponse, error) {
	response, err := s.project.SearchGrantedProjects(ctx, grantedProjectSearchRequestsToModel(in))
	if err != nil {
		return nil, err
	}
	return grantedProjectSearchResponseFromModel(response), nil
}

func (s *Server) ProjectByID(ctx context.Context, id *ProjectID) (*Project, error) {
	project, err := s.project.ProjectByID(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return projectFromModel(project), nil
}

func (s *Server) GetGrantedProjectGrantByID(ctx context.Context, in *ProjectGrantID) (*ProjectGrant, error) {
	project, err := s.project.ProjectGrantByID(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return projectGrantFromModel(project), nil
}

func (s *Server) AddProjectRole(ctx context.Context, in *ProjectRoleAdd) (*ProjectRole, error) {
	role, err := s.project.AddProjectRole(ctx, projectRoleAddToModel(in))
	if err != nil {
		return nil, err
	}
	return projectRoleFromModel(role), nil
}
func (s *Server) ChangeProjectRole(ctx context.Context, in *ProjectRoleChange) (*ProjectRole, error) {
	role, err := s.project.ChangeProjectRole(ctx, projectRoleChangeToModel(in))
	if err != nil {
		return nil, err
	}
	return projectRoleFromModel(role), nil
}

func (s *Server) RemoveProjectRole(ctx context.Context, in *ProjectRoleRemove) (*empty.Empty, error) {
	err := s.project.RemoveProjectRole(ctx, in.Id, in.Key)
	return &empty.Empty{}, err
}

func (s *Server) SearchProjectRoles(ctx context.Context, in *ProjectRoleSearchRequest) (*ProjectRoleSearchResponse, error) {
	response, err := s.project.SearchProjectRoles(ctx, projectRoleSearchRequestsToModel(in))
	if err != nil {
		return nil, err
	}
	return projectRoleSearchResponseFromModel(response), nil
}

func (s *Server) ProjectChanges(ctx context.Context, changesRequest *ChangeRequest) (*Changes, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-mci3f", "Not implemented")
}
