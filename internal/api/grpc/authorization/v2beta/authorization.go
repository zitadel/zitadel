package authorization

import (
	"context"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	authorization "github.com/zitadel/zitadel/pkg/grpc/authorization/v2beta"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateAuthorization(ctx context.Context, req *authorization.CreateAuthorizationRequest) (*authorization.CreateAuthorizationResponse, error) {
	grant := &domain.UserGrant{
		UserID:    req.UserId,
		ProjectID: req.ProjectId,
		RoleKeys:  req.RoleKeys,
	}
	grant, err := s.command.AddUserGrant(ctx, grant, req.OrganizationId, s.command.NewPermissionCheckUserGrantWrite(ctx))
	if err != nil {
		return nil, err
	}
	return &authorization.CreateAuthorizationResponse{
		Id:           grant.AggregateID,
		CreationDate: timestamppb.New(grant.ChangeDate),
	}, nil
}

func (s *Server) UpdateAuthorization(ctx context.Context, request *authorization.UpdateAuthorizationRequest) (*authorization.UpdateAuthorizationResponse, error) {
	changedUserGrant, err := s.command.ChangeUserGrant(ctx, &domain.UserGrant{
		ObjectRoot: models.ObjectRoot{AggregateID: request.Id},
		RoleKeys:   request.RoleKeys,
	}, nil, true, s.command.NewPermissionCheckUserGrantWrite(ctx))
	if err != nil {
		return nil, err
	}
	return &authorization.UpdateAuthorizationResponse{
		ChangeDate: timestamppb.New(changedUserGrant.ChangeDate),
	}, nil
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
