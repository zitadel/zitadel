package handlers

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const NotifyUserID = "NOTIFICATION" //TODO: system?

func HandlerContext(parent context.Context, event *eventstore.Aggregate) context.Context {
	ctx := authz.WithInstanceID(parent, event.InstanceID)
	return authz.SetCtxData(ctx, authz.CtxData{UserID: NotifyUserID, OrgID: event.ResourceOwner})
}

func ContextWithNotifier(ctx context.Context, aggregate *eventstore.Aggregate) context.Context {
	return authz.WithInstanceID(authz.SetCtxData(ctx, authz.CtxData{UserID: NotifyUserID, OrgID: aggregate.ResourceOwner}), aggregate.InstanceID)
}

func (n *NotificationQueries) HandlerContext(parent context.Context, event *eventstore.Aggregate) (context.Context, error) {
	instance, err := n.InstanceByID(parent, event.InstanceID)
	if err != nil {
		return nil, err
	}
	ctx := authz.WithInstance(parent, instance)
	return authz.SetCtxData(ctx, authz.CtxData{UserID: NotifyUserID, OrgID: event.ResourceOwner}), nil
}
