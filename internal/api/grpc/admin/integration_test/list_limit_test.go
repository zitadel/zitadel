//go:build integration

package admin_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/object"
)

func TestServer_ListOrgs_LimitExceeded(t *testing.T) {
	t.Parallel()

	_, err := Client.ListOrgs(AdminCTX, &admin_pb.ListOrgsRequest{
		Query: &object.ListQuery{Limit: 1001},
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}
