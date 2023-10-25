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
	ai := info.ActivityInfoFromContext(ctx)
	triggerLog(authz.GetInstance(ctx).InstanceID(), orgID, userID, http_utils.ComposedOrigin(ctx), trigger, ai.Method, ai.Path, ai.RequestMethod)
}

func TriggerWithContext(ctx context.Context, trigger TriggerMethod) {
	data := authz.GetCtxData(ctx)
	ai := info.ActivityInfoFromContext(ctx)
	// if GRPC call, path is prefilled with the grpc fullmethod and method is empty
	triggerLog(authz.GetInstance(ctx).InstanceID(), data.OrgID, data.UserID, http_utils.ComposedOrigin(ctx), trigger, ai.Path, "", ai.RequestMethod)
}

func triggerLog(instanceID, orgID, userID, domain string, trigger TriggerMethod, method, path, requestMethod string) {
	logging.WithFields(
		"instance", instanceID,
		"org", orgID,
		"user", userID,
		"domain", domain,
		"trigger", trigger.String(),
		"method", method,
		"path", path,
		"requestMethod", requestMethod,
	).Info(Activity)
}
