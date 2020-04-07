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

func (s *Server) SearchProjects(ctx context.Context, in *ProjectSearchRequest) (*ProjectSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-2sFvd", "Not implemented")
}

func (s *Server) ProjectByID(ctx context.Context, id *ProjectID) (*Project, error) {
	project, err := s.project.ProjectByID(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return projectFromModel(project), nil
}

func (s *Server) GetGrantedProjectGrantByID(ctx context.Context, request *GrantedGrantID) (*ProjectGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-974vd", "Not implemented")
}

func (s *Server) AddProjectRole(ctx context.Context, in *ProjectRoleAdd) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-0ow2C", "Not implemented")
}
func (s *Server) RemoveProjectRole(ctx context.Context, in *ProjectRoleRemove) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-bm6iB", "Not implemented")
}

func (s *Server) SearchProjectRoles(ctx context.Context, in *ProjectRoleSearchRequest) (*ProjectRoleSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-plV56", "Not implemented")
}

func (s *Server) ProjectChanges(ctx context.Context, changesRequest *ChangeRequest) (*Changes, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-mci3f", "Not implemented")
}
