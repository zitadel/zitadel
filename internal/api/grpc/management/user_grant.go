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

func (s *Server) BulkCreateUserGrant(ctx context.Context, in *management.UserGrantCreateBulk) (*empty.Empty, error) {
	err := s.usergrant.BulkAddUserGrant(ctx, userGrantCreateBulkToModel(in)...)
	return &empty.Empty{}, err
}

func (s *Server) BulkUpdateUserGrant(ctx context.Context, in *management.UserGrantUpdateBulk) (*empty.Empty, error) {
	err := s.usergrant.BulkChangeUserGrant(ctx, userGrantUpdateBulkToModel(in)...)
	return &empty.Empty{}, err
}

func (s *Server) BulkRemoveUserGrant(ctx context.Context, in *management.UserGrantRemoveBulk) (*empty.Empty, error) {
	err := s.usergrant.BulkRemoveUserGrant(ctx, userGrantRemoveBulkToModel(in)...)
	return &empty.Empty{}, err
}

func (s *Server) SearchProjectUserGrants(ctx context.Context, in *management.ProjectUserGrantSearchRequest) (*management.UserGrantSearchResponse, error) {
	request := projectUserGrantSearchRequestsToModel(in)
	request.AppendMyOrgQuery(authz.GetCtxData(ctx).OrgID)
	request.AppendProjectIDQuery(in.ProjectId)
	response, err := s.usergrant.SearchUserGrants(ctx, request)
	if err != nil {
		return nil, err
	}
	return userGrantSearchResponseFromModel(response), nil
}

func (s *Server) ProjectUserGrantByID(ctx context.Context, request *management.ProjectUserGrantID) (*management.UserGrantView, error) {
	user, err := s.usergrant.UserGrantByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return userGrantViewFromModel(user), nil
}

func (s *Server) CreateProjectUserGrant(ctx context.Context, in *management.UserGrantCreate) (*management.UserGrant, error) {
	user, err := s.usergrant.AddUserGrant(ctx, userGrantCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}
func (s *Server) UpdateProjectUserGrant(ctx context.Context, in *management.ProjectUserGrantUpdate) (*management.UserGrant, error) {
	user, err := s.usergrant.ChangeUserGrant(ctx, projectUserGrantUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) DeactivateProjectUserGrant(ctx context.Context, in *management.ProjectUserGrantID) (*management.UserGrant, error) {
	user, err := s.usergrant.DeactivateUserGrant(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) ReactivateProjectUserGrant(ctx context.Context, in *management.ProjectUserGrantID) (*management.UserGrant, error) {
	user, err := s.usergrant.ReactivateUserGrant(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) SearchProjectGrantUserGrants(ctx context.Context, in *management.ProjectGrantUserGrantSearchRequest) (*management.UserGrantSearchResponse, error) {
	grant, err := s.project.ProjectGrantByID(ctx, in.ProjectGrantId)
	if err != nil {
		return nil, err
	}
	request := projectGrantUserGrantSearchRequestsToModel(in)
	request.AppendMyOrgQuery(authz.GetCtxData(ctx).OrgID)
	request.AppendProjectIDQuery(grant.ProjectID)
	response, err := s.usergrant.SearchUserGrants(ctx, request)
	if err != nil {
		return nil, err
	}
	return userGrantSearchResponseFromModel(response), nil
}

func (s *Server) ProjectGrantUserGrantByID(ctx context.Context, request *management.ProjectGrantUserGrantID) (*management.UserGrantView, error) {
	user, err := s.usergrant.UserGrantByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return userGrantViewFromModel(user), nil
}

func (s *Server) CreateProjectGrantUserGrant(ctx context.Context, in *management.ProjectGrantUserGrantCreate) (*management.UserGrant, error) {
	user, err := s.usergrant.AddUserGrant(ctx, projectGrantUserGrantCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}
func (s *Server) UpdateProjectGrantUserGrant(ctx context.Context, in *management.ProjectGrantUserGrantUpdate) (*management.UserGrant, error) {
	user, err := s.usergrant.ChangeUserGrant(ctx, projectGrantUserGrantUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) DeactivateProjectGrantUserGrant(ctx context.Context, in *management.ProjectGrantUserGrantID) (*management.UserGrant, error) {
	user, err := s.usergrant.DeactivateUserGrant(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}

func (s *Server) ReactivateProjectGrantUserGrant(ctx context.Context, in *management.ProjectGrantUserGrantID) (*management.UserGrant, error) {
	user, err := s.usergrant.ReactivateUserGrant(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return usergrantFromModel(user), nil
}
