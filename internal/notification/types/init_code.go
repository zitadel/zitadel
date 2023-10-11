package types

import (
	"context"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func (notify Notify) SendUserInitCode(ctx context.Context, user *query.NotifyUser, code string) error {
	url := login.InitUserLink(http_utils.ComposedOrigin(ctx), user.ID, user.PreferredLoginName, code, user.ResourceOwner, user.PasswordSet)
	args := make(map[string]interface{})
	args["Code"] = code
	return notify(url, args, domain.InitCodeMessageType, true)
}
