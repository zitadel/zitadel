package handlers

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const NotifyUserID = "NOTIFICATION" //TODO: system?

func HandlerContext(event *eventstore.Aggregate) context.Context {
	ctx := authz.WithInstanceID(context.Background(), event.InstanceID)
	return authz.SetCtxData(ctx, authz.CtxData{UserID: NotifyUserID, OrgID: event.ResourceOwner})
}

func (n *NotificationQueries) HandlerContext(event *eventstore.Aggregate) (context.Context, error) {
	ctx := context.Background()
	instance, err := n.InstanceByID(ctx, event.InstanceID)
	if err != nil {
		return nil, err
	}
	ctx = authz.WithInstance(ctx, instance)
	return authz.SetCtxData(ctx, authz.CtxData{UserID: NotifyUserID, OrgID: event.ResourceOwner}), nil
}
