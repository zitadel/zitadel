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
	//TODO: Check explicit Permissions
	user, err := s.command.AddUserGrant(ctx, userGrantCreateToDomain(in), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return userGrantFromDomain(user), nil
}

func (s *Server) UpdateUserGrant(ctx context.Context, in *management.UserGrantUpdate) (*management.UserGrant, error) {
	//TODO: Check explicit Permissions
	user, err := s.command.ChangeUserGrant(ctx, userGrantUpdateToDomain(in), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return userGrantFromDomain(user), nil
}

func (s *Server) DeactivateUserGrant(ctx context.Context, in *management.UserGrantID) (*empty.Empty, error) {
	//TODO: Check explicit Permissions
	err := s.command.DeactivateUserGrant(ctx, in.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}
func (s *Server) ReactivateUserGrant(ctx context.Context, in *management.UserGrantID) (*empty.Empty, error) {
	//TODO: Check explicit Permissions
	err := s.command.ReactivateUserGrant(ctx, in.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) RemoveUserGrant(ctx context.Context, in *management.UserGrantID) (*empty.Empty, error) {
	//TODO: Check explicit Permissions
	err := s.command.RemoveUserGrant(ctx, in.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) BulkRemoveUserGrant(ctx context.Context, in *management.UserGrantRemoveBulk) (*empty.Empty, error) {
	//TODO: Check explicit Permissions
	err := s.command.BulkRemoveUserGrant(ctx, userGrantRemoveBulkToModel(in), authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}
