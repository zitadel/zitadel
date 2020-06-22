package grpc

import (
	"context"

	"github.com/caos/zitadel/internal/api"
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) SearchProjectGrants(ctx context.Context, in *ProjectGrantSearchRequest) (*ProjectGrantSearchResponse, error) {
	request := projectGrantSearchRequestsToModel(in)
	orgID := grpc_util.GetHeader(ctx, api.ZitadelOrgID)
	request.AppendMyResourceOwnerQuery(orgID)
	request.AppendNotMyOrgQuery(orgID)
	response, err := s.project.SearchProjectGrants(ctx, request)
	if err != nil {
		return nil, err
	}
	return projectGrantSearchResponseFromModel(response), nil
}

func (s *Server) ProjectGrantByID(ctx context.Context, in *ProjectGrantID) (*ProjectGrantView, error) {
	grant, err := s.project.ProjectGrantByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return projectGrantFromGrantedProjectModel(grant), nil
}

func (s *Server) CreateProjectGrant(ctx context.Context, in *ProjectGrantCreate) (*ProjectGrant, error) {
	grant, err := s.project.AddProjectGrant(ctx, projectGrantCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return projectGrantFromModel(grant), nil
}
func (s *Server) UpdateProjectGrant(ctx context.Context, in *ProjectGrantUpdate) (*ProjectGrant, error) {
	grant, err := s.project.ChangeProjectGrant(ctx, projectGrantUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return projectGrantFromModel(grant), nil
}
func (s *Server) DeactivateProjectGrant(ctx context.Context, in *ProjectGrantID) (*ProjectGrant, error) {
	grant, err := s.project.DeactivateProjectGrant(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return projectGrantFromModel(grant), nil
}

func (s *Server) ReactivateProjectGrant(ctx context.Context, in *ProjectGrantID) (*ProjectGrant, error) {
	grant, err := s.project.ReactivateProjectGrant(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return projectGrantFromModel(grant), nil
}

func (s *Server) RemoveProjectGrant(ctx context.Context, in *ProjectGrantID) (*empty.Empty, error) {
	err := s.project.RemoveProjectGrant(ctx, in.ProjectId, in.Id)
	return &empty.Empty{}, err
}
