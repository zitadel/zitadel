package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
)

func (s *Server) SearchAuthGrant(ctx context.Context, grantSearch *AuthGrantSearchRequest) (*AuthGrantSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-dkwd5", "Not implemented")
}
