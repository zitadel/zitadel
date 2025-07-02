package authorization

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	authorization "github.com/zitadel/zitadel/pkg/grpc/authorization/v2beta"
)

func (s *Server) CreateAuthorization(ctx context.Context, req *authorization.CreateAuthorizationRequest) (*authorization.CreateAuthorizationResponse, error) {
	grant := &domain.UserGrant{
		UserID:    req.UserId,
		ProjectID: req.ProjectId,
		RoleKeys:  req.RoleKeys,
		ObjectRoot: models.ObjectRoot{
			ResourceOwner: req.GetOrganizationId(),
		},
	}
	grant, err := s.command.AddUserGrant(ctx, grant, s.command.NewPermissionCheckUserGrantWrite(ctx))
	if err != nil {
		return nil, err
	}
	return &authorization.CreateAuthorizationResponse{
		Id:           grant.AggregateID,
		CreationDate: timestamppb.New(grant.ChangeDate),
	}, nil
}

func (s *Server) UpdateAuthorization(ctx context.Context, request *authorization.UpdateAuthorizationRequest) (*authorization.UpdateAuthorizationResponse, error) {
	userGrant, err := s.command.ChangeUserGrant(ctx, &domain.UserGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: request.Id,
		},
		RoleKeys: request.RoleKeys,
	}, true, true, s.command.NewPermissionCheckUserGrantWrite(ctx))
	if err != nil {
		return nil, err
	}
	return &authorization.UpdateAuthorizationResponse{
		ChangeDate: timestamppb.New(userGrant.ChangeDate),
	}, nil
}

func (s *Server) DeleteAuthorization(ctx context.Context, request *authorization.DeleteAuthorizationRequest) (*authorization.DeleteAuthorizationResponse, error) {
	details, err := s.command.RemoveUserGrant(ctx, request.Id, "", true, s.command.NewPermissionCheckUserGrantDelete(ctx))
	if err != nil {
		return nil, err
	}
	return &authorization.DeleteAuthorizationResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}, nil
}

func (s *Server) ActivateAuthorization(ctx context.Context, request *authorization.ActivateAuthorizationRequest) (*authorization.ActivateAuthorizationResponse, error) {
	details, err := s.command.ReactivateUserGrant(ctx, request.Id, "", s.command.NewPermissionCheckUserGrantWrite(ctx))
	if err != nil {
		return nil, err
	}
	return &authorization.ActivateAuthorizationResponse{
		ChangeDate: timestamppb.New(details.EventDate),
	}, nil
}

func (s *Server) DeactivateAuthorization(ctx context.Context, request *authorization.DeactivateAuthorizationRequest) (*authorization.DeactivateAuthorizationResponse, error) {
	details, err := s.command.DeactivateUserGrant(ctx, request.Id, "", s.command.NewPermissionCheckUserGrantWrite(ctx))
	if err != nil {
		return nil, err
	}
	return &authorization.DeactivateAuthorizationResponse{
		ChangeDate: timestamppb.New(details.EventDate),
	}, nil
}
