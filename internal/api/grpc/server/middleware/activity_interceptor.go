package middleware

import (
	"context"
	"slices"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/zitadel/internal/activity"
	"github.com/zitadel/zitadel/internal/api/grpc/gerrors"
	ainfo "github.com/zitadel/zitadel/internal/api/info"
)

func ActivityInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = activityInfoFromGateway(ctx).SetMethod(info.FullMethod).IntoContext(ctx)
		resp, err := handler(ctx, req)
		if isResourceAPI(info.FullMethod) {
			code, _, _, _ := gerrors.ExtractZITADELError(err)
			ctx = ainfo.ActivityInfoFromContext(ctx).SetGRPCStatus(code).IntoContext(ctx)
			activity.TriggerGRPCWithContext(ctx, activity.ResourceAPI)
		}
		return resp, err
	}
}

var resourcePrefixes = []string{
	"/zitadel.management.v1.ManagementService/",
	"/zitadel.admin.v1.AdminService/",
	"/zitadel.user.v2.UserService/",
	"/zitadel.settings.v2.SettingsService/",
	"/zitadel.auth.v1.AuthService/",
}

func isResourceAPI(method string) bool {
	return slices.ContainsFunc(resourcePrefixes, func(prefix string) bool {
		return strings.HasPrefix(method, prefix)
	})
}

func activityInfoFromGateway(ctx context.Context) *ainfo.ActivityInfo {
	info := ainfo.ActivityInfoFromContext(ctx)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return info
	}
	path := md.Get(activity.PathKey)
	if len(path) != 1 {
		return info
	}
	requestMethod := md.Get(activity.RequestMethodKey)
	if len(requestMethod) != 1 {
		return info
	}
	return info.SetPath(path[0]).SetRequestMethod(requestMethod[0])
}
