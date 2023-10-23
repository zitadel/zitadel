package activity

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/info"
)

const (
	Activity = "activity"
)

type TriggerMethod int

const (
	Unspecified TriggerMethod = iota
	ResourceAPI
	OIDCAccessToken
	OIDCRefreshToken
	SessionAPI
	SAMLResponse
)

func (t TriggerMethod) String() string {
	switch t {
	case Unspecified:
		return "unspecified"
	case ResourceAPI:
		return "resourceAPI"
	case OIDCRefreshToken:
		return "refreshToken"
	case OIDCAccessToken:
		return "accessToken"
	case SessionAPI:
		return "sessionAPI"
	case SAMLResponse:
		return "samlResponse"
	default:
		return "unknown"
	}
}

func Trigger(ctx context.Context, orgID, userID string, trigger TriggerMethod) {
	path, _ := info.HTTPPathFromContext()(ctx)
	reqMethod, _ := info.RequestMethodFromContext()(ctx)
	method, _ := info.RPCMethodFromContext()(ctx)
	logging.WithFields(
		"instance", authz.GetInstance(ctx).InstanceID(),
		"org", orgID,
		"user", userID,
		"domain", http_utils.ComposedOrigin(ctx),
		"trigger", trigger.String(),
		"method", method,
		"path", path,
		"requestMethod", reqMethod,
	).Info(Activity)
}

func TriggerWithContext(ctx context.Context, trigger TriggerMethod) {
	data := authz.GetCtxData(ctx)
	path, _ := info.HTTPPathFromContext()(ctx)
	reqMethod, _ := info.RequestMethodFromContext()(ctx)
	method, _ := info.RPCMethodFromContext()(ctx)
	// if GRPC call, path is prefilled with the grpc fullmethod and method is empty
	if method == "" {
		method = path
		path = ""
	}
	logging.WithFields(
		"instance", authz.GetInstance(ctx).InstanceID(),
		"org", data.OrgID,
		"user", data.UserID,
		"domain", http_utils.ComposedOrigin(ctx),
		"trigger", trigger.String(),
		"method", method,
		"path", path,
		"requestMethod", reqMethod,
	).Info(Activity)
}
