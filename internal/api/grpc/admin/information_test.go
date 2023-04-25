package admin_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

func TestServer_Healthz(t *testing.T) {
	client := admin.NewAdminServiceClient(Tester.ClientConn)
	_, err := client.Healthz(context.TODO(), &admin.HealthzRequest{})
	require.NoError(t, err)
}
