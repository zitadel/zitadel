//go:build integration

package session_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2beta"
)

func TestServer_ListSessions_LimitExceeded(t *testing.T) {
	t.Parallel()

	_, err := Client.ListSessions(IAMOwnerCTX, &session.ListSessionsRequest{
		Query: &object.ListQuery{Limit: 1001},
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}
