package types

import (
	"context"
	"strings"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func (notify Notify) SendPasswordCode(ctx context.Context, user *query.NotifyUser, code, urlTmpl, authRequestID string, loginPolicy *query.LoginPolicy) error {
	var url string

	if loginPolicy != nil && loginPolicy.DefaultRedirectURI != "" && loginPolicy.UseDefaultUriForNotificationLinks {
		url = loginPolicy.DefaultRedirectURI
	} else {
		if urlTmpl == "" {
			url = login.InitPasswordLink(http_utils.ComposedOrigin(ctx), user.ID, code, user.ResourceOwner, authRequestID)
		} else {
			var buf strings.Builder
			if err := domain.RenderConfirmURLTemplate(&buf, urlTmpl, user.ID, code, user.ResourceOwner); err != nil {
				return err
			}
			url = buf.String()
		}
	}

	args := make(map[string]interface{})
	args["Code"] = code
	return notify(url, args, domain.PasswordResetMessageType, true)
}
