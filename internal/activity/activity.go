package activity

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
)

const (
	Activity = "activity"
)

type TriggerMethod int

const (
	Unspecified TriggerMethod = iota
	ResourceAPI
	Authentication
	RefreshToken
	SessionAPI
)

func (t TriggerMethod) String() string {
	switch t {
	case Unspecified:
		return "unspecified"
	case ResourceAPI:
		return "resourceAPI"
	case RefreshToken:
		return "refreshToken"
	case Authentication:
		return "authentication"
	case SessionAPI:
		return "sessionAPI"
	default:
		return "unknown"
	}
}

func Trigger(ctx context.Context, orgID, userID string, trigger TriggerMethod) {
	logging.WithFields(
		"instance", authz.GetInstance(ctx).InstanceID(),
		"org", orgID,
		"user", userID,
		"domain", http_utils.ComposedOrigin(ctx),
		"trigger", trigger.String(),
	).Info(Activity)
}

func TriggerWithContext(ctx context.Context, method string, trigger TriggerMethod) {
	data := authz.GetCtxData(ctx)
	logging.WithFields(
		"instance", authz.GetInstance(ctx).InstanceID(),
		"org", data.OrgID,
		"user", data.UserID,
		"domain", http_utils.ComposedOrigin(ctx),
		"trigger", trigger.String(),
		"method", method,
	).Info(Activity)
}
