//go:build integration

package system_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/object"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func TestServer_ListInstances_LimitExceeded(t *testing.T) {
	t.Parallel()

	_, err := integration.SystemClient().ListInstances(CTX, &system.ListInstancesRequest{
		Query: &object.ListQuery{Limit: 1001},
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}
