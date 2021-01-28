package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"

	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) AddMachineKey(ctx context.Context, req *management.AddMachineKeyRequest) (*management.AddMachineKeyResponse, error) {
	key, err := s.command.AddUserMachineKey(ctx, addMachineKeyToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return addMachineKeyFromDomain(key), nil
}

func (s *Server) DeleteMachineKey(ctx context.Context, req *management.MachineKeyIDRequest) (*empty.Empty, error) {
	err := s.command.RemoveUserMachineKey(ctx, req.UserId, req.KeyId, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) GetMachineKey(ctx context.Context, req *management.MachineKeyIDRequest) (*management.MachineKeyView, error) {
	key, err := s.user.GetMachineKey(ctx, req.UserId, req.KeyId)
	if err != nil {
		return nil, err
	}
	return machineKeyViewFromModel(key), nil
}

func (s *Server) SearchMachineKeys(ctx context.Context, req *management.MachineKeySearchRequest) (*management.MachineKeySearchResponse, error) {
	result, err := s.user.SearchMachineKeys(ctx, machineKeySearchRequestToModel(req))
	if err != nil {
		return nil, err
	}
	return machineKeySearchResponseFromModel(result), nil
}
