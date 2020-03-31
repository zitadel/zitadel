package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) CreateProject(ctx context.Context, in *ProjectCreateRequest) (*Project, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-mo34X", "Not implemented")
}
func (s *Server) UpdateProject(ctx context.Context, in *ProjectUpdateRequest) (*Project, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-0o4fB", "Not implemented")
}
func (s *Server) DeactivateProject(ctx context.Context, in *ProjectID) (*Project, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-4Sck8", "Not implemented")
}
func (s *Server) ReactivateProject(ctx context.Context, in *ProjectID) (*Project, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-0oVre", "Not implemented")
}

func (s *Server) SearchProjects(ctx context.Context, in *ProjectSearchRequest) (*ProjectSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-2sFvd", "Not implemented")
}

func (s *Server) ProjectByID(ctx context.Context, id *ProjectID) (*Project, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-plV5x", "Not implemented")
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
