package handlers

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
)

const NotifyUserID = "NOTIFICATION" //TODO: system?

func HandlerContext(event *eventstore.Aggregate) context.Context {
	ctx := authz.WithInstanceID(context.Background(), event.InstanceID)
	return authz.SetCtxData(ctx, authz.CtxData{UserID: NotifyUserID, OrgID: event.ResourceOwner})
}
