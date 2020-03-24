package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetUserAgentSessions(ctx context.Context, id *UserAgentID) (_ *UserSessions, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-ue45f", "Not implemented")
}

func (s *Server) GetUserSession(ctx context.Context, id *UserSessionID) (_ *UserSession, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-lor5h", "Not implemented")
}

func (s *Server) TerminateUserSession(ctx context.Context, id *UserSessionID) (_ *empty.Empty, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-bnmt6", "Not implemented")
}

func (s *Server) GetMyUserSessions(ctx context.Context, _ *empty.Empty) (_ *UserSessionViews, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-nc52s", "Not implemented")
}
