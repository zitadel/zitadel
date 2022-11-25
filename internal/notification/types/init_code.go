package types

import (
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func (notify Notify) SendUserInitCode(user *query.NotifyUser, origin, code string) error {
	url := login.InitUserLink(origin, user.ID, user.PreferredLoginName, code, user.ResourceOwner, user.PasswordSet)
	args := make(map[string]interface{})
	args["Code"] = code
	return notify(url, args, domain.InitCodeMessageType, true)
}
