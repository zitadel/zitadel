package v2

import (
	"context"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func ListSessions(ctx context.Context, request *connect.Request[session.ListSessionsRequest]) (*connect.Response[session.ListSessionsResponse], error) {
	return defaultServer.ListSessions(ctx, request)
}

// ListSessions implements [sessionconnect.SessionServiceHandler].
func (s *server) ListSessions(ctx context.Context, request *connect.Request[session.ListSessionsRequest]) (*connect.Response[session.ListSessionsResponse], error) {
	panic("unimplemented")
}

// func listRequestToDomain(request *session.ListSessionsRequest) (domain.session)
