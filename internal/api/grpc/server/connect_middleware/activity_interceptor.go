package connect_middleware

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/activity"
	"github.com/zitadel/zitadel/internal/api/grpc/gerrors"
	ainfo "github.com/zitadel/zitadel/internal/api/info"
)

func ActivityInterceptor() connect.UnaryInterceptorFunc {
	return func(handler connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			ctx = activityInfoFromGateway(ctx, req.Header()).SetMethod(req.Spec().Procedure).IntoContext(ctx)
			resp, err := handler(ctx, req)
			if isResourceAPI(req.Spec().Procedure) {
				code, _, _, _ := gerrors.ExtractZITADELError(err)
				ctx = ainfo.ActivityInfoFromContext(ctx).SetGRPCStatus(code).IntoContext(ctx)
				activity.TriggerGRPCWithContext(ctx, activity.ResourceAPI)
			}
			return resp, err
		}
	}
}

var resourcePrefixes = []string{
	"/zitadel.management.v1.ManagementService/",
	"/zitadel.admin.v1.AdminService/",
	"/zitadel.user.v2.UserService/",
	"/zitadel.settings.v2.SettingsService/",
	"/zitadel.user.v2beta.UserService/",
	"/zitadel.settings.v2beta.SettingsService/",
	"/zitadel.auth.v1.AuthService/",
}

func isResourceAPI(method string) bool {
	return slices.ContainsFunc(resourcePrefixes, func(prefix string) bool {
		return strings.HasPrefix(method, prefix)
	})
}

func activityInfoFromGateway(ctx context.Context, headers http.Header) *ainfo.ActivityInfo {
	info := ainfo.ActivityInfoFromContext(ctx)
	path := headers.Get(activity.PathKey)
	requestMethod := headers.Get(activity.RequestMethodKey)
	return info.SetPath(path).SetRequestMethod(requestMethod)
}
