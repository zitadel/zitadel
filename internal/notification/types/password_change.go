package types

import (
	"github.com/zitadel/zitadel/internal/api/ui/console"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func (notify Notify) SendPasswordChange(user *query.NotifyUser, origin, username string) error {
	url := console.LoginHintLink(origin, username)
	args := make(map[string]interface{})
	return notify(url, args, domain.PasswordChangeMessageType, true)
}
