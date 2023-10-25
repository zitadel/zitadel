package middleware

import (
	"context"
	"slices"
	"strings"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/activity"
)

func ActivityInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if isResourceAPI(info.FullMethod) {
			activity.TriggerWithContext(ctx, activity.ResourceAPI)
		}
		return resp, err
	}
}

var resourcePrefixes = []string{
	"/zitadel.management.v1.ManagementService/",
	"/zitadel.admin.v1.AdminService/",
	"/zitadel.user.v2beta.UserService/",
	"/zitadel.settings.v2beta.SettingsService/",
	"/zitadel.auth.v1.AuthService/",
}

func isResourceAPI(method string) bool {
	return slices.ContainsFunc(resourcePrefixes, func(prefix string) bool {
		return strings.HasPrefix(method, prefix)
	})
}
