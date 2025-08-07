package user

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) AddKey(ctx context.Context, req *connect.Request[user.AddKeyRequest]) (*connect.Response[user.AddKeyResponse], error) {
	newMachineKey := &command.MachineKey{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.Msg.GetUserId(),
		},
		ExpirationDate:  req.Msg.GetExpirationDate().AsTime(),
		Type:            domain.AuthNKeyTypeJSON,
		PermissionCheck: s.command.NewPermissionCheckUserWrite(ctx),
	}
	newMachineKey.PublicKey = req.Msg.GetPublicKey()

	pubkeySupplied := len(newMachineKey.PublicKey) > 0
	details, err := s.command.AddUserMachineKey(ctx, newMachineKey)
	if err != nil {
		return nil, err
	}
	// Return key details only if the pubkey wasn't supplied, otherwise the user already has
	// private key locally
	var keyDetails []byte
	if !pubkeySupplied {
		var err error
		keyDetails, err = newMachineKey.Detail()
		if err != nil {
			return nil, err
		}
	}
	return connect.NewResponse(&user.AddKeyResponse{
		KeyId:        newMachineKey.KeyID,
		KeyContent:   keyDetails,
		CreationDate: timestamppb.New(details.EventDate),
	}), nil
}

func (s *Server) RemoveKey(ctx context.Context, req *connect.Request[user.RemoveKeyRequest]) (*connect.Response[user.RemoveKeyResponse], error) {
	machineKey := &command.MachineKey{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.Msg.GetUserId(),
		},
		PermissionCheck: s.command.NewPermissionCheckUserWrite(ctx),
		KeyID:           req.Msg.GetKeyId(),
	}
	objectDetails, err := s.command.RemoveUserMachineKey(ctx, machineKey)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.RemoveKeyResponse{
		DeletionDate: timestamppb.New(objectDetails.EventDate),
	}), nil
}
