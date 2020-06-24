package management

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/pkg/management/grpc"
)

func (s *Server) SearchProjectGrants(ctx context.Context, in *grpc.ProjectGrantSearchRequest) (*grpc.ProjectGrantSearchResponse, error) {
	request := projectGrantSearchRequestsToModel(in)
	orgID := grpc_util.GetHeader(ctx, http.ZitadelOrgID)
	request.AppendMyResourceOwnerQuery(orgID)
	request.AppendNotMyOrgQuery(orgID)
	response, err := s.project.SearchProjectGrants(ctx, request)
	if err != nil {
		return nil, err
	}
	return projectGrantSearchResponseFromModel(response), nil
}

func (s *Server) ProjectGrantByID(ctx context.Context, in *grpc.ProjectGrantID) (*grpc.ProjectGrantView, error) {
	grant, err := s.project.ProjectGrantByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return projectGrantFromGrantedProjectModel(grant), nil
}

func (s *Server) CreateProjectGrant(ctx context.Context, in *grpc.ProjectGrantCreate) (*grpc.ProjectGrant, error) {
	grant, err := s.project.AddProjectGrant(ctx, projectGrantCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return projectGrantFromModel(grant), nil
}
func (s *Server) UpdateProjectGrant(ctx context.Context, in *grpc.ProjectGrantUpdate) (*grpc.ProjectGrant, error) {
	grant, err := s.project.ChangeProjectGrant(ctx, projectGrantUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return projectGrantFromModel(grant), nil
}
func (s *Server) DeactivateProjectGrant(ctx context.Context, in *grpc.ProjectGrantID) (*grpc.ProjectGrant, error) {
	grant, err := s.project.DeactivateProjectGrant(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return projectGrantFromModel(grant), nil
}

func (s *Server) ReactivateProjectGrant(ctx context.Context, in *grpc.ProjectGrantID) (*grpc.ProjectGrant, error) {
	grant, err := s.project.ReactivateProjectGrant(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return projectGrantFromModel(grant), nil
}

func (s *Server) RemoveProjectGrant(ctx context.Context, in *grpc.ProjectGrantID) (*empty.Empty, error) {
	err := s.project.RemoveProjectGrant(ctx, in.ProjectId, in.Id)
	return &empty.Empty{}, err
}
