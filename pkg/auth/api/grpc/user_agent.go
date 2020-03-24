package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
)

func (s *Server) GetUserAgent(ctx context.Context, request *UserAgentID) (_ *UserAgent, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-iu7jF", "Not implemented")
}

func (s *Server) CreateUserAgent(ctx context.Context, request *UserAgentCreation) (_ *UserAgent, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-sdfk3", "Not implemented")
}

func (s *Server) RevokeUserAgent(ctx context.Context, id *UserAgentID) (_ *UserAgent, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-09HjK", "Not implemented")
}
