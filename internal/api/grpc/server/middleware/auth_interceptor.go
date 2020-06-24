package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/caos/zitadel/internal/api/authz"
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/http"
)

func AuthorizationInterceptor(verifier authz.TokenVerifier, authConfig *authz.Config, authMethods authz.MethodMapping) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return authorize(ctx, req, info, handler, verifier, authConfig, authMethods)
	}
}

func authorize(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, verifier authz.TokenVerifier, authConfig *authz.Config, authMethods authz.MethodMapping) (interface{}, error) {
	authOpt, needsToken := authMethods[info.FullMethod]
	if !needsToken {
		return handler(ctx, req)
	}

	authToken := grpc_util.GetAuthorizationHeader(ctx)
	if authToken == "" {
		return nil, status.Error(codes.Unauthenticated, "auth header missing")
	}

	orgID := grpc_util.GetHeader(ctx, http.ZitadelOrgID)

	ctx, err := authz.CheckUserAuthorization(ctx, req, authToken, orgID, verifier, authConfig, authOpt)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}
