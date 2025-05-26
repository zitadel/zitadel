package user

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) AddKey(ctx context.Context, req *user.AddKeyRequest) (*user.AddKeyResponse, error) {
	newMachineKey := &command.MachineKey{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.UserId,
		},
		ExpirationDate:  req.GetExpirationDate().AsTime(),
		Type:            domain.AuthNKeyTypeJSON,
		PermissionCheck: s.command.NewPermissionCheckUserWrite(ctx),
	}
	newMachineKey.PublicKey = req.PublicKey

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
	return &user.AddKeyResponse{
		KeyId:        newMachineKey.KeyID,
		KeyContent:   keyDetails,
		CreationDate: timestamppb.New(details.EventDate),
	}, nil
}

func (s *Server) RemoveKey(ctx context.Context, req *user.RemoveKeyRequest) (*user.RemoveKeyResponse, error) {
	machineKey := &command.MachineKey{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.UserId,
		},
		PermissionCheck: s.command.NewPermissionCheckUserWrite(ctx),
		KeyID:           req.KeyId,
	}
	objectDetails, err := s.command.RemoveUserMachineKey(ctx, machineKey)
	if err != nil {
		return nil, err
	}
	return &user.RemoveKeyResponse{
		DeletionDate: timestamppb.New(objectDetails.EventDate),
	}, nil
}
