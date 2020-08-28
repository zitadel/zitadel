package management

import (
	"context"

	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) AddMachineKey(ctx context.Context, req *management.AddMachineKeyRequest) (*management.AddMachineKeyResponse, error) {
	key, err := s.user.AddMachineKey(ctx, addMachineKeyToModel(req))
	if err != nil {
		return nil, err
	}
	return addMachineKeyFromModel(key), nil
}

func (s *Server) DeleteServiceAccountKey(ctx context.Context, req *management.MachineKeyIDRequest) (*empty.Empty, error) {
	err := s.user.RemoveMachineKey(ctx, req.UserId, req.KeyId)
	return &empty.Empty{}, err
}
