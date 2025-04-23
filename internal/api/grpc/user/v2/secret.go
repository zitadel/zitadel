package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) AddSecret(ctx context.Context, req *user.AddSecretRequest) (*user.AddSecretResponse, error) {
	return nil, zerrors.ThrowUnimplemented(nil, "", "not implemented")
}

func (s *Server) RemoveSecret(ctx context.Context, req *user.RemoveSecretRequest) (*user.RemoveSecretResponse, error) {
	return nil, zerrors.ThrowUnimplemented(nil, "", "not implemented")
}
