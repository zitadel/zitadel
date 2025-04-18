package user

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) AddSecret(ctx context.Context, req *user.AddSecretRequest) (*user.AddSecretResponse, error) {
	newSecret := new(command.GenerateMachineSecret)
	details, err := s.command.GenerateMachineSecret(ctx, req.UserId, "", false, newSecret)
	if err != nil {
		return nil, err
	}
	return &user.AddSecretResponse{
		CreationDate: timestamppb.New(details.EventDate),
		ClientSecret: newSecret.ClientSecret,
	}, nil
}

func (s *Server) RemoveSecret(ctx context.Context, req *user.RemoveSecretRequest) (*user.RemoveSecretResponse, error) {
	details, err := s.command.RemoveMachineSecret(ctx, req.UserId, "", false)
	if err != nil {
		return nil, err
	}
	return &user.RemoveSecretResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}, nil
}
