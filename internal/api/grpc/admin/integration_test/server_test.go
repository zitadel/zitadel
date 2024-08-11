//go:build integration

package admin_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

var (
	CTX, AdminCTX, SystemCTX context.Context
	Instance                 *integration.Instance
	Client                   admin_pb.AdminServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var err error
		Instance, err = integration.GetInstance(ctx)
		if err != nil {
			panic(err)
		}

		CTX = ctx
		AdminCTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		SystemCTX = Instance.WithAuthorization(ctx, integration.UserTypeSystem)
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
		100*time.Millisecond,
		"awaiting successful callback failed",
	)
}

var _ assert.TestingT = (*noopAssertionT)(nil)

type noopAssertionT struct{}

func (*noopAssertionT) FailNow() {}

func (*noopAssertionT) Errorf(string, ...interface{}) {}
