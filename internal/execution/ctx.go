package execution

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const ExecutionUserID = "EXECUTION"

func ContextWithExecuter(ctx context.Context, aggregate *eventstore.Aggregate) context.Context {
	ctx = authz.WithInstanceID(ctx, aggregate.InstanceID)
	return authz.SetCtxData(ctx, authz.CtxData{UserID: ExecutionUserID, OrgID: aggregate.ResourceOwner})
}
