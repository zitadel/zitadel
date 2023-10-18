package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/activity"
)

func ActivityInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if isResourceAPI(info.FullMethod) {
			activity.TriggerWithContext(ctx, info.FullMethod, activity.ResourceAPI)
		}
		return resp, err
	}
}

func isResourceAPI(method string) bool {
	if strings.HasPrefix(method, "/zitadel.management.v1.ManagementService/") ||
		strings.HasPrefix(method, "/zitadel.admin.v1.AdminService/") {
		return true
	}
	return false
}
