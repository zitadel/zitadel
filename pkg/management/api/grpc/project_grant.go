package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetProjectGrantMemberRoles(ctx context.Context, _ *empty.Empty) (*ProjectGrantMemberRoles, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-mGo89", "Not implemented")
}

func (s *Server) SearchProjectGrants(ctx context.Context, request *ProjectGrantSearchRequest) (*ProjectGrantSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-po9fs", "Not implemented")
}

func (s *Server) ProjectGrantByID(ctx context.Context, request *ProjectGrantID) (*ProjectGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-nmr54", "Not implemented")
}

func (s *Server) CreateProjectGrant(ctx context.Context, in *ProjectGrantCreate) (*ProjectGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-fi45f", "Not implemented")
}
func (s *Server) UpdateProjectGrant(ctx context.Context, in *ProjectGrantUpdate) (*ProjectGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-nm7Ds", "Not implemented")
}
func (s *Server) DeactivateProjectGrant(ctx context.Context, in *ProjectGrantID) (*ProjectGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-xkwpr", "Not implemented")
}
func (s *Server) ReactivateProjectGrant(ctx context.Context, in *ProjectGrantID) (*ProjectGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-mdk23", "Not implemented")
}
