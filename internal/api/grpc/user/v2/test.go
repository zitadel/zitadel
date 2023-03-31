package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func (s *Server) TestGet(ctx context.Context, req *user.TestGetRequest) (*user.TestGetResponse, error) {
	return &user.TestGetResponse{
		Ctx: req.Ctx.String(),
	}, nil
}

func (s *Server) TestPost(ctx context.Context, req *user.TestPostRequest) (*user.TestPostResponse, error) {
	return &user.TestPostResponse{
		Ctx: req.Ctx.String(),
	}, nil
}

func (s *Server) TestAuth(ctx context.Context, _ *user.TestAuthRequest) (*user.TestAuthResponse, error) {
	return &user.TestAuthResponse{
		User: &user.User{Id: authz.GetCtxData(ctx).UserID},
	}, nil
}
