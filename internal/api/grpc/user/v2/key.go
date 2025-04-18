package user

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) AddKey(ctx context.Context, req *user.AddKeyRequest) (*user.AddKeyResponse, error) {
	machineKey := &command.MachineKey{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.UserId,
		},
		PublicKey:      req.PublicKey,
		Type:           domain.AuthNKeyTypeJSON,
		ExpirationDate: req.ExpirationDate.AsTime(),
	}
	pubkeySupplied := len(machineKey.PublicKey) > 0
	details, err := s.command.AddUserMachineKey(ctx, machineKey, false)
	if err != nil {
		return nil, err
	}
	// Return key details only if the pubkey wasn't supplied, otherwise the user already has
	// private key locally
	var keyDetails []byte
	if !pubkeySupplied {
		var err error
		keyDetails, err = machineKey.Detail()
		if err != nil {
			return nil, err
		}
	}
	return &user.AddKeyResponse{
		KeyId:        machineKey.KeyID,
		KeyContent:   keyDetails,
		CreationDate: timestamppb.New(details.EventDate),
	}, nil
}

func (s *Server) RemoveKey(ctx context.Context, req *user.RemoveKeyRequest) (*user.RemoveKeyResponse, error) {
	objectDetails, err := s.command.RemoveUserMachineKey(ctx, &command.MachineKey{KeyID: req.KeyId}, false, false)
	if err != nil {
		return nil, err
	}
	return &user.RemoveKeyResponse{
		DeletionDate: timestamppb.New(objectDetails.EventDate),
	}, nil
}

func (s *Server) ListKeys(ctx context.Context, req *user.ListKeysRequest) (*user.ListKeysResponse, error) {
	return nil, zerrors.ThrowUnimplemented(nil, "", "not implemented")
}
