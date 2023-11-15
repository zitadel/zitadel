package types

import (
	"context"
	"strings"

	http_utils "github.com/zitadel/zitadel/v2/internal/api/http"
	"github.com/zitadel/zitadel/v2/internal/api/ui/login"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/query"
)

func (notify Notify) SendPasswordCode(ctx context.Context, user *query.NotifyUser, code, urlTmpl string) error {
	var url string
	if urlTmpl == "" {
		url = login.InitPasswordLink(http_utils.ComposedOrigin(ctx), user.ID, code, user.ResourceOwner)
	} else {
		var buf strings.Builder
		if err := domain.RenderConfirmURLTemplate(&buf, urlTmpl, user.ID, code, user.ResourceOwner); err != nil {
			return err
		}
		url = buf.String()
	}
	args := make(map[string]interface{})
	args["Code"] = code
	return notify(url, args, domain.PasswordResetMessageType, true)
}
