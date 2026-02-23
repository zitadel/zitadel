package types

import (
	"context"
	"strings"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func (notify Notify) SendDomainClaimed(ctx context.Context, user *query.NotifyUser, username, urlTemplate string) error {
	index := strings.LastIndex(user.LastEmail, "@")
	domainSuffix := user.LastEmail[index+1:]
	var url string
	if urlTemplate == "" {
		url = login.LoginLink(http_utils.DomainContext(ctx).Origin(), user.ResourceOwner)
	} else {
		var buf strings.Builder
		if err := domain.RenderURLTemplate(&buf, urlTemplate, &DomainClaimedData{
			Domain:       domainSuffix,
			TempUsername: username,
			OrgID:        user.ResourceOwner,
		}); err != nil {
			return err
		}
		url = buf.String()
	}
	args := make(map[string]interface{})
	args["TempUsername"] = username
	args["Domain"] = domainSuffix
	return notify(url, args, domain.DomainClaimedMessageType, true)
}

type DomainClaimedData struct {
	TempUsername string
	Domain       string
	OrgID        string
}
