package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/caos/zitadel/internal/api"
	"github.com/caos/zitadel/internal/api/auth"
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
)

func AuthorizationInterceptor(verifier auth.TokenVerifier, authConfig *auth.Config, authMethods auth.MethodMapping) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		authOpt, needsToken := authMethods[info.FullMethod]
		if !needsToken {
			return handler(ctx, req)
		}

		authToken := grpc_util.GetAuthorizationHeader(ctx)
		if authToken == "" {
			return nil, status.Error(codes.Unauthenticated, "auth header missing")
		}

		orgID := grpc_util.GetHeader(ctx, api.ZitadelOrgID)

		ctx, err := auth.CheckUserAuthorization(ctx, req, authToken, orgID, verifier, authConfig, authOpt)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
