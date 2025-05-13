package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) AddKey(ctx context.Context, req *user.AddKeyRequest) (*user.AddKeyResponse, error) {
	return nil, zerrors.ThrowUnimplemented(nil, "", "not implemented")
}

func (s *Server) RemoveKey(ctx context.Context, req *user.RemoveKeyRequest) (*user.RemoveKeyResponse, error) {
	return nil, zerrors.ThrowUnimplemented(nil, "", "not implemented")
}

func (s *Server) ListKeys(ctx context.Context, req *user.ListKeysRequest) (*user.ListKeysResponse, error) {
	return nil, zerrors.ThrowUnimplemented(nil, "", "not implemented")
}
