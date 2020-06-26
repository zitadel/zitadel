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

func AuthorizationInterceptor(verifier *authz.TokenVerifier2, authConfig authz.Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return authorize(ctx, req, info, handler, verifier, authConfig)
	}
}

func authorize(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, verifier *authz.TokenVerifier2, authConfig authz.Config) (interface{}, error) {
	authOpt, needsToken := verifier.CheckAuthMethod(info.FullMethod)
	if !needsToken {
		return handler(ctx, req)
	}

	authToken := grpc_util.GetAuthorizationHeader(ctx)
	if authToken == "" {
		return nil, status.Error(codes.Unauthenticated, "auth header missing")
	}

	orgID := grpc_util.GetHeader(ctx, http.ZitadelOrgID)

	ctx, err := authz.CheckUserAuthorization(ctx, req, authToken, orgID, verifier, authConfig, authOpt, info.FullMethod)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}
