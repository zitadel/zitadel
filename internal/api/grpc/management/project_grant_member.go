package management

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/management/grpc"
)

func (s *Server) GetProjectGrantMemberRoles(ctx context.Context, _ *empty.Empty) (*grpc.ProjectGrantMemberRoles, error) {
	return &grpc.ProjectGrantMemberRoles{Roles: s.project.GetProjectGrantMemberRoles()}, nil
}

func (s *Server) SearchProjectGrantMembers(ctx context.Context, in *grpc.ProjectGrantMemberSearchRequest) (*grpc.ProjectGrantMemberSearchResponse, error) {
	response, err := s.project.SearchProjectGrantMembers(ctx, projectGrantMemberSearchRequestsToModel(in))
	if err != nil {
		return nil, err
	}
	return projectGrantMemberSearchResponseFromModel(response), nil
}

func (s *Server) AddProjectGrantMember(ctx context.Context, in *grpc.ProjectGrantMemberAdd) (*grpc.ProjectGrantMember, error) {
	member, err := s.project.AddProjectGrantMember(ctx, projectGrantMemberAddToModel(in))
	if err != nil {
		return nil, err
	}
	return projectGrantMemberFromModel(member), nil
}

func (s *Server) ChangeProjectGrantMember(ctx context.Context, in *grpc.ProjectGrantMemberChange) (*grpc.ProjectGrantMember, error) {
	member, err := s.project.ChangeProjectGrantMember(ctx, projectGrantMemberChangeToModel(in))
	if err != nil {
		return nil, err
	}
	return projectGrantMemberFromModel(member), nil
}

func (s *Server) RemoveProjectGrantMember(ctx context.Context, in *grpc.ProjectGrantMemberRemove) (*empty.Empty, error) {
	err := s.project.RemoveProjectGrantMember(ctx, in.ProjectId, in.GrantId, in.UserId)
	return &empty.Empty{}, err
}
