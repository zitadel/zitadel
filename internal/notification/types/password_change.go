package types

import (
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func (notify Notify) SendPasswordChange(user *query.NotifyUser, origin string) error {
	url := login.LoginLink(origin, user.ResourceOwner)
	args := make(map[string]interface{})
	return notify(url, args, domain.PasswordChangeMessageType, true)
}
