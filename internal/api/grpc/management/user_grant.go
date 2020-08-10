package management

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) SearchUserGrants(ctx context.Context, in *management.UserGrantSearchRequest) (*management.UserGrantSearchResponse, error) {
	request := userGrantSearchRequestsToModel(in)
	request.AppendMyOrgQuery(authz.GetCtxData(ctx).OrgID)
	response, err := s.usergrant.SearchUserGrants(ctx, request)
	if err != nil {
		return nil, err
	}
	return userGrantSearchResponseFromModel(response), nil
}

func (s *Server) UserGrantByID(ctx context.Context, request *management.UserGrantID) (*management.UserGrantView, error) {
	user, err := s.usergrant.UserGrantByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return userGrantViewFromModel(user), nil
}

func (s *Server) CreateUserGrant(ctx context.Context, in *management.UserGrantCreate) (*management.UserGrant, error) {
	user, err := s.usergrant.AddUserGrant(ctx, userGrantCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) UpdateUserGrant(ctx context.Context, in *management.UserGrantUpdate) (*management.UserGrant, error) {
	user, err := s.usergrant.ChangeUserGrant(ctx, userGrantUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) DeactivateUserGrant(ctx context.Context, in *management.UserGrantID) (*management.UserGrant, error) {
	user, err := s.usergrant.DeactivateUserGrant(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}
func (s *Server) ReactivateUserGrant(ctx context.Context, in *management.UserGrantID) (*management.UserGrant, error) {
	user, err := s.usergrant.ReactivateUserGrant(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) RemoveUserGrant(ctx context.Context, in *management.UserGrantID) (*empty.Empty, error) {
	err := s.usergrant.RemoveUserGrant(ctx, in.Id)
	return &empty.Empty{}, err
}

func (s *Server) BulkRemoveUserGrant(ctx context.Context, in *management.UserGrantRemoveBulk) (*empty.Empty, error) {
	err := s.usergrant.BulkRemoveUserGrant(ctx, userGrantRemoveBulkToModel(in)...)
	return &empty.Empty{}, err
}
