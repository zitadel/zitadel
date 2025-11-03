package execution

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const ExecutionUserID = "EXECUTION"

func HandlerContext(parent context.Context, event *eventstore.Aggregate) context.Context {
	ctx := authz.WithInstanceID(parent, event.InstanceID)
	return authz.SetCtxData(ctx, authz.CtxData{UserID: ExecutionUserID, OrgID: event.ResourceOwner})
}

func ContextWithExecuter(ctx context.Context, aggregate *eventstore.Aggregate) context.Context {
	return authz.SetCtxData(ctx, authz.CtxData{UserID: ExecutionUserID, OrgID: aggregate.ResourceOwner})
}
