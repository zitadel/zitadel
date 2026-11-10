//go:build integration

package user_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_ListUsers_LimitExceeded(t *testing.T) {
	t.Parallel()

	_, err := Client.ListUsers(IamCTX, &user.ListUsersRequest{
		Query: &object.ListQuery{Limit: 1001},
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestServer_ListIDPLinks_LimitExceeded(t *testing.T) {
	t.Parallel()

	_, err := Client.ListIDPLinks(IamCTX, &user.ListIDPLinksRequest{
		UserId: Instance.Users.Get(integration.UserTypeLogin).ID,
		Query:  &object.ListQuery{Limit: 1001},
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}
