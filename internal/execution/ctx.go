package execution

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const ExecutionUserID = "EXECUTION"

func HandlerContext(event *eventstore.Aggregate) context.Context {
	ctx := authz.WithInstanceID(context.Background(), event.InstanceID)
	return authz.SetCtxData(ctx, authz.CtxData{UserID: ExecutionUserID, OrgID: event.ResourceOwner})
}

func ContextWithExecuter(ctx context.Context, aggregate *eventstore.Aggregate) context.Context {
	return authz.SetCtxData(ctx, authz.CtxData{UserID: ExecutionUserID, OrgID: aggregate.ResourceOwner})
}

func (q *ExecutionsQueries) HandlerContext(event *eventstore.Aggregate) (context.Context, error) {
	ctx := context.Background()
	instance, err := q.InstanceByID(ctx, event.InstanceID)
	if err != nil {
		return nil, err
	}
	ctx = authz.WithInstance(ctx, instance)
	return authz.SetCtxData(ctx, authz.CtxData{UserID: ExecutionUserID, OrgID: event.ResourceOwner}), nil
}
