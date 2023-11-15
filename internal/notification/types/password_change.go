package types

import (
	"context"

	http_utils "github.com/zitadel/zitadel/v2/internal/api/http"
	"github.com/zitadel/zitadel/v2/internal/api/ui/console"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/query"
)

func (notify Notify) SendPasswordChange(ctx context.Context, user *query.NotifyUser) error {
	url := console.LoginHintLink(http_utils.ComposedOrigin(ctx), user.PreferredLoginName)
	args := make(map[string]interface{})
	return notify(url, args, domain.PasswordChangeMessageType, true)
}
