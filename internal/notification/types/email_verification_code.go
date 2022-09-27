package types

import (
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func (notify Notify) SendEmailVerificationCode(user *query.NotifyUser, origin, code string) error {
	url := login.MailVerificationLink(origin, user.ID, code, user.ResourceOwner)
	args := make(map[string]interface{})
	args["Code"] = code
	return notify(url, args, domain.VerifyEmailMessageType, true)
}
