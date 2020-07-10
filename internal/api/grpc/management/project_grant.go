package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) SearchProjectGrants(ctx context.Context, in *management.ProjectGrantSearchRequest) (*management.ProjectGrantSearchResponse, error) {
	request := projectGrantSearchRequestsToModel(in)
	ctxData := authz.GetCtxData(ctx)
	request.AppendMyResourceOwnerQuery(ctxData.OrgID)
	response, err := s.project.SearchProjectGrants(ctx, request)
	if err != nil {
		return nil, err
	}
	return projectGrantSearchResponseFromModel(response), nil
}

func (s *Server) ProjectGrantByID(ctx context.Context, in *management.ProjectGrantID) (*management.ProjectGrantView, error) {
	grant, err := s.project.ProjectGrantByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return projectGrantFromGrantedProjectModel(grant), nil
}

func (s *Server) CreateProjectGrant(ctx context.Context, in *management.ProjectGrantCreate) (*management.ProjectGrant, error) {
	grant, err := s.project.AddProjectGrant(ctx, projectGrantCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return projectGrantFromModel(grant), nil
}
func (s *Server) UpdateProjectGrant(ctx context.Context, in *management.ProjectGrantUpdate) (*management.ProjectGrant, error) {
	grant, err := s.project.ChangeProjectGrant(ctx, projectGrantUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return projectGrantFromModel(grant), nil
}
func (s *Server) DeactivateProjectGrant(ctx context.Context, in *management.ProjectGrantID) (*management.ProjectGrant, error) {
	grant, err := s.project.DeactivateProjectGrant(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return projectGrantFromModel(grant), nil
}

func (s *Server) ReactivateProjectGrant(ctx context.Context, in *management.ProjectGrantID) (*management.ProjectGrant, error) {
	grant, err := s.project.ReactivateProjectGrant(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return projectGrantFromModel(grant), nil
}

func (s *Server) RemoveProjectGrant(ctx context.Context, in *management.ProjectGrantID) (*empty.Empty, error) {
	err := s.project.RemoveProjectGrant(ctx, in.ProjectId, in.Id)
	return &empty.Empty{}, err
}
