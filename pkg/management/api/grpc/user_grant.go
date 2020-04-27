package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) SearchUserGrants(ctx context.Context, request *UserGrantSearchRequest) (*UserGrantSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-dk3ds", "Not implemented")
}

func (s *Server) UserGrantByID(ctx context.Context, request *UserGrantID) (*UserGrant, error) {
	user, err := s.user.UserGrantByID(ctx, request.UserId, request.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) CreateUserGrant(ctx context.Context, in *UserGrantCreate) (*UserGrant, error) {
	user, err := s.user.AddUserGrant(ctx, userGrantCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}
func (s *Server) UpdateUserGrant(ctx context.Context, in *UserGrantUpdate) (*UserGrant, error) {
	user, err := s.user.ChangeUserGrant(ctx, userGrantUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}
func (s *Server) DeactivateUserGrant(ctx context.Context, in *UserGrantID) (*UserGrant, error) {
	user, err := s.user.DeactivateUserGrant(ctx, in.UserId, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}
func (s *Server) ReactivateUserGrant(ctx context.Context, in *UserGrantID) (*UserGrant, error) {
	user, err := s.user.ReactivateUserGrant(ctx, in.UserId, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) RemoveUserGrant(ctx context.Context, in *UserGrantID) (*empty.Empty, error) {
	err := s.user.RemoveUserGrant(ctx, in.UserId, in.Id)
	return &empty.Empty{}, err
}

func (s *Server) SearchProjectUserGrants(ctx context.Context, request *ProjectUserGrantSearchRequest) (*UserGrantSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-8jdSw", "Not implemented")
}

func (s *Server) ProjectUserGrantByID(ctx context.Context, request *ProjectUserGrantID) (*UserGrant, error) {
	user, err := s.user.UserGrantByID(ctx, request.UserId, request.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) CreateProjectUserGrant(ctx context.Context, in *UserGrantCreate) (*UserGrant, error) {
	user, err := s.user.AddUserGrant(ctx, userGrantCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}
func (s *Server) UpdateProjectUserGrant(ctx context.Context, in *ProjectUserGrantUpdate) (*UserGrant, error) {
	user, err := s.user.ChangeUserGrant(ctx, projectUserGrantUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) DeactivateProjectUserGrant(ctx context.Context, in *ProjectUserGrantID) (*UserGrant, error) {
	user, err := s.user.DeactivateUserGrant(ctx, in.UserId, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) ReactivateProjectUserGrant(ctx context.Context, in *ProjectUserGrantID) (*UserGrant, error) {
	user, err := s.user.ReactivateUserGrant(ctx, in.UserId, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) SearchProjectGrantUserGrants(ctx context.Context, request *ProjectGrantUserGrantSearchRequest) (*UserGrantSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-32sFs", "Not implemented")
}

func (s *Server) ProjectGrantUserGrantByID(ctx context.Context, request *ProjectGrantUserGrantID) (*UserGrant, error) {
	user, err := s.user.UserGrantByID(ctx, request.UserId, request.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) CreateProjectGrantUserGrant(ctx context.Context, in *ProjectGrantUserGrantCreate) (*UserGrant, error) {
	user, err := s.user.ChangeUserGrant(ctx, projectGrantUserGrantCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}
func (s *Server) UpdateProjectGrantUserGrant(ctx context.Context, in *ProjectGrantUserGrantUpdate) (*UserGrant, error) {
	user, err := s.user.ChangeUserGrant(ctx, projectGrantUserGrantUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) DeactivateProjectGrantUserGrant(ctx context.Context, in *ProjectGrantUserGrantID) (*UserGrant, error) {
	user, err := s.user.DeactivateUserGrant(ctx, in.UserId, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) ReactivateProjectGrantUserGrant(ctx context.Context, in *ProjectGrantUserGrantID) (*UserGrant, error) {
	user, err := s.user.ReactivateUserGrant(ctx, in.UserId, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}
