package types

import (
	"context"
	"strings"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func (notify Notify) SendPasswordlessRegistrationLink(ctx context.Context, user *query.NotifyUser, code, codeID, urlTmpl string) error {
	var url string
	if urlTmpl == "" {
		url = domain.PasswordlessInitCodeLink(http_utils.DomainContext(ctx).Origin()+login.HandlerPrefix+login.EndpointPasswordlessRegistration, user.ID, user.ResourceOwner, codeID, code)
	} else {
		var buf strings.Builder
		if err := domain.RenderPasskeyURLTemplate(&buf, urlTmpl, user.ID, user.ResourceOwner, codeID, code); err != nil {
			return err
		}
		url = buf.String()
	}
	return notify(url, nil, domain.PasswordlessRegistrationMessageType, true)
}
