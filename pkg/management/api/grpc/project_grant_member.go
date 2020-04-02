package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) SearchProjectGrantMembers(ctx context.Context, request *ProjectGrantMemberSearchRequest) (*ProjectGrantMemberSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-pldE4", "Not implemented")
}

func (s *Server) AddProjectGrantMember(ctx context.Context, in *ProjectGrantMemberAdd) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-po8r3", "Not implemented")
}

func (s *Server) ChangeProjectGrantMember(ctx context.Context, in *ProjectGrantMemberChange) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-asd3c", "Not implemented")
}

func (s *Server) RemoveProjectGrantMember(ctx context.Context, in *ProjectGrantMemberRemove) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-04kfs", "Not implemented")
}
