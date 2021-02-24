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
	grant, err := s.command.AddProjectGrant(ctx, projectGrantCreateToDomain(in), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return projectGrantFromDomain(grant), nil
}
func (s *Server) UpdateProjectGrant(ctx context.Context, in *management.ProjectGrantUpdate) (*management.ProjectGrant, error) {
	userGrants, err := s.usergrant.UserGrantsByProjectAndGrantID(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	grant, err := s.command.ChangeProjectGrant(ctx, projectGrantUpdateToDomain(in), authz.GetCtxData(ctx).OrgID, userGrantsToIDs(userGrants)...)
	if err != nil {
		return nil, err
	}
	return projectGrantFromDomain(grant), nil
}
func (s *Server) DeactivateProjectGrant(ctx context.Context, in *management.ProjectGrantID) (*empty.Empty, error) {
	err := s.command.DeactivateProjectGrant(ctx, in.ProjectId, in.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) ReactivateProjectGrant(ctx context.Context, in *management.ProjectGrantID) (*empty.Empty, error) {
	err := s.command.ReactivateProjectGrant(ctx, in.ProjectId, in.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) RemoveProjectGrant(ctx context.Context, in *management.ProjectGrantID) (*empty.Empty, error) {
	err := s.command.RemoveProjectGrant(ctx, in.ProjectId, in.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}
