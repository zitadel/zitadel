package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/caos/zitadel/internal/api/authz"
)

type InstanceVerifier interface {
	GetInstance(ctx context.Context)
}

func InstanceInterceptor(verifier authz.InstanceVerifier, headerName string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return setInstance(ctx, req, info, handler, verifier, headerName)
	}
}

func setInstance(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, verifier authz.InstanceVerifier, headerName string) (_ interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {

	}
	host, ok := md[headerName]
	if !ok {

	}
	if len(host) != 1 {

	}
	instance, err := verifier.InstanceByHost(ctx, host[0])
	if err != nil {
		return nil, err
	}
	return handler(authz.WithInstance(ctx, instance), req)
	//authOpt, needsToken := verifier.CheckAuthMethod(info.FullMethod)
	//if !needsToken {
	//	return handler(ctx, req)
	//}
	//
	//authCtx, span := tracing.NewServerInterceptorSpan(ctx)
	//defer func() { span.EndWithError(err) }()
	//
	//authToken := grpc_util.GetAuthorizationHeader(authCtx)
	//if authToken == "" {
	//	return nil, status.Error(codes.Unauthenticated, "auth header missing")
	//}
	//
	//orgID := grpc_util.GetHeader(authCtx, http.ZitadelOrgID)
	//
	//ctxSetter, err := authz.CheckUserAuthorization(authCtx, req, authToken, orgID, verifier, authConfig, authOpt, info.FullMethod)
	//if err != nil {
	//	return nil, err
	//}
	//span.End()
	//return handler(ctxSetter(ctx), req)
}
