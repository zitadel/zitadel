package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) AddPersonalAccessToken(ctx context.Context, req *user.AddPersonalAccessTokenRequest) (*user.AddPersonalAccessTokenResponse, error) {
	return nil, zerrors.ThrowUnimplemented(nil, "", "not implemented")
}

func (s *Server) RemovePersonalAccessToken(ctx context.Context, req *user.RemovePersonalAccessTokenRequest) (*user.RemovePersonalAccessTokenResponse, error) {
	return nil, zerrors.ThrowUnimplemented(nil, "", "not implemented")
}

func (s *Server) ListPersonalAccessTokens(ctx context.Context, req *user.ListPersonalAccessTokensRequest) (*user.ListPersonalAccessTokensResponse, error) {
	return nil, zerrors.ThrowUnimplemented(nil, "", "not implemented")
}
