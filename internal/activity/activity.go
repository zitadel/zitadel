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
	data := authz.GetCtxData(ctx)
	triggerLog(
		authz.GetInstance(ctx).InstanceID(),
		orgID,
		userID,
		http_utils.ComposedOrigin(ctx),
		trigger,
		ai.Method,
		ai.Path,
		ai.RequestMethod,
		data.SystemMemberships != nil,
	)
}

func TriggerWithContext(ctx context.Context, trigger TriggerMethod) {
	Trigger(ctx, authz.GetCtxData(ctx).OrgID, authz.GetCtxData(ctx).UserID, trigger)
}

func triggerLog(instanceID, orgID, userID, domain string, trigger TriggerMethod, method, path, requestMethod string, isSystemUser bool) {
	logging.WithFields(
		"instance", instanceID,
		"org", orgID,
		"user", userID,
		"domain", domain,
		"trigger", trigger.String(),
		"method", method,
		"path", path,
		"requestMethod", requestMethod,
		"isSystemUser", isSystemUser,
	).Info(Activity)
}
