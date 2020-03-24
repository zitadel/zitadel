package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
)

func (s *Server) CreateAuthSession(ctx context.Context, request *AuthSessionCreation) (_ *AuthSessionResponse, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-dh3Rt", "Not implemented")
}

func (s *Server) GetAuthSession(ctx context.Context, id *AuthSessionID) (_ *AuthSessionResponse, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-dk56g", "Not implemented")
}

func (s *Server) SelectUser(ctx context.Context, request *SelectUserRequest) (_ *AuthSessionResponse, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-dl5gs", "Not implemented")
}

func (s *Server) VerifyUser(ctx context.Context, request *VerifyUserRequest) (_ *AuthSessionResponse, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-39dGs", "Not implemented")
}

func (s *Server) VerifyPassword(ctx context.Context, password *VerifyPasswordRequest) (_ *AuthSessionResponse, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-tu9j2", "Not implemented")
}

func (s *Server) VerifyMfa(ctx context.Context, mfa *VerifyMfaRequest) (_ *AuthSessionResponse, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-oi9GB", "Not implemented")
}

func (s *Server) GetAuthSessionByTokenID(ctx context.Context, id *TokenID) (_ *AuthSessionView, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-dk56z", "Not implemented")
}
