package authorization

import (
	"context"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	authorization "github.com/zitadel/zitadel/pkg/grpc/authorization/v2beta"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateAuthorization(ctx context.Context, req *authorization.CreateAuthorizationRequest) (*authorization.CreateAuthorizationResponse, error) {
	grant := &domain.UserGrant{
		UserID:    req.UserId,
		ProjectID: req.ProjectId,
		RoleKeys:  req.RoleKeys,
	}
	grant, err := s.command.AddUserGrant(ctx, grant, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &authorization.CreateAuthorizationResponse{
		Id:           grant.AggregateID,
		CreationDate: timestamppb.New(grant.CreationDate),
	}, nil
}

func (s *Server) UpdateAuthorization(ctx context.Context, request *authorization.UpdateAuthorizationRequest) (*authorization.UpdateAuthorizationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) DeleteAuthorization(ctx context.Context, request *authorization.DeleteAuthorizationRequest) (*authorization.DeleteAuthorizationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) ActivateAuthorization(ctx context.Context, request *authorization.ActivateAuthorizationRequest) (*authorization.ActivateAuthorizationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) DeactivateAuthorization(ctx context.Context, request *authorization.DeactivateAuthorizationRequest) (*authorization.DeactivateAuthorizationResponse, error) {
	//TODO implement me
	panic("implement me")
}
