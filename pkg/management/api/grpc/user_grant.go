package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
)

func (s *Server) SearchUserGrants(ctx context.Context, request *UserGrantSearchRequest) (*UserGrantSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-dk3ds", "Not implemented")
}

func (s *Server) UserGrantByID(ctx context.Context, request *UserGrantID) (*UserGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-9dksF", "Not implemented")
}

func (s *Server) CreateUserGrant(ctx context.Context, in *UserGrantCreate) (*UserGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-2kdl2", "Not implemented")
}
func (s *Server) UpdateUserGrant(ctx context.Context, in *UserGrantUpdate) (*UserGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-83jsF", "Not implemented")
}
func (s *Server) DeactivateUserGrant(ctx context.Context, in *UserGrantID) (*UserGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-93dj3", "Not implemented")
}
func (s *Server) ReactivateUserGrant(ctx context.Context, in *UserGrantID) (*UserGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-2kSfs", "Not implemented")
}

func (s *Server) SearchProjectUserGrants(ctx context.Context, request *ProjectUserGrantSearchRequest) (*UserGrantSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-8jdSw", "Not implemented")
}

func (s *Server) ProjectUserGrantByID(ctx context.Context, request *ProjectUserGrantID) (*UserGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-dk32s", "Not implemented")
}

func (s *Server) CreateProjectUserGrant(ctx context.Context, in *UserGrantCreate) (*UserGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-0or5G", "Not implemented")
}
func (s *Server) UpdateProjectUserGrant(ctx context.Context, in *ProjectUserGrantUpdate) (*UserGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-asl4D", "Not implemented")
}

func (s *Server) DeactivateProjectUserGrant(ctx context.Context, in *ProjectUserGrantID) (*UserGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-2fG6h", "Not implemented")
}

func (s *Server) ReactivateProjectUserGrant(ctx context.Context, in *ProjectUserGrantID) (*UserGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-03kSc", "Not implemented")
}

func (s *Server) SearchProjectGrantUserGrants(ctx context.Context, request *ProjectGrantUserGrantSearchRequest) (*UserGrantSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-32sFs", "Not implemented")
}

func (s *Server) ProjectGrantUserGrantByID(ctx context.Context, request *ProjectGrantUserGrantID) (*UserGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-9kfSc", "Not implemented")
}

func (s *Server) CreateProjectGrantUserGrant(ctx context.Context, in *ProjectGrantUserGrantCreate) (*UserGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-293md", "Not implemented")
}
func (s *Server) UpdateProjectGrantUserGrant(ctx context.Context, in *ProjectGrantUserGrantUpdate) (*UserGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-76fGe", "Not implemented")
}

func (s *Server) DeactivateProjectGrantUserGrant(ctx context.Context, in *ProjectGrantUserGrantID) (*UserGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-sFsi3", "Not implemented")
}

func (s *Server) ReactivateProjectGrantUserGrant(ctx context.Context, in *ProjectGrantUserGrantID) (*UserGrant, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-ckr56", "Not implemented")
}
