package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetProjectMemberRoles(ctx context.Context, _ *empty.Empty) (*management.ProjectMemberRoles, error) {
	roles, err := s.project.GetProjectMemberRoles(ctx)
	if err != nil {
		return nil, err
	}
	return &management.ProjectMemberRoles{Roles: roles}, nil
}

func (s *Server) SearchProjectMembers(ctx context.Context, in *management.ProjectMemberSearchRequest) (*management.ProjectMemberSearchResponse, error) {
	request := projectMemberSearchRequestsToModel(in)
	request.AppendProjectQuery(in.ProjectId)
	response, err := s.project.SearchProjectMembers(ctx, request)
	if err != nil {
		return nil, err
	}
	return projectMemberSearchResponseFromModel(response), nil
}

func (s *Server) AddProjectMember(ctx context.Context, in *management.ProjectMemberAdd) (*management.ProjectMember, error) {
	member, err := s.command.AddProjectMember(ctx, projectMemberAddToDomain(in), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return projectMemberFromDomain(member), nil
}

func (s *Server) ChangeProjectMember(ctx context.Context, in *management.ProjectMemberChange) (*management.ProjectMember, error) {
	member, err := s.command.ChangeProjectMember(ctx, projectMemberChangeToDomain(in), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return projectMemberFromDomain(member), nil
}

func (s *Server) RemoveProjectMember(ctx context.Context, in *management.ProjectMemberRemove) (*empty.Empty, error) {
	err := s.command.RemoveProjectMember(ctx, in.Id, in.UserId, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}
