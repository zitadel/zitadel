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
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(
		t,
		func(tt *assert.CollectT) {
			defer func() {
				// Panics are not recovered and don't mark the test as failed, so we need to do that ourselves
				assert.Nil(tt, recover(), "panic in await callback")
			}()
			cb(tt)
		},
		retryDuration,
		tick,
		"awaiting successful callback failed",
	)
}

var _ assert.TestingT = (*noopAssertionT)(nil)

type noopAssertionT struct{}

func (*noopAssertionT) FailNow() {}

func (*noopAssertionT) Errorf(string, ...interface{}) {}
