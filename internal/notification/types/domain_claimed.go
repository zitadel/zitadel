package types

import (
	"strings"

	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func (notify Notify) SendDomainClaimed(user *query.NotifyUser, origin, username string) error {
	url := login.LoginLink(origin, user.ResourceOwner)
	index := strings.LastIndex(user.LastEmail, "@")
	args := make(map[string]interface{})
	args["TempUsername"] = username
	args["Domain"] = user.LastEmail[index+1:]
	return notify(url, args, domain.DomainClaimedMessageType, true)
}
