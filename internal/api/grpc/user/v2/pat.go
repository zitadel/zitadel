package user

import (
	"context"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"google.golang.org/protobuf/types/known/timestamppb"

	z_oidc "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) AddPersonalAccessToken(ctx context.Context, req *user.AddPersonalAccessTokenRequest) (*user.AddPersonalAccessTokenResponse, error) {
	pat := &command.PersonalAccessToken{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.UserId,
		},
		AllowedUserType: domain.UserTypeMachine,
		ExpirationDate:  req.ExpirationDate.AsTime(),
		Scopes:          []string{oidc.ScopeOpenID, oidc.ScopeProfile, z_oidc.ScopeUserMetaData, z_oidc.ScopeResourceOwner},
	}
	details, err := s.command.AddPersonalAccessToken(ctx, pat, false)
	if err != nil {
		return nil, err
	}
	return &user.AddPersonalAccessTokenResponse{
		CreationDate: timestamppb.New(details.EventDate),
		TokenId:      pat.TokenID,
		Token:        pat.Token,
	}, nil
}

func (s *Server) RemovePersonalAccessToken(ctx context.Context, req *user.RemovePersonalAccessTokenRequest) (*user.RemovePersonalAccessTokenResponse, error) {
	objectDetails, err := s.command.RemovePersonalAccessToken(ctx, &command.PersonalAccessToken{TokenID: req.TokenId}, false, false)
	if err != nil {
		return nil, err
	}
	return &user.RemovePersonalAccessTokenResponse{
		DeletionDate: timestamppb.New(objectDetails.EventDate),
	}, nil
}
