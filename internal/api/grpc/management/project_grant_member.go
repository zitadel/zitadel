package management

import (
	"context"

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
	member, err := s.project.AddProjectGrantMember(ctx, projectGrantMemberAddToModel(in))
	if err != nil {
		return nil, err
	}
	return projectGrantMemberFromModel(member), nil
}

func (s *Server) ChangeProjectGrantMember(ctx context.Context, in *management.ProjectGrantMemberChange) (*management.ProjectGrantMember, error) {
	member, err := s.project.ChangeProjectGrantMember(ctx, projectGrantMemberChangeToModel(in))
	if err != nil {
		return nil, err
	}
	return projectGrantMemberFromModel(member), nil
}

func (s *Server) RemoveProjectGrantMember(ctx context.Context, in *management.ProjectGrantMemberRemove) (*empty.Empty, error) {
	err := s.project.RemoveProjectGrantMember(ctx, in.ProjectId, in.GrantId, in.UserId)
	return &empty.Empty{}, err
}
