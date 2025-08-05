package user

import (
	"context"

	"connectrpc.com/connect"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"google.golang.org/protobuf/types/known/timestamppb"

	z_oidc "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) AddPersonalAccessToken(ctx context.Context, req *connect.Request[user.AddPersonalAccessTokenRequest]) (*connect.Response[user.AddPersonalAccessTokenResponse], error) {
	newPat := &command.PersonalAccessToken{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.Msg.GetUserId(),
		},
		PermissionCheck: s.command.NewPermissionCheckUserWrite(ctx),
		ExpirationDate:  req.Msg.GetExpirationDate().AsTime(),
		Scopes: []string{
			oidc.ScopeOpenID,
			oidc.ScopeProfile,
			z_oidc.ScopeUserMetaData,
			z_oidc.ScopeResourceOwner,
		},
		AllowedUserType: domain.UserTypeMachine,
	}
	details, err := s.command.AddPersonalAccessToken(ctx, newPat)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.AddPersonalAccessTokenResponse{
		CreationDate: timestamppb.New(details.EventDate),
		TokenId:      newPat.TokenID,
		Token:        newPat.Token,
	}), nil
}

func (s *Server) RemovePersonalAccessToken(ctx context.Context, req *connect.Request[user.RemovePersonalAccessTokenRequest]) (*connect.Response[user.RemovePersonalAccessTokenResponse], error) {
	objectDetails, err := s.command.RemovePersonalAccessToken(ctx, &command.PersonalAccessToken{
		TokenID: req.Msg.GetTokenId(),
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.Msg.GetUserId(),
		},
		PermissionCheck: s.command.NewPermissionCheckUserWrite(ctx),
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.RemovePersonalAccessTokenResponse{
		DeletionDate: timestamppb.New(objectDetails.EventDate),
	}), nil
}
