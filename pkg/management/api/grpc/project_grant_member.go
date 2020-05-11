package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) SearchProjectGrantMembers(ctx context.Context, request *ProjectGrantMemberSearchRequest) (*ProjectGrantMemberSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-pldE4", "Not implemented")
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
