package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetProjectMemberRoles(ctx context.Context, _ *empty.Empty) (*ProjectMemberRoles, error) {
	return &ProjectMemberRoles{Roles: s.project.GetProjectMemberRoles()}, nil
}

func (s *Server) SearchProjectMembers(ctx context.Context, in *ProjectMemberSearchRequest) (*ProjectMemberSearchResponse, error) {
	request := projectMemberSearchRequestsToModel(in)
	request.AppendProjectQuery(in.ProjectId)
	response, err := s.project.SearchProjectMembers(ctx, request)
	if err != nil {
		return nil, err
	}
	return projectMemberSearchResponseFromModel(response), nil
}

func (s *Server) AddProjectMember(ctx context.Context, in *ProjectMemberAdd) (*ProjectMember, error) {
	member, err := s.project.AddProjectMember(ctx, projectMemberAddToModel(in))
	if err != nil {
		return nil, err
	}
	return projectMemberFromModel(member), nil
}

func (s *Server) ChangeProjectMember(ctx context.Context, in *ProjectMemberChange) (*ProjectMember, error) {
	member, err := s.project.ChangeProjectMember(ctx, projectMemberChangeToModel(in))
	if err != nil {
		return nil, err
	}
	return projectMemberFromModel(member), nil
}

func (s *Server) RemoveProjectMember(ctx context.Context, in *ProjectMemberRemove) (*empty.Empty, error) {
	err := s.project.RemoveProjectMember(ctx, in.Id, in.UserId)
	return &empty.Empty{}, err
}
