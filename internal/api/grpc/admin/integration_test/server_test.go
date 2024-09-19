//go:build integration

package admin_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/v2/internal/integration"
	admin_pb "github.com/zitadel/zitadel/v2/pkg/grpc/admin"
)

var (
	CTX, AdminCTX context.Context
	Instance      *integration.Instance
	Client        admin_pb.AdminServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)
		CTX = ctx
		AdminCTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		Client = Instance.Client.Admin
		return m.Run()
	}())
}

func await(t *testing.T, ctx context.Context, cb func(*assert.CollectT)) {
	deadline, ok := ctx.Deadline()
	require.True(t, ok, "context must have deadline")
	require.EventuallyWithT(
		t,
		func(tt *assert.CollectT) {
			defer func() {
				// Panics are not recovered and don't mark the test as failed, so we need to do that ourselves
				require.Nil(t, recover(), "panic in await callback")
			}()
			cb(tt)
		},
		time.Until(deadline),
		time.Second,
		"awaiting successful callback failed",
	)
}

var _ assert.TestingT = (*noopAssertionT)(nil)

type noopAssertionT struct{}

func (*noopAssertionT) FailNow() {}

func (*noopAssertionT) Errorf(string, ...interface{}) {}
