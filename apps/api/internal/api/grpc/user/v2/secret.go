package user

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) AddSecret(ctx context.Context, req *connect.Request[user.AddSecretRequest]) (*connect.Response[user.AddSecretResponse], error) {
	newSecret := &command.GenerateMachineSecret{
		PermissionCheck: s.command.NewPermissionCheckUserWrite(ctx),
	}
	details, err := s.command.GenerateMachineSecret(ctx, req.Msg.GetUserId(), "", newSecret)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.AddSecretResponse{
		CreationDate: timestamppb.New(details.EventDate),
		ClientSecret: newSecret.ClientSecret,
	}), nil
}

func (s *Server) RemoveSecret(ctx context.Context, req *connect.Request[user.RemoveSecretRequest]) (*connect.Response[user.RemoveSecretResponse], error) {
	details, err := s.command.RemoveMachineSecret(
		ctx,
		req.Msg.GetUserId(),
		"",
		s.command.NewPermissionCheckUserWrite(ctx),
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.RemoveSecretResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}), nil
}
