package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
)

func (s *Server) CreateToken(ctx context.Context, request *CreateTokenRequest) (_ *Token, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-86dG3", "Not implemented")
}
