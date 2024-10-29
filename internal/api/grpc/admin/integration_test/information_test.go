//go:build integration

package admin_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

func TestServer_Healthz(t *testing.T) {
	ctx, cancel := context.WithTimeout(AdminCTX, time.Minute)
	defer cancel()
	_, err := Instance.Client.Admin.Healthz(ctx, &admin.HealthzRequest{})
	require.NoError(t, err)
}
