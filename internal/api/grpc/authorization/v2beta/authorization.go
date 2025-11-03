package authorization

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	authorization "github.com/zitadel/zitadel/pkg/grpc/authorization/v2beta"
)

func (s *Server) CreateAuthorization(ctx context.Context, req *connect.Request[authorization.CreateAuthorizationRequest]) (*connect.Response[authorization.CreateAuthorizationResponse], error) {
	grant := &domain.UserGrant{
		UserID:    req.Msg.UserId,
		ProjectID: req.Msg.ProjectId,
		RoleKeys:  req.Msg.RoleKeys,
		ObjectRoot: models.ObjectRoot{
			ResourceOwner: req.Msg.GetOrganizationId(),
		},
	}
	grant, err := s.command.AddUserGrant(ctx, grant, s.command.NewPermissionCheckUserGrantWrite(ctx))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&authorization.CreateAuthorizationResponse{
		Id:           grant.AggregateID,
		CreationDate: timestamppb.New(grant.ChangeDate),
	}), nil
}

func (s *Server) UpdateAuthorization(ctx context.Context, request *connect.Request[authorization.UpdateAuthorizationRequest]) (*connect.Response[authorization.UpdateAuthorizationResponse], error) {
	userGrant, err := s.command.ChangeUserGrant(ctx, &domain.UserGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: request.Msg.Id,
		},
		RoleKeys: request.Msg.RoleKeys,
	}, true, true, s.command.NewPermissionCheckUserGrantWrite(ctx))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&authorization.UpdateAuthorizationResponse{
		ChangeDate: timestamppb.New(userGrant.ChangeDate),
	}), nil
}

func (s *Server) DeleteAuthorization(ctx context.Context, request *connect.Request[authorization.DeleteAuthorizationRequest]) (*connect.Response[authorization.DeleteAuthorizationResponse], error) {
	details, err := s.command.RemoveUserGrant(ctx, request.Msg.Id, "", true, s.command.NewPermissionCheckUserGrantDelete(ctx))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&authorization.DeleteAuthorizationResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}), nil
}

func (s *Server) ActivateAuthorization(ctx context.Context, request *connect.Request[authorization.ActivateAuthorizationRequest]) (*connect.Response[authorization.ActivateAuthorizationResponse], error) {
	details, err := s.command.ReactivateUserGrant(ctx, request.Msg.Id, "", s.command.NewPermissionCheckUserGrantWrite(ctx))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&authorization.ActivateAuthorizationResponse{
		ChangeDate: timestamppb.New(details.EventDate),
	}), nil
}

func (s *Server) DeactivateAuthorization(ctx context.Context, request *connect.Request[authorization.DeactivateAuthorizationRequest]) (*connect.Response[authorization.DeactivateAuthorizationResponse], error) {
	details, err := s.command.DeactivateUserGrant(ctx, request.Msg.Id, "", s.command.NewPermissionCheckUserGrantWrite(ctx))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&authorization.DeactivateAuthorizationResponse{
		ChangeDate: timestamppb.New(details.EventDate),
	}), nil
}
