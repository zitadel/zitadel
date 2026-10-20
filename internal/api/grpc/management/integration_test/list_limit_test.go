//go:build integration

package management_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/pkg/grpc/change"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/object"
)

func TestServer_ListUsers_LimitExceeded(t *testing.T) {
	t.Parallel()

	_, err := Client.ListUsers(OrgCTX, &mgmt_pb.ListUsersRequest{
		Query: &object.ListQuery{Limit: 1001},
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestServer_ListProjects_LimitExceeded(t *testing.T) {
	t.Parallel()

	_, err := Client.ListProjects(OrgCTX, &mgmt_pb.ListProjectsRequest{
		Query: &object.ListQuery{Limit: 1001},
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestServer_ListOrgChanges_LimitExceeded(t *testing.T) {
	t.Parallel()

	_, err := Client.ListOrgChanges(OrgCTX, &mgmt_pb.ListOrgChangesRequest{
		Query: &change.ChangeQuery{Limit: 1001},
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}
