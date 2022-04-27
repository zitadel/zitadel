package middleware

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/api/authz"
)

type InstanceVerifier interface {
	GetInstance(ctx context.Context)
}

func InstanceInterceptor(verifier authz.InstanceVerifier, headerName string, ignoredServices ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return setInstance(ctx, req, info, handler, verifier, headerName, ignoredServices...)
	}
}

func setInstance(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, verifier authz.InstanceVerifier, headerName string, ignoredServices ...string) (_ interface{}, err error) {
	for _, service := range ignoredServices {
		if strings.HasPrefix(info.FullMethod, service) {
			return handler(ctx, req)
		}
	}

	host, err := hostNameFromContext(ctx, headerName)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	instance, err := verifier.InstanceByHost(ctx, host)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	return handler(authz.WithInstance(ctx, instance), req)
}

func hostNameFromContext(ctx context.Context, headerName string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("cannot read metadata")
	}
	host, ok := md[headerName]
	if !ok {
		return "", fmt.Errorf("cannot find header: %v", headerName)
	}
	if len(host) != 1 {
		return "", fmt.Errorf("invalid host header: %v", host)
	}
	return host[0], nil
}
