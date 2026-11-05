package types

import (
	"context"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/ui/console"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

// SendEmailChange notifies the user that their email address was changed. Uses the
// verified email channel (allowUnverifiedNotificationChannel: false), deliberately not
// NotifyUser.LastEmail -- the user projection updates LastEmail to the *new* address as
// part of this very event, so allowing the unverified channel would send the security
// alert to the new (and potentially attacker-controlled) address instead of the still-
// verified old one that actually needs to see it.
func (notify Notify) SendEmailChange(ctx context.Context, user *query.NotifyUser) error {
	url := console.LoginHintLink(http_utils.DomainContext(ctx).Origin(), user.PreferredLoginName)
	args := make(map[string]interface{})
	return notify(url, args, domain.EmailChangeMessageType, false)
}
