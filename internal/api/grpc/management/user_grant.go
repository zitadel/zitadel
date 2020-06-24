package grpc

import (
	"context"

	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/pkg/management/grpc"

	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) SearchUserGrants(ctx context.Context, in *grpc.UserGrantSearchRequest) (*grpc.UserGrantSearchResponse, error) {
	request := userGrantSearchRequestsToModel(in)
	request.AppendMyOrgQuery(auth.GetCtxData(ctx).OrgID)
	response, err := s.usergrant.SearchUserGrants(ctx, request)
	if err != nil {
		return nil, err
	}
	return userGrantSearchResponseFromModel(response), nil
}

func (s *Server) UserGrantByID(ctx context.Context, request *grpc.UserGrantID) (*grpc.UserGrantView, error) {
	user, err := s.usergrant.UserGrantByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return userGrantViewFromModel(user), nil
}

func (s *Server) CreateUserGrant(ctx context.Context, in *grpc.UserGrantCreate) (*grpc.UserGrant, error) {
	user, err := s.usergrant.AddUserGrant(ctx, userGrantCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) UpdateUserGrant(ctx context.Context, in *grpc.UserGrantUpdate) (*grpc.UserGrant, error) {
	user, err := s.usergrant.ChangeUserGrant(ctx, userGrantUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}
func (s *Server) DeactivateUserGrant(ctx context.Context, in *grpc.UserGrantID) (*grpc.UserGrant, error) {
	user, err := s.usergrant.DeactivateUserGrant(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}
func (s *Server) ReactivateUserGrant(ctx context.Context, in *grpc.UserGrantID) (*grpc.UserGrant, error) {
	user, err := s.usergrant.ReactivateUserGrant(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) RemoveUserGrant(ctx context.Context, in *grpc.UserGrantID) (*empty.Empty, error) {
	err := s.usergrant.RemoveUserGrant(ctx, in.Id)
	return &empty.Empty{}, err
}

func (s *Server) BulkCreateUserGrant(ctx context.Context, in *grpc.UserGrantCreateBulk) (*empty.Empty, error) {
	err := s.usergrant.BulkAddUserGrant(ctx, userGrantCreateBulkToModel(in)...)
	return &empty.Empty{}, err
}

func (s *Server) BulkUpdateUserGrant(ctx context.Context, in *grpc.UserGrantUpdateBulk) (*empty.Empty, error) {
	err := s.usergrant.BulkChangeUserGrant(ctx, userGrantUpdateBulkToModel(in)...)
	return &empty.Empty{}, err
}

func (s *Server) BulkRemoveUserGrant(ctx context.Context, in *grpc.UserGrantRemoveBulk) (*empty.Empty, error) {
	err := s.usergrant.BulkRemoveUserGrant(ctx, userGrantRemoveBulkToModel(in)...)
	return &empty.Empty{}, err
}

func (s *Server) SearchProjectUserGrants(ctx context.Context, request *grpc.ProjectUserGrantSearchRequest) (*grpc.UserGrantSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-8jdSw", "Not implemented")
}

func (s *Server) ProjectUserGrantByID(ctx context.Context, request *grpc.ProjectUserGrantID) (*grpc.UserGrantView, error) {
	user, err := s.usergrant.UserGrantByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return userGrantViewFromModel(user), nil
}

func (s *Server) CreateProjectUserGrant(ctx context.Context, in *grpc.UserGrantCreate) (*grpc.UserGrant, error) {
	user, err := s.usergrant.AddUserGrant(ctx, userGrantCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}
func (s *Server) UpdateProjectUserGrant(ctx context.Context, in *grpc.ProjectUserGrantUpdate) (*grpc.UserGrant, error) {
	user, err := s.usergrant.ChangeUserGrant(ctx, projectUserGrantUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) DeactivateProjectUserGrant(ctx context.Context, in *grpc.ProjectUserGrantID) (*grpc.UserGrant, error) {
	user, err := s.usergrant.DeactivateUserGrant(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) ReactivateProjectUserGrant(ctx context.Context, in *grpc.ProjectUserGrantID) (*grpc.UserGrant, error) {
	user, err := s.usergrant.ReactivateUserGrant(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) SearchProjectGrantUserGrants(ctx context.Context, request *grpc.ProjectGrantUserGrantSearchRequest) (*grpc.UserGrantSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-32sFs", "Not implemented")
}

func (s *Server) ProjectGrantUserGrantByID(ctx context.Context, request *grpc.ProjectGrantUserGrantID) (*grpc.UserGrantView, error) {
	user, err := s.usergrant.UserGrantByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return userGrantViewFromModel(user), nil
}

func (s *Server) CreateProjectGrantUserGrant(ctx context.Context, in *grpc.ProjectGrantUserGrantCreate) (*grpc.UserGrant, error) {
	user, err := s.usergrant.ChangeUserGrant(ctx, projectGrantUserGrantCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}
func (s *Server) UpdateProjectGrantUserGrant(ctx context.Context, in *grpc.ProjectGrantUserGrantUpdate) (*grpc.UserGrant, error) {
	user, err := s.usergrant.ChangeUserGrant(ctx, projectGrantUserGrantUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) DeactivateProjectGrantUserGrant(ctx context.Context, in *grpc.ProjectGrantUserGrantID) (*grpc.UserGrant, error) {
	user, err := s.usergrant.DeactivateUserGrant(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) ReactivateProjectGrantUserGrant(ctx context.Context, in *grpc.ProjectGrantUserGrantID) (*grpc.UserGrant, error) {
	user, err := s.usergrant.ReactivateUserGrant(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}
