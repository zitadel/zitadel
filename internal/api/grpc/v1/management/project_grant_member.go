package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetProjectGrantMemberRoles(ctx context.Context, _ *empty.Empty) (*management.ProjectGrantMemberRoles, error) {
	return &management.ProjectGrantMemberRoles{Roles: s.project.GetProjectGrantMemberRoles()}, nil
}

func (s *Server) SearchProjectGrantMembers(ctx context.Context, in *management.ProjectGrantMemberSearchRequest) (*management.ProjectGrantMemberSearchResponse, error) {
	response, err := s.project.SearchProjectGrantMembers(ctx, projectGrantMemberSearchRequestsToModel(in))
	if err != nil {
		return nil, err
	}
	return projectGrantMemberSearchResponseFromModel(response), nil
}

func (s *Server) AddProjectGrantMember(ctx context.Context, in *management.ProjectGrantMemberAdd) (*management.ProjectGrantMember, error) {
	member, err := s.command.AddProjectGrantMember(ctx, projectGrantMemberAddToDomain(in), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return projectGrantMemberFromDomain(member), nil
}

func (s *Server) ChangeProjectGrantMember(ctx context.Context, in *management.ProjectGrantMemberChange) (*management.ProjectGrantMember, error) {
	member, err := s.command.ChangeProjectGrantMember(ctx, projectGrantMemberChangeToDomain(in), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return projectGrantMemberFromDomain(member), nil
}

func (s *Server) RemoveProjectGrantMember(ctx context.Context, in *management.ProjectGrantMemberRemove) (*empty.Empty, error) {
	err := s.command.RemoveProjectGrantMember(ctx, in.ProjectId, in.UserId, in.GrantId, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}
