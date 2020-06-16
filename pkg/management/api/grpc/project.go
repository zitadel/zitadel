package grpc

import (
	"context"

	"github.com/caos/zitadel/internal/api"
	"github.com/caos/zitadel/internal/api/auth"
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
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
	if s.IsZitadel(ctx, in.Id) {
		return nil, errors.ThrowInvalidArgument(nil, "GRPC-SFH8d", "Zitadel Project should not be updated")
	}
	project, err := s.project.UpdateProject(ctx, projectUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return projectFromModel(project), nil
}
func (s *Server) DeactivateProject(ctx context.Context, in *ProjectID) (*Project, error) {
	if s.IsZitadel(ctx, in.Id) {
		return nil, errors.ThrowInvalidArgument(nil, "GRPC-PS9cs", "Zitadel Project should not be deactivated")
	}
	project, err := s.project.DeactivateProject(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return projectFromModel(project), nil
}
func (s *Server) ReactivateProject(ctx context.Context, in *ProjectID) (*Project, error) {
	if s.IsZitadel(ctx, in.Id) {
		return nil, errors.ThrowInvalidArgument(nil, "GRPC-LSpe2", "Zitadel Project should not be reactivated")
	}
	project, err := s.project.ReactivateProject(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return projectFromModel(project), nil
}

func (s *Server) SearchProjects(ctx context.Context, in *ProjectSearchRequest) (*ProjectSearchResponse, error) {
	request := projectSearchRequestsToModel(in)
	request.AppendMyResourceOwnerQuery(grpc_util.GetHeader(ctx, api.ZitadelOrgID))
	response, err := s.project.SearchProjects(ctx, request)
	if err != nil {
		return nil, err
	}
	return projectSearchResponseFromModel(response), nil
}

func (s *Server) ProjectByID(ctx context.Context, id *ProjectID) (*Project, error) {
	project, err := s.project.ProjectByID(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return projectFromModel(project), nil
}

func (s *Server) SearchGrantedProjects(ctx context.Context, in *GrantedProjectSearchRequest) (*ProjectGrantSearchResponse, error) {
	request := grantedProjectSearchRequestsToModel(in)
	request.AppendMyOrgQuery(grpc_util.GetHeader(ctx, api.ZitadelOrgID))
	response, err := s.project.SearchProjectGrants(ctx, request)
	if err != nil {
		return nil, err
	}
	return projectGrantSearchResponseFromModel(response), nil
}

func (s *Server) GetGrantedProjectByID(ctx context.Context, in *ProjectGrantID) (*ProjectGrantView, error) {
	project, err := s.project.ProjectGrantViewByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return projectGrantFromGrantedProjectModel(project), nil
}

func (s *Server) AddProjectRole(ctx context.Context, in *ProjectRoleAdd) (*ProjectRole, error) {
	if s.IsZitadel(ctx, in.Id) {
		return nil, errors.ThrowInvalidArgument(nil, "GRPC-PS9cs", "Zitadel Project should not get new role")
	}
	role, err := s.project.AddProjectRole(ctx, projectRoleAddToModel(in))
	if err != nil {
		return nil, err
	}
	return projectRoleFromModel(role), nil
}
func (s *Server) ChangeProjectRole(ctx context.Context, in *ProjectRoleChange) (*ProjectRole, error) {
	if s.IsZitadel(ctx, in.Id) {
		return nil, errors.ThrowInvalidArgument(nil, "GRPC-LASj8", "Zitadel Project should not change roles")
	}
	role, err := s.project.ChangeProjectRole(ctx, projectRoleChangeToModel(in))
	if err != nil {
		return nil, err
	}
	return projectRoleFromModel(role), nil
}

func (s *Server) RemoveProjectRole(ctx context.Context, in *ProjectRoleRemove) (*empty.Empty, error) {
	if s.IsZitadel(ctx, in.Id) {
		return nil, errors.ThrowInvalidArgument(nil, "GRPC-Psn7s", "do not remove roles from Zitadel Project")
	}
	err := s.project.RemoveProjectRole(ctx, in.Id, in.Key)
	return &empty.Empty{}, err
}

func (s *Server) SearchProjectRoles(ctx context.Context, in *ProjectRoleSearchRequest) (*ProjectRoleSearchResponse, error) {
	request := projectRoleSearchRequestsToModel(in)
	request.AppendMyOrgQuery(auth.GetCtxData(ctx).OrgID)
	request.AppendProjectQuery(in.ProjectId)
	response, err := s.project.SearchProjectRoles(ctx, request)
	if err != nil {
		return nil, err
	}
	return projectRoleSearchResponseFromModel(response), nil
}

func (s *Server) ProjectChanges(ctx context.Context, changesRequest *ChangeRequest) (*Changes, error) {
	response, err := s.project.ProjectChanges(ctx, changesRequest.Id, 0, 0)
	if err != nil {
		return nil, err
	}
	return projectChangesToResponse(response, changesRequest.GetSequenceOffset(), changesRequest.GetLimit()), nil
}

func (s *Server) IsZitadel(ctx context.Context, projectID string) bool {
	iam, err := s.iam.IamByID(ctx, s.systemDefaults.IamID)
	if err != nil {
		return false
	}
	if iam.IamProjectID == projectID {
		return true
	}
	return false
}
