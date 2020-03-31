package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetProjectMemberRoles(ctx context.Context, _ *empty.Empty) (*ProjectMemberRoles, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-qw34d", "Not implemented")
}

func (s *Server) SearchProjectMembers(ctx context.Context, request *ProjectMemberSearchRequest) (*ProjectMemberSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-PLr84", "Not implemented")
}

func (s *Server) AddProjectMember(ctx context.Context, in *ProjectMemberAdd) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-c2dks", "Not implemented")
}

func (s *Server) ChangeProjectMember(ctx context.Context, in *ProjectMemberChange) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-cms47", "Not implemented")
}

func (s *Server) RemoveProjectMember(ctx context.Context, in *ProjectMemberRemove) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-olw21", "Not implemented")
}
