package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetProjectGrantMemberRoles(ctx context.Context, _ *empty.Empty) (*ProjectGrantMemberRoles, error) {
	return &ProjectGrantMemberRoles{Roles: s.getProjectGrantMemberRoles()}, nil
}

func (s *Server) SearchProjectGrantMembers(ctx context.Context, in *ProjectGrantMemberSearchRequest) (*ProjectGrantMemberSearchResponse, error) {
	response, err := s.project.SearchProjectGrantMembers(ctx, projectGrantMemberSearchRequestsToModel(in))
	if err != nil {
		return nil, err
	}
	return projectGrantMemberSearchResponseFromModel(response), nil
}

func (s *Server) AddProjectGrantMember(ctx context.Context, in *ProjectGrantMemberAdd) (*ProjectGrantMember, error) {
	member, err := s.project.AddProjectGrantMember(ctx, projectGrantMemberAddToModel(in))
	if err != nil {
		return nil, err
	}
	return projectGrantMemberFromModel(member), nil
}

func (s *Server) ChangeProjectGrantMember(ctx context.Context, in *ProjectGrantMemberChange) (*ProjectGrantMember, error) {
	member, err := s.project.ChangeProjectGrantMember(ctx, projectGrantMemberChangeToModel(in))
	if err != nil {
		return nil, err
	}
	return projectGrantMemberFromModel(member), nil
}

func (s *Server) RemoveProjectGrantMember(ctx context.Context, in *ProjectGrantMemberRemove) (*empty.Empty, error) {
	err := s.project.RemoveProjectGrantMember(ctx, in.ProjectId, in.GrantId, in.UserId)
	return &empty.Empty{}, err
}
