package types

import (
	"context"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/ui/console"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

// SendPhoneChange notifies the user that their phone number was changed. This
// notification is always delivered by email (not SMS), so -- as with SendEmailChange --
// allowUnverifiedNotificationChannel is deliberately false to always use the verified
// email channel.
func (notify Notify) SendPhoneChange(ctx context.Context, user *query.NotifyUser) error {
	url := console.LoginHintLink(http_utils.DomainContext(ctx).Origin(), user.PreferredLoginName)
	args := make(map[string]interface{})
	return notify(url, args, domain.PhoneChangeMessageType, false)
}
