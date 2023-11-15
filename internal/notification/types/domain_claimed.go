package types

import (
	"context"
	"strings"

	http_utils "github.com/zitadel/zitadel/v2/internal/api/http"
	"github.com/zitadel/zitadel/v2/internal/api/ui/login"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/query"
)

func (notify Notify) SendDomainClaimed(ctx context.Context, user *query.NotifyUser, username string) error {
	url := login.LoginLink(http_utils.ComposedOrigin(ctx), user.ResourceOwner)
	index := strings.LastIndex(user.LastEmail, "@")
	args := make(map[string]interface{})
	args["TempUsername"] = username
	args["Domain"] = user.LastEmail[index+1:]
	return notify(url, args, domain.DomainClaimedMessageType, true)
}
