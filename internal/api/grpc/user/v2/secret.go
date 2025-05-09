package user

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) AddSecret(ctx context.Context, req *user.AddSecretRequest) (*user.AddSecretResponse, error) {
	newSecret := &command.GenerateMachineSecret{
		PermissionCheck: s.command.NewPermissionCheckUserWrite(ctx),
	}
	details, err := s.command.GenerateMachineSecret(ctx, req.UserId, "", newSecret)
	if err != nil {
		return nil, err
	}
	return &user.AddSecretResponse{
		CreationDate: timestamppb.New(details.EventDate),
		ClientSecret: newSecret.ClientSecret,
	}, nil
}

func (s *Server) RemoveSecret(ctx context.Context, req *user.RemoveSecretRequest) (*user.RemoveSecretResponse, error) {
	details, err := s.command.RemoveMachineSecret(
		ctx,
		req.UserId,
		"",
		s.command.NewPermissionCheckUserWrite(ctx),
	)
	if err != nil {
		return nil, err
	}
	return &user.RemoveSecretResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}, nil
}
