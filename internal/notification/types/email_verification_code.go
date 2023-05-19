package types

import (
	"strings"

	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func (notify Notify) SendEmailVerificationCode(user *query.NotifyUser, origin, code string, urlTmpl string) error {
	var url string
	if urlTmpl == "" {
		url = login.MailVerificationLink(origin, user.ID, code, user.ResourceOwner)
	} else {
		var buf strings.Builder
		if err := domain.RenderConfirmURLTemplate(&buf, urlTmpl, user.ID, code, user.ResourceOwner); err != nil {
			return err
		}
		url = buf.String()
	}

	args := make(map[string]interface{})
	args["Code"] = code
	return notify(url, args, domain.VerifyEmailMessageType, true)
}
