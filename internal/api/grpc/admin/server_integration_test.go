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
)

var (
	AdminCTX, SystemCTX context.Context
	Tester              *integration.Tester
	// NoopAssertionT is useful in combination with assert.Eventuallyf to use testify assertions in a callback
	NoopAssertionT = new(noopAssertionT)
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, _, cancel := integration.Contexts(3 * time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		AdminCTX = Tester.WithAuthorization(ctx, integration.IAMOwner)
		SystemCTX = Tester.WithAuthorization(ctx, integration.SystemUser)

		return m.Run()
	}())
}

func await(t *testing.T, ctx context.Context, cb func() bool) {
	deadline, ok := ctx.Deadline()
	require.True(t, ok, "context must have deadline")
	assert.Eventuallyf(
		t,
		func() bool {
			defer func() {
				// Panics are not recovered and don't mark the test as failed, so we need to do that ourselves
				require.Nil(t, recover(), "panic in await callback")
			}()
			return cb()
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
