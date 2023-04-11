package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/errors"
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

func (s *Server) TestAuth(ctx context.Context, req *user.TestAuthRequest) (*user.TestAuthResponse, error) {
	reqCtx, err := authDemo(ctx, req.Ctx)
	if err != nil {
		return nil, err
	}
	return &user.TestAuthResponse{
		User: &user.User{Id: authz.GetCtxData(ctx).UserID},
		Ctx:  reqCtx,
	}, nil
}

func authDemo(ctx context.Context, reqCtx *user.Context) (*user.Context, error) {
	ro := authz.GetCtxData(ctx).ResourceOwner
	if reqCtx == nil {
		return &user.Context{Ctx: &user.Context_OrgId{OrgId: ro}}, nil
	}
	switch c := reqCtx.Ctx.(type) {
	case *user.Context_OrgId:
		if c.OrgId == ro {
			return reqCtx, nil
		}
		return nil, errors.ThrowPermissionDenied(nil, "USER-dg4g", "Errors.User.NotAllowedOrg")
	case *user.Context_OrgDomain:
		if c.OrgDomain == "forbidden.com" {
			return nil, errors.ThrowPermissionDenied(nil, "USER-SDg4g", "Errors.User.NotAllowedOrg")
		}
		return reqCtx, nil
	case *user.Context_Instance:
		return reqCtx, nil
	default:
		return reqCtx, nil
	}
}
