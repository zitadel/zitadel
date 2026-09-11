package execution_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/execution"
)

func TestContextWithExecuter(t *testing.T) {
	ctx := execution.ContextWithExecuter(context.Background(), &eventstore.Aggregate{
		InstanceID:    instanceID,
		ResourceOwner: orgID,
	})

	assert.Equal(t, instanceID, authz.GetInstance(ctx).InstanceID())
	assert.Equal(t, orgID, authz.GetCtxData(ctx).OrgID)
	assert.Equal(t, execution.ExecutionUserID, authz.GetCtxData(ctx).UserID)
}
