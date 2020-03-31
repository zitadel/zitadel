package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetMyUserSessions(ctx context.Context, _ *empty.Empty) (_ *UserSessionViews, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-nc52s", "Not implemented")
}
